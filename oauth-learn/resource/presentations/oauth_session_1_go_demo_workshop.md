# OAuth Session 1 Workshop Deck

**Deck purpose**  
Teach OAuth Session 1 with the runnable Go demo app as the main example.

**Primary lab app**  
`resource/labs/day1-go-oauth-demo/`

**Export target**  
`export/presentations/oauth_session1_go_demo_workshop.pptx` (rebuilt with pptxgenjs; navy/teal brand matching the lesson PDFs).

# Demo Ground Truth

The deck must match the real app. Key facts:

- The app runs **three separate servers**: Authorization Server `:8080` (`op.go`, go-oidc), Resource Server API `:8081` (`api.go`), Client App `:8082` (`client.go`). **Start at http://localhost:8082.**
- `/authorize`, `/token`, `/introspect`, `/device_authorization` are served by **go-oidc** (`op.Handler()`), not hand-written functions.
- Client routes are unprefixed: `/start`, `/callback`, `/call-api`, `/refresh`, `/client-credentials`, `/device/start`, `/device/poll`.
- Real client handlers: `StartAuthCode`, `Callback`, `Refresh`, `CallAPI`, `ClientCredentials`, `DeviceStart`, `DevicePoll` (+ helpers `pkceChallenge`, `randomString`, `postForm`).
- Login/consent is a go-oidc `Policy` (`goidc.NewPolicy("simple_login", ...)`) in `op.go`, which sets `Subject="alya"` and grants scopes. go-oidc creates/stores the auth code internally.
- Resource server validates by calling `/introspect` (`introspect()` in `api.go`); handlers `apiProfile`, `apiFraudReport`, `apiDeploy`.
- The auth-code client requests `openid profile profile.read offline_access`, so go-oidc also issues an **id_token** — flag this as OpenID Connect (Session 2), per project decision.

# Teaching Flow

Session 1 is taught as a live walkthrough:

1. Explain the concept briefly.
2. Show the diagram.
3. Run the Go demo.
4. Open the relevant Go function in `client.go` / `op.go` / `api.go`.
5. Ask participants to explain the request, response, and token use.

# Slide Plan

| Slide | Title | Main Point | Demo Anchor (real) |
|---:|---|---|---|
| 1 | OAuth 2.0 — the real flow, in Go | Learn OAuth through a runnable go-oidc demo. | `go run .` → http://localhost:8082 |
| 2 | Explain the flow without hiding behind JWT | Session 1 outcomes. | Workshop plan |
| 3 | One Go program, three separate servers | Real network separation per role. | `main.go`; :8080 op.go, :8081 api.go, :8082 client.go |
| 4 | How we teach each flow | Concept, diagram, Go code, browser, checkpoint. | All flows |
| 5 | Limited access without the password | RepoBoard needs scoped access, not a password. | RepoBoard scenario |
| 6 | Map the four actors onto the code | Route prefixes/ports make actors visible. | :8082 client.go, :8080 op.go, :8081 api.go |
| 7 | /authorize and /token are not interchangeable | Front-channel vs back-channel; both served by go-oidc. | `op.Handler()` |
| 8 | Authorization Code + PKCE overview | Browser gets a code; client gets tokens from `/token`. | Click “Start Auth Code + PKCE flow” |
| 9 | Client builds the request before redirecting | state + PKCE created client-side. | `StartAuthCode` (client.go) |
| 10 | Authorization request parameters | Explain every query parameter. | `/authorize?...` |
| 11 | Login & consent is a go-oidc policy | AS approval is a `Policy`, not hand-written endpoints. | `goidc.NewPolicy` (op.go) |
| 12 | Verify state, then exchange the code | Callback checks `state`, then posts to `/token`. | `Callback` (client.go) |
| 13 | go-oidc verifies the exchange for you | PKCE/code checks are server config, not client code. | `WithAuthCodeGrant`, `WithPKCE` (op.go) |
| 14 | Why redirect_uri goes to /token | Binds the exchange to the original request. | conceptual |
| 15 | PKCE mental model | Challenge first, verifier later. | `pkceChallenge` (client.go) |
| 16 | Access token at the API | Resource server validates via introspection. | `apiProfile` + `introspect` (api.go) |
| 17 | Refresh token flow | Renew access; rotation on; needs `offline_access`. | `Refresh` (client.go) |
| 18 | Client Credentials flow | Service-to-service, no user/redirect. | `ClientCredentials` (client.go) |
| 19 | Device Authorization flow | Approve on a second device while CLI polls. | `DeviceStart`, `DevicePoll` (client.go) |
| 20 | Choose the flow | Match app shape to flow. | Flow comparison |
| 21 | OAuth vs OpenID Connect | OAuth = API access; OIDC = login identity. | Heads-up: demo issues an id_token (OIDC, Session 2) |
| 22 | Common Session 1 mistakes | Don’t mix tokens, actors, endpoints. | Review |
| 23 | Instructor demo checklist | Exact live-run order. | Buttons on :8082 |
| 24 | Session 1 checkpoint | Explain the chain from memory. | Final discussion |

# Notes

Keep slide copy short. The instructor should spend time in the running app (`:8082`) and the Go source, not reading paragraphs from slides. Code shown on slides 9, 11, 12, 13, 16 is taken directly from the demo so the editor matches the slide.
