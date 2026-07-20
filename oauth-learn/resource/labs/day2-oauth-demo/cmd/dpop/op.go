package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"html"
	"net/http"

	"day2-oauth-demo/internal/base"

	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
)

var clientRepoboard = &goidc.Client{
	ID: clientID,
	ClientMeta: goidc.ClientMeta{
		RedirectURIs:     []string{clientBaseURL + "/callback"},
		ResponseTypes:    []goidc.ResponseType{goidc.ResponseTypeCode},
		GrantTypes:       []goidc.GrantType{goidc.GrantAuthorizationCode, goidc.GrantRefreshToken},
		ScopeIDs:         "openid profile profile.read offline_access",
		TokenAuthnMethod: goidc.AuthnMethodNone,
	},
}

// newProvider builds the Authorization Server. The Day 2 / Lesson 7 line is WithDPoP:
// when the client presents a DPoP proof on the token request, the AS binds the issued
// access token to that key (embedding cnf.jkt in the JWT). The token is then only usable
// by whoever can produce a fresh proof signed by the matching private key.
func newProvider(store *base.MemoryStore) (*provider.Provider, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating signing key: %w", err)
	}
	jwks := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{KeyID: signingKeyID, Key: key, Algorithm: "RS256"}},
	}

	loginPolicy := goidc.NewPolicy(
		"simple_login",
		func(_ *http.Request, _ *goidc.AuthnSession, _ *goidc.Client) bool { return true },
		func(w http.ResponseWriter, r *http.Request, as *goidc.AuthnSession, _ *goidc.Client) (goidc.Status, error) {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprintf(w, `<!DOCTYPE html><html><body>
<h2>Mock Login &amp; Consent</h2>
<p>Client: <strong>%s</strong></p>
<p>Scopes requested: <strong>%s</strong></p>
<form method="POST" action="/authorize/%s">
  <input type="hidden" name="approve" value="true">
  <button type="submit">Approve for Alya Developer</button>
</form>
</body></html>`, html.EscapeString(as.ClientID), html.EscapeString(as.Scopes), html.EscapeString(as.ID))
				return goidc.StatusPending, nil
			}
			if r.FormValue("approve") == "true" {
				as.Subject = "alya"
				as.GrantedScopes = as.Scopes
				return goidc.StatusSuccess, nil
			}
			return goidc.StatusFailure, nil
		},
	)

	return provider.New(
		opBaseURL,
		store,
		func(_ context.Context) (goidc.JSONWebKeySet, error) { return jwks, nil },
		provider.WithStaticClients(clientRepoboard),
		provider.WithAuthCodeGrant(store, goidc.ResponseTypeCode),
		provider.WithPKCERequired(goidc.CodeChallengeMethodSHA256),
		provider.WithRefreshTokenGrant(store),
		provider.WithRefreshTokenRotation(),
		provider.WithTokenAuthnMethods(goidc.AuthnMethodNone),
		provider.WithTokenOptions(func(_ context.Context, _ *goidc.Grant, _ *goidc.Client) goidc.TokenOptions {
			return goidc.NewJWTTokenOptions("RS256", 300)
		}),
		// >>> Lesson 7: accept DPoP proofs and bind tokens to the proof key (RS256). <<<
		provider.WithDPoP(goidc.RS256),
		provider.WithHandleErrorFunc(func(_ context.Context, err error) {
			fmt.Printf("OIDC ERROR: %v\n", err)
		}),
		provider.WithPolicies(loginPolicy),
		provider.WithScopes(
			goidc.ScopeOpenID, goidc.ScopeProfile, goidc.ScopeOfflineAccess,
			goidc.NewScope("profile.read"),
		),
	)
}
