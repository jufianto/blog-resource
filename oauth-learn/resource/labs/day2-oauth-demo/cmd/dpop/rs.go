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

// The Resource Server does two checks per request:
//  1. the access token is a valid JWT (signature via JWKS, issuer, expiry, scope);
//  2. the caller proves possession of the key the token is bound to (cnf.jkt) by sending
//     a fresh DPoP proof for this exact method+URL+token.
// A token presented without a valid proof is rejected — that is what "sender-constrained"
// means, and why a leaked DPoP token cannot be replayed.

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

func validateAccessToken(raw, requiredScope string) (*jwt.Claims, *accessClaims, error) {
	tok, err := jwt.ParseSigned(raw, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, nil, fmt.Errorf("not a well-formed RS256 JWT: %w", err)
	}
	if len(tok.Headers) == 0 {
		return nil, nil, fmt.Errorf("JWT has no header")
	}
	key, err := jwksForKID(tok.Headers[0].KeyID)
	if err != nil {
		return nil, nil, err
	}
	var std jwt.Claims
	var custom accessClaims
	if err := tok.Claims(key.Key, &std, &custom); err != nil {
		return nil, nil, fmt.Errorf("signature check failed: %w", err)
	}
	if err := std.Validate(jwt.Expected{Issuer: opBaseURL, Time: time.Now()}); err != nil {
		return nil, nil, fmt.Errorf("claim validation failed: %w", err)
	}
	if !slices.Contains(strings.Fields(custom.Scope), requiredScope) {
		return nil, nil, fmt.Errorf("token lacks required scope %q (has %q)", requiredScope, custom.Scope)
	}
	return &std, &custom, nil
}

func apiProfile(w http.ResponseWriter, r *http.Request) {
	// DPoP tokens are presented with the "DPoP" auth scheme, not "Bearer".
	auth := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(auth, "DPoP ")
	if !ok || token == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error":  "token rejected",
			"reason": "this API requires a DPoP-bound token (Authorization: DPoP <token>) — a plain Bearer token is not accepted",
		})
		return
	}

	std, custom, err := validateAccessToken(token, "profile.read")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": err.Error()})
		return
	}
	if custom.Cnf.JKT == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token rejected", "reason": "token is not DPoP-bound (no cnf.jkt)"})
		return
	}

	proof := r.Header.Get("DPoP")
	if proof == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error":  "token rejected",
			"reason": "missing DPoP proof header — the token is bound to a key you must prove you hold",
		})
		return
	}
	// htu is the request URL without query/fragment; here it is the fixed profile endpoint.
	if err := validateDPoPProof(proof, http.MethodGet, profileURL, token, custom.Cnf.JKT); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "DPoP proof rejected", "reason": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"api":           "profile",
		"validated":     "JWT verified locally AND DPoP proof-of-possession confirmed",
		"subject":       std.Subject,
		"client_id":     custom.ClientID,
		"bound_key_jkt": custom.Cnf.JKT,
		"scope":         custom.Scope,
		"message":       "A stolen copy of this token is useless without the client's private key.",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
