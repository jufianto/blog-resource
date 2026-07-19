# Day 1 OAuth 2.0 Demo - Flow Testing Guide

This guide details the step-by-step execution of each OAuth 2.0 flow implemented in our demo. It is designed to help instructors demonstrate the concepts and can be directly used as a storyboard for generating workshop presentation slides.

## Architecture Overview

The demo simulates real-world network isolation by running three distinct servers concurrently:

*   **Client Application (`http://localhost:8082`)**: The frontend or backend app acting on behalf of the user or itself.
*   **Authorization Server (`http://localhost:8080`)**: The `go-oidc` OpenID Provider responsible for authenticating users and issuing tokens.
*   **Resource Server (`http://localhost:8081`)**: The API hosting protected resources, validating tokens via introspection.

---

## Flow 1: Authorization Code + PKCE

**Scenario**: A public web application (RepoBoard) wants to access the user's GitHub-like profile.

### Step-by-step Walkthrough:

1.  **User Initiates Login**
    *   **Action:** Click "Start Auth Code + PKCE flow" on `http://localhost:8082`.
    *   **Behind the scenes:** The Client App generates a cryptographic `state` (for CSRF protection), a `code_verifier` (a random secret), and a `code_challenge` (the SHA-256 hash of the verifier).
2.  **Browser Redirect**
    *   **Action:** The browser physically navigates away from Port 8082 to Port 8080.
    *   **Behind the scenes:** The redirect URL contains `response_type=code`, `client_id`, `scope`, `state`, and the `code_challenge`. The `code_verifier` is kept securely in the client's session.
3.  **Authentication & Consent**
    *   **Action:** The Authorization Server displays a login/consent screen. User reviews the requested scopes and clicks "Approve".
    *   **Behind the scenes:** The Authorization Server registers that the user "Alya" has consented.
4.  **Authorization Code Issuance**
    *   **Action:** The Authorization Server redirects the browser back to the Client App (`http://localhost:8082/callback`).
    *   **Behind the scenes:** The URL contains a short-lived, single-use `code` and the original `state`.
5.  **Token Exchange (Back-channel)**
    *   **Behind the scenes:** The Client App verifies the `state` matches. It then makes a backend `POST` request to the Authorization Server's `/token` endpoint, sending the `code` and the secret `code_verifier`.
    *   **Security:** The Authorization Server hashes the `code_verifier` and compares it to the original `code_challenge`. If they match, it proves the client exchanging the code is the same one that initiated the request.
6.  **Token Issuance**
    *   **Action:** The UI updates to show the JSON Token Response.
    *   **Behind the scenes:** The Authorization Server issues an `access_token`, an `id_token` (since `openid` scope was requested), and a `refresh_token`.
7.  **Resource Access**
    *   **Action:** Click "Call resource API with current access token".
    *   **Behind the scenes:** The Client App sends a `GET` request to `http://localhost:8081/api/profile` with `Authorization: Bearer <access_token>`. The Resource Server introspects the token and returns the profile data.

---

## Flow 2: Client Credentials

**Scenario**: A backend worker service (Order Service) needs to call an internal Fraud API. No end-user is involved.

### Step-by-step Walkthrough:

1.  **Client Authentication**
    *   **Action:** Click "Run Client Credentials demo" on `http://localhost:8082`.
    *   **Behind the scenes:** There is no browser redirect. The Client App makes a direct backend `POST` request to the Authorization Server (`http://localhost:8080/token`). It authenticates itself using Basic Auth (`client_id` and `client_secret`) and requests `grant_type=client_credentials`.
2.  **Token Issuance**
    *   **Behind the scenes:** The Authorization Server verifies the credentials and returns an `access_token`. Notice that **no refresh token** or **ID token** is issued, as there is no user session to refresh or identify.
3.  **Resource Access**
    *   **Action:** The UI displays the token and the API response instantly.
    *   **Behind the scenes:** The Client App automatically sends a `GET` request to `http://localhost:8081/api/fraud-report` with the `access_token`. The Resource Server validates the token and returns the report.

---

## Flow 3: Refresh Token

**Scenario**: The user's access token has expired (or the client wants a fresh one), and the client needs to maintain access without forcing the user to log in again.

### Step-by-step Walkthrough:

1.  **Prerequisite**
    *   You must have successfully completed Flow 1 (Auth Code + PKCE) to obtain a `refresh_token` in your current session.
2.  **Refresh Request**
    *   **Action:** Click "Use refresh token" on `http://localhost:8082`.
    *   **Behind the scenes:** The Client App makes a backend `POST` request to the Authorization Server's `/token` endpoint, sending `grant_type=refresh_token` and the actual `refresh_token`.
3.  **Token Rotation & Issuance**
    *   **Action:** The UI displays a new Token Response.
    *   **Behind the scenes:** The Authorization Server validates the refresh token. It returns a *new* `access_token` and a *new* `refresh_token` (this is called Refresh Token Rotation). The old refresh token is immediately invalidated.

---

## Flow 4: Device Authorization

**Scenario**: A CLI tool or Smart TV (Deploy CLI) needs to authenticate a user, but it lacks a suitable web browser or keyboard for a normal login flow.

### Step-by-step Walkthrough:

1.  **Device Initiation**
    *   **Action:** Click "Start Device Authorization demo" on `http://localhost:8082` (representing the CLI/TV).
    *   **Behind the scenes:** The CLI makes a backend `POST` request to the Authorization Server's `/device_authorization` endpoint.
2.  **Code Generation**
    *   **Action:** The UI displays the Authorization Response.
    *   **Behind the scenes:** The Authorization Server returns a `device_code` (long, kept secret by the CLI), a `user_code` (short, to be shown to the user), and a `verification_uri`.
3.  **User Action (Secondary Device)**
    *   **Action:** Click "Open verification page" (simulating the user opening the link on their smartphone).
    *   **Behind the scenes:** The user visits `http://localhost:8080/device`, enters the short `user_code`, and clicks "Verify". They are then prompted to approve the access.
4.  **Client Polling**
    *   **Behind the scenes:** Meanwhile, the CLI tool is designed to repeatedly poll the Authorization Server (`POST http://localhost:8080/token` with `grant_type=urn:ietf:params:oauth:grant-type:device_code` and the `device_code`). If the user hasn't approved yet, the server returns an `authorization_pending` error.
5.  **Token Issuance**
    *   **Action:** Switch back to the Client App UI and click "Poll token endpoint" *after* you have approved the code on the verification page.
    *   **Behind the scenes:** Because the user has now approved the session, the polling request succeeds. The Authorization Server returns an `access_token`.
6.  **Resource Access**
    *   **Action:** The UI displays the successful Token Response and API Response.
    *   **Behind the scenes:** The CLI automatically calls the Resource Server (`http://localhost:8081/api/deploy`) with the new token.