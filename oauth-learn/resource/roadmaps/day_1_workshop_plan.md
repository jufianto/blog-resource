# Day 1 Workshop Plan

**Workshop format**  
2-day OAuth/JWT/FAPI workshop.

**Day 1 focus**  
Build the OAuth mental model before JWT, JOSE, and advanced security controls.

# Day 1 Outcome

By the end of Day 1, participants should be able to:

- Explain the OAuth problem without mentioning JWT.
- Map actors and endpoints in common scenarios.
- Draw Authorization Code Flow and Authorization Code + PKCE.
- Explain what travels through the browser and what happens back-channel.
- Explain why `state`, `redirect_uri`, `code_challenge`, and `code_verifier` matter.
- Choose between Authorization Code, Client Credentials, Refresh Token, and Device Authorization Flow.

# Teaching Style

Day 1 should not be pure theory.

Use this rhythm:

1. Explain the concept.
2. Show the diagram.
3. Run the Go demo app.
4. Show the real browser redirect, callback, token request, and API call.
5. Ask participants to explain what happened.

Participants do not need to create code on Day 1. The instructor shares and runs the Go demo app, then explains the implementation.

# Schedule

| Time | Session | Material | Lab |
|---|---|---|---|
| 09:00-09:20 | Opening and mental model | Deck 1 | Demo overview |
| 09:20-10:20 | OAuth problem, actors, endpoints | Deck 1 | Go app actor mapping |
| 10:20-10:35 | Break |  |  |
| 10:35-12:00 | Authorization Code Flow | Deck 2 | Auth Code + PKCE in Go app |
| 12:00-13:00 | Lunch |  |  |
| 13:00-14:20 | PKCE deep dive | Deck 2 | Code verifier/challenge walkthrough |
| 14:20-14:35 | Break |  |  |
| 14:35-15:50 | OAuth flow comparison | Deck 3 | Client Credentials, Refresh Token, Device Flow in Go app |
| 15:50-16:30 | Day 1 review | Decks 1-3 | Real flow recap |
| 16:30-17:00 | Questions and recap | All Day 1 material | Instructor-led discussion |

# Day 1 Files

Presentations:

- `resource/presentations/oauth_deck_1_core_model.md`
- `resource/presentations/oauth_deck_2_auth_code_pkce.md`
- `resource/presentations/oauth_deck_3_flow_comparison.md`

Lab app:

- `resource/labs/day1-go-oauth-demo/`
- `resource/labs/day1-go-oauth-demo/main.go`
- `resource/labs/day1-go-oauth-demo/LAB.md`
- `export/labs/day1_go_oauth_demo_lab.pdf`

# Implementation Walkthrough Examples

Day 1 implementation examples come from the Go demo app:

- Authorization URL construction.
- Callback URL containing `code` and `state`.
- Token request using `authorization_code`.
- PKCE `code_verifier` and `code_challenge`.
- Client Credentials token request.
- Refresh Token token request.
- Device Authorization polling request.

Keep the examples small and visible. The goal is recognition and reasoning, not production implementation.
