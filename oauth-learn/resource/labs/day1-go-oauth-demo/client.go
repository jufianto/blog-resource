package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type ClientApp struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

type Session struct {
	ID           string
	State        string
	CodeVerifier string
	AccessToken  string
	RefreshToken string
	LastEvent    string
}

func NewClientApp() *ClientApp {
	return &ClientApp{sessions: make(map[string]*Session)}
}

func (c *ClientApp) ResetSessions() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessions = make(map[string]*Session)
}

func (c *ClientApp) Home(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	body := fmt.Sprintf(`
<p>This single Go app simulates three OAuth parts using <strong>go-oidc</strong>:
<strong>Client App</strong>, <strong>Authorization Server</strong>, and <strong>Resource Server API</strong>.</p>

<h2>1. Authorization Code + PKCE</h2>
<p>User (Alya Developer) connects RepoBoard to read a GitHub-like profile.</p>
<p><a class="button" href="/start">Start Auth Code + PKCE flow</a></p>
<p><a class="button secondary" href="/call-api">Call resource API with current access token</a></p>
<p><a class="button secondary" href="/refresh">Use refresh token</a></p>

<h2>2. Client Credentials</h2>
<p>Backend service gets its own token. No user involved.</p>
<p><a class="button" href="/client-credentials">Run Client Credentials demo</a></p>

<h2>3. Device Authorization</h2>
<p>CLI or TV-like app starts authorization; user approves in browser.</p>
<p><a class="button" href="/device/start">Start Device Authorization demo</a></p>

<h2>Current Client Session</h2>
<pre>%s</pre>
<form method="POST" action="/admin/reset" style="margin-top:20px;">
	<button type="submit" class="button secondary">Reset All Data (Clear DB & Sessions)</button>
</form>
`, sessionSummary(s))
	render(w, "Day 1 OAuth Demo (go-oidc)", body)
}

func (c *ClientApp) StartAuthCode(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)

	state := randomString(18)
	codeVerifier := randomString(48)
	codeChallenge := pkceChallenge(codeVerifier)

	c.mu.Lock()
	s.State = state
	s.CodeVerifier = codeVerifier
	s.LastEvent = "Created state + PKCE verifier; redirected browser to /authorize."
	c.mu.Unlock()

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", "repoboard-web")
	q.Set("redirect_uri", clientBaseURL+"/callback")
	q.Set("scope", "openid profile profile.read offline_access")
	q.Set("state", state)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")
	http.Redirect(w, r, opBaseURL+"/authorize?"+q.Encode(), http.StatusFound)
}

func (c *ClientApp) Callback(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	c.mu.Lock()
	expectedState := s.State
	codeVerifier := s.CodeVerifier
	c.mu.Unlock()

	if state == "" || state != expectedState {
		http.Error(w, "state mismatch — possible CSRF", http.StatusBadRequest)
		return
	}

	tokenResp, status, err := postForm(opBaseURL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {"repoboard-web"},
		"code":          {code},
		"redirect_uri":  {clientBaseURL + "/callback"},
		"code_verifier": {codeVerifier},
	}, "", "")
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("token exchange failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}

	var tr tokenResponse
	json.Unmarshal([]byte(tokenResp), &tr)

	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	s.RefreshToken = tr.RefreshToken
	s.LastEvent = "Callback: verified state, exchanged code + code_verifier at /token."
	c.mu.Unlock()

	render(w, "Callback Complete", fmt.Sprintf(`
<p>The client verified <code>state</code>, then posted the code and <code>code_verifier</code> to <code>/token</code>.</p>
<h2>Token response</h2>
<pre>%s</pre>
<p><a class="button" href="/call-api">Call Resource API</a></p>
<p><a href="/">Home</a></p>
`, esc(tokenResp)))
}

func (c *ClientApp) Refresh(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.RefreshToken == "" {
		render(w, "No Refresh Token", `<p>Run the Auth Code + PKCE flow first.</p><p><a href="/">Home</a></p>`)
		return
	}

	tokenResp, status, err := postForm(opBaseURL+"/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {s.RefreshToken},
		"client_id":     {"repoboard-web"},
	}, "", "")
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("refresh failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}

	var tr tokenResponse
	json.Unmarshal([]byte(tokenResp), &tr)

	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	if tr.RefreshToken != "" {
		s.RefreshToken = tr.RefreshToken
	}
	s.LastEvent = "Used refresh_token grant to get a new access token."
	c.mu.Unlock()

	render(w, "Refresh Token Flow", fmt.Sprintf(`
<p>The client exchanged its refresh token for a new access token.</p>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, esc(tokenResp)))
}

func (c *ClientApp) CallAPI(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.AccessToken == "" {
		render(w, "No Access Token", `<p>Run the Auth Code + PKCE flow first.</p><p><a href="/">Home</a></p>`)
		return
	}

	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	render(w, "Resource API Call", fmt.Sprintf(`
<p>Client called <code>/api/profile</code> with <code>Authorization: Bearer ...</code></p>
<h2>Response (HTTP %d)</h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, resp.StatusCode, esc(string(b))))
}

