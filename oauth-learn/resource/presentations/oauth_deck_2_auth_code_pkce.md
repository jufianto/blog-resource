# Deck 2: Authorization Code Flow + PKCE

**Workshop position**  
Day 1, Sessions 2-3.

**Source paper**  
Lesson 2: Authorization Code Flow and PKCE.

**Estimated teaching time**  
120 minutes including Labs 2 and 3.

# Teaching Goal

Participants should be able to draw Authorization Code Flow and Authorization Code + PKCE from memory, explain browser/front-channel vs back-channel requests, and explain why `state`, `redirect_uri`, `code_challenge`, and `code_verifier` are not interchangeable.

# Audience Assumption

Participants may know that OAuth uses redirects, but they may not know which request returns the code, which request returns the token, or why `/token` receives `redirect_uri`.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | The Main User-Based OAuth Flow | Authorization Code Flow is the foundation for user-delegated access. | Title |
| 2 | Flow | Browser vs Back-Channel | The code travels through the browser; the token comes from `/token`. | `auth_code_sequence.png` |
| 3 | Endpoint | Authorization Request | `/authorize` starts the user authorization. | HTTP request |
| 4 | Parameter | `state` | The client verifies the callback belongs to a request it started. | Callback example |
| 5 | Parameter | `redirect_uri` | It must match what was registered and what was used. | Redirect URI comparison |
| 6 | Endpoint | Token Request | `/token` returns data directly; it does not redirect. | HTTP POST |
| 7 | Question | Why Send `redirect_uri` to `/token`? | It binds the code exchange to the original redirect target. | Code exchange check |
| 8 | Demo | Build Authorization URL | Show tiny route that redirects to authorization endpoint. | JavaScript snippet |
| 9 | Checkpoint | Map the Parameters | Participants map endpoint and parameters. | Lab 2 |
| 10 | Transition | What If the Code Is Stolen? | A stolen code should not be enough. | Attack setup |
| 11 | PKCE Concept | Challenge First, Verifier Later | Server stores challenge, later checks verifier. | `auth_code_pkce_sequence.png` |
| 12 | Crypto | SHA-256 Is One-Way | The server hashes the verifier and compares to challenge. | Comparison diagram |
| 13 | Demo | Create PKCE Values | Show code for verifier and challenge. | JavaScript snippet |
| 14 | Token Request | `code_verifier` Goes to `/token` | The verifier is not sent in the browser authorization request. | HTTP POST |
| 15 | Comparison | `state` vs PKCE | They solve different problems and should both exist. | Two-column table |
| 16 | Checkpoint | Code Interception Scenario | Participants explain why attacker fails. | Lab 3 |
| 17 | Summary | Draw It From Memory | Participants redraw Auth Code + PKCE. | Blank flow |

# Implementation Walkthrough Notes

Show the route that starts authorization:

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

Then show PKCE value creation:

```js
const codeVerifier = base64url(crypto.randomBytes(32));
const codeChallenge = base64url(
  crypto.createHash("sha256").update(codeVerifier).digest()
);
```

Instructor emphasis:

- `state` is stored by the client and checked on callback.
- `code_challenge` goes to `/authorize`.
- `code_verifier` goes to `/token`.
- `/token` returns JSON-like data to the client; it does not redirect the browser.

# Checkpoint Questions

1. Which request is browser/front-channel?
2. Which request is back-channel?
3. Why is `redirect_uri` included in the token request?
4. If an attacker steals only the code, why should PKCE stop token theft?

# Speaker Notes

Use the phrase:

> The code is not the token. The code is a short-lived ticket that must be exchanged correctly.

For PKCE:

> The authorization request says, "later I will prove I know the original verifier."
