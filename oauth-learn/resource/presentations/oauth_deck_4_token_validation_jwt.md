# Deck 4: Token Validation and JWT Basics

**Workshop position**  
Day 2, Session 1.

**Source paper**  
Lesson 4: Token Validation and JWT Basics.

**Estimated teaching time**  
75 minutes including Exercise 1.

# Teaching Goal

Participants should understand that a JWT is only a data format, that security comes from
*validation*, and that a resource server can validate a JWT locally against a JWKS instead of
calling the authorization server. Decoding is not validating.

# Audience Assumption

Participants often think "the token is encrypted" or "if I can read the claims it must be
trusted". They may not know the difference between an opaque token (validated by
introspection) and a JWT (validated by signature).

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | A Token Is a Claim, Not a Guarantee | Anyone can send a string; validation is what matters. | Title |
| 2 | Concept | Opaque vs JWT Access Tokens | Opaque → introspect; JWT → verify signature locally. | Two-column compare |
| 3 | Diagram | Anatomy of a JWT | Header, payload, signature; Base64URL is readable, not secret. | `jwt_structure.png` |
| 4 | Concept | Who Validates the Token? | The resource server, on every request. | Actor recap |
| 5 | Diagram | Local Validation Against JWKS | Fetch keys once, verify signature, then claims. | `token_validation_sequence.png` |
| 6 | Concept | The Four Checks | Signature, issuer, expiry, scope — in that order. | Checklist |
| 7 | Demo | Decode vs Verify | Reading the payload proves nothing; the signature does. | `cmd/jwt-validation` |
| 8 | Demo | Tampered Token Rejected | Flip one byte → signature check fails, no AS call. | `cmd/jwt-validation` tamper |
| 9 | Risk | Trusting `alg` | Never trust the token's own `alg`; pin RS256. | `alg: none` warning |
| 10 | Comparison | Local Validation vs Introspection | Speed vs revocation/rotation trade-off. | Table |
| 11 | Checkpoint | Validate This Token | Participants list what to check and why. | Exercise 1 |
| 12 | Summary | Decoding Is Not Validation | Signature + issuer + expiry + scope. | Recap |

# Implementation Walkthrough Notes

Run `go run ./cmd/jwt-validation` from `resource/labs/day2-oauth-demo`, open
`http://localhost:8082`.

1. Run Authorization Code + PKCE; show the decoded JWT (header/payload/signature).
2. Call the API with the valid token → 200, "validated locally, via JWKS".
3. Tamper with the token → 401 "signature check failed" — caught with no call to the AS.

Emphasise: the resource server fetched public keys from `/jwks` and did the maths itself.

# Checkpoint Questions

1. If a JWT's payload is readable by anyone, why is it safe to send? What makes it trusted?
2. What does local validation give up compared to introspection?
3. Why must the server pin the accepted algorithm instead of trusting the token's `alg`?

# Speaker Notes

Analogy: a JWT is a signed letter. Anyone can read it through the envelope window; what makes
it trustworthy is the signature they cannot forge — verified against a public key you already
trust (the JWKS). Reading the letter (decoding) and confirming who wrote it (validating) are
different acts. Stress the ordering of checks: never inspect claims you have not yet verified.
