package main

import (
	"bytes"
	"crypto/rsa"
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
	jarKey   *rsa.PrivateKey
}

type Session struct {
	ID           string
	State        string
	CodeVerifier string
	AccessToken  string
	LastEvent    string
}

func NewClientApp(jarKey *rsa.PrivateKey) *ClientApp {
	return &ClientApp{sessions: make(map[string]*Session), jarKey: jarKey}
}

func (c *ClientApp) ResetSessions() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessions = make(map[string]*Session)
}

func (c *ClientApp) Home(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	body := fmt.Sprintf(`
<p>Day 2, Lesson 7. The client sends its authorization request as a <strong>signed request
object (JAR)</strong>, <strong>pushed</strong> to the AS back-channel first (PAR). The browser
only ever carries an opaque <code>request_uri</code> — no scopes, redirect URI, or PKCE
challenge are exposed in front-channel URLs, and none can be tampered with.</p>

<h2>1. Push the request and authorize</h2>
<p><a class="button" href="/start">Build request object → PAR → authorize</a></p>

<h2>2. Use the token</h2>
<p><a class="button secondary" href="/call-api">Call resource API</a></p>

<h2>Current Client Session</h2>
<pre>%s</pre>
<form method="POST" action="/admin/reset" style="margin-top:20px;">
  <button type="submit" class="button secondary">Reset</button>
</form>
`, base.Esc(sessionSummary(s)))
	base.Render(w, "Day 2 — PAR + JAR (go-oidc)", body)
}

func (c *ClientApp) StartAuthCode(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	state := base.RandomString(18)
	verifier := base.RandomString(48)
	challenge := base.PKCEChallenge(verifier)

	c.mu.Lock()
	s.State = state
	s.CodeVerifier = verifier
	s.LastEvent = "Built signed request object; pushed to /par."
	c.mu.Unlock()

	// 1. Build the signed request object (JAR) with all authorization parameters.
	requestObject, err := makeRequestObject(c.jarKey, map[string]any{
		"response_type":         "code",
		"client_id":             clientID,
		"redirect_uri":          clientBaseURL + "/callback",
		"scope":                 "openid profile profile.read offline_access",
		"state":                 state,
		"code_challenge":        challenge,
		"code_challenge_method": "S256",
	})
	if err != nil {
		http.Error(w, "building request object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Push it to the PAR endpoint. Response contains an opaque request_uri.
	parResp, status, err := postForm(opBaseURL+"/par", url.Values{
		"client_id": {clientID},
		"request":   {requestObject},
	})
	if err != nil || status != http.StatusCreated {
		http.Error(w, fmt.Sprintf("PAR failed (%d): %s", status, parResp), http.StatusBadGateway)
		return
	}
	var pr parResponse
	_ = json.Unmarshal([]byte(parResp), &pr)

	// 3. Send the browser to /authorize with the opaque request_uri. For OpenID requests,
	// response_type and client_id must also appear as plain params (they still can't be
	// tampered with — the AS cross-checks them against the signed, pushed request).
	authorizeURL := opBaseURL + "/authorize?" + url.Values{
		"client_id":     {clientID},
		"request_uri":   {pr.RequestURI},
		"response_type": {"code"},
		"scope":         {"openid profile profile.read offline_access"},
	}.Encode()

	base.Render(w, "Pushed Authorization Request", fmt.Sprintf(`
<p>The signed request object (JAR) — its payload is integrity-protected by the client signature:</p>
<pre>%s</pre>
<p>PAR response from the Authorization Server:</p>
<pre>%s</pre>
<p>The browser now goes to <code>/authorize</code> carrying only <code>client_id</code> and the
opaque <code>request_uri</code> — nothing sensitive is in this URL:</p>
<pre>%s</pre>
<p><a class="button" href="%s">Continue to /authorize</a></p>
<p><a href="/">Home</a></p>
`, base.Esc(decodeJWTPayload(requestObject)), base.Esc(parResp), base.Esc(authorizeURL), base.Esc(authorizeURL)))
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

	tokenResp, status, err := postForm(opBaseURL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"code":          {code},
		"redirect_uri":  {clientBaseURL + "/callback"},
		"code_verifier": {verifier},
	})
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("token exchange failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}
	var tr tokenResponse
	_ = json.Unmarshal([]byte(tokenResp), &tr)

	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	s.LastEvent = "Authorized via request_uri; exchanged code for a token."
	c.mu.Unlock()

	base.Render(w, "Token Received", fmt.Sprintf(`
<p>Authorization completed from the pushed, signed request. Token response:</p>
<pre>%s</pre>
<p><a class="button" href="/call-api">Call Resource API</a></p>
<p><a href="/">Home</a></p>
`, base.Esc(tokenResp)))
}

func (c *ClientApp) CallAPI(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.AccessToken == "" {
		base.Render(w, "No Access Token", `<p>Run the flow first.</p><p><a href="/">Home</a></p>`)
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

	base.Render(w, "Resource API Call", fmt.Sprintf(`
<h2>Response (HTTP %d)</h2><pre>%s</pre>
<p><a href="/">Home</a></p>`, resp.StatusCode, base.Esc(string(b))))
}

type parResponse struct {
	RequestURI string `json:"request_uri"`
	ExpiresIn  int    `json:"expires_in"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func postForm(target string, form url.Values) (string, int, error) {
	req, _ := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b), resp.StatusCode, nil
}

func decodeJWTPayload(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "(not a JWT)"
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "(could not decode)"
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
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
	return fmt.Sprintf("session_id:   %s\nstate:        %s\naccess_token: %s\nlast_event:   %s",
		short(s.ID), short(s.State), short(s.AccessToken), s.LastEvent)
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
