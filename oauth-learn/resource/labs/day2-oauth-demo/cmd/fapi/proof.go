package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"day2-oauth-demo/internal/base"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// DPoP proof: built by the client, validated by the Resource Server. Same mechanism as the
// dpop example — included here so the FAPI capstone is self-contained.

func makeDPoPProof(key *rsa.PrivateKey, htm, htu, accessToken string) (string, error) {
	opts := (&jose.SignerOptions{EmbedJWK: true}).WithType("dpop+jwt")
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, opts)
	if err != nil {
		return "", err
	}
	claims := map[string]any{
		"jti": base.RandomString(16),
		"htm": htm,
		"htu": htu,
		"iat": time.Now().Unix(),
	}
	if accessToken != "" {
		claims["ath"] = sha256b64(accessToken)
	}
	return jwt.Signed(signer).Claims(claims).Serialize()
}

type dpopClaims struct {
	HTM string `json:"htm"`
	HTU string `json:"htu"`
	ATH string `json:"ath"`
}

func validateDPoPProof(proof, htm, htu, accessToken, boundJKT string) error {
	tok, err := jwt.ParseSigned(proof, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return fmt.Errorf("proof is not a well-formed RS256 JWT: %w", err)
	}
	if len(tok.Headers) == 0 {
		return fmt.Errorf("proof has no header")
	}
	hdr := tok.Headers[0]
	if hdr.ExtraHeaders["typ"] != "dpop+jwt" {
		return fmt.Errorf(`proof header typ must be "dpop+jwt"`)
	}
	if hdr.JSONWebKey == nil {
		return fmt.Errorf("proof header has no embedded jwk")
	}
	var c dpopClaims
	if err := tok.Claims(hdr.JSONWebKey.Key, &c); err != nil {
		return fmt.Errorf("proof signature invalid: %w", err)
	}
	if c.HTM != htm {
		return fmt.Errorf("proof htm %q does not match request method %q", c.HTM, htm)
	}
	if c.HTU != htu {
		return fmt.Errorf("proof htu %q does not match request URL %q", c.HTU, htu)
	}
	tp, err := hdr.JSONWebKey.Thumbprint(crypto.SHA256)
	if err != nil {
		return fmt.Errorf("computing jwk thumbprint: %w", err)
	}
	if jkt := base64.RawURLEncoding.EncodeToString(tp); jkt != boundJKT {
		return fmt.Errorf("proof key thumbprint %q is not the token's bound key %q", jkt, boundJKT)
	}
	if c.ATH != sha256b64(accessToken) {
		return fmt.Errorf("proof ath does not match the presented access token")
	}
	return nil
}

func sha256b64(s string) string {
	sum := sha256.Sum256([]byte(s))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
