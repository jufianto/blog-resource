package main

import (
	"log"
	"net/http"
)

const opBaseURL = "http://localhost:8080"
const apiBaseURL = "http://localhost:8081"
const clientBaseURL = "http://localhost:8082"

func main() {
	store := NewMemoryStore()

	op, err := newProvider(store)
	if err != nil {
		log.Fatalf("failed to init OpenID Provider: %v", err)
	}

	clientApp := NewClientApp()

	// Authorization Server — go-oidc serves /authorize, /token, /introspect,
	// /device_authorization, /.well-known/openid-configuration, /jwks, etc.
	opMux := http.NewServeMux()
	opMux.Handle("/", logMiddleware(op.Handler(), "AuthServer (8080)"))

	// Resource Server
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/profile", apiProfile)
	apiMux.HandleFunc("/api/fraud-report", apiFraudReport)
	apiMux.HandleFunc("/api/deploy", apiDeploy)

	// Client App
	clientMux := http.NewServeMux()
	clientMux.HandleFunc("/", clientApp.Home)
	clientMux.HandleFunc("/start", clientApp.StartAuthCode)
	clientMux.HandleFunc("/callback", clientApp.Callback)
	clientMux.HandleFunc("/call-api", clientApp.CallAPI)
	clientMux.HandleFunc("/refresh", clientApp.Refresh)
	clientMux.HandleFunc("/client-credentials", clientApp.ClientCredentials)
	clientMux.HandleFunc("/device/start", clientApp.DeviceStart)
	clientMux.HandleFunc("/device/poll", clientApp.DevicePoll)
	clientMux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		store.Reset()
		clientApp.ResetSessions()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Println("Starting 3 separate servers for the Day 1 OAuth Demo:")
	log.Println("  1. Authorization Server -> http://localhost:8080")
	log.Println("  2. Resource Server API  -> http://localhost:8081")
	log.Println("  3. Client Application   -> http://localhost:8082")
	log.Println(" ")
	log.Println("👉 START HERE: http://localhost:8082")

	// Run them all concurrently
	go func() {
		log.Fatal(http.ListenAndServe(":8080", opMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8081", logMiddleware(apiMux, "ResourceAPI (8081)")))
	}()

	log.Fatal(http.ListenAndServe(":8082", logMiddleware(clientMux, "ClientApp (8082)")))
}

func logMiddleware(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", name, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
