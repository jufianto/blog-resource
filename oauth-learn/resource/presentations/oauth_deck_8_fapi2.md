# Deck 8: FAPI 2.0 (Capstone) — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. The live demo is a terminal you switch to at the marked point.
> This deck also owns the workshop-closing slide (page 33). Dense detail = Lesson 8 paper (handout).

**Slides:** PDF pages 30–33 · **Teach in:** Day 2, ~60 min (incl. Exercise 4 + review) · demo = `cmd/fapi`

## The one thing
FAPI 2.0 is a **security profile — a required *combination* of controls for high-value APIs — not
a new protocol and not one toggle.** Its strength is that all the controls are mandatory together.

---

## Walk the slides

### Page 30 — Divider: "PART 08 · FAPI 2.0"
Nearly empty; set the capstone frame:
> "When the API guards money or health records, 'we did OAuth' isn't enough. Any single control
> has an edge case. FAPI's answer isn't a new protocol — it's 'you must use PAR *and* PKCE(S256)
> *and* a sender-constrained token *and* strong client auth, all at once.' Defense in depth, made
> non-optional."

### Page 31 — "Stricter OAuth, combined" (what FAPI is)
On screen: **FAPI 2.0 combines** (strict config, PAR + sender-constrained tokens, strong client
auth, strict validation, optional message signing) next to **Where it's used** (open banking,
healthcare, regulated APIs).
Say: "Notice you already know every ingredient — Decks 4–7. FAPI just says *use the strict option,
every time*." Land the slide line: "normal OAuth gives options; FAPI chooses the strict ones for
high-value APIs."

### Page 32 — "FAPI controls mapped to risk" (controls → attacks)
On screen: table — request tampering → PAR/JAR; bearer replay → DPoP/mTLS; weak client auth →
private-key/cert auth; wrong-audience token → strict audience validation; payment dispute →
message signing + audit.
Say: go row by row — "this is the payoff of the whole workshop: every risk you met this week, and
the FAPI control that removes it." This table *is* Exercise 4.

### Page 33 — Closing: "Ready to teach"
On screen (dark): the ✓ checklist — draw Auth Code + PKCE and Client Credentials from memory;
explain code vs access vs refresh vs ID token; list the JWT validation checks; connect attacks to
mitigations; explain FAPI as strict OAuth.
Say: this is the workshop wrap for **both days**. Read the checklist as "here's what you can now
do" and close on the through-line: "from 'don't share your password' all the way to a
bank-grade profile — same four actors, same two channels, tighter every step."

---

## Run the example (after page 32, folded into Exercise 4)
From `resource/labs/day2-oauth-demo`, run `go run ./cmd/fapi`, open `http://localhost:8082`:
1. Start the flow → show the decoded **`private_key_jwt` client assertion** (client authenticates
   with a signed JWT, no shared secret) and the **PAR push**.
2. Approve → show **`token_type: DPoP`** — client authenticated *and* proved possession of its key.
3. Call the API → **200**: "FAPI 2.0: JWT verified locally AND DPoP proof-of-possession confirmed."

Drive home: the server uses **`WithProfile(ProfileFAPI2)`**, which **refuses to start** if any
required control (PAR, PKCE, DPoP, strong client auth, short-lived codes) is missing — **fail-closed**,
so you can't accidentally ship a weakened FAPI server. → **Exercise 4**: hand the room an
architecture sketch, map each FAPI control to the threat it removes (this mirrors page 32).

---

## Say it like this
> "FAPI is a combination lock, not a single key. PAR secures the request, `private_key_jwt` proves
> the client, DPoP binds the token — and FAPI says you don't get to skip any of them."

> "What FAPI is *not*: a login system, a product, or one config flag."

## Check they got it
1. Name the FAPI 2.0 baseline controls and the attack each one addresses.
2. Why does FAPI forbid public clients and shared secrets in favour of `private_key_jwt`?
3. Why is a plain bearer token never acceptable for a high-value API?

## They can now
Explain FAPI 2.0 as a required stack of controls, name which control stops which attack, and
recognize that the profile enforces itself (fail-closed).
