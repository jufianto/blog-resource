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

// newProvider registers the client WITH its public key (so the AS can verify the client's
// signed request objects) and requires both PAR and JAR. clientPub is the public half of the
// key the client signs request objects with.
func newProvider(store *base.MemoryStore, clientPub *rsa.PublicKey) (*provider.Provider, error) {
	asKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating AS signing key: %w", err)
	}
	asJWKS := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{KeyID: signingKeyID, Key: asKey, Algorithm: "RS256"}},
	}

	// The client's registered public JWKS — how the AS verifies JAR request objects.
	clientJWKS := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{
			KeyID:     clientKeyID,
			Key:       clientPub,
			Algorithm: "RS256",
			Use:       "sig",
		}},
	}

	clientRepoboard := &goidc.Client{
		ID: clientID,
		ClientMeta: goidc.ClientMeta{
			RedirectURIs:     []string{clientBaseURL + "/callback"},
			ResponseTypes:    []goidc.ResponseType{goidc.ResponseTypeCode},
			GrantTypes:       []goidc.GrantType{goidc.GrantAuthorizationCode, goidc.GrantRefreshToken},
			ScopeIDs:         "openid profile profile.read offline_access",
			TokenAuthnMethod: goidc.AuthnMethodNone,
			JWKS:             &clientJWKS,
			JARSigAlg:        goidc.RS256,
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
		provider.WithStaticClients(clientRepoboard),
		provider.WithAuthCodeGrant(store, goidc.ResponseTypeCode),
		provider.WithPKCERequired(goidc.CodeChallengeMethodSHA256),
		provider.WithRefreshTokenGrant(store),
		provider.WithRefreshTokenRotation(),
		provider.WithTokenAuthnMethods(goidc.AuthnMethodNone),
		provider.WithTokenOptions(func(_ context.Context, _ *goidc.Grant, _ *goidc.Client) goidc.TokenOptions {
			return goidc.NewJWTTokenOptions("RS256", 300)
		}),
		// >>> Lesson 7: require PAR (push the request back-channel) and JAR (sign it). <<<
		provider.WithPARRequired(store),
		provider.WithJARRequired(goidc.RS256),
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
