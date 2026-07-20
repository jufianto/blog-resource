# Deck 5: JWS, JWE, JWK and JWKS

**Workshop position**  
Day 2, Session 2.

**Source paper**  
Lesson 5: JWS, JWE, JWK/JWKS.

**Estimated teaching time**  
75 minutes including JWKS validation and rotation.

# Teaching Goal

Participants should distinguish JWS (signing), JWE (encryption), and JWK/JWKS (key
publication), and understand how a resource server finds the right key via `kid` and how key
rotation works. Signed does not mean private.

# Audience Assumption

Participants routinely conflate "signed" with "encrypted" and assume a JWT hides its contents.
They may not know where the verifying key comes from or what happens when it changes.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | JOSE Is a Toolbox, Not One Thing | Different jobs: sign, encrypt, publish keys. | Title |
| 2 | Concept | JWS: Signing | Integrity + authenticity; contents still readable. | Definition |
| 3 | Concept | JWE: Encryption | Confidentiality; contents hidden. | Definition |
| 4 | Diagram | JWS vs JWE Side by Side | Signed ≠ private; choose by the threat. | `jws_jwe_comparison.png` |
| 5 | Concept | JWK and JWKS | A JWK is one key; a JWKS is the published set. | Definition |
| 6 | Concept | Finding the Key: `kid` | Token header `kid` selects the JWKS key. | Header → key mapping |
| 7 | Diagram | JWKS Key Rotation | Old and new keys coexist; verifiers refetch by `kid`. | `jwks_key_rotation_sequence.png` |
| 8 | Demo | Verify Against the JWKS | The RS fetched keys from `/jwks` and matched `kid`. | `cmd/jwt-validation` |
| 9 | Concept | Rotation Without Downtime | Publish new key before signing with it; retire old later. | Timeline |
| 10 | Risk | Encrypted vs Signed Confusion | A signed token leaks its claims if you put secrets in it. | Warning |
| 11 | Checkpoint | Which JOSE Piece? | Match scenario → JWS, JWE, or JWKS. | Quiz |
| 12 | Summary | Sign, Encrypt, Publish | JWS signs, JWE encrypts, JWKS publishes keys. | Recap |

# Implementation Walkthrough Notes

Reuse `cmd/jwt-validation`. Show `GET http://localhost:8080/jwks` — the published key set the
resource server fetches. Point out the `kid` in the token header and the matching `kid` in the
JWKS. Explain that when the RS sees an unknown `kid` it refetches the JWKS — the hook that makes
rotation work. JWE is taught from the diagram (no runnable JWE demo in this module).

# Checkpoint Questions

1. You must hide the token contents from the client. JWS or JWE?
2. How does a resource server pick the correct key to verify a token?
3. During rotation, why do the old and new keys both appear in the JWKS for a while?

# Speaker Notes

The single most common mistake here is "it's a JWT so it's encrypted". Hammer the JWS/JWE
split. Analogy: JWS is a wax seal (you can read the letter, but not forge the seal); JWE is a
locked box (you cannot read it at all). JWKS is the public noticeboard of seals people can
check against. `kid` is the label that says which seal to compare.