func (c *ClientApp) ClientCredentials(w http.ResponseWriter, r *http.Request) {
	tokenResp, status, err := postForm(opBaseURL+"/token", url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"fraud.check"},
	}, "order-service", "order-service-secret")
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("token request failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}

	var tr tokenResponse
	json.Unmarshal([]byte(tokenResp), &tr)

	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/api/fraud-report", nil)
	req.Header.Set("Authorization", "Bearer "+tr.AccessToken)
	apiResp, _ := http.DefaultClient.Do(req)
	defer apiResp.Body.Close()
	b, _ := io.ReadAll(apiResp.Body)

	render(w, "Client Credentials Flow", fmt.Sprintf(`
<p><strong>No browser redirect. No user approval.</strong> The service authenticated directly.</p>
<h2>Token response</h2>
<pre>%s</pre>
<h2>Internal API response (HTTP %d)</h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, esc(tokenResp), apiResp.StatusCode, esc(string(b))))
}

func (c *ClientApp) DeviceStart(w http.ResponseWriter, r *http.Request) {
	resp, status, err := postForm(opBaseURL+"/device_authorization", url.Values{
		"client_id": {"deploy-cli"},
		"scope":     {"deploy.write"},
	}, "", "")
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("device authorization failed (%d): %s", status, resp), http.StatusBadGateway)
		return
	}

	var dr deviceResponse
	json.Unmarshal([]byte(resp), &dr)

	render(w, "Device Authorization Started", fmt.Sprintf(`
<p>The CLI asked the AS for a device code. The user must now visit the verification URL and enter the code.</p>
<h2>Device authorization response</h2>
<pre>%s</pre>
<p>Open verification page and enter code <strong>%s</strong>:</p>
<p><a class="button" href="%s" target="_blank">Open verification page</a></p>
<p>Then poll the token endpoint:</p>
<p><a class="button secondary" href="/device/poll?device_code=%s">Poll token endpoint</a></p>
<p><a href="/">Home</a></p>
`, esc(resp), esc(dr.UserCode), esc(dr.VerificationURI), esc(dr.DeviceCode)))
}

func (c *ClientApp) DevicePoll(w http.ResponseWriter, r *http.Request) {
	dc := r.URL.Query().Get("device_code")
	resp, status, err := postForm(opBaseURL+"/token", url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {dc},
		"client_id":   {"deploy-cli"},
	}, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if status != http.StatusOK {
		render(w, "Polling…", fmt.Sprintf(`
<p>Token endpoint returned HTTP %d:</p>
<pre>%s</pre>
<p>If this says <code>authorization_pending</code>, approve the device first, then poll again.</p>
<p><a class="button" href="/device/poll?device_code=%s">Poll again</a></p>
<p><a href="/">Home</a></p>
`, status, esc(resp), esc(dc)))
		return
	}

	var tr tokenResponse
	json.Unmarshal([]byte(resp), &tr)

	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/api/deploy", nil)
	req.Header.Set("Authorization", "Bearer "+tr.AccessToken)
	apiResp, _ := http.DefaultClient.Do(req)
	defer apiResp.Body.Close()
	b, _ := io.ReadAll(apiResp.Body)

	render(w, "Device Flow Complete", fmt.Sprintf(`
<p>Device authorization complete. The CLI received an access token and called the deploy API.</p>
<h2>Token response</h2>
<pre>%s</pre>
<h2>Deploy API response (HTTP %d)</h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, esc(resp), apiResp.StatusCode, esc(string(b))))
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type deviceResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
}

func (c *ClientApp) getSession(w http.ResponseWriter, r *http.Request) *Session {
	if cookie, err := r.Cookie("session_id"); err == nil {
		c.mu.Lock()
		if s, ok := c.sessions[cookie.Value]; ok {
			c.mu.Unlock()
			return s
		}
		c.mu.Unlock()
	}
	id := randomString(24)
	s := &Session{ID: id}
	c.mu.Lock()
	c.sessions[id] = s
	c.mu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: id, Path: "/", HttpOnly: true})
	return s
}

func sessionSummary(s *Session) string {
	return fmt.Sprintf(
		"session_id:    %s\nstate:         %s\ncode_verifier: %s\naccess_token:  %s\nrefresh_token: %s\nlast_event:    %s",
		short(s.ID), short(s.State), short(s.CodeVerifier), short(s.AccessToken), short(s.RefreshToken), s.LastEvent,
	)
}

func short(v string) string {
	if v == "" {
		return "(empty)"
	}
	if len(v) <= 12 {
		return v
	}
	return v[:8] + "..."
}

func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func postForm(target string, form url.Values, basicUser, basicPass string) (string, int, error) {
	req, _ := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if basicUser != "" {
		req.SetBasicAuth(basicUser, basicPass)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b), resp.StatusCode, nil
}

func esc(v string) string { return template.HTMLEscapeString(v) }

func render(w http.ResponseWriter, title, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	const page = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; line-height: 1.5; max-width: 980px; margin: 40px auto; padding: 0 20px; color: #172033; }
    h1, h2 { line-height: 1.2; }
    pre { background: #f5f7fb; border: 1px solid #d8deea; padding: 14px; overflow-x: auto; border-radius: 6px; }
    code { background: #f5f7fb; padding: 1px 4px; border-radius: 4px; }
    .button { display: inline-block; background: #1b5cff; color: white; border-radius: 6px; padding: 10px 14px; text-decoration: none; font-weight: 700; margin: 4px 0; }
    .secondary { background: #32415c; }
  </style>
</head>
<body>
  <h1>{{.Title}}</h1>
  {{.Body}}
</body>
</html>`
	t := template.Must(template.New("p").Parse(page))
	t.Execute(w, map[string]any{"Title": title, "Body": template.HTML(body)})
}
