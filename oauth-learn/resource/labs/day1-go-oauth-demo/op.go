package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/luikyv/go-oidc/pkg/goidc"
	"github.com/luikyv/go-oidc/pkg/provider"
)

var (
	clientRepoboard = &goidc.Client{
		ID: "repoboard-web",
		ClientMeta: goidc.ClientMeta{
			RedirectURIs:     []string{clientBaseURL + "/callback"},
			ResponseTypes:    []goidc.ResponseType{goidc.ResponseTypeCode},
			GrantTypes:       []goidc.GrantType{goidc.GrantAuthorizationCode, goidc.GrantRefreshToken},
			ScopeIDs:         "openid profile profile.read offline_access",
			TokenAuthnMethod: goidc.AuthnMethodNone,
		},
	}
	clientOrderService = &goidc.Client{
		ID:     "order-service",
		Secret: "order-service-secret",
		ClientMeta: goidc.ClientMeta{
			GrantTypes:                    []goidc.GrantType{goidc.GrantClientCredentials},
			ScopeIDs:                      "fraud.check",
			TokenAuthnMethod:              goidc.AuthnMethodSecretBasic,
			TokenIntrospectionAuthnMethod: goidc.AuthnMethodSecretBasic,
		},
	}
	clientDeployCLI = &goidc.Client{
		ID: "deploy-cli",
		ClientMeta: goidc.ClientMeta{
			GrantTypes:       []goidc.GrantType{goidc.GrantDeviceCode},
			ScopeIDs:         "deploy.write",
			TokenAuthnMethod: goidc.AuthnMethodNone,
		},
	}
)

func newProvider(store *MemoryStore) (*provider.Provider, error) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwks := goidc.JSONWebKeySet{
		Keys: []goidc.JSONWebKey{{
			KeyID:     "key_id",
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
				
				// Determine callback path based on flow
				callbackPath := "/authorize/" + as.ID
				if as.DeviceCode != "" {
					callbackPath = "/device/" + as.ID
				}

				fmt.Fprintf(w, `<!DOCTYPE html><html><body>
<h2>Mock Login &amp; Consent</h2>
<p>Client: <strong>%s</strong></p>
<p>Scopes requested: <strong>%s</strong></p>
<form method="POST" action="%s">
  <input type="hidden" name="approve" value="true">
  <button type="submit">Approve for Alya Developer</button>
</form>
</body></html>`, html.EscapeString(as.ClientID), html.EscapeString(as.Scopes), html.EscapeString(callbackPath))
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
		provider.WithStaticClients(clientRepoboard, clientOrderService, clientDeployCLI),
		provider.WithAuthCodeGrant(store, goidc.ResponseTypeCode),
		provider.WithPKCE(goidc.CodeChallengeMethodSHA256),
		provider.WithRefreshTokenGrant(store),
		provider.WithRefreshTokenRotation(),
		provider.WithClientCredentialsGrant(),
		provider.WithTokenAuthnMethods(goidc.AuthnMethodNone, goidc.AuthnMethodSecretPost, goidc.AuthnMethodSecretBasic),
		provider.WithDeviceGrant(store, renderUserCodePage, renderDeviceApprovedPage),
		provider.WithTokenIntrospection(func(_ context.Context, _ *goidc.Client, _ goidc.TokenInfo) bool { return true }),
		provider.WithRefreshTokenShouldIssueFunc(func(_ context.Context, _ *goidc.Client, grant *goidc.Grant) bool {
			// Require offline_access scope to issue refresh tokens
			return strings.Contains(grant.Scopes, goidc.ScopeOfflineAccess.ID)
		}),
		provider.WithHandleErrorFunc(func(_ context.Context, err error) {
			fmt.Printf("OIDC ERROR: %v\n", err)
		}),
		provider.WithPolicies(loginPolicy),
		provider.WithScopes(
			goidc.ScopeOpenID, goidc.ScopeProfile, goidc.ScopeOfflineAccess,
			goidc.NewScope("profile.read"), goidc.NewScope("fraud.check"), goidc.NewScope("deploy.write"),
		),
	)
}

func renderUserCodePage(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!DOCTYPE html><html><body>
<h2>Device Verification</h2>
<form method="GET" action="/device">
  <label>Enter the code shown on your device:</label><br>
  <input type="text" name="user_code" required style="font-size:1.4em;letter-spacing:.2em">
  <button type="submit">Verify</button>
</form>
</body></html>`)
	return nil
}

func renderDeviceApprovedPage(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!DOCTYPE html><html><body>
<h2>Device approved!</h2>
<p>You may now return to your device.</p>
</body></html>`)
	return nil
}
