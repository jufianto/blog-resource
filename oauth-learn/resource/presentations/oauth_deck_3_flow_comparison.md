# Deck 3: OAuth Flow Comparison

**Workshop position**  
Day 1, Session 4.

**Source paper**  
Lesson 3: Client Credentials, Refresh Token, and Device Authorization Flow.

**Estimated teaching time**  
75 minutes including Lab 4.

# Teaching Goal

Participants should choose the correct OAuth flow for common application types and understand what changes in the token request.

# Audience Assumption

Participants may think every OAuth scenario is "login". They may also think refresh tokens are a separate login method rather than a token lifecycle mechanism.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | Not Every OAuth Flow Is Login | Different application shapes need different flows. | Title |
| 2 | Decision | First Question: Is There a User? | User-based access and service access are different. | Decision tree |
| 3 | Flow | Authorization Code + PKCE | User-delegated access for web, SPA, mobile. | Small recap |
| 4 | Flow | Client Credentials | Service acts on its own behalf. | `client_credentials_sequence.png` |
| 5 | Demo | Client Credentials Token Request | No browser redirect, no user. | HTTP POST |
| 6 | Flow | Refresh Token | Renew access after initial authorization. | `refresh_token_sequence.png` |
| 7 | Demo | Refresh Token Request | Refresh token is more sensitive than access token. | HTTP POST |
| 8 | Flow | Device Authorization | Good for CLI, TV, limited-input devices. | `device_authorization_sequence.png` |
| 9 | Demo | Device Code Start and Polling | User code on one device, authorization in browser. | HTTP requests |
| 10 | Boundary | Login with Google | Use OpenID Connect when the app needs identity login. | OAuth vs OIDC reminder |
| 11 | Comparison | Flow Table | Compare user, redirect, secret, token purpose. | Table |
| 12 | Checkpoint | Choose the Correct Flow | Participants choose flow for scenarios. | Lab 4 |
| 13 | Summary | Flow Selection Rules | Ask: user? browser? public client? token renewal? limited input? | Five-rule recap |

# Implementation Walkthrough Notes

Show how token requests differ.

Client Credentials:

```http
POST https://auth.example.com/token
Authorization: Basic base64(client_id:client_secret)
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&
scope=fraud.check
```

Refresh Token:

```http
POST https://auth.example.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&
refresh_token=old_refresh_token&
client_id=web-client
```

Device Authorization polling:

```http
POST https://auth.example.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=urn:ietf:params:oauth:grant-type:device_code&
device_code=device-code-from-first-response&
client_id=cli-tool
```

# Checkpoint Questions

1. Which flow has no user?
2. Which flow renews access after the initial authorization?
3. Which flow is best for a CLI or TV app?
4. Which flow should a SPA/mobile app use for user-delegated access?
5. Why is Client Credentials not "service login as a user"?

# Speaker Notes

Use this decision rule:

> If a user is approving access, think Authorization Code + PKCE. If a service acts on its own behalf, think Client Credentials. If access is being renewed, think Refresh Token. If input is limited, think Device Authorization.
