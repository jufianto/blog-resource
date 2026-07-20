# Presenter Transcript — Deck 3: OAuth Flow Comparison

**Day 1, Session 4 · Source: Lesson 3 · ~75 min (incl. Lab 4)**

Format: one block per slide. Bullets = points to make. `»` = stage directions. Decision rule to repeat: *"User approving? → Auth Code + PKCE. Service on its own behalf? → Client Credentials. Renewing access? → Refresh Token. Limited input? → Device."*

---

## Slide 1 — Not Every OAuth Flow Is Login · ~2 min
- Yesterday's flow assumed a user in a browser. That's not every case.
- Different *application shapes* need different flows; this session is how to pick.

## Slide 2 — First Question: Is There a User? · ~5 min
- The one question that forks everything: **is a human approving access, or is a service acting on its own?**
- User present → user-delegated flows (Auth Code + PKCE). No user → service flow (Client Credentials).
- Everything else is a refinement of these two.
» Draw the decision tree; keep it on screen as the spine for the session.

## Slide 3 — Authorization Code + PKCE (recap) · ~4 min
- The user-delegated case, from Deck 2: web, SPA, mobile.
- One-line recap: browser gets a code, backend exchanges it (with PKCE) for a token.
- Anchor: this is the default whenever a *user* is delegating access.

## Slide 4 — Client Credentials · ~7 min
- **No user, no browser, no consent.** A service authenticates as *itself* to get a token.
- Use when your backend calls another API on its *own* behalf (e.g., a worker calling an internal fraud API).
- The client authenticates with its own credential (secret or key) directly at `/token`.
» Walk `client_credentials_sequence.png` — note there's no `/authorize`, no redirect.

## Slide 5 — Client Credentials Token Request (Demo) · ~5 min
» Show:
```http
POST https://auth.example.com/token
Authorization: Basic base64(client_id:client_secret)
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&scope=fraud.check
```
- Narrate: no user step; the service authenticates and gets a token immediately.
- Point out: **no refresh token, no ID token** — there's no user session to refresh or identify.
» In `day1-go-oauth-demo`, run the Client Credentials demo — instant token, no browser.

## Slide 6 — Refresh Token · ~6 min
- Not a login method — a **token lifecycle** mechanism. Access tokens are short-lived on purpose.
- When the access token expires, the client exchanges a refresh token for a new one — no user re-prompt.
- The refresh token is **more sensitive** than the access token (longer-lived, renews access).
» Walk `refresh_token_sequence.png`.

## Slide 7 — Refresh Token Request (Demo) · ~5 min
» Show:
```http
POST https://auth.example.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&refresh_token=old_refresh_token&client_id=web-client
```
- Narrate: swap old refresh token for a fresh access token (and often a new refresh token — rotation).
» In `day1-go-oauth-demo`, use the refresh demo; show the old token being invalidated on rotation.

## Slide 8 — Device Authorization · ~6 min
- For devices with **limited input**: CLI, TV, IoT — no good browser/keyboard.
- The device shows a short **user code**; the user approves on their *phone/laptop* browser; the device **polls** `/token` until approved.
- Splits the flow across two devices on purpose.
» Walk `device_authorization_sequence.png`.

## Slide 9 — Device Code Start & Polling (Demo) · ~6 min
» Show the polling request:
```http
POST https://auth.example.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=urn:ietf:params:oauth:grant-type:device_code&
device_code=device-code-from-first-response&client_id=cli-tool
```
- Narrate: before approval you get `authorization_pending`; after approval, a token.
» In `day1-go-oauth-demo`, run Device: show the user code, approve in the browser, poll again → token.

## Slide 10 — Login with Google (OIDC boundary) · ~4 min
- Reminder from Deck 1: if the app needs **identity/login**, that's **OpenID Connect**, not plain OAuth.
- These flows deliver *access*; OIDC adds the ID token on top.
- Keep the split clean: today's flows = access; login = OIDC.

## Slide 11 — Flow Comparison Table · ~6 min
- Put them side by side across: user present? browser redirect? client secret? token purpose? refresh?
- Let the table make the pattern obvious — the differences are exactly the "is there a user / is there a browser / is it renewing" axes.
» This table is the artifact they photograph. Linger on it.

## Slide 12 — Checkpoint: Choose the Correct Flow (Lab 4) · ~8 min
» Run Lab 4. Give scenarios (SPA, backend cron job, smart TV app, mobile app, service-to-service) → they pick the flow and justify.
- Watch for the classic error: calling Client Credentials "the service logging in as the user." It isn't — no user is involved.

## Slide 13 — Flow Selection Rules · ~4 min
- Five questions decide it: user? browser? public client? renewing? limited input?
- Recite the decision rule one more time so it sticks.
- Close Day 1: "You can now map actors, draw Auth Code + PKCE, and choose a flow. Tomorrow: what makes the *tokens* themselves trustworthy."
» Transition to Day 2 / Deck 4.
