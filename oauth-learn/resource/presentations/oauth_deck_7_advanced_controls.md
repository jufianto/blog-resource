# Deck 7: Advanced Controls — PAR, JAR, JARM, DPoP, mTLS — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. The live demos are terminals you switch to at the marked points.
> Dense detail = Lesson 7 paper (handout).

**Slides:** PDF pages 27–29 · **Teach in:** Day 2, ~90 min (incl. Exercises 2 & 3) · demos = `cmd/par-jar`, `cmd/dpop`

## The one thing
Three different weak points, three families of control: harden the **request** (PAR/JAR), the
**response** (JARM), and the **token** (DPoP/mTLS). A stolen token that's sender-constrained is
useless to the thief.

---

## Walk the slides

### Page 27 — Divider: "PART 07 · Advanced: PAR, JAR, JARM, DPoP, mTLS"
Nearly empty; name the two gaps that remain after the basic flow:
> "The basic flow protects the *code*. Two things are still exposed. First, the authorization
> request rides in a browser URL — visible, tamperable, loggable. Second, and worse: once a token
> is issued, whoever holds it can use it. Steal it, replay it, you're in. HTTPS fixes neither."

### Page 28 — "PAR, JAR, JARM" (request & response)
On screen: table — PAR (push request back-channel), JAR (sign the request), JARM (sign the
response), each with what it protects.
Say: "Harden the request and the response. **PAR** pushes the request details to a back-channel
endpoint first, so only an opaque handle rides in the browser. **JAR** signs the request object so
its parameters can't be altered. **JARM** signs the *response* too, which also stops mix-up
attacks." Land the slide line: "PAR pushes · JAR signs the request · JARM signs the response."

### Page 29 — "DPoP vs mTLS" (stop token replay)
On screen: two cards — **DPoP** (application-layer signed proof, key-bound token) and **mTLS**
(TLS-layer client cert, cert-bound token).
Say: "Both **sender-constrain** the token — bind it to something the client must prove it holds.
**DPoP** is an application-layer signed proof per request; fits public clients and APIs. **mTLS**
binds to a client TLS certificate; fits confidential, high-security clients." Land the slide line:
"a stolen bearer token alone is no longer enough."

---

## Run the examples (folded into Exercises 2 & 3)
From `resource/labs/day2-oauth-demo`:

**`go run ./cmd/par-jar`** — show the decoded **signed request object**, the PAR response, and
that the `/authorize` URL now carries only `client_id`, `response_type`, `scope`, and an opaque
`request_uri`. Then show a plain `/authorize` *without* PAR being **rejected (400)**.

**`go run ./cmd/dpop`** — show `token_type: DPoP`; a correct call *with* a proof → **200**; then
**replay the same token as a plain bearer → 401.** That replay is the moment the session turns on:
the token that worked one second ago is refused the instant it's presented without the key.

(JARM and mTLS are taught from page 28/29 — there's no runnable demo for them in this module. Say
so plainly; don't imply a demo you're not running.)

---

## Say it like this
> "Three targets: the request, the response, and the token. PAR/JAR seal the request, JARM seals
> the response, DPoP/mTLS bind the token to a holder."

> "DPoP vs mTLS — same goal, different plumbing: application-layer signed proof vs transport-layer
> client certificate."

## Check they got it
1. What can leak or be tampered with in a classic front-channel authorization URL? How do PAR and
   JAR each remove that risk?
2. How does a DPoP-bound token defeat replay, and what does the resource server check?
3. When would mTLS be a better sender-constraining choice than DPoP?

## They can now
Match a request/response/token threat to the control that removes it, and explain why a
DPoP-bound token can't be replayed.
