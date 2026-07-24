# Deck 3: Choosing the OAuth Flow — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This walks the **actual slides in that PDF, in order** — each heading is the
> slide's on-screen title. The live demo is a terminal you switch to at the marked point.
> Dense detail = Lesson 3 paper (handout).

**Slides:** PDF pages 14–17 · **Teach in:** Day 1, ~75 min (incl. Lab 4)

## The one thing
**Not every OAuth flow is login.** The flow you pick is decided by a few questions — mainly "is
there a human approving this?" — and only the token request really changes.

---

## Walk the slides

### Page 14 — Divider: "PART 03 · Client Credentials, Refresh & Device"
Nearly empty; frame the problem:
> "People learn Authorization Code and try to use it for *everything* — a nightly batch job, a
> CLI, a smart TV. It doesn't fit: no browser, no user to click approve. OAuth has different
> flows for different shapes of app. Let's pick the right tool."

### Page 15 — "Three more flows, three problems"
On screen: three cards — **Client Credentials** (machine-to-machine, no user), **Refresh Token**
(renew without re-login), **Device Authorization** (CLI/TV, approve on a second device).
Say: give each a one-liner and a real example. "Client Credentials = your fraud-check service
calling another API at 3am. Refresh = keep a session alive without bouncing the user through
login. Device = logging into Netflix on a TV by typing a code on your phone."

### Page 16 — "Rotation makes reuse detectable" (refresh tokens)
On screen: table — Issue → Refresh → Rotate → Reuse seen.
Say: "Refresh tokens are more sensitive than access tokens — they *mint* new ones. Rotation is
the safety net: each use hands back a new refresh token and burns the old. If an old one shows
up again, that's theft — revoke the whole family." Land the slide line: "an access token opens
the door once; a refresh token keeps making new keys — store it well."

### Page 17 — "Choose the flow" (cheat sheet)
On screen: the big decision table (User? Redirect? Main use? Today?).
Say: this is the takeaway slide — read it as rules, not rows:
> "User approving? → **Auth Code + PKCE**. Service on its own behalf? → **Client Credentials**.
> Renewing access? → **Refresh Token**. Limited input? → **Device**. Implicit / Password? →
> **Avoid** — that's Lesson 6."

---

## Run the example (folded into Lab 4)
Switch to a terminal. Show that only the **token request** changes between flows — same endpoint,
different `grant_type`:
- **Client Credentials:** `grant_type=client_credentials` + Basic auth, no browser, no user.
- **Refresh:** `grant_type=refresh_token` + the refresh token.
- **Device:** `grant_type=urn:ietf:params:oauth:grant-type:device_code` while the client polls.

The **Day 1 Go demo** can show Client Credentials issuing a token with no redirect at all.
→ hand off to **Lab 4** (participants pick the correct flow for given scenarios).

---

## Say it like this
> "User approving access → Auth Code + PKCE. Service acting as itself → Client Credentials.
> Renewing access → Refresh Token. Limited input → Device Authorization."

And the boundary again: need to *log the user in*? That's OpenID Connect, not a new grant.

## Check they got it
1. Which flow has no user?
2. Which flow renews access after the initial authorization?
3. Which flow fits a CLI or TV app?
4. Which flow should a SPA/mobile app use for user-delegated access?
5. Why is Client Credentials not "service login as a user"?

## They can now
Given an app type, pick the correct flow and say what its token request looks like.
