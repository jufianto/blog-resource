package main

import (
	"crypto/rand"
	"crypto/rsa"
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
	// dpopKey is the client's own key. It never leaves the client; the public half is
	// embedded in each proof, and the token is bound to its thumbprint.
	dpopKey *rsa.PrivateKey
}

type Session struct {
	ID           string
	State        string
	CodeVerifier string
	AccessToken  string
	TokenType    string
	LastEvent    string
}

func NewClientApp() (*ClientApp, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &ClientApp{sessions: make(map[string]*Session), dpopKey: key}, nil
}

func (c *ClientApp) ResetSessions() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessions = make(map[string]*Session)
}

func (c *ClientApp) Home(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	body := fmt.Sprintf(`
<p>Day 2, Lesson 7. The access token is <strong>DPoP sender-constrained</strong>: the
Authorization Server binds it to this client's key (<code>cnf.jkt</code>). Every API call
must carry a fresh DPoP proof signed by that key, so a stolen token alone is worthless.</p>

<h2>1. Get a DPoP-bound token</h2>
<p><a class="button" href="/start">Run Authorization Code + PKCE (with DPoP)</a></p>

<h2>2. Use it correctly</h2>
<p><a class="button secondary" href="/call-api">Call API with token + DPoP proof</a></p>

<h2>3. Attacker replay</h2>
<p><a class="button bad" href="/replay">Replay the token WITHOUT the key (as a plain Bearer)</a></p>

<h2>Current Client Session</h2>
<pre>%s</pre>
<form method="POST" action="/admin/reset" style="margin-top:20px;">
  <button type="submit" class="button secondary">Reset</button>
</form>
`, base.Esc(sessionSummary(s)))
	base.Render(w, "Day 2 — DPoP (go-oidc)", body)
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

	// The token request carries a DPoP proof for POST /token. The AS reads the proof's key
	// and binds the issued token to it.
	proof, err := makeDPoPProof(c.dpopKey, http.MethodPost, opBaseURL+"/token", "")
	if err != nil {
		http.Error(w, "building DPoP proof: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tokenResp, status, err := postFormWithDPoP(opBaseURL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"code":          {code},
		"redirect_uri":  {clientBaseURL + "/callback"},
		"code_verifier": {verifier},
	}, proof)
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("token exchange failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}

	var tr tokenResponse
	_ = json.Unmarshal([]byte(tokenResp), &tr)

	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	s.TokenType = tr.TokenType
	s.LastEvent = "Exchanged code with a DPoP proof; received a DPoP-bound token."
	c.mu.Unlock()

	base.Render(w, "DPoP Token Received", fmt.Sprintf(`
<p>Note <code>"token_type": "DPoP"</code> — this token is bound to the client's key. The
access token's <code>cnf.jkt</code> claim is the thumbprint of that key.</p>
<h2>Token response</h2><pre>%s</pre>
<p><a class="button" href="/call-api">Call Resource API (correctly)</a></p>
<p><a href="/">Home</a></p>
`, base.Esc(tokenResp)))
}

func (c *ClientApp) CallAPI(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.AccessToken == "" {
		base.Render(w, "No Access Token", `<p>Run the flow first.</p><p><a href="/">Home</a></p>`)
		return
	}
	// A fresh proof for THIS request: GET the profile URL, tied to this token via ath.
	proof, err := makeDPoPProof(c.dpopKey, http.MethodGet, profileURL, s.AccessToken)
	if err != nil {
		http.Error(w, "building DPoP proof: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req, _ := http.NewRequest(http.MethodGet, profileURL, nil)
	req.Header.Set("Authorization", "DPoP "+s.AccessToken)
	req.Header.Set("DPoP", proof)
	status, respBody := do(req)

	base.Render(w, "Resource API Call", fmt.Sprintf(`
<p>Sent <code>Authorization: DPoP &lt;token&gt;</code> plus a fresh <code>DPoP</code> proof
signed by the client key.</p>
<h2>Response (HTTP %d) <span class="%s">%s</span></h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, status, okClass(status), okLabel(status), base.Esc(respBody)))
}

// Replay simulates an attacker who copied the token but does NOT have the client's private
// key: they present it as a plain Bearer token, with no proof. The Resource Server rejects it.
func (c *ClientApp) Replay(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	if s.AccessToken == "" {
		base.Render(w, "No Access Token", `<p>Run the flow first.</p><p><a href="/">Home</a></p>`)
		return
	}
	req, _ := http.NewRequest(http.MethodGet, profileURL, nil)
	req.Header.Set("Authorization", "Bearer "+s.AccessToken) // wrong scheme, no DPoP proof
	status, respBody := do(req)

	base.Render(w, "Attacker Replay", fmt.Sprintf(`
<p>The stolen token is presented as a plain Bearer token, with no DPoP proof — exactly what
an attacker who exfiltrated the token (but not the key) could do.</p>
<h2>Response (HTTP %d) <span class="%s">%s</span></h2>
<pre>%s</pre>
<p><a href="/">Home</a></p>
`, status, okClass(status), okLabel(status), base.Esc(respBody)))
}

func do(req *http.Request) (int, string) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func postFormWithDPoP(target string, form url.Values, proof string) (string, int, error) {
	req, _ := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("DPoP", proof)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b), resp.StatusCode, nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
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
	return fmt.Sprintf("session_id:   %s\ntoken_type:   %s\naccess_token: %s\nlast_event:   %s",
		short(s.ID), orEmpty(s.TokenType), short(s.AccessToken), s.LastEvent)
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

func orEmpty(v string) string {
	if v == "" {
		return "(empty)"
	}
	return v
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
