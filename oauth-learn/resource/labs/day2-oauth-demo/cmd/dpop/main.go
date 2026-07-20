// Command dpop is the Day 2 (Lesson 7) example: DPoP sender-constrained access tokens.
// The Authorization Server binds each access token to a key the client proves it holds,
// so a stolen token is useless without the matching private key. Run it with
// `go run ./cmd/dpop` from the module root, then open http://localhost:8082.
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

	profileURL = apiBaseURL + "/api/profile"
)

func main() {
	store := base.NewMemoryStore()

	op, err := newProvider(store)
	if err != nil {
		log.Fatalf("failed to init OpenID Provider: %v", err)
	}

	clientApp, err := NewClientApp()
	if err != nil {
		log.Fatalf("failed to init client (DPoP key): %v", err)
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
	clientMux.HandleFunc("/replay", clientApp.Replay)
	clientMux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		store.Reset()
		clientApp.ResetSessions()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Println("Day 2 — DPoP sender-constrained tokens:")
	log.Println("  1. Authorization Server -> http://localhost:8080 (binds tokens to the client's DPoP key)")
	log.Println("  2. Resource Server API  -> http://localhost:8081 (requires a valid DPoP proof per request)")
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
