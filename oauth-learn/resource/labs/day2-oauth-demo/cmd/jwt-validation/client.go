package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"day2-oauth-demo/internal/base"
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

func NewClientApp() *ClientApp { return &ClientApp{sessions: make(map[string]*Session)} }

func (c *ClientApp) ResetSessions() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessions = make(map[string]*Session)
}

func (c *ClientApp) Home(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	body := fmt.Sprintf(`
<p>Day 2, Lessons 4-5. The Authorization Server issues <strong>JWT</strong> access tokens.
The Resource Server verifies them <strong>locally</strong> against the AS's published JWKS —
it never calls the AS per request (contrast with Day 1's introspection).</p>

<h2>1. Get a JWT access token</h2>
<p><a class="button" href="/start">Run Authorization Code + PKCE</a></p>

<h2>2. Use it</h2>
<p><a class="button secondary" href="/call-api">Call resource API (valid token)</a></p>
<p><a class="button bad" href="/tamper">Tamper with the token, then call the API</a></p>
<p><a class="button secondary" href="/refresh">Use refresh token</a></p>

<h2>Current Client Session</h2>
<pre>%s</pre>
<form method="POST" action="/admin/reset" style="margin-top:20px;">
  <button type="submit" class="button secondary">Reset</button>
</form>
`, base.Esc(sessionSummary(s)))
	base.Render(w, "Day 2 — JWT Validation (go-oidc)", body)
}

func (c *ClientApp) StartAuthCode(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)

	state := base.RandomString(18)
	verifier := base.RandomString(48)

	c.mu.Lock()
	s.State = state
	s.CodeVerifier = verifier
	s.LastEvent = "Created state + PKCE verifier; redirected to /authorize."
	c.mu.Unlock()

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", clientBaseURL+"/callback")
	q.Set("scope", "openid profile profile.read offline_access")
	q.Set("state", state)
	q.Set("code_challenge", base.PKCEChallenge(verifier))
	q.Set("code_challenge_method", "S256")
	http.Redirect(w, r, opBaseURL+"/authorize?"+q.Encode(), http.StatusFound)
}

func (c *ClientApp) Callback(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	c.mu.Lock()
	expected := s.State
	verifier := s.CodeVerifier
	c.mu.Unlock()

	if state == "" || state != expected {
		http.Error(w, "state mismatch — possible CSRF", http.StatusBadRequest)
		return
	}

	tokenResp, status, err := base.PostForm(opBaseURL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"code":          {code},
		"redirect_uri":  {clientBaseURL + "/callback"},
		"code_verifier": {verifier},
	}, "", "")
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("token exchange failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}

	var tr tokenResponse
	_ = json.Unmarshal([]byte(tokenResp), &tr)

	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	s.RefreshToken = tr.RefreshToken
	s.LastEvent = "Exchanged code + verifier; received a JWT access token."
	c.mu.Unlock()

	base.Render(w, "Token Received", fmt.Sprintf(`
<p>The access token is a JWT. Its header and payload are Base64URL — anyone can read them,
which is why the signature (verified against the JWKS) is what makes it trustworthy.</p>
<h2>Decoded JWT</h2>
<pre>%s</pre>
<h2>Raw token response</h2>
<pre>%s</pre>
<p><a class="button" href="/call-api">Call Resource API</a></p>
<p><a href="/">Home</a></p>
`, base.Esc(decodeJWT(tr.AccessToken)), base.Esc(tokenResp)))
}

func (c *ClientApp) Refresh(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.RefreshToken == "" {
		base.Render(w, "No Refresh Token", `<p>Run the Auth Code + PKCE flow first.</p><p><a href="/">Home</a></p>`)
		return
	}
	tokenResp, status, err := base.PostForm(opBaseURL+"/token", url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {s.RefreshToken},
		"client_id":     {clientID},
	}, "", "")
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("refresh failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}
	var tr tokenResponse
	_ = json.Unmarshal([]byte(tokenResp), &tr)
	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	if tr.RefreshToken != "" {
		s.RefreshToken = tr.RefreshToken
	}
	s.LastEvent = "Refreshed; received a new JWT access token."
	c.mu.Unlock()
	base.Render(w, "Refreshed", fmt.Sprintf(`
<h2>Decoded new JWT</h2><pre>%s</pre>
<h2>Raw response</h2><pre>%s</pre>
<p><a href="/">Home</a></p>`, base.Esc(decodeJWT(tr.AccessToken)), base.Esc(tokenResp)))
}

func (c *ClientApp) CallAPI(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.AccessToken == "" {
		base.Render(w, "No Access Token", `<p>Run the Auth Code + PKCE flow first.</p><p><a href="/">Home</a></p>`)
		return
	}
	status, respBody := callProfile(s.AccessToken)
	base.Render(w, "Resource API Call", fmt.Sprintf(`
<p>Called <code>/api/profile</code> with the JWT. The Resource Server verified it locally.</p>
<h2>Response (HTTP %d) <span class="%s">%s</span></h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, status, okClass(status), okLabel(status), base.Esc(respBody)))
}

// Tamper flips one character in the JWT payload, leaving the signature intact. Because the
// signature no longer matches the altered payload, local validation must reject it.
func (c *ClientApp) Tamper(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.AccessToken == "" {
		base.Render(w, "No Access Token", `<p>Run the Auth Code + PKCE flow first.</p><p><a href="/">Home</a></p>`)
		return
	}
	tampered := tamperPayload(s.AccessToken)
	status, respBody := callProfile(tampered)
	base.Render(w, "Tampered Token", fmt.Sprintf(`
<p>We altered one byte of the token payload without re-signing it, then called the API.
Local JWKS validation catches this — no round trip to the Authorization Server needed.</p>
<h2>Response (HTTP %d) <span class="%s">%s</span></h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, status, okClass(status), okLabel(status), base.Esc(respBody)))
}

func callProfile(token string) (int, string) {
	req, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	id := base.RandomString(24)
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

// decodeJWT pretty-prints the header and payload of a JWT for display. It does NOT verify
// anything — decoding and verifying are different acts, which is the lesson.
func decodeJWT(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "(not a JWT: expected three dot-separated parts)"
	}
	header := prettyJSONSegment(parts[0])
	payload := prettyJSONSegment(parts[1])
	return fmt.Sprintf("HEADER:\n%s\n\nPAYLOAD:\n%s\n\nSIGNATURE:\n%s", header, payload, short(parts[2]))
}

func prettyJSONSegment(seg string) string {
	raw, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return "(could not base64url-decode)"
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
}

func tamperPayload(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || len(parts[1]) == 0 {
		return token
	}
	b := []byte(parts[1])
	// Flip the last character to a different valid base64url character.
	if b[len(b)-1] == 'A' {
		b[len(b)-1] = 'B'
	} else {
		b[len(b)-1] = 'A'
	}
	parts[1] = string(b)
	return strings.Join(parts, ".")
}

func okClass(status int) string {
	if status == http.StatusOK {
		return "ok"
	}
	return "fail"
}

func okLabel(status int) string {
	if status == http.StatusOK {
		return "ACCEPTED"
	}
	return "REJECTED"
}
