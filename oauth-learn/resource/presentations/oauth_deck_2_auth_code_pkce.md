# Deck 2: Authorization Code Flow + PKCE — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. The live demo is a terminal you switch to at the marked point.
> Dense detail = Lesson 2 paper (handout).

**Slides:** PDF pages 9–13 · **Teach in:** Day 1, ~120 min (incl. Labs 2 & 3)

## The one thing
The code travels through the **browser (front channel)**; the token is fetched from
**`/token` (back channel)**. The code is a short-lived ticket, **not** the token — and PKCE
makes a stolen ticket worthless.

---

## Walk the slides

### Page 9 — Divider: "PART 02 · Authorization Code + PKCE"
Nearly empty; you set the stakes:
> "This is the flow behind almost every 'Log in with…' button. It looks like one redirect, but
> two different things happen: something comes back through the *browser*, something else is
> fetched *server to server*. Mix them up and you build an insecure app. Let's watch them travel
> on different roads."

### Page 10 — "Authorization Code, step by step"
On screen: 6 numbered steps (redirect to `/authorize` → user approves → code + state back →
client posts code to `/token` → gets access token → calls API).
Say: trace the steps with your finger. Hit the key beat: "Step 3 returns a **code** in the
browser. Step 4 exchanges it **server-to-server**. The code is not the token."

### Page 11 — "Authorization request parameters"
On screen: table of `/authorize` params (`response_type`, `client_id`, `redirect_uri`, `scope`,
`state`, `code_challenge`).
Say: "Everything on this slide rides in the browser URL. Two of these are safety checks, not
plumbing — `state` and `code_challenge`. Remember them; the next two slides are all about them."

### Page 12 — "A stolen code should not be enough" (why PKCE)
On screen: **Without PKCE** (red) vs **With PKCE** (teal), and the formula panel
`code_challenge = BASE64URL(SHA256(code_verifier))`.
Say: "Here's the attack — steal the code off the redirect, send it to `/token`, and on a public
client you might get tokens. PKCE closes it: the client sent a **hash** up front, and must later
present the **original** to redeem. SHA-256 is one-way, so seeing the hash tells the attacker
nothing." Land it: "Attacker has the code but not the verifier → fails."

### Page 13 — "state vs redirect_uri vs PKCE" (not interchangeable)
On screen: table — who checks each and what it protects.
Say: "These get lumped together and they shouldn't. `state` is checked by the **client**
(anti-CSRF). `redirect_uri` is checked by the **authorization server** (binds the code to the
original target). `code_verifier` is checked at the **token endpoint** (proves you started this).
A modern browser flow uses **all three**."

---

## Run the example (after page 13, folded into Labs 2 & 3)
Switch to a terminal and run the **Day 1 Go demo**. Show the real Authorization Code + PKCE
round-trip: the code arriving in the browser, then the token coming back from `/token`. Then
break it — tamper the `state` on the callback and show it rejected. Point at the split:
- `state` → stored by client, checked on callback.
- `code_challenge` → goes to `/authorize` (in the browser).
- `code_verifier` → goes to `/token` (never touched the browser).
→ hand off to **Lab 2** (map endpoints & params) and **Lab 3** (explain why the interceptor fails).

---

## Say it like this
> "The code is not the token. The code is a short-lived ticket that must be exchanged correctly."

> "The authorization request says: *later I'll prove I knew the original verifier.*"

## Check they got it
1. Which request is front-channel (browser)? Which is back-channel?
2. Why is `redirect_uri` included in the **token** request?
3. If an attacker steals only the code, why does PKCE stop them getting a token?

## They can now
Draw Authorization Code + PKCE from memory and explain why `state`, `redirect_uri`,
`code_challenge`, and `code_verifier` each exist and are **not** interchangeable.
