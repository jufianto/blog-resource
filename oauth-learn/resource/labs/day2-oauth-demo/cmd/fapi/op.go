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

// newProvider builds a FAPI 2.0 baseline Authorization Server. WithProfile(ProfileFAPI2)
// makes go-oidc reject any configuration that is not FAPI-compliant, so this option set is
// essentially the minimum FAPI 2.0 requires: PAR + PKCE(S256) + DPoP-required +
// private_key_jwt client auth + short-lived authorization codes.
func newProvider(store *base.MemoryStore, clientPub *rsa.PublicKey) (*provider.Provider, error) {
	asKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating AS signing key: %w", err)
	}
	asJWKS := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{KeyID: signingKeyID, Key: asKey, Algorithm: "RS256"}},
	}

	// The client authenticates with private_key_jwt, so it registers a public key (JWKS)
	// and no secret.
	clientJWKS := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{KeyID: clientKeyID, Key: clientPub, Algorithm: "RS256", Use: "sig"}},
	}
	clientRepoboard := &goidc.Client{
		ID: clientID,
		ClientMeta: goidc.ClientMeta{
			RedirectURIs:     []string{clientBaseURL + "/callback"},
			ResponseTypes:    []goidc.ResponseType{goidc.ResponseTypeCode},
			GrantTypes:       []goidc.GrantType{goidc.GrantAuthorizationCode},
			ScopeIDs:         "profile.read",
			TokenAuthnMethod: goidc.AuthnMethodPrivateKeyJWT,
			JWKS:             &clientJWKS,
		},
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
		func(_ context.Context) (goidc.JSONWebKeySet, error) { return asJWKS, nil },
		provider.WithProfile(goidc.ProfileFAPI2),
		provider.WithStaticClients(clientRepoboard),
		provider.WithAuthCodeGrant(store, goidc.ResponseTypeCode),
		provider.WithAuthCodeLifetime(30), // FAPI 2.0: authorization code lifetime < 60s
		provider.WithPKCERequired(goidc.CodeChallengeMethodSHA256),
		provider.WithPARRequired(store),
		provider.WithPARLifetime(60),
		provider.WithDPoPRequired(goidc.RS256),
		provider.WithTokenAuthnMethods(goidc.AuthnMethodPrivateKeyJWT),
		provider.WithPrivateKeyJWTSignatureAlgs(goidc.RS256),
		provider.WithTokenOptions(func(_ context.Context, _ *goidc.Grant, _ *goidc.Client) goidc.TokenOptions {
			return goidc.NewJWTTokenOptions("RS256", 300)
		}),
		provider.WithHandleErrorFunc(func(_ context.Context, err error) {
			fmt.Printf("OIDC ERROR: %v\n", err)
		}),
		provider.WithPolicies(loginPolicy),
		provider.WithScopes(goidc.NewScope("profile.read")),
	)
}
