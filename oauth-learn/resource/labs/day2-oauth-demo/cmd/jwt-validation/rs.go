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

// This Resource Server never calls the Authorization Server per-request. It fetches the
// AS's public keys once from the JWKS endpoint, then verifies every incoming JWT access
// token locally: signature, issuer, expiry, and scope. That is the whole point of Lesson
// 4-5 — a JWT is self-contained and independently verifiable.

var (
	jwksMu    sync.Mutex
	jwksCache *jose.JSONWebKeySet
)

// jwksForKID returns the AS public keys, fetching (or refetching, on an unknown kid —
// the key-rotation case) from the JWKS endpoint. In production you would cache with a TTL
// and honour Cache-Control; here we refetch only when a kid is missing.
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

// validateJWT performs local validation and returns the standard + custom claims.
// Every returned error is a distinct, teachable failure reason.
func validateJWT(raw, requiredScope string) (*jwt.Claims, *accessClaims, error) {
	// 1. Parse, restricting the accepted signing algorithm (never trust the token's own alg blindly).
	tok, err := jwt.ParseSigned(raw, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, nil, fmt.Errorf("not a well-formed RS256 JWT: %w", err)
	}
	if len(tok.Headers) == 0 {
		return nil, nil, fmt.Errorf("JWT has no header")
	}

	// 2. Find the signing key by kid in the AS JWKS.
	key, err := jwksForKID(tok.Headers[0].KeyID)
	if err != nil {
		return nil, nil, err
	}

	// 3. Verify the signature and decode claims (Claims() fails if the signature is invalid).
	var std jwt.Claims
	var custom accessClaims
	if err := tok.Claims(key.Key, &std, &custom); err != nil {
		return nil, nil, fmt.Errorf("signature check failed: %w", err)
	}

	// 4. Validate issuer and expiry.
	if err := std.Validate(jwt.Expected{Issuer: opBaseURL, Time: time.Now()}); err != nil {
		return nil, nil, fmt.Errorf("claim validation failed: %w", err)
	}

	// 5. Enforce the scope the endpoint requires.
	if !hasScope(custom.Scope, requiredScope) {
		return nil, nil, fmt.Errorf("token lacks required scope %q (has %q)", requiredScope, custom.Scope)
	}

	return &std, &custom, nil
}

type accessClaims struct {
	Scope    string `json:"scope"`
	ClientID string `json:"client_id"`
}

func hasScope(scopes, required string) bool {
	return slices.Contains(strings.Fields(scopes), required)
}

func apiProfile(w http.ResponseWriter, r *http.Request) {
	raw := bearer(r)
	if raw == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
		return
	}
	std, custom, err := validateJWT(raw, "profile.read")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error":  "token rejected",
			"reason": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"api":                 "profile",
		"validated":           "locally, via JWKS — no call to the Authorization Server",
		"subject":             std.Subject,
		"issuer":              std.Issuer,
		"client_id":           custom.ClientID,
		"scope":               custom.Scope,
		"expires":             std.Expiry.Time().Format(time.RFC3339),
		"message":             "RepoBoard read the profile resource with a self-verified JWT.",
	})
}

func bearer(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
