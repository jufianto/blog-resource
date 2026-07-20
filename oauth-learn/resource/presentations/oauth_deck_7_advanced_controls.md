# Deck 7: Advanced Controls — PAR, JAR, JARM, DPoP, mTLS

**Workshop position**  
Day 2, Session 4.

**Source paper**  
Lesson 7: Advanced OAuth Controls.

**Estimated teaching time**  
90 minutes including Exercises 2 and 3.

# Teaching Goal

Participants should understand why the request, the response, and the token itself each need
hardening: PAR and JAR protect the request, JARM protects the response, and DPoP or mTLS
constrain the token to a holder so a stolen token cannot be replayed.

# Audience Assumption

Participants know the basic Authorization Code + PKCE flow. They may think HTTPS alone protects
the request, and they usually assume any valid access token can be used by anyone who holds it.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | Harden the Request, Response, and Token | Three different weak points, three families of control. | Title |
| 2 | Concept | Problem: Front-Channel Requests | Params in the browser URL can leak or be tampered with. | Threat |
| 3 | Concept | PAR: Push the Request | Send the request to the AS back-channel first; get a `request_uri`. | Definition |
| 4 | Concept | JAR: Sign the Request | The request object is a signed JWT; params can't be altered. | Definition |
| 5 | Concept | JARM: Sign the Response | Authorization response parameters are signed too. | Definition |
| 6 | Diagram | PAR / JAR / JARM Map | Where each one sits in the flow. | `par_jar_jarm_map.png` |
| 7 | Demo | PAR + JAR | Only an opaque `request_uri` reaches the browser. | `cmd/par-jar` |
| 8 | Concept | Problem: Bearer Token Replay | A stolen bearer token just works. | Threat |
| 9 | Diagram | Bearer vs DPoP | DPoP binds the token to a client key. | `bearer_vs_dpop_sequence.png` |
| 10 | Demo | DPoP Sender-Constraint | Proof required per request; replay as bearer is rejected. | `cmd/dpop` |
| 11 | Diagram | mTLS Sender-Constraint | Bind the token to a client TLS certificate instead. | `mtls_sender_constraint_sequence.png` |
| 12 | Comparison | DPoP vs mTLS | Application-layer proof vs transport-layer certificate. | Table |
| 13 | Concept | What Changes at the Resource Server | RS now checks proof-of-possession, not just the token. | RS checklist |
| 14 | Checkpoint | Choose the Control | Match request/response/token threat to control. | Exercises 2-3 |
| 15 | Summary | Request, Response, Token | PAR/JAR, JARM, DPoP/mTLS. | Recap |

# Implementation Walkthrough Notes

From `resource/labs/day2-oauth-demo`:

- `go run ./cmd/par-jar` — show the decoded signed request object, the PAR response, and that
  the `/authorize` URL carries only `client_id`, `response_type`, `scope`, and an opaque
  `request_uri`. Then show that a plain `/authorize` without PAR is rejected (400).
- `go run ./cmd/dpop` — show `token_type: DPoP`, a correct call with a proof (200), and the
  replay of the token as a plain bearer being rejected (401).

JARM and mTLS are taught from their diagrams; there is no runnable demo for them in this module.

# Checkpoint Questions

1. What can leak or be tampered with in a classic front-channel authorization URL? How do PAR
   and JAR each remove that risk?
2. How does a DPoP-bound token defeat replay, and what does the resource server check?
3. When would mTLS be a better sender-constraining choice than DPoP?

# Speaker Notes

Organise the whole deck around three targets: the request (PAR/JAR), the response (JARM), and
the token (DPoP/mTLS). For DPoP, the memorable point is the replay demo: the same token that
worked a moment ago is refused the instant it is presented without the key. For mTLS vs DPoP,
frame it as transport-layer (certificate) vs application-layer (signed proof) — same goal,
different plumbing.
