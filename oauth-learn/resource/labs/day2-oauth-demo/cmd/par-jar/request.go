package main

import (
	"crypto/rsa"
	"maps"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// makeRequestObject builds a JAR: the authorization request as a JWT signed with the
// client's key. The AS verifies it against the client's registered public key, so the
// request parameters cannot be tampered with in transit or by the browser.
func makeRequestObject(key *rsa.PrivateKey, params map[string]any) (string, error) {
	opts := (&jose.SignerOptions{}).
		WithType("oauth-authz-req+jwt").
		WithHeader("kid", clientKeyID)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, opts)
	if err != nil {
		return "", err
	}

	claims := map[string]any{
		"iss": clientID,   // the request object is issued by the client
		"aud": opBaseURL,  // and intended for this Authorization Server
	}
	maps.Copy(claims, params)
	return jwt.Signed(signer).Claims(claims).Serialize()
}
