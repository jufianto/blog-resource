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
    ├── dpop/                       # Lesson 7      (DONE)
    ├── par-jar/                    # Lesson 7      (DONE)
    └── fapi/                       # Lesson 8      (DONE — capstone)
```

A participant lab guide covering the built examples is in `LAB.md`
(build to PDF with `make lab2-pdf` from the `oauth-learn/` root).

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

## `cmd/dpop` — Lesson 7 (sender-constrained tokens)

**The one idea:** the access token is bound to a key the client holds (`cnf.jkt`), so a
stolen token is useless without the private key.

- AS: `WithDPoP(goidc.RS256)` — binds the token to the DPoP proof key (`op.go`).
- Client: generates its own key, sends a fresh DPoP proof on the token request and on each API
  call (`proof.go`, `client.go`).
- RS: validates the JWT locally, then validates the DPoP proof and that its key thumbprint
  matches the token's `cnf.jkt` (`proof.go`, `rs.go`).
- UI shows: `token_type: DPoP`; a correct call (token + proof) → 200; a **replay** as a plain
  Bearer with no proof → 401.

## `cmd/par-jar` — Lesson 7 (request hardening)

**The one idea:** the authorization request is a **signed request object (JAR)** **pushed**
to the AS back-channel first (**PAR**); the browser only carries an opaque `request_uri`.

- AS: `WithPARRequired` + `WithJARRequired(goidc.RS256)`, and the client is registered with
  its public JWKS so the AS can verify the signed request (`op.go`).
- Client: signs the request object with its key (`request.go`), pushes it to `/par`, then
  sends the browser to `/authorize` with the opaque `request_uri` (`client.go`).
- UI shows the decoded signed request, the PAR response, and the clean `/authorize` URL. A
  plain `/authorize` without PAR is rejected (`400`).

## `cmd/fapi` — Lesson 8 (capstone)

**The one idea:** FAPI 2.0 is not one feature — it is a *profile* that requires several
controls together, and `WithProfile(goidc.ProfileFAPI2)` rejects any non-compliant config.

- AS: `WithProfile(ProfileFAPI2)` + `WithPARRequired` + `WithPKCERequired(S256)` +
  `WithDPoPRequired` + `private_key_jwt` client auth + auth-code lifetime < 60s (`op.go`).
- Client: authenticates with a signed `private_key_jwt` assertion (`assertion.go`), pushes the
  request via PAR, and sends DPoP proofs on the token request and API calls
  (`proof.go`, `client.go`).
- RS: validates the JWT locally *and* the DPoP proof/binding (`rs.go`).
- Verified end-to-end: PAR → authorize → `private_key_jwt` token exchange → DPoP-bound token →
  API call `200`.

## Not demonstrated as code (see the lesson papers)

To keep each example focused, these Lesson 5/7 topics are covered in the papers but not as
runnable demos here: **JARM** (signed authorization responses), **mTLS** sender-constraining
(the DPoP alternative), and **JWE** (encrypted JOSE).
