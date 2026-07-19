# Day 1 — OAuth 2.0 Go Demo

A single runnable Go application that shows **three real OAuth 2.0 flows** side by side. It is split into five focused files so you can read each flow in isolation without getting lost in one giant file.

## Quick Start

```bash
go run .
```

Then open **<http://localhost:8080/client>** in your browser.

> Requires Go 1.21+. No Docker. No external services. Everything runs in memory.

---

## What Is Running

One process serves three logical roles on the same port:

| Role | Path prefix | Description |
|---|---|---|
| **Authorization Server** | `/authorize`, `/token`, `/introspect`, `/device_authorization`, `/.well-known/*` | Powered by [go-oidc](https://github.com/luikyv/go-oidc) — a spec-compliant OpenID Provider |
| **Client Application** | `/client/*` | Simulates a web app (RepoBoard), a backend service, and a CLI tool |
| **Resource Server** | `/api/*` | Three protected API endpoints that validate tokens via introspection |

Keeping all three in one process lets you trace a complete flow end-to-end without configuring multiple services. In production each role would be a separate service.

---

## File Structure

```
day1-go-oauth-demo/
├── main.go    — entry point, route wiring, HTTP server
├── op.go      — Authorization Server: provider config, clients, login/consent policy
├── store.go   — in-memory storage implementing go-oidc's manager interfaces
├── client.go  — Client App: all four OAuth flows + session management + HTML rendering
└── api.go     — Resource Server: three protected endpoints + token introspection helper
```

---

## File-by-File Explanation

### `main.go` — Entry Point

Initialises the three parts and wires their HTTP handlers onto a single `ServeMux`:

```
store  →  op (Authorization Server)
                ↘
       mux ←── client (Client App)
                ↗
          api (Resource Server)
```

The go-oidc `Handler()` is mounted at `/` to catch all standard OAuth endpoints. Client and API routes are registered on top and take priority because of Go's `ServeMux` longest-match rule.

---

### `op.go` — Authorization Server

Uses **go-oidc** to spin up a real OpenID Provider in ~80 lines.

**Three pre-registered clients:**

| Client ID | Type | Grant types | Used for |
|---|---|---|---|
| `repoboard-web` | Public (no secret) | `authorization_code`, `refresh_token` | Auth Code + PKCE demo |
| `order-service` | Confidential | `client_credentials` | Client Credentials demo |
| `deploy-cli` | Public (no secret) | `device_code` | Device Authorization demo |

**Login / Consent Policy:**

go-oidc does not hard-code how users authenticate. You plug in an `AuthnPolicy` that receives the HTTP request and the active `AuthnSession`. This demo uses a single policy that:

1. On `GET /authorize` — renders a simple HTML form showing which client is requesting access and which scopes it wants.
2. On `POST` (form submission) — sets `Subject = "alya"` and `GrantedScopes` to approve the request, returning `StatusSuccess`.

This is exactly the integration point where a real app would check a username/password, call an MFA service, or verify a session cookie.

**Enabled capabilities:**

- `WithAuthCodeGrant` — authorization code flow with PKCE support
- `WithRefreshTokenGrant` — refresh token issuance and rotation
- `WithClientCredentialsGrant` — machine-to-machine tokens
- `WithDeviceGrant` — device authorization flow (CLI/TV devices)
- `WithTokenIntrospection` — allows the Resource Server to validate tokens

---

### `store.go` — In-Memory Storage

go-oidc is storage-agnostic. It defines interfaces (`GrantManager`, `AuthManager`, `RefreshTokenManager`, `DeviceAuthManager`, `DCRManager`) and you provide the implementation.

`MemoryStore` implements all of them using plain Go maps protected by `sync.RWMutex`. This is intentionally simple — for a real deployment you would replace this with a database (Postgres, Redis, etc.) while keeping the same interface.

**Key lookups implemented:**

| Method | Used by |
|---|---|
| `GrantByAuthCode` | Token endpoint exchanges the code for tokens |
| `GrantByRefreshToken` | Refresh token grant looks up the original grant |
| `SessionByUserCode` | Device verification page finds the pending session |
| `SessionByDeviceCode` | Device token polling resolves the grant |
| `GrantByDeviceCode` | Device token endpoint issues the access token |

---

### `client.go` — Client Application

Simulates three different OAuth clients in one file. Each flow is a self-contained set of HTTP handlers.

**Flow 1 — Authorization Code + PKCE** (`/client/start` → `/client/callback`)

```
Browser                     Client App               Authorization Server
   |                            |                            |
   |-- GET /client/start ------>|                            |
   |                            |-- generate state, PKCE --> |
   |<-- 302 /authorize?... -----|                            |
   |                                                         |
   |-- GET /authorize ---------------------------------------->
   |<-- 200 login/consent form --------------------------------
   |-- POST (approve) ---------------------------------------->
   |<-- 302 /client/callback?code=...&state=... -------------
   |                            |                            |
   |-- GET /client/callback --->|                            |
   |                            |-- POST /token (code + verifier) -->
   |                            |<-- access_token, refresh_token ---
   |<-- 200 token shown --------|
```

What travels through the browser: the authorization code and state in the redirect URL.  
What is back-channel: the code-for-token exchange at `/token`.  
PKCE protects this exchange — the `code_verifier` is never sent to the browser.

**Flow 2 — Client Credentials** (`/client/client-credentials`)

```
order-service                Authorization Server
      |-- POST /token (client_id + secret) -->
      |<-- access_token -----------------------
      |
      |-- GET /api/fraud-report (Bearer token) --> Resource Server
      |<-- 200 JSON response -------------------
```

No browser. No user. The service authenticates with its own credentials and gets a token scoped to `fraud.check`.

**Flow 3 — Refresh Token** (`/client/refresh`)

```
Client App                   Authorization Server
      |-- POST /token (refresh_token) -->
      |<-- new access_token, new refresh_token --
```

The demo uses the token stored in the browser session from the Auth Code flow. Refresh tokens are single-use — the old one is invalidated and a new one is issued.

**Flow 4 — Device Authorization** (`/client/device/start` → `/client/device/poll`)

```
CLI (deploy-cli)             Authorization Server           User Browser
      |-- POST /device_authorization -->                        |
      |<-- device_code, user_code, verification_uri ---         |
      |                                                         |
      |   (display user_code to user)                          |
      |                               <-- GET /device_verify ---
      |                               --> show code input form -|
      |                               <-- POST user_code -------
      |                               --> approved             |
      |                                                         |
      |-- POST /token (device_code, poll) -->                  |
      |<-- access_token ------------------                      |
```

The CLI never touches the browser. The user approves on a separate device. The CLI polls until the token is ready.

**Session management:** A cookie (`session_id`) maps the browser to an in-memory `Session` struct that holds the state, code verifier, and current tokens between requests.

---

### `api.go` — Resource Server

Three protected endpoints. Each one:

1. Reads the `Authorization: Bearer <token>` header.
2. Calls `introspect()` which posts the token to `/introspect` using `order-service` credentials.
3. Checks `active: true` and that the required scope is present.
4. Returns the JSON response.

This is the correct decoupled pattern — the Resource Server does **not** share memory with the Authorization Server. It validates tokens by calling a standard endpoint, exactly as a real separate service would.

```
Client → GET /api/profile (Bearer token)
Resource Server → POST /introspect (token)
Authorization Server → { active: true, scope: "profile.read", sub: "alya" }
Resource Server → 200 { api: "profile", subject: "alya", ... }
```

---

## Discussion Questions

- What would happen if you removed `code_challenge` from the authorization request?
- Why does the Resource Server use `order-service` credentials to call `/introspect`? Could it use no credentials?
- What is the difference between the `sub` claim in a user token vs a client credentials token?
- What can an attacker do if they intercept the authorization code in the redirect URL? Why does PKCE mitigate this?
- What happens when the refresh token expires? What should the client do?
