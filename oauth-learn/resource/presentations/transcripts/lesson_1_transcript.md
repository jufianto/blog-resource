# Presenter Transcript — Deck 1: OAuth Core Model

**Day 1, Sessions 1–2 · Source: Lesson 1 · ~60 min (incl. Lab 1)**

Format: one block per slide. Bullets = the points you make (expand in your own words). `»` = stage directions (click, demo, ask the room). Anchor line to repeat all session: *"OAuth lets an app use an API with limited permission — without ever getting the user's password."*

---

## Slide 1 — OAuth Starts With a Problem · ~2 min
- One problem to hold in your head: your app needs to act in *another* system on a user's behalf — read their GitHub, add to their calendar.
- Naive answer: take their password and log in as them. OAuth is the clean answer to that.
- For the next hour, OAuth is *only* this: limited access to someone's account elsewhere, without their password.
» "We are not talking about JWT yet — it's just a token format that shows up later. Park it."

## Slide 2 — Password Sharing Is the Wrong Shape · ~4 min
- Not just "insecure" — the *wrong shape*. A password gives all-or-nothing access.
- Three breakages: (1) too much power — everything, including delete; (2) can't scope — no "read but not push"; (3) can't revoke without changing the password and breaking every other app.
- Image to plant: **master key vs. guest badge**. OAuth is the guest-badge system.
» Ask: "Who's revoked a third-party app in GitHub/Google? That per-app cut-off screen exists *because* of OAuth."

## Slide 3 — RepoBoard Connects GitHub · ~4 min
- Running example all day: **RepoBoard** (a dashboard) wants to read **Alya's** GitHub repos.
- Requirements: read-only, revocable anytime, never sees her password.
- GitHub on purpose — everyone can picture it; banking examples wait for Day 2.
» On screen: Alya → RepoBoard → GitHub API. Leave it up while you name the actors.

## Slide 4 — The Four Actors · ~7 min
- **Resource Owner** — the human who owns the data (Alya).
- **Client** — the *app* wanting access (RepoBoard). Stress: client = the app, not the user.
- **Authorization Server** — authenticates the user, issues tokens; the *only* one that sees the password.
- **Resource Server** — the API holding the data (GitHub's API); receives and honors tokens.
- Why split them: keep the password with the one actor that should have it; everyone else gets a scoped, revocable token.
» Point at `oauth_actor_general.png`; say each as a verb: owns / asks / issues / serves. Ask: "In 'Login with Google to Figma,' who's the client?" (Figma.)

## Slide 5 — Two Endpoints: `/authorize` and `/token` · ~6 min
- Confusing these is the #1 source of OAuth bugs.
- **`/authorize`** — on the Auth Server, reached *through the browser*; job = log in + get consent. Human-facing, front-channel.
- **`/token`** — on the Auth Server, reached *server-to-server*; job = hand out tokens. Machine-facing, browser never touches it.
- Contrast line: `/authorize` is where the *user* approves; `/token` is where the *app* collects.

## Slide 6 — What the Access Token Means · ~5 min
- It is the guest badge made real: *not* the password, *not* proof of identity.
- It means one thing: "bearer may do these actions, on this API, for a limited time." Permission, scoped, expiring.
- Two consequences: API never relearns the password; a leaked token is *time-boxed*, not permanent.
- Pocket phrase: **"permission, not identity"** — we use it to split OAuth from OIDC later.

## Slide 7 — Scope Is Permission · ~5 min
- Scope is the dial a password never had. RepoBoard asks for `repo:read`, not "everything."
- The Auth Server shows Alya exactly those scopes; she consents; the token carries only those.
- Try to push code later → API refuses; the token lacks that permission.
- Scope does double duty: how the app *requests* least privilege, and how the user *sees & consents* to it.
» "First question in any OAuth review: what scopes does this actually need?"

## Slide 8 — Audience Is the API · ~4 min
- Scope = *what* the token may do. Audience = *where* it may be used.
- A token minted for GitHub's API must not be accepted by your payments API, even if the signature checks out.
- Without an audience check, a token leaked from one service is a skeleton key for another.
- One-liner: **scope is *what*, audience is *where*.** (We make it concrete Day 2.)

## Slide 9 — User Approval (Consent) · ~3 min
- The consent screen = least privilege as a human decision: "RepoBoard wants to read your repositories → Approve."
- Nothing issued without it; the user approves *specific scopes*, not a blanket "allow app."
- Served by the **Auth Server** (the one trusted with the password), never by the client.
» "If a client renders its own 'enter your GitHub password' screen, that's phishing, not OAuth."

## Slide 10 — OAuth vs OpenID Connect · ~6 min
- The field's most confused boundary — you'll be the one who gets it right.
- **OAuth = authorization**: "may this app do this thing?" It was never meant to say *who* the user is.
- **OIDC = authentication**, layered on OAuth; it adds the **ID token** (identity). Access token stays about permission.
- Clean split: **OAuth = access, OIDC = login.** "Login with Google" = OIDC; "let this app read my Drive" = OAuth.
- Same machinery, different questions — which is why people fuse them.
» Slow down here. Ask: "Which gives an ID token?" (OIDC) "Which is about permission?" (OAuth)

## Slide 11 — Real Configuration Values (Demo) · ~5 min
» Put config on screen:
```text
CLIENT_ID=repo-board-web
AUTHORIZATION_ENDPOINT=https://github.com/login/oauth/authorize
TOKEN_ENDPOINT=https://github.com/login/oauth/access_token
RESOURCE_API=https://api.github.com
REDIRECT_URI=https://repoboard.example.com/oauth/callback
SCOPES=repo:read
```
- Map each line to an actor/endpoint: `CLIENT_ID` = the Client; the two endpoints = the Auth Server (the browser-one and the backend-one from slide 5); `RESOURCE_API` = Resource Server; `SCOPES` = least privilege in one line; `REDIRECT_URI` = where the browser returns (becomes a security control tomorrow).
- Point: OAuth isn't abstract — every actor/endpoint is a literal config value.

## Slide 12 — Checkpoint: Actor Mapping (Lab 1) · ~6 min
» Run Lab 1. Four scenarios; name all four actors for each:
- (a) RepoBoard ↔ GitHub repos; (b) app ↔ Google Calendar; (c) CLI ↔ backend deploy; (d) backend service ↔ internal API.
- On (d): let them notice there's *no human* — teaser for Client Credentials in Session 4.
» Don't rush — the whole day rests on naming actors fluently.

## Slide 13 — The Mental Model · ~3 min
- Four questions collapse everything: **Who asks?** (client) **Who approves?** (resource owner) **Who issues?** (auth server) **Who validates?** (resource server).
- Answer those four for any system and you understand its OAuth model before reading its code.
- Leave-with line: *OAuth lets an app use an API with limited permission, without ever getting the user's password.*
» Transition: "Let's stop talking about it and run it." → Deck 2.
