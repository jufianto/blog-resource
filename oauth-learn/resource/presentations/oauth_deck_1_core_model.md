# Deck 1: OAuth Core Model — Presenter Script

> **How to use this:** open `export/presentations/oauth_complete_workshop_lessons_1_8.pdf`
> and present. This script walks the **actual slides in that PDF, in order**. Each heading
> is the slide's on-screen title, so you always know where you are. The live demo is not a
> slide — it's a terminal you switch to at the marked point. Dense detail = Lesson 1 paper (handout).

**Slides:** PDF pages 3–8 · **Teach in:** Day 1, ~60 min (incl. Lab 1)

## The one thing
OAuth lets an app use an API on your behalf **with limited permission and without ever
getting your password** — delegated authorization, not "JWT login".

---

## Walk the slides

### Page 3 — Divider: "PART 01 · OAuth Problem, Actors & Endpoints"
Big "01" on dark. This is your hook moment — the slide is nearly empty, so *you* talk:
> "RepoBoard wants to show your GitHub repos. The lazy way: ask for your GitHub password.
> Now a repo dashboard can delete everything, read everything, forever — until you change
> that password. That's the wrong shape. OAuth exists to hand out *permission*, not *passwords*."

### Page 4 — "A framework, not a login or a token format"
On screen: two cards — **OAuth 2.0 IS** (teal) vs **is NOT** (red).
Say: "Before anything else, kill three myths. OAuth is a *framework* for authorization. It is
**not** password sharing, **not** a single token format, and **not** a login protocol —
login is OpenID Connect, which we'll bound at the end." Read the red card aloud; that's where
people are wrong.

### Page 5 — "Password sharing is the wrong shape"
On screen: **Give away the password** (red) vs **Give a scoped token** (teal).
Say: walk the red column as pain (total power, can't scope, revoke = change password, breach
exposes everything), then the teal column as relief. Land the bottom line on the slide:
> "The client receives **permission**, not the user's password."

### Page 6 — "The four actors"
On screen: four cards — Resource Owner, Client, Authorization Server, Resource Server.
Say: point at each. "Resource Owner = you. Client = the app that wants in. Authorization
Server = issues tokens. Resource Server = the API that checks them." Use the slide's own
recall trick aloud: *"who owns the data, who wants access, who grants it, who serves the API?"*

### Page 7 — "Two endpoints, two channels"
On screen: a table (`/authorize` = front-channel/browser; `/token` = back-channel/server) plus
an access-token card.
Say: "Two endpoints, two very different roads. `/authorize` runs **in the browser** — exposed.
`/token` is **server-to-server** — protected. This front-channel vs back-channel split is the
idea behind *every* security control later in the workshop." Then the token card: "scoped,
expiring, revocable — checked on every call."

### Page 8 — "OAuth vs OpenID Connect"
On screen: comparison table (Purpose, Question, Main token, Example).
Say: draw the boundary in one line and repeat it:
> "Want to **call an API** → OAuth. Want to **know who the user is** → OpenID Connect."
"Login with Google" is OIDC. Don't let anyone leave thinking OAuth is a login system.

---

## Run the example (after page 8, before Lab 1)
Switch away from the slides to a terminal/editor. Put real config on screen and map each line
to an actor you just taught:

```text
CLIENT_ID=repo-board-web                                          # the Client
AUTHORIZATION_ENDPOINT=https://github.com/login/oauth/authorize   # Authorization Server
TOKEN_ENDPOINT=https://github.com/login/oauth/access_token        # Authorization Server
RESOURCE_API=https://api.github.com                               # Resource Server
REDIRECT_URI=https://repoboard.example.com/oauth/callback         # where the browser returns
SCOPES=repo:read                                                  # the permission requested
```

Say: "The token this produces is a permission card for `repo:read` at `api.github.com` — and
useless anywhere else. That 'useless anywhere else' is **audience**; that 'only repo:read' is
**scope**." → hand off to **Lab 1** (participants label the actors in GitHub / Calendar / CLI
/ backend scenarios).

---

## Say it like this
> "OAuth lets an app access an API with limited permission without receiving the user's
> password." (Repeat it — it's the spine of the whole workshop.)

Do **not** open with JWT. JWT is just a format that shows up later; leading with it is the
number-one way people misunderstand OAuth.

## Check they got it
1. Is OAuth primarily authentication or authorization?
2. Which actor *issues* the access token? Which actor *validates* it?
3. When does OpenID Connect enter the story?

## They can now
Point at any OAuth integration and name the four actors, the two endpoints, and what the token
is permission *to do* and *where*.
