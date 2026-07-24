# Deck 6: OAuth Security Controls — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. This session is threat-modeling, not a new demo — keep the room asking
> "what leaks, and what stops it?" Dense detail = Lesson 6 paper (handout).

**Slides:** PDF pages 24–26 · **Teach in:** Day 2, ~75 min (incl. threat-model discussion)

## The one thing
Every control exists to stop a **specific attack.** "We use OAuth" is not a security statement —
say the *threat* before you name the *control*.

---

## Walk the slides

### Page 24 — Divider: "PART 06 · OAuth Security Controls"
Nearly empty; reframe the room:
> "People treat OAuth controls as a checklist of acronyms and can't say what any of them actually
> *prevent*. That's how you tick every box and still get owned. Today we go the other way: start
> from the attack, then reveal the control that kills it."

### Page 25 — "Avoid Implicit and Password grant" (legacy flows)
On screen: two red cards — **Implicit flow** and **Password credentials** — both bad.
Say: "Two flows you'll still see in old code, both retired. **Implicit** returned tokens directly
in the browser redirect — front-channel token exposure in history, logs, scripts. **Password
grant** had the client collect the user's real password — it breaks the entire point of OAuth and
blocks MFA." Land the slide line: "Use **Auth Code + PKCE** for users, **Client Credentials** for
services."

### Page 26 — "Attacks mapped to mitigations" (threat model)
On screen: the anchor table — attack/mistake on the left, mitigation on the right (missing state,
weak redirect matching, code interception, long-lived tokens, refresh theft, bearer replay).
Say: this is the slide participants photograph. Go **row by row as attack → defense**:
> "Missing `state` → CSRF; generate and verify it. Loose redirect matching → open redirect;
> exact pre-registered URIs. Code interception → PKCE. Long-lived token → short lifetimes.
> Refresh theft → rotation + reuse detection. Bearer replay → sender-constrained tokens — which is
> the whole next session."
Note the last row is the bridge to Deck 7.

---

## Run the example (no new demo — point at what they've already seen)
Don't spin up anything new; make the controls concrete by pointing back:
- **Day 1 Go demo** — PKCE required, `state` checked on the callback.
- **Day 2 `cmd/jwt-validation`** — local token validation, scope enforced, tamper → 401.

For each step of a real flow, ask the room: **"What leaks here, and what stops it?"** Keep them in
threat-modeling mode rather than lecturing. The last table row (bearer replay) sets up Deck 7.

---

## Say it like this
> "A bearer token is a password-equivalent for as long as it's valid — which is exactly why the
> next session's sender-constraining matters."

Drive every point from the **attack first**, then reveal the control. Never lead with the acronym.

## Check they got it
1. An attacker intercepts the authorization code. Which control makes it useless, and why?
2. What's the difference between restricting a token by **scope** and by **audience**?
3. Which control would you add *first* to a legacy OAuth integration?

## They can now
Threat-model a flow: name what travels through the browser, what can leak at each step, and which
control closes each gap.
