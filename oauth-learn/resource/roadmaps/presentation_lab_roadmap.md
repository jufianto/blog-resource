# OAuth Teaching Roadmap: Presentations and Labs

**Purpose**  
Turn the eight participant papers into teachable slide decks and practical labs.

This roadmap is not a participant handout yet. It is the build plan for the instructor materials, decks, exercises, and later classroom delivery.

# 1. Teaching Goal

Participants should finish the course able to:

- Explain OAuth without starting from JWT.
- Map actors, endpoints, tokens, scopes, audience, and consent.
- Choose the correct OAuth flow for a scenario.
- Draw Authorization Code + PKCE from memory.
- Validate JWT access tokens conceptually.
- Explain JWS, JWE, and JWKS without mixing them up.
- Identify common OAuth security mistakes and mitigations.
- Explain why PAR, JAR, JARM, DPoP, mTLS, and FAPI 2.0 exist.

# 2. Recommended Teaching Format

Use three layers:

1. **Paper**: detailed reference material for reading before or after class.
2. **Presentation**: shorter visual teaching path for live explanation.
3. **Lab**: active practice so participants prove they understand the flow.

The paper explains deeply.  
The presentation teaches visually.  
The lab makes the participant do the thinking.

# 3. Suggested Course Structure

## Chosen Format: 2-Day Workshop

Use the 2-day workshop as the main target.

The 1-day and 4-week structures are still useful references, but all deck and lab work should prioritize the 2-day version first.

## Option A: 1-Day Workshop

Best when participants already have API/security background.

| Session | Topic | Source Lessons | Time |
|---|---|---:|---:|
| 1 | OAuth mental model, actors, endpoints | Lesson 1 | 60 min |
| 2 | Authorization Code + PKCE | Lesson 2 | 90 min |
| 3 | Client Credentials, Refresh Token, Device Flow | Lesson 3 | 75 min |
| 4 | Token Validation, JWT, JOSE | Lessons 4-5 | 120 min |
| 5 | OAuth Security Controls | Lesson 6 | 75 min |
| 6 | Advanced OAuth and FAPI overview | Lessons 7-8 | 90 min |
| 7 | Final scenario review | All | 60 min |

## Option B: 2-Day Workshop

Best for mixed participants.

| Day | Session | Topic | Source Lessons | Main Exercise |
|---|---|---|---|---|
| 1 | 1 | OAuth problem, actors, endpoints | Lesson 1 | Actor mapping |
| 1 | 2 | Auth Code Flow | Lesson 2 | Endpoint and parameter mapping |
| 1 | 3 | PKCE deep dive | Lesson 2 | Auth Code + PKCE walkthrough |
| 1 | 4 | Client Credentials, Refresh Token, Device Flow | Lesson 3 | Choose the correct flow |
| 2 | 1 | Token Validation and JWT Basics | Lesson 4 | Token validation drills |
| 2 | 2 | JWS, JWE, JWK/JWKS | Lesson 5 | JOSE and JWKS rotation |
| 2 | 3 | OAuth Security Controls | Lesson 6 | Threat-model the flow |
| 2 | 4 | PAR, JAR, JARM, DPoP, mTLS, FAPI 2.0 | Lessons 7-8 | Advanced controls decision |
| 2 | 5 | Capstone lab and review | All | Architecture defense |

## Option C: 4-Week Learning Series

Best when participants need time between topics.

| Week | Focus | Output |
|---|---|---|
| 1 | OAuth core model and flows | Draw actors and Auth Code + PKCE from memory |
| 2 | Other flows and token lifecycle | Choose correct flow and design token lifetime rules |
| 3 | JWT, JOSE, validation, keys | Validate token trust and explain key rotation |
| 4 | Security, advanced controls, FAPI | Threat-model a high-value API architecture |

# 4. Deck Plan

Build multiple focused decks instead of one huge deck.

| Deck | Title | Source | Suggested Slides |
|---|---|---|---:|
| 1 | OAuth Core Model | Lesson 1 | 12-16 |
| 2 | Authorization Code + PKCE | Lesson 2 | 14-18 |
| 3 | OAuth Flow Comparison | Lesson 3 | 14-18 |
| 4 | Token Validation, JWT, and JOSE | Lessons 4-5 | 18-24 |
| 5 | OAuth Security Controls | Lesson 6 | 14-18 |
| 6 | Advanced OAuth and FAPI 2.0 | Lessons 7-8 | 20-26 |

Deck rule:

> One slide should teach one idea, one decision, or one diagram.

# 5. Slide Style

Use this structure for most slides:

- **Concept slide**: one statement, one diagram or table.
- **Flow slide**: actors and arrows, with browser/front-channel vs back-channel marked.
- **Decision slide**: compare choices and when to use each.
- **Risk slide**: attack, impact, mitigation.
- **Checkpoint slide**: scenario question before giving the answer.
- **Implementation walkthrough slide**: small request, response, or code snippet that the instructor explains.

Avoid:

- Long paragraphs copied from the paper.
- Too many standards on one slide.
- Raw RFC wording unless needed.
- Participant-facing labels that expose the creation process.

