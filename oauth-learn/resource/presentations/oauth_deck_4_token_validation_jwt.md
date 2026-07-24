# Deck 4: Token Validation and JWT Basics — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. The live demo is a terminal you switch to at the marked point.
> Dense detail = Lesson 4 paper (handout).

**Slides:** PDF pages 18–20 · **Teach in:** Day 2, ~75 min (incl. Exercise 1) · demo = `cmd/jwt-validation`

## The one thing
A JWT is just a **data format**. Security comes from **validation**, not from the token being a
JWT. **Decoding is not validating.**

---

## Walk the slides

### Page 18 — Divider: "PART 04 · Token Validation & JWT Basics"
Day 2 opens here. Nearly empty; set the trap:
> "Everyone can read a JWT — paste it into jwt.io, there are the claims. So people assume: if I
> can read `scope: admin`, it must be true. Wrong. Anyone can *write* a string that says admin.
> The question is never 'what does it say?' — it's 'who signed it, and can I verify that myself?'"

### Page 19 — "JWT: header . payload . signature" (the format)
On screen: a decoded payload (`iss`, `sub`, `aud`, `exp`, `scope`) next to a card of the claims
that drive validation.
Say: "Three parts, dot-separated. The payload is just **base64url** — readable, **not** secret.
So read the claims *as the inputs to checks*: `iss` — do I trust who issued this? `aud` — is it
meant for *this* API? `exp` — is it still valid? `scope` — enough permission?" Say the slide's
own line out loud: "anyone holding it can read the payload."

### Page 20 — "What the resource server must check" (decoding ≠ validation)
On screen: **Validation checklist** (signature, alg/kid, iss/aud, exp/nbf, scope) next to
**Reject for the right reason** (401 vs 403; opaque token → introspect).
Say: "This slide is the whole lesson. Order matters — **signature first**, never inspect a claim
you haven't verified. Then trust the issuer, check the audience, check time, check scope."
Then the status codes: "**401** = I can't trust this token (missing/expired/bad signature).
**403** = I trust it, but it lacks the scope. Don't blur those." And the escape hatch: "if the
token is opaque, not a JWT, you can't verify locally — you call introspection."

---

## Run the example (after page 20, folded into Exercise 1)
From `resource/labs/day2-oauth-demo`, run `go run ./cmd/jwt-validation`, open
`http://localhost:8082`:
1. Run Authorization Code + PKCE → show the decoded JWT (header / payload / signature).
2. Call the API with the valid token → **200**, "validated locally, via JWKS".
3. **Tamper**: flip one byte of the payload → call again → **401 "signature check failed"** — and
   stress **no call was made to the authorization server**. The resource server caught it alone,
   with maths, using the public key it already had.

Then the `alg` trap: "never trust the token's *own* `alg` — pin RS256 on the server, or an
attacker sends `alg: none`." → hand off to **Exercise 1** (list what to check, and why).

---

## Say it like this
> "A JWT is a signed letter. Anyone can read it through the envelope window. What makes it
> trustworthy is the signature you can't forge — checked against a public key you already trust
> (the JWKS). Reading the letter and confirming who wrote it are different acts."

## Check they got it
1. If a JWT's payload is readable by anyone, why is it safe to send? What makes it *trusted*?
2. What does local validation give up compared to introspection?
3. Why must the server pin the accepted algorithm instead of trusting the token's `alg`?

## They can now
Explain why a readable token can still be trustworthy, list the checks in order (signature first),
and say when to validate locally (JWT + JWKS) vs call introspection (opaque token).
