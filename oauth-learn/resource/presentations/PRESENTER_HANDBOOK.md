# Presenter Handbook — OAuth 2.0 / JWT / FAPI Workshop

> **Who this is for:** the person delivering the workshop. The `oauth_deck_*.md` scripts tell
> you *what to say* on each slide. This handbook is everything else you need to actually stand
> up and run the room without falling over — the demos, the timing, the questions, the words.
> Read it once end-to-end before you present, then keep it open beside the deck on the day.

**What sits where**
- **Slides (project this):** `export/presentations/oauth_complete_workshop_lessons_1_8.pdf` — 33 pages.
- **Per-slide scripts (your lines):** `oauth_deck_1..8_*.md`, keyed to the PDF by page number + title.
- **Handouts (give to participants):** the Lesson PDFs in `export/pdf/`.
- **Labs (participants run):** `resource/labs/day1-go-oauth-demo/` and `resource/labs/day2-oauth-demo/`.
- **This handbook:** how to run the demos, the clock, the Q&A, pronunciation, and lab facilitation.

---

## 1. Pre-flight checklist — do this BEFORE the room fills up

Nothing rattles a junior like a broken laptop at 9:00. Run this 30 minutes early.

- [ ] **Go 1.25 installed:** `go version` → prints `go1.25.x`.
- [ ] **Deck open** in Keynote/PowerPoint (or the PDF), full-screen tested on the projector.
- [ ] **Ports free:** `for p in 8080 8081 8082; do lsof -ti tcp:$p; done` → prints **nothing**.
      If it prints a PID, kill it (see §3 reset).
- [ ] **Warm up every demo once** so the first `go run` compile delay doesn't happen live:
      from each lab dir, run the command, open `http://localhost:8082`, click through once,
      then **Ctrl-C**. Do this for all five (Day-1 demo + the four Day-2 `cmd/`s).
- [ ] **Fallback recordings ready** (see §6) — screenshots/GIF of each demo's money moment,
      in case the live demo dies.
- [ ] **This handbook + the day's deck scripts** open in a second window / on a tablet.
- [ ] **Water, timer/clock visible, phone on silent.**

---

## 2. The clock (run-of-show)

Times are **durations** — shift the start to match your room. Teaching time comes from the deck
scripts; labs and breaks are added here. Day 2 is long — if you're tight, the trim targets are
marked ✂.

### Day 1 — Core OAuth (Lessons 1–3, PDF pages 1–17)
| Block | Mins | What |
|---|---:|---|
| Welcome + pre-flight with the room | 15 | Everyone runs the Day-1 setup, confirms `:8082` loads |
| **Deck 1** — pages 3–8 | 25 | OAuth core model |
| **Lab 1** (Auth Code + PKCE) | 20 | Day-1 demo, Exercise 1 |
| Break | 10 | |
| **Deck 2** — pages 9–13 | 50 | Auth Code + PKCE |
| **Labs 2 & 3** (params, interception) | 35 | Day-1 demo, Exercises 1–3 |
| Lunch | 45 | |
| **Deck 3** — pages 14–17 | 35 | Choosing the flow |
| **Lab 4** (Device / Client Creds) | 20 | Day-1 demo, Exercises 2 & 4 |
| Day-1 wrap + questions | 15 | |

### Day 2 — Tokens & hardening (Lessons 4–8, PDF pages 18–33)
| Block | Mins | What |
|---|---:|---|
| Recap Day 1 + pre-flight | 10 | Ports free, Day-2 module builds |
| **Deck 4** — pages 18–20 | 40 | Token validation & JWT |
| **Exercise 1** (`cmd/jwt-validation`) | 20 | tamper → 401 |
| Break | 10 | |
| **Deck 5** — pages 21–23 | 45 | JWS/JWE/JWK/JWKS (reuses Ex 1's `/jwks`) |
| **Deck 6** — pages 24–26 | 45 | Security controls (discussion, no new demo) |
| Lunch | 45 | |
| **Deck 7** — pages 27–29 | 50 | PAR/JAR/JARM, DPoP/mTLS |
| **Exercises 2 & 3** (`cmd/dpop`, `cmd/par-jar`) | 35 | replay → 401; no-PAR → 400 |
| **Deck 8** — pages 30–33 | 45 | FAPI 2.0 capstone + close |
| **Exercise 4** (`cmd/fapi`) + wrap | 25 | full stack → 200 |

✂ **If you're behind:** Deck 6 can run as pure discussion off the page-26 table (skip the
callbacks to earlier demos); Deck 5's JWE half is diagram-only anyway.

