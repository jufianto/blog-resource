# Deck 1: OAuth Core Model

**Workshop position**  
Day 1, Sessions 1-2.

**Source paper**  
Lesson 1: OAuth Problem, Actors, and Endpoints.

**Estimated teaching time**  
60 minutes including Lab 1.

# Teaching Goal

Participants should understand OAuth as delegated API access, not as "JWT login". They should be able to identify actors, endpoints, tokens, scopes, audience, and the OpenID Connect boundary.

# Audience Assumption

Participants may have used "Login with Google" or API tokens before, but they may mix up authentication, authorization, OAuth, OIDC, JWT, and API keys.

# Slide Outline

| Slide | Type | Title | Main Teaching Point | Visual or Demo |
|---:|---|---|---|---|
| 1 | Opening | OAuth Starts With a Problem | Apps need limited API access without user passwords. | One-sentence problem |
| 2 | Concept | Password Sharing Is the Wrong Shape | Sharing a password gives too much power and is hard to revoke. | Before/after comparison |
| 3 | Scenario | RepoBoard Connects GitHub | Use a general GitHub example before banking. | GitHub app/API sketch |
| 4 | Actors | The Four Actors | Resource Owner, Client, Authorization Server, Resource Server. | `oauth_actor_general.png` |
| 5 | Endpoint | Where Requests Go | `/authorize` and `/token` have different jobs. | Endpoint map |
| 6 | Token | What the Access Token Means | Permission to call an API, not the user's password. | Token as permission card |
| 7 | Scope | Scope Is Permission | Scope limits what the client asks for. | Scope examples |
| 8 | Audience | Audience Is the API | The API must know whether the token is meant for it. | API audience diagram |
| 9 | Consent | User Approval | The user approves limited access. | Consent prompt mock |
| 10 | OIDC Boundary | OAuth vs OpenID Connect | OAuth is authorization; OIDC adds login and ID token. | Split table |
| 11 | Demo | Real Configuration Values | Config reveals actors and endpoints. | Show config snippet |
| 12 | Checkpoint | Actor Mapping | Participants identify actors in GitHub, Calendar, CLI, backend scenarios. | Lab 1 |
| 13 | Summary | The Mental Model | Who asks, who approves, who issues, who validates. | Four-question recap |

# Implementation Walkthrough Notes

Show this configuration and map it with participants:

```text
CLIENT_ID=repo-board-web
AUTHORIZATION_ENDPOINT=https://github.com/login/oauth/authorize
TOKEN_ENDPOINT=https://github.com/login/oauth/access_token
RESOURCE_API=https://api.github.com
REDIRECT_URI=https://repoboard.example.com/oauth/callback
SCOPES=repo:read
```

Instructor script:

- `CLIENT_ID` identifies the Client.
- `AUTHORIZATION_ENDPOINT` and `TOKEN_ENDPOINT` belong to the Authorization Server.
- `RESOURCE_API` is the Resource Server.
- `SCOPES` are requested permissions.
- `REDIRECT_URI` is where the browser comes back.

# Checkpoint Questions

1. Is OAuth primarily authentication or authorization?
2. Which actor issues the access token?
3. Which actor validates the access token?
4. When does OpenID Connect enter the story?

# Speaker Notes

Keep repeating:

> OAuth lets an app access an API with limited permission without receiving the user's password.

Do not start with JWT. JWT is only a token format that may appear later.