# 6. Lab Plan

For this workshop, labs should be based on runnable example code. Participants do not need to create the code from scratch, but they should see the real app and understand the real flow.

Day 1 starts with one Go demo app:

- `resource/labs/day1-go-oauth-demo/`

The Day 1 Go demo covers:

- Authorization Code + PKCE.
- Refresh Token Flow.
- Client Credentials Flow.
- Device Authorization Flow.

Later labs can add more Go demos for JWT validation, JWKS, security controls, and advanced OAuth.

| Lab | Title | Main Skill |
|---|---|---|
| 1 | Day 1 Go OAuth Demo | Run and explain real OAuth flows in Go |
| 5 | Token Validation Drills | Reject invalid JWT access tokens for the right reason |
| 6 | JOSE and JWKS Key Rotation | Explain JWS, JWE, JWKS, and key rollover |
| 7 | OAuth Threat Model | Match attacks to mitigations |
| 8 | Advanced Controls Decision Lab | Choose PAR, JAR, JARM, DPoP, or mTLS for scenarios |
| 9 | FAPI Architecture Review | Review a high-value API design and find missing controls |
| 10 | Capstone | Design and defend a secure OAuth architecture |

# 7. Lab Format

Each code lab should use the same structure:

1. **What the app demonstrates**
2. **How to run it**
3. **What routes/endpoints exist**
4. **What flow to click through**
5. **What code functions to explain**
6. **What questions to ask participants**
7. **What mistakes to warn about**

For this workshop, "lab" means instructor-led implementation walkthrough with runnable example code. Participants can inspect and run the code, but they are not required to build it from scratch.

# 8. Recommended Folder Plan

```text
resource/
├── presentations/
│   ├── oauth_deck_1_core_model.md
│   ├── oauth_deck_2_auth_code_pkce.md
│   ├── oauth_deck_3_flow_comparison.md
│   ├── oauth_deck_4_token_jwt_jose.md
│   ├── oauth_deck_5_security_controls.md
│   └── oauth_deck_6_advanced_fapi.md
├── labs/
│   ├── day1-go-oauth-demo/
│   ├── lab_05_token_validation.md
│   ├── lab_06_jose_jwks_rotation.md
│   ├── lab_07_oauth_threat_model.md
│   ├── lab_08_advanced_controls_decision.md
│   ├── lab_09_fapi_architecture_review.md
│   └── lab_10_capstone.md
└── roadmaps/
    └── presentation_lab_roadmap.md
```

Generated deck exports can later live under:

```text
export/presentations/
export/labs/
```

# 9. Build Order for the 2-Day Workshop

Important distinction:

- **Build order** means the order we create files.
- **Milestone** means the first useful teaching package after those files exist.

For the 2-day workshop, build in this order:

1. Create Day 1 Go demo app:
   - Authorization Code + PKCE
   - Refresh Token Flow
   - Client Credentials Flow
   - Device Authorization Flow
2. Create Day 1 lab guide and PDF.
3. Create Day 1 decks:
   - Deck 1: OAuth Core Model
   - Deck 2: Authorization Code + PKCE
   - Deck 3: OAuth Flow Comparison
4. Review Day 1 as one teachable block.
5. Create Day 2 code/demo labs:
   - Lab 5: Token Validation Drills
   - Lab 6: JOSE and JWKS Key Rotation
   - Lab 7: OAuth Threat Model
   - Lab 8: Advanced Controls Decision Lab
   - Lab 9: FAPI Architecture Review
   - Lab 10: Capstone
6. Create Day 2 decks:
   - Deck 4: Token Validation, JWT, and JOSE
   - Deck 5: OAuth Security Controls
   - Deck 6: Advanced OAuth and FAPI 2.0
7. Export decks and lab PDFs.
8. Run a final classroom simulation pass.

# 10. First Build Package

Build this first:

- Day 1 Go OAuth demo app
- Day 1 Go OAuth demo lab guide
- Day 1 Go OAuth demo PDF
- Deck 1: OAuth Core Model
- Deck 2: Authorization Code + PKCE
- Deck 3: OAuth Flow Comparison

This is the **Day 1 foundation package**.

The reason the Go demo comes first is practical: the real flow decides what the slides must explain. After the demo is working, the deck can teach exactly what participants will see in the browser and source code.

# 11. Quality Checklist

Before a deck or lab is considered ready:

- Every diagram is readable on projector and PDF.
- Every exercise has an answer key.
- Every lab has at least one common mistake section.
- OAuth and OpenID Connect are not mixed without explanation.
- Every flow identifies browser/front-channel and back-channel requests.
- Every token discussion identifies who issues it and who validates it.
- Advanced topics explain the problem first, then the standard/control.
- No participant material uses AI-branded labels.

# 12. Final Course Output

Target final artifact set:

- Six teaching decks.
- Ten labs.
- Ten instructor answer keys or instructor-note sections.
- Existing eight participant papers.
- Reusable Mermaid diagrams and exported images.
- Optional final quiz or certification-style review.
