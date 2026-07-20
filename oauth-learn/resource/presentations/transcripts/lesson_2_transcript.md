# Presenter Transcript — Deck 2: Authorization Code Flow + PKCE

**Day 1, Sessions 2–3 · Source: Lesson 2 · ~120 min (incl. Labs 2 & 3)**

Format: one block per slide. Bullets = points to make. `»` = stage directions. Anchor lines: *"The code is not the token — it's a short-lived ticket that must be exchanged correctly."* and for PKCE: *"The authorization request says: later I'll prove I know the original secret."*

---

## Slide 1 — The Main User-Based OAuth Flow · ~2 min
- This is *the* flow for user-delegated access — web, SPA, mobile. Everything else on Day 1 is a variation or a special case.
- Goal for this session: you can draw it from memory and defend every parameter.

## Slide 2 — Browser vs Back-Channel · ~8 min
- The whole flow is one idea: **the code travels through the browser; the token comes back-channel.**
- Front-channel (browser): the user goes to `/authorize`, approves, gets redirected back with a **code**.
- Back-channel (server-to-server): the app posts that code to `/token` and gets the **token**. The browser never sees the token.
- Why split it: the browser is visible and tamperable; the token is too sensitive to expose there.
» Walk `auth_code_sequence.png` left to right. Say each hop out loud: "browser to /authorize… redirect back with code… backend to /token… token returned."

## Slide 3 — Authorization Request · ~6 min
- The trip to `/authorize` carries: `response_type=code`, `client_id`, `redirect_uri`, `scope`, `state` (and later `code_challenge`).
- It's a *browser redirect*, not an API call — the user lands on the Auth Server's login/consent.
- Output of this step is only a **code** in the redirect back — not a token.
» Show the URL being built (slide 8 has the snippet); preview it here.

## Slide 4 — `state` · ~6 min
- `state` = a random value the client generates, stashes in its session, and checks on the callback.
- Purpose: **CSRF protection** — proves the callback belongs to a request *this* client actually started.
- If the returned `state` doesn't match what you stored → reject; someone's forging a callback.
» "state is about the *request's integrity*, not about tokens — keep it separate from PKCE (slide 15)."

## Slide 5 — `redirect_uri` · ~5 min
- Must match what was **registered** with the Auth Server *and* what was used in the request.
- It's where the browser comes back with the code — so if an attacker could change it, they'd redirect your code to themselves.
- Exact-match, pre-registered. No wildcards, no "close enough."

## Slide 6 — Token Request · ~6 min
- `/token` is a **direct POST** from the backend; it returns data, it does **not** redirect the browser.
- Sends: `grant_type=authorization_code`, the `code`, `client_id`, `redirect_uri` (and later `code_verifier`).
- Returns: the access token (and often a refresh token). This is the back-channel half.

## Slide 7 — Why Send `redirect_uri` to `/token`? · ~5 min
- It **binds the code exchange to the original redirect target** — the Auth Server checks it matches the one used at `/authorize`.
- Stops a stolen code from being redeemed against a different redirect URI.
- Point: the same parameter appears twice on purpose; the second time it's a consistency check.

## Slide 8 — Build the Authorization URL (Demo) · ~7 min
» Show the route:
```js
app.get("/connect/github", (req, res) => {
  const state = createRandomState();
  saveStateInSession(req, state);
  const url = new URL("https://github.com/login/oauth/authorize");
  url.searchParams.set("response_type", "code");
  url.searchParams.set("client_id", "repo-board-web");
  url.searchParams.set("redirect_uri", "https://repoboard.example.com/oauth/callback");
  url.searchParams.set("scope", "repo:read");
  url.searchParams.set("state", state);
  res.redirect(url.toString());
});
```
- Narrate: generate `state`, save it, build the URL, redirect. "This is all slide 3 was describing."
» Or run `day1-go-oauth-demo` and click "Start Auth Code + PKCE" — show the real redirect in the address bar.

## Slide 9 — Checkpoint: Map the Parameters (Lab 2) · ~8 min
» Run Lab 2. Participants place each parameter on the right endpoint: which go to `/authorize`, which to `/token`.
- Confirm they can answer: which request is front-channel, which is back-channel.

## Slide 10 — What If the Code Is Stolen? · ~4 min
- Set up the attack: the code rides through the browser (redirects, logs, referrer headers) — assume it *can* leak.
- A stolen code alone should **not** be enough to get a token. Everything after this slide is why.
» "This is the gap PKCE closes. Watch."

## Slide 11 — PKCE: Challenge First, Verifier Later · ~8 min
- At `/authorize` the client sends a **`code_challenge`** (a hash). The Auth Server stores it with the code.
- At `/token` the client sends the **`code_verifier`** (the original secret). The server hashes it and compares.
- So the authorization request effectively says: *"later I'll prove I know the original secret."*
- A thief who grabs the code from the browser never saw the verifier → can't complete the exchange.
» Walk `auth_code_pkce_sequence.png`; point to where challenge is stored vs. where verifier is checked.

## Slide 12 — SHA-256 Is One-Way · ~5 min
- The challenge is `SHA-256(verifier)`. You can go verifier → challenge, never back.
- So even seeing the challenge (which *does* travel through the browser) tells an attacker nothing useful.
- The server just re-hashes the presented verifier and checks it equals the stored challenge.

## Slide 13 — Create PKCE Values (Demo) · ~6 min
» Show:
```js
const codeVerifier = base64url(crypto.randomBytes(32));
const codeChallenge = base64url(
  crypto.createHash("sha256").update(codeVerifier).digest()
);
```
- Narrate: random verifier → hash → challenge. Challenge goes to `/authorize`, verifier is kept secret for `/token`.
» In `day1-go-oauth-demo`, show the session panel: `code_verifier` is stored client-side and never in a redirect URL.

## Slide 14 — `code_verifier` Goes to `/token` · ~4 min
- The verifier is sent **only** on the back-channel token request — never in the browser authorization request.
- That's the whole trick: the secret proving "same client" never travels where it could be stolen.

## Slide 15 — `state` vs PKCE · ~6 min
- They solve *different* problems; you need both.
- **`state`** — CSRF: is this callback tied to a request I started?
- **PKCE** — code interception: is the party redeeming the code the same one that requested it?
- Two-column table on screen; the takeaway: not interchangeable, both mandatory for public clients.

## Slide 16 — Code Interception Scenario (Lab 3) · ~8 min
» Run Lab 3. Participants explain, step by step, why an attacker who steals only the code fails.
- Expected answer: no `code_verifier`, so the `/token` exchange fails the challenge check.
» Optionally demo: in `day1-go-oauth-demo`, an `/authorize` request without a challenge is now rejected (PKCE required).

## Slide 17 — Draw It From Memory · ~6 min
- Have them redraw Auth Code + PKCE on a blank sheet: browser → `/authorize` (challenge) → code back → backend → `/token` (verifier) → token.
- If they can draw it and name what each parameter defends, this session succeeded.
» Transition: "Not every app is a user in a browser, though. Next: the other flows." → Deck 3.