---

## 3. The one demo rule that saves the day

**Only one demo runs at a time.** Every demo — the Day-1 app and all four Day-2 `cmd/`s — binds
the **same three ports**:

| Port | Role | 
|---|---|
| **8080** | Authorization Server (issues tokens) |
| **8081** | Resource Server / API (validates tokens) |
| **8082** | Client app — **always your entry point: `http://localhost:8082`** |

So you **must fully stop one demo before starting the next**, or the new one dies with
*"address already in use."*

**Reset between demos:**
1. In the terminal running the demo, press **Ctrl-C**. That stops `go run` and its child.
2. If a port is still stuck (backgrounded run, crash):
   ```bash
   for p in 8080 8081 8082; do lsof -ti tcp:$p | xargs -r kill -9; done
   ```
3. Confirm clear: `for p in 8080 8081 8082; do lsof -ti tcp:$p; done` → prints nothing.

Every demo prints a startup banner ending in **`👉 START HERE: http://localhost:8082`** (or the
equivalent) — when you see it, the servers are up.

---

## 4. Demo runbook (per demo: command → click → money moment → reset)

> Run each from its lab dir. All open at **http://localhost:8082**. The "money moment" is the one
> thing the room must *see*; say the landing line from the deck script as it happens.

### Day 1 — the full OAuth app (Labs 1–4)
- **Dir:** `resource/labs/day1-go-oauth-demo`
- **Run:** `go run .`
- **Open:** `http://localhost:8082`
- **What it is:** one app running AS (8080) + API (8081) + Client (8082). The four Day-1 lab
  exercises are all buttons in this one running app — you do **not** restart between them.
- **Money moments:**
  - *Auth Code + PKCE* → after approving, the client shows a real access token and calls the API successfully.
  - *Client Credentials* → a token with **no user** (`sub` is the service).
  - *Refresh* → old access token swapped for a new one (needs Auth Code done first).
  - *Device* → user code entered on a "second device," client polls and then succeeds.
- **Reset:** Ctrl-C only when Day 1 labs are fully done.

### Day 2 · Exercise 1 — JWT validation (Decks 4 & 5)
- **Dir:** `resource/labs/day2-oauth-demo` · **Run:** `go run ./cmd/jwt-validation` · **Open:** `:8082`
- **Click:** run the login flow → call the API with the token → then use the **tamper** action.
- **Money moment:** valid token → **200**, *"JWT verified locally, via JWKS."* Tamper one byte →
  **401 "signature check failed"** — and stress **no call was made to the AS**. Also show
  `http://localhost:8080/jwks` for Deck 5 (the published keys; point at the matching `kid`).
- **Reset:** Ctrl-C, then run the next `cmd/`.

### Day 2 · Exercise 2 — DPoP (Deck 7)
- **Run:** `go run ./cmd/dpop` · **Open:** `:8082`
- **Click:** do the flow → note **`token_type: DPoP`** → make a correct call → then hit the
  **`/replay`** action (sends the same token as a plain bearer).
- **Money moment:** correct call → **200**, *"JWT verified locally AND DPoP proof-of-possession
  confirmed."* Replay → **401 "token rejected / DPoP proof rejected."** *The token that worked one
  second ago is refused the instant it's presented without the key.*
- **Reset:** Ctrl-C.

### Day 2 · Exercise 3 — PAR + JAR (Deck 7)
- **Run:** `go run ./cmd/par-jar` · **Open:** `:8082`
- **Click:** start the flow → show the decoded **signed request object** and the PAR push → point
  out the `/authorize` URL now carries only `client_id`, `response_type`, `scope`, and an opaque
  **`request_uri`**. Then try a plain `/authorize` **without** PAR.
- **Money moment:** the browser URL has nothing to tamper with; plain no-PAR request → **400 rejected**.
- **Reset:** Ctrl-C.

### Day 2 · Exercise 4 — FAPI 2.0 capstone (Deck 8)
- **Run:** `go run ./cmd/fapi` · **Open:** `:8082`
- **Click:** run the flow → show the decoded **`private_key_jwt` client assertion** and the PAR
  push → approve → note **`token_type: DPoP`** → call the API.
