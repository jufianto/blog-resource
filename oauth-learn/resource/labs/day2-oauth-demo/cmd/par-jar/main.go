// Command par-jar is the Day 2 (Lesson 7) example: Pushed Authorization Requests (PAR)
// and JWT-Secured Authorization Requests (JAR). The client sends its authorization request
// as a signed request object, pushed directly to the AS back-channel; the browser only ever
// carries an opaque request_uri. Run `go run ./cmd/par-jar`, then open http://localhost:8082.
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
	signingKeyID = "as_key"      // Authorization Server's token-signing key id
	clientKeyID  = "client_key"  // Client's request-object-signing key id
)

func main() {
	// The client's signing key: the private half signs request objects (JAR); the public
	// half is registered with the AS so it can verify them. Same process, so we share it.
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating client JAR key: %v", err)
	}

	store := base.NewMemoryStore()
	op, err := newProvider(store, &clientKey.PublicKey)
	if err != nil {
		log.Fatalf("failed to init OpenID Provider: %v", err)
	}
	clientApp := NewClientApp(clientKey)

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

	log.Println("Day 2 — PAR + JAR:")
	log.Println("  1. Authorization Server -> http://localhost:8080 (PAR + signed request objects required)")
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
