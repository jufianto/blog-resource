// Command jwt-validation is the Day 2 (Lessons 4-5) example: the Authorization Server
// issues JWT access tokens, and the Resource Server validates them locally against the
// AS's JWKS — no introspection round-trip. Run it with `go run ./cmd/jwt-validation`
// from the day2-oauth-demo module root, then open http://localhost:8082.
package main

import (
	"log"
	"net/http"

	"day2-oauth-demo/internal/base"
)

const (
	opBaseURL     = "http://localhost:8080"
	apiBaseURL    = "http://localhost:8081"
	clientBaseURL = "http://localhost:8082"

	clientID     = "repoboard-web"
	signingKeyID = "key_id"
)

func main() {
	store := base.NewMemoryStore()

	op, err := newProvider(store)
	if err != nil {
		log.Fatalf("failed to init OpenID Provider: %v", err)
	}

	clientApp := NewClientApp()

	// Authorization Server (:8080) — go-oidc serves /authorize, /token, /jwks,
	// /.well-known/openid-configuration, etc.
	opMux := http.NewServeMux()
	opMux.Handle("/", logMiddleware(op.Handler(), "AuthServer (8080)"))

	// Resource Server (:8081) — validates JWTs locally.
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/profile", apiProfile)

	// Client Application (:8082).
	clientMux := http.NewServeMux()
	clientMux.HandleFunc("/", clientApp.Home)
	clientMux.HandleFunc("/start", clientApp.StartAuthCode)
	clientMux.HandleFunc("/callback", clientApp.Callback)
	clientMux.HandleFunc("/call-api", clientApp.CallAPI)
	clientMux.HandleFunc("/tamper", clientApp.Tamper)
	clientMux.HandleFunc("/refresh", clientApp.Refresh)
	clientMux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		store.Reset()
		clientApp.ResetSessions()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Println("Day 2 — JWT validation demo:")
	log.Println("  1. Authorization Server -> http://localhost:8080 (issues JWT access tokens)")
	log.Println("  2. Resource Server API  -> http://localhost:8081 (validates JWTs locally via JWKS)")
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