- **Money moment:** *"Authenticated with private_key_jwt; received a DPoP-bound token."* → API
  returns **200** confirming JWT + DPoP. A non-DPoP call → **401 "FAPI 2.0 requires a DPoP-bound
  token."** Say: the server uses `WithProfile(ProfileFAPI2)` and **refuses to start** if any
  control is missing — **fail-closed**.
- **Reset:** Ctrl-C.

---

## 5. Which lab is which (terminology bridge)

The deck scripts say **"Lab N"**; the lab files say **"Exercise N."** Same thing:

| Day | Deck says | Lab file (`LAB.md`) | Demo |
|---|---|---|---|
| 1 | Lab 1 | Exercise 1 — Auth Code + PKCE | Day-1 app |
| 1 | Lab 2 / 3 | Exercises 1–3 (params, interception, refresh) | Day-1 app |
| 1 | Lab 4 | Exercises 2 & 4 — Client Credentials / Device | Day-1 app |
| 2 | Exercise 1 | Exercise 1 — JWT validation | `cmd/jwt-validation` |
| 2 | Exercise 2 | Exercise 2 — DPoP | `cmd/dpop` |
| 2 | Exercise 3 | Exercise 3 — PAR + JAR | `cmd/par-jar` |
| 2 | Exercise 4 | Exercise 4 — FAPI 2.0 | `cmd/fapi` |

---

## 6. When a demo fails (it will, once)

Live demos break. The difference between a pro and a panic is having a plan.

1. **Don't debug on stage.** Give it one retry (usually a leftover process — do the §3 reset).
2. **If it doesn't come back in ~20 seconds, switch to the fallback recording** and narrate over
   it exactly as if it were live. The room barely notices.
3. Say it lightly: *"The demo gods are testing us — here's the recording of this exact run,"* and
   move on. Never apologize twice.

**Prepare the fallbacks in advance** (part of pre-flight): a screenshot or short screen-recording
of each money moment above — the **401** for jwt-validation and dpop, the **400** for par-jar, the
**200** for fapi. Store them next to the deck.

**Most common real failures & fixes:**
- *"address already in use"* → a previous demo is still up. §3 reset.
- *Blank page at `:8082`* → the server isn't up yet; wait for the `START HERE` banner (first run
  compiles — that's why you warm up in pre-flight).
- *Refresh exercise "no refresh token"* → Auth Code (Exercise 1) must be done first in the same session.

---

## 7. Q&A survival kit

You do **not** have to know everything. A crisp answer or a clean punt both build trust.

**Punt phrasing (memorize one):** *"Great question — let me take that at the break so I don't
derail the group."* Then actually follow up at the break.

**The questions this material always draws:**

1. **Why not just use API keys?** — API keys are one static secret with no per-user scope, no
   expiry, no easy revoke. OAuth issues scoped, expiring, revocable tokens without sharing the
   user's password.
2. **Isn't a JWT encrypted?** — No. A signed JWT (JWS) is base64url — readable by anyone.
   Encryption is a separate thing (JWE). Never put secrets in a signed-only JWT.
3. **Do I still need PKCE if I have a client secret?** — Yes. Modern guidance (OAuth 2.1) is PKCE
   for all clients. The secret doesn't stop a stolen *code* on the redirect; PKCE does.
4. **Access vs refresh vs ID token?** — Access = call APIs (short-lived). Refresh = get new access
   tokens (sensitive, back-channel only). ID = who the user is (OIDC; for the client, not the API).
5. **OAuth vs OIDC in one line?** — OAuth = authorization (access an API). OIDC = authentication
   (log in), adds the ID token on top.
6. **Where does the resource server get the key to validate a JWT?** — From the AS's published
   **JWKS** (a trusted `jwks_uri` set in config), selecting by `kid`. Never from the token itself.
7. **What stops `alg: none`?** — Pin the accepted algorithm in server config and reject anything
   else. Never trust the token's own `alg`.
8. **Local validation vs introspection?** — Local (JWT + JWKS) is fast and offline but can't see
   instant revocation. Introspection calls the AS every time (sees revocation, slower). JWT →
   local; opaque → introspect.
9. **Is a bearer token really like a password?** — Yes — while valid, anyone holding it can use
   it. That's exactly why sender-constraining (DPoP/mTLS) exists.
10. **DPoP or mTLS — which?** — mTLS for confidential/high-security clients with certificate
    infrastructure; DPoP (app-layer proof) for public clients and APIs. Same goal, different layer.
11. **Is FAPI a new protocol?** — No. It's a **profile**: a mandatory *combination* of existing
    OAuth controls for high-value APIs.
