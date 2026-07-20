package main

import (
	"bytes"
	"crypto/rand"
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

const assertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

type ClientApp struct {
	mu       sync.Mutex
	sessions map[string]*Session
	authKey  *rsa.PrivateKey // private_key_jwt client authentication
	dpopKey  *rsa.PrivateKey // DPoP proof-of-possession
}

type Session struct {
	ID           string
	State        string
	CodeVerifier string
	AccessToken  string
	LastEvent    string
}

func NewClientApp(authKey *rsa.PrivateKey) (*ClientApp, error) {
	dpopKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &ClientApp{sessions: make(map[string]*Session), authKey: authKey, dpopKey: dpopKey}, nil
}

func (c *ClientApp) ResetSessions() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessions = make(map[string]*Session)
}

func (c *ClientApp) Home(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	body := fmt.Sprintf(`
<p>Day 2, Lesson 8 — the FAPI 2.0 capstone. Every control from the earlier examples is on at
once, enforced by <code>WithProfile(FAPI2)</code>:</p>
<ul>
  <li><strong>private_key_jwt</strong> — the client authenticates with a signed assertion, no secret.</li>
  <li><strong>PAR</strong> — the request is pushed to the back-channel first.</li>
  <li><strong>PKCE (S256)</strong> — required.</li>
  <li><strong>DPoP</strong> — access tokens are sender-constrained.</li>
</ul>

<h2>Run the full FAPI 2.0 flow</h2>
<p><a class="button" href="/start">1. Assertion → PAR → authorize</a></p>
<p><a class="button secondary" href="/call-api">2. Call API (DPoP-bound token + proof)</a></p>

<h2>Current Client Session</h2>
<pre>%s</pre>
<form method="POST" action="/admin/reset" style="margin-top:20px;">
  <button type="submit" class="button secondary">Reset</button>
</form>
`, base.Esc(sessionSummary(s)))
	base.Render(w, "Day 2 — FAPI 2.0 (go-oidc)", body)
}

func (c *ClientApp) StartAuthCode(w http.ResponseWriter, r *http.Request) {
	s := c.getSession(w, r)
	state := base.RandomString(18)
	verifier := base.RandomString(48)

	c.mu.Lock()
	s.State = state
	s.CodeVerifier = verifier
	s.LastEvent = "Built client assertion; pushed request to /par."
	c.mu.Unlock()

	assertion, err := makeClientAssertion(c.authKey)
	if err != nil {
		http.Error(w, "building client assertion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// PAR, authenticated with private_key_jwt (client_assertion). No DPoP on PAR itself.
	parResp, status, err := postForm(opBaseURL+"/par", url.Values{
		"client_id":             {clientID},
		"client_assertion_type": {assertionType},
		"client_assertion":      {assertion},
		"response_type":         {"code"},
		"redirect_uri":          {clientBaseURL + "/callback"},
		"scope":                 {"profile.read"},
		"state":                 {state},
		"code_challenge":        {base.PKCEChallenge(verifier)},
		"code_challenge_method": {"S256"},
	}, "")
	if err != nil || status != http.StatusCreated {
		http.Error(w, fmt.Sprintf("PAR failed (%d): %s", status, parResp), http.StatusBadGateway)
		return
	}
	var pr parResponse
	_ = json.Unmarshal([]byte(parResp), &pr)

	authorizeURL := opBaseURL + "/authorize?" + url.Values{
		"client_id":   {clientID},
		"request_uri": {pr.RequestURI},
	}.Encode()

	base.Render(w, "Pushed Authorization Request", fmt.Sprintf(`
<p>Client assertion (private_key_jwt) — signed with the client's registered key:</p>
<pre>%s</pre>
<p>PAR response:</p><pre>%s</pre>
<p><a class="button" href="%s">Continue to /authorize</a></p>
<p><a href="/">Home</a></p>
`, base.Esc(decodeJWTPayload(assertion)), base.Esc(parResp), base.Esc(authorizeURL)))
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

	assertion, err := makeClientAssertion(c.authKey)
	if err != nil {
		http.Error(w, "building client assertion: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Token request: private_key_jwt auth + a DPoP proof for POST /token.
	proof, err := makeDPoPProof(c.dpopKey, http.MethodPost, opBaseURL+"/token", "")
	if err != nil {
		http.Error(w, "building DPoP proof: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tokenResp, status, err := postForm(opBaseURL+"/token", url.Values{
		"grant_type":            {"authorization_code"},
		"client_id":             {clientID},
		"code":                  {code},
		"redirect_uri":          {clientBaseURL + "/callback"},
		"code_verifier":         {verifier},
		"client_assertion_type": {assertionType},
		"client_assertion":      {assertion},
	}, proof)
	if err != nil || status != http.StatusOK {
		http.Error(w, fmt.Sprintf("token exchange failed (%d): %s", status, tokenResp), http.StatusBadGateway)
		return
	}
	var tr tokenResponse
	_ = json.Unmarshal([]byte(tokenResp), &tr)

	c.mu.Lock()
	s.AccessToken = tr.AccessToken
	s.LastEvent = "Authenticated with private_key_jwt; received a DPoP-bound token."
	c.mu.Unlock()

	base.Render(w, "Token Received", fmt.Sprintf(`
<p>Note <code>"token_type": "DPoP"</code>. The client authenticated with a signed assertion
(no secret), the request was pushed via PAR, and the token is sender-constrained.</p>
<h2>Token response</h2><pre>%s</pre>
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
	proof, err := makeDPoPProof(c.dpopKey, http.MethodGet, profileURL, s.AccessToken)
	if err != nil {
		http.Error(w, "building DPoP proof: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req, _ := http.NewRequest(http.MethodGet, profileURL, nil)
	req.Header.Set("Authorization", "DPoP "+s.AccessToken)
	req.Header.Set("DPoP", proof)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	label := "REJECTED"
	cls := "fail"
	if resp.StatusCode == http.StatusOK {
		label, cls = "ACCEPTED", "ok"
	}
	base.Render(w, "Resource API Call", fmt.Sprintf(`
<h2>Response (HTTP %d) <span class="%s">%s</span></h2><pre>%s</pre>
<p><a href="/">Home</a></p>`, resp.StatusCode, cls, label, base.Esc(string(b))))
}

type parResponse struct {
	RequestURI string `json:"request_uri"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// postForm posts a urlencoded form, optionally with a DPoP proof header.
func postForm(target string, form url.Values, dpopProof string) (string, int, error) {
	req, _ := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if dpopProof != "" {
		req.Header.Set("DPoP", dpopProof)
	}
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
