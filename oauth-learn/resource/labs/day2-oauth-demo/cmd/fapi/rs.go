package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// The Resource Server validates the JWT locally AND the DPoP proof-of-possession — the same
// checks as the dpop example. Under FAPI 2.0 every access token is sender-constrained.

var (
	jwksMu    sync.Mutex
	jwksCache *jose.JSONWebKeySet
)

func jwksForKID(kid string) (*jose.JSONWebKey, error) {
	jwksMu.Lock()
	defer jwksMu.Unlock()
	if jwksCache != nil {
		if keys := jwksCache.Key(kid); len(keys) > 0 {
			return &keys[0], nil
		}
	}
	resp, err := http.Get(opBaseURL + "/jwks")
	if err != nil {
		return nil, fmt.Errorf("fetching JWKS: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var set jose.JSONWebKeySet
	if err := json.Unmarshal(body, &set); err != nil {
		return nil, fmt.Errorf("parsing JWKS: %w", err)
	}
	jwksCache = &set
	keys := set.Key(kid)
	if len(keys) == 0 {
		return nil, fmt.Errorf("no key in JWKS matches kid %q", kid)
	}
	return &keys[0], nil
}

type accessClaims struct {
	Scope    string `json:"scope"`
	ClientID string `json:"client_id"`
	Cnf      struct {
		JKT string `json:"jkt"`
	} `json:"cnf"`
}

func apiProfile(w http.ResponseWriter, r *http.Request) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "DPoP ")
	if !ok || token == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error":  "token rejected",
			"reason": "FAPI 2.0 requires a DPoP-bound token (Authorization: DPoP <token>)",
		})
		return
	}

	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil || len(tok.Headers) == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": "malformed JWT"})
		return
	}
	key, err := jwksForKID(tok.Headers[0].KeyID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": err.Error()})
		return
	}
	var std jwt.Claims
	var custom accessClaims
	if err := tok.Claims(key.Key, &std, &custom); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": err.Error()})
		return
	}
	if err := std.Validate(jwt.Expected{Issuer: opBaseURL, Time: time.Now()}); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": err.Error()})
		return
	}
	if !slices.Contains(strings.Fields(custom.Scope), "profile.read") {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "insufficient scope"})
		return
	}
	if custom.Cnf.JKT == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": "token is not DPoP-bound"})
		return
	}

	proof := r.Header.Get("DPoP")
	if proof == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": "missing DPoP proof"})
		return
	}
	if err := validateDPoPProof(proof, http.MethodGet, profileURL, token, custom.Cnf.JKT); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "DPoP proof rejected", "reason": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"api":           "profile",
		"validated":     "FAPI 2.0: JWT verified locally AND DPoP proof-of-possession confirmed",
		"subject":       std.Subject,
		"client_id":     custom.ClientID,
		"bound_key_jkt": custom.Cnf.JKT,
		"scope":         custom.Scope,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