12. **What if the refresh token is stolen?** — Rotation + reuse detection: if an already-used
    refresh token appears again, treat it as theft and revoke the whole token family.
13. **Can I use Client Credentials to act as a user?** — No. That flow has no user; the service is
    acting as itself. Use Auth Code + PKCE when a user is involved.
14. **Why send `redirect_uri` to the token endpoint too?** — It binds the code exchange to the
    original redirect target, so a code can't be redeemed against a different one.
15. **If I use PKCE, can I drop `state`?** — No. They solve different problems (CSRF vs code
    interception). Use both.

---

## 8. Say-it-right card (pronunciation)

Mispronouncing the core terms is the fastest way to lose a technical room. Pick one and be
consistent.

| Term | Say | Notes |
|---|---|---|
| OAuth | "oh-auth" | |
| JWT | "jot" *or* spell "J-W-T" | "jot" is the RFC-blessed pronunciation |
| JOSE | "JOH-zay" (some say "joe-zee") | the umbrella family |
| JWS / JWE | "J-W-S" / "J-W-E" | signing / encryption |
| JWK / JWKS | "J-W-K" / "J-W-K-S" | one key / the published set |
| PKCE | "pixie" | Proof Key for Code Exchange |
| DPoP | "dee-pop" | Demonstrating Proof-of-Possession |
| PAR / JAR / JARM | "P-A-R" / "jar" / "jarm" | pushed request / signed request / signed response |
| FAPI | "FAP-ee" (community also says "fappy") | Financial-grade API |
| mTLS | "em-T-L-S" | mutual TLS |
| nonce | "nonss" | number used once |
| `kid` | "kid" | key ID (in the JWT header) |
| `iss` / `aud` / `exp` / `sub` / `nbf` | "issuer / audience / expiry / subject / not-before" | say the word, not the letters |

---

## 9. Depth guardrails — where NOT to go

Juniors lose the room by over-explaining or chasing a rabbit hole. When these come up: **one
sentence, then move on** (or punt to the break).

- **Lesson 1:** Don't recite OAuth 2.0-vs-2.1 RFC history, or enumerate every grant here — that's Lesson 3.
- **Lesson 2:** Don't derive base64url or SHA-256 internals. Don't introduce `nonce` (that's OIDC, not this flow).
- **Lesson 3:** Don't go deep on device-flow polling intervals/backoff math.
- **Lesson 4:** Don't teach every JOSE header parameter. Don't turn "JWT vs opaque" into a philosophy debate.
- **Lesson 5:** Don't go into JWE algorithms (RSA-OAEP, AES-GCM…) or a full key-ceremony story.
- **Lesson 6:** Don't drift into a general web-security lecture (deep XSS/CSRF). Stay on the OAuth attack→control table.
- **Lesson 7:** **Do not attempt a live JARM or mTLS demo — there isn't one.** Teach them from the
  slide, say so plainly, move on. The runnable demos are DPoP and PAR/JAR only.
- **Lesson 8:** Don't recite FAPI spec clause numbers. Keep it to "a required combination,
  fail-closed."

---

## 10. Running the labs

- **Setup once per day.** Day 1: everyone runs `go run .` in `day1-go-oauth-demo` — all four
  exercises live in that one app at `:8082`; nobody restarts between them. Day 2: one
  `go run ./cmd/<name>` per exercise, **stopped before the next** (§3).
- **Give a clear target.** Tell the room the money moment to reach ("you're done when you see the
  401 signature failure"). "Done" = they saw the expected output, not "they read the file."
- **Circulate.** Walk the room; most stuck participants hit one of the §6 failures (port conflict,
  or skipping Exercise 1 before Refresh).
- **Timebox out loud.** Announce the clock ("~15 minutes, I'll call time"), then debrief with the
  Questions at the end of each exercise in `LAB.md`.
- **Answer key:** the `## Questions` blocks in each `LAB.md`, backed by the Q&A in §7 of this handbook.

---

## 11. Delivery notes

- **Repeat the one-liners.** Each deck script has a "Say it like this" line — say it more than
  once. Repetition is what participants carry home.
- **Problem before acronym, every time.** Make the room feel the pain first; the control is the relief.
- **Pause after the money moment.** Let the 401 land. Silence is fine.
- **You are the expert in the room.** You don't need every edge case — you need the mental model,
  which these scripts and demos give you. If you don't know, punt cleanly (§7) and follow up.
