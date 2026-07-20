// Command fapi is the Day 2 (Lesson 8) capstone: a FAPI 2.0 baseline profile that combines
// everything from the earlier examples under go-oidc's strict WithProfile(FAPI2):
//   - PAR required (request pushed to the back-channel)
//   - PKCE S256 required
//   - DPoP-sender-constrained access tokens required
//   - private_key_jwt client authentication (no client secret, no public client)
// Run `go run ./cmd/fapi`, then open http://localhost:8082.
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"

	"day2-oauth-demo/internal/base"
)

const (
	opBaseURL     = "http://localhost:8080"
	apiBaseURL    = "http://localhost:8081"
	clientBaseURL = "http://localhost:8082"

	clientID     = "repoboard-web"
	signingKeyID = "as_key"     // AS token-signing key
	clientKeyID  = "client_key" // client's private_key_jwt signing key

	profileURL = apiBaseURL + "/api/profile"
)

func main() {
	// The client's authentication key: private half signs client_assertion JWTs; public half
	// is registered with the AS. Same process, so we share it.
	authKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating client auth key: %v", err)
	}

	store := base.NewMemoryStore()
	op, err := newProvider(store, &authKey.PublicKey)
	if err != nil {
		log.Fatalf("failed to init OpenID Provider: %v", err)
	}
	clientApp, err := NewClientApp(authKey)
	if err != nil {
		log.Fatalf("failed to init client: %v", err)
	}

	opMux := http.NewServeMux()
	opMux.Handle("/", logMiddleware(op.Handler(), "AuthServer (8080)"))

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/profile", apiProfile)

	clientMux := http.NewServeMux()
	clientMux.HandleFunc("/", clientApp.Home)
	clientMux.HandleFunc("/start", clientApp.StartAuthCode)
	clientMux.HandleFunc("/callback", clientApp.Callback)
	clientMux.HandleFunc("/call-api", clientApp.CallAPI)
	clientMux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		store.Reset()
		clientApp.ResetSessions()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Println("Day 2 — FAPI 2.0 capstone (PAR + PKCE + DPoP + private_key_jwt):")
	log.Println("  1. Authorization Server -> http://localhost:8080")
	log.Println("  2. Resource Server API  -> http://localhost:8081")
	log.Println("  3. Client Application   -> http://localhost:8082")
	log.Println(" ")
	log.Println("👉 START HERE: http://localhost:8082")

	go func() { log.Fatal(http.ListenAndServe(":8080", opMux)) }()
	go func() { log.Fatal(http.ListenAndServe(":8081", logMiddleware(apiMux, "ResourceAPI (8081)"))) }()
	log.Fatal(http.ListenAndServe(":8082", logMiddleware(clientMux, "ClientApp (8082)")))
}

func logMiddleware(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", name, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
