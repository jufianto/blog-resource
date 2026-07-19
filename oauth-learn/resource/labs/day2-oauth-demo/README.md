# Day 2 — Advanced OAuth / JOSE / FAPI Demos

Day 2 covers Lessons 4–8: how tokens are actually validated, JOSE (JWS/JWE/JWK/JWKS),
security controls, and the sender-constrained / FAPI 2.0 controls (PAR, JAR, JARM, DPoP,
mTLS). Where Day 1 taught the **flows**, Day 2 teaches what makes them **trustworthy**.

## Structure — one module, one focused command per concept

This is a single Go module. Each concept is its own runnable entrypoint under `cmd/`, so
every example concentrates on exactly one idea. Shared, concept-neutral plumbing (the
in-memory go-oidc store and small web/OAuth-client helpers) lives in `internal/base`, so
the concept-specific code in each `cmd/` stays short and readable.

```
day2-oauth-demo/
├── go.mod  go.sum  vendor/         # deps vendored (same as Day 1); no network needed
├── internal/base/                  # shared: MemoryStore + web/PKCE/HTTP helpers
└── cmd/
    ├── jwt-validation/             # Lessons 4–5  (DONE)
    ├── dpop/                       # Lesson 7      (planned)
    ├── par-jar/                    # Lesson 7      (planned)
    └── fapi/                       # Lesson 8      (planned)
```

Run any example from the module root:

```bash
cd resource/labs/day2-oauth-demo
go run ./cmd/jwt-validation      # then open http://localhost:8082
```

Each command starts the same three-actor setup as Day 1 — Authorization Server (`:8080`),
Resource Server (`:8081`), Client (`:8082`) — so the mental model carries straight over.
Enter the flow at **<http://localhost:8082>**.

## `cmd/jwt-validation` — Lessons 4 & 5

**The one idea:** a JWT access token is self-contained and independently verifiable, so the
Resource Server validates it **locally** instead of introspecting it at the Authorization
Server (the Day 1 approach).

What changes from Day 1:

- The AS issues **JWT** access tokens — one line: `WithTokenOptions` returning
  `goidc.NewJWTTokenOptions("RS256", 300)` (`cmd/jwt-validation/op.go`).
- The RS (`cmd/jwt-validation/rs.go`) fetches the AS's public keys from `/jwks`, then for
  each request verifies **signature → issuer → expiry → scope** with go-jose. No round trip.

Walk participants through, in the UI:

1. **Run Authorization Code + PKCE** → the token page decodes the JWT header/payload so they
   can see the claims are readable (Base64URL), and that the **signature** is what matters.
2. **Call resource API (valid token)** → `200`, "validated locally, via JWKS".
3. **Tamper with the token** → one payload byte is flipped without re-signing; local
   validation rejects it (`401`, "signature check failed") with no call to the AS.

Discussion questions:

- If anyone can read a JWT's payload, why is it safe to send? What actually makes it trusted?
- The RS never contacts the AS to validate. What does it gain, and what does it give up
  (think revocation and key rotation)?
- Why must the RS pin the accepted algorithm (`RS256`) instead of trusting the token's own
  `alg` header?
- When would you still prefer introspection (Day 1) over local JWT validation?

> Note: this demo keeps validation to signature/issuer/expiry/scope to stay focused.
> Audience (`aud`) binding and JWKS key **rotation** are natural extensions — the RS already
> refetches the JWKS when it sees an unknown `kid`.
