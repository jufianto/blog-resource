# Deck 5: JWS, JWE, JWK and JWKS — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. The live demo is a terminal you switch to at the marked point.
> Dense detail = Lesson 5 paper (handout).

**Slides:** PDF pages 21–23 · **Teach in:** Day 2, ~75 min · demo = `cmd/jwt-validation`

## The one thing
JOSE is a **toolbox**, not one thing: **JWS signs, JWE encrypts, JWK/JWKS publish keys.** And
**signed does not mean private.**

---

## Walk the slides

### Page 21 — Divider: "PART 05 · JWS, JWE, JWK / JWKS"
Nearly empty; open with the myth you're about to kill:
> "Yesterday someone said 'the JWT is encrypted, so the secret's safe inside.' It is not. A
> normal signed JWT is *readable by anyone* — signing proves who wrote it, it does **not** hide
> it. Put a secret in a signed-only token and you've published the secret."

### Page 22 — "Five pieces that fit together" (the JOSE family)
On screen: table — JWT (container), JWS (signature), JWE (encryption), JWK (one key), JWKS
(published set).
Say: "One family, different jobs. JWT is the envelope. JWS is the wax seal — proves it wasn't
modified, still readable. JWE is the locked box — hides the contents. JWK is one key as JSON;
JWKS is the public noticeboard of keys." Land the slide's line: "**JWS proves the content wasn't
modified. JWE hides the content.**"

### Page 23 — "Where the verifier gets the key" (keys & rotation)
On screen: **How verification works** (AS signs with private key → publishes public key in JWKS →
RS fetches from trusted `jwks_uri` → selects by `kid` → verifies) next to **Common alg/kid
mistakes** (red).
Say: walk the four-step flow left-to-right. "The `kid` in the token header points at which key in
the JWKS to use. When the RS sees an **unknown `kid`**, it refetches the JWKS — that refetch is
the hook that makes **rotation** work: publish the new key first, sign with it later, retire the
old one after old tokens expire." Then read the red card as the don't-do list: "never accept
`alg: none`, never trust the token's algorithm, never fetch keys from a URL the token controls,
never drop an old key too early." Land it: "**allowed algorithms and key sources come from
configuration — never from the token.**"

---

## Run the example (after page 23)
Reuse `cmd/jwt-validation`. In the browser open `http://localhost:8080/jwks` and show the
published key set the resource server fetches. Point at:
- the `kid` in the **token header** and the matching `kid` in the **JWKS** — that's key selection;
- the refetch-on-unknown-`kid` behavior — that's what makes zero-downtime rotation possible.

(JWE has no runnable demo here — teach it from page 22. Point is just: same family, different job.)

---

## Say it like this
> "JWS is a wax seal — you can read the letter, but you can't forge the seal. JWE is a locked box
> — you can't read it at all. JWKS is the public noticeboard of seals. `kid` is the label saying
> which seal to compare."

Hammer the one sentence people get wrong: **a signed token still leaks its claims.**

## Check they got it
1. You must hide the token contents from the client. JWS or JWE?
2. How does a resource server pick the correct key to verify a token?
3. During rotation, why do old and new keys both appear in the JWKS for a while?

## They can now
Choose JWS vs JWE for a requirement, trace how a `kid` selects a JWKS key, and explain rotation
without downtime.
