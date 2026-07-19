package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func apiProfile(w http.ResponseWriter, r *http.Request) {
	info, ok := introspect(r, "profile.read")
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid_token_or_scope"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"api":       "profile",
		"subject":   info.Sub,
		"client_id": info.ClientID,
		"scopes":    info.Scope,
		"message":   "RepoBoard can read this user's profile-like resource.",
	})
}

func apiFraudReport(w http.ResponseWriter, r *http.Request) {
	info, ok := introspect(r, "fraud.check")
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "service_token_required"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"api":       "fraud-report",
		"client_id": info.ClientID,
		"scopes":    info.Scope,
		"message":   "Order service called fraud API using Client Credentials.",
	})
}

func apiDeploy(w http.ResponseWriter, r *http.Request) {
	info, ok := introspect(r, "deploy.write")
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid_token_or_scope"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"api":       "deploy",
		"subject":   info.Sub,
		"client_id": info.ClientID,
		"scopes":    info.Scope,
		"message":   "Deploy API accepted token from Device Authorization Flow.",
	})
}

type introspectResponse struct {
	Active   bool   `json:"active"`
	Scope    string `json:"scope"`
	ClientID string `json:"client_id"`
	Sub      string `json:"sub"`
}

func introspect(r *http.Request, requiredScope string) (*introspectResponse, bool) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, false
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")

	form := url.Values{}
	form.Set("token", tokenStr)
	req, _ := http.NewRequest(http.MethodPost, opBaseURL+"/introspect", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("order-service", "order-service-secret")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, false
	}
	defer resp.Body.Close()

	var ir introspectResponse
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil || !ir.Active {
		return nil, false
	}

	for _, s := range strings.Fields(ir.Scope) {
		if s == requiredScope {
			return &ir, true
		}
	}
	return nil, false
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
