# Deck 6: OAuth Security Controls

**Workshop position**  
Day 2, Session 3.

**Source paper**  
Lesson 6: OAuth Security Controls.

**Estimated teaching time**  
75 minutes including a threat-modeling discussion.

# Teaching Goal

Participants should connect common OAuth attacks to the specific controls that stop them, and
be able to threat-model a flow: what travels through the browser, what can leak, and which
control closes each gap.

# Audience Assumption

Participants may know the controls as a checklist of acronyms without knowing which attack each
one prevents. They may think "we use OAuth" is itself a security statement.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | Controls Exist Because of Attacks | Every control answers a specific threat. | Title |
| 2 | Concept | The Threat Surfaces | Front channel, back channel, token at rest, token in use. | Four surfaces |
| 3 | Attack | Authorization Code Interception | Stolen code in the redirect. | Attack sketch |
| 4 | Control | PKCE | The verifier the attacker never saw. | Recap from Day 1 |
| 5 | Attack | CSRF on the Callback | Forged callback without `state`. | Attack sketch |
| 6 | Control | `state` and `redirect_uri` | Bind and restrict the redirect. | Recap |
| 7 | Attack | Token Leakage and Replay | A stolen bearer token just works. | Attack sketch |
| 8 | Control | Sender-Constraining (preview) | DPoP/mTLS bind the token to a key — see Deck 7. | Forward pointer |
| 9 | Diagram | Controls Map | Attacks mapped to controls. | `oauth_security_controls_map.png` |
| 10 | Concept | Least Privilege: Scope and Audience | Narrow what a token can do and where. | Scope/aud |
| 11 | Checkpoint | Threat-Model This Flow | Participants name leaks and matching controls. | Group exercise |
| 12 | Summary | Attack → Control Mapping | Say the threat before naming the control. | Recap |

# Implementation Walkthrough Notes

This session is analysis, not a new demo. Reuse the Day 1 `day1-go-oauth-demo` and the Day 2
`cmd/jwt-validation` to point at concrete controls already seen: PKCE required, `state` checked
on the callback, local token validation. Sender-constraining (DPoP/mTLS) is previewed here and
demonstrated in Deck 7. Keep the room in threat-modeling mode: for each step, ask "what leaks,
and what stops it".

# Checkpoint Questions

1. An attacker intercepts the authorization code. Which control makes it useless, and why?
2. What is the difference between restricting a token by scope and by audience?
3. Which of today's controls would you add first to a legacy OAuth integration?

# Speaker Notes

The goal is to stop treating controls as a checklist. Drive every slide from the attack, then
reveal the control. Emphasise that a bearer token is a password-equivalent while it is valid —
which is exactly why Deck 7's sender-constraining matters. Use the controls map as the anchor
diagram participants take away.
