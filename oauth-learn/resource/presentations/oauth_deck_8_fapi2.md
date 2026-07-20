# Deck 8: FAPI 2.0

**Workshop position**  
Day 2, Session 4-5 (capstone).

**Source paper**  
Lesson 8: FAPI 2.0.

**Estimated teaching time**  
60 minutes including Exercise 4 and the capstone review.

# Teaching Goal

Participants should understand that FAPI 2.0 is a security *profile* — a required combination
of controls for high-value APIs — not a new protocol, and be able to name which control stops
which attack. It ties together everything from Decks 4-7.

# Audience Assumption

Participants may think FAPI is "a different OAuth" or a single feature. They may not see why a
high-value API needs several controls at once rather than any one of them.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | FAPI 2.0 Is a Profile, Not a Protocol | It constrains OAuth; it does not replace it. | Title |
| 2 | Concept | Why One Control Is Not Enough | High-value APIs need defence in depth. | Motivation |
| 3 | Concept | The FAPI 2.0 Baseline | PAR + PKCE(S256) + sender-constrained + strong client auth. | Requirement list |
| 4 | Diagram | FAPI 2.0 Controls Map | Each requirement and what it defends. | `fapi_controls_map.png` |
| 5 | Concept | Client Authentication | `private_key_jwt` or mTLS — no secrets, no public clients. | Definition |
| 6 | Concept | Attacker Model | What FAPI assumes the attacker can already do. | Threat model |
| 7 | Diagram | High-Value API Architecture | Where the AS, RS, and clients sit. | `fapi_high_value_architecture.png` |
| 8 | Demo | FAPI 2.0 End to End | Assertion → PAR → DPoP-bound token → API. | `cmd/fapi` |
| 9 | Concept | The Profile Enforces Itself | `WithProfile(FAPI2)` rejects non-compliant config. | Fail-closed |
| 10 | Concept | What FAPI Is Not | Not a login system, not a product, not one toggle. | Myth-busting |
| 11 | Checkpoint | Map Control to Attack | Participants defend an architecture. | Exercise 4 + review |
| 12 | Summary | Defence in Depth, Enforced | PAR + PKCE + DPoP + private_key_jwt together. | Recap |

# Implementation Walkthrough Notes

From `resource/labs/day2-oauth-demo`, run `go run ./cmd/fapi` and open
`http://localhost:8082`.

1. Start the flow: show the decoded `private_key_jwt` client assertion and the PAR push.
2. Approve; show `token_type: DPoP` — the client authenticated with the assertion and proved
   possession of its DPoP key.
3. Call the API → 200, "FAPI 2.0: JWT verified locally AND DPoP proof-of-possession confirmed".

Emphasise that the server is configured with `WithProfile(ProfileFAPI2)`, which refuses to
start if any required control (PAR, PKCE, DPoP, strong client auth, short-lived codes) is
missing — the profile is fail-closed.

# Checkpoint Questions

1. Name the FAPI 2.0 baseline controls and the attack each one addresses.
2. Why does FAPI forbid public clients and shared secrets in favour of `private_key_jwt`?
3. Why is a plain bearer token never acceptable for a high-value API?

# Speaker Notes

The core message: FAPI is a *combination*, and its strength is that all the controls are
required together. Use the capstone demo to show the layers stacking — client auth, request
integrity via PAR, and sender-constraining via DPoP — in one flow. Close the workshop with the
architecture defense: give participants a system sketch and have them map each FAPI control to
the threat it removes.
