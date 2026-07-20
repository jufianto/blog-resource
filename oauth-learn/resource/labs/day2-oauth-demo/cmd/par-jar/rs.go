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

// The Resource Server here just validates the JWT locally (as in the jwt-validation
// example). PAR/JAR harden the front channel; token validation is unchanged.

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

func apiProfile(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	raw, ok := strings.CutPrefix(auth, "Bearer ")
	if !ok || raw == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
		return
	}

	tok, err := jwt.ParseSigned(raw, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": err.Error()})
		return
	}
	if len(tok.Headers) == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": "no header"})
		return
	}
	key, err := jwksForKID(tok.Headers[0].KeyID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": err.Error()})
		return
	}
	var std jwt.Claims
	var custom struct {
		Scope    string `json:"scope"`
		ClientID string `json:"client_id"`
	}
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

	writeJSON(w, http.StatusOK, map[string]any{
		"api":       "profile",
		"subject":   std.Subject,
		"client_id": custom.ClientID,
		"scope":     custom.Scope,
		"message":   "Authorized via a pushed, signed request (PAR + JAR).",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
