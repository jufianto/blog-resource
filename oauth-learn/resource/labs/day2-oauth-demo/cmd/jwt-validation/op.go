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

// The one client this example needs: a public web app using Authorization Code + PKCE.
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

// newProvider builds the Authorization Server. The one line that matters for this
// example is WithTokenOptions returning a JWT token format: access tokens are now
// signed JWTs the Resource Server can verify on its own, instead of opaque strings
// it has to introspect (the Day 1 approach).
func newProvider(store *base.MemoryStore) (*provider.Provider, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	jwks := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{
			KeyID:     signingKeyID,
			Key:       key,
			Algorithm: "RS256",
		}},
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
		provider.WithPKCE(goidc.CodeChallengeMethodSHA256),
		provider.WithRefreshTokenGrant(store),
		provider.WithRefreshTokenRotation(),
		provider.WithTokenAuthnMethods(goidc.AuthnMethodNone),
		// >>> The Day 2 change: issue JWT access tokens (RS256), valid 5 minutes. <<<
		provider.WithTokenOptions(func(_ context.Context, _ *goidc.Grant, _ *goidc.Client) goidc.TokenOptions {
			return goidc.NewJWTTokenOptions("RS256", 300)
		}),
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
