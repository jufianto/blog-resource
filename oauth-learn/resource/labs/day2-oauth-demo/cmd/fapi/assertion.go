package main

import (
	"crypto/rsa"
	"time"

	"day2-oauth-demo/internal/base"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// makeClientAssertion builds a private_key_jwt client assertion: a short-lived JWT the client
// signs with its registered key to authenticate to the token and PAR endpoints. No shared
// secret ever crosses the wire. FAPI 2.0 requires exactly one audience; the issuer is always
// accepted, so we use it.
func makeClientAssertion(key *rsa.PrivateKey) (string, error) {
	opts := (&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", clientKeyID)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, opts)
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := jwt.Claims{
		Issuer:   clientID,
		Subject:  clientID,
		Audience: jwt.Audience{opBaseURL},
		ID:       base.RandomString(16),
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(60 * time.Second)),
	}
	return jwt.Signed(signer).Claims(claims).Serialize()
}
