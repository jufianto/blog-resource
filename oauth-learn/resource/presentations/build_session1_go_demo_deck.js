// Regenerates the Session 1 Go-demo walkthrough deck.
//   npm install pptxgenjs && node build_session1_go_demo_deck.js
//   mv oauth_session1_go_demo_workshop.pptx ../../export/presentations/

const P = require("pptxgenjs");
const pres = new P();
pres.layout = "LAYOUT_WIDE";            // 13.33 x 7.5
pres.author = "OAuth/JWT/FAPI Workshop";
pres.title = "OAuth Session 1 — Go Demo Walkthrough";

const W = 13.33, H = 7.5;
const INK = "1F3A5F", INK2 = "2F5A8F", ACCENT = "0F8A7E", GREY = "5A6B7B";
const LIGHT = "EEF3F8", CODEBG = "F5F7FB", CODEBORDER = "D8DEEA", WHITE = "FFFFFF";
const BAD = "C0492F", INKSOFT = "DCE7F2";
const HEAD = "Arial", BODY = "Arial", MONO = "Courier New";
const MX = 0.7;                          // left margin

let N = 0;
const shadow = () => ({ type: "outer", color: "1F3A5F", blur: 7, offset: 3, angle: 135, opacity: 0.13 });

// ---------- content slide scaffold (white bg) ----------
function content(kicker, titleText) {
  N++;
  const s = pres.addSlide();
  s.background = { color: WHITE };
  if (kicker) s.addText(kicker.toUpperCase(), { x: MX, y: 0.42, w: 12, h: 0.3, margin: 0,
    fontFace: HEAD, fontSize: 12, bold: true, color: ACCENT, charSpacing: 2 });
  s.addText(titleText, { x: MX, y: 0.72, w: 12, h: 0.85, margin: 0,
    fontFace: HEAD, fontSize: 28, bold: true, color: INK });
  // footer (plain text, no bar)
  s.addText("OAuth 2.0  ·  Session 1  ·  Go demo walkthrough", { x: MX, y: H - 0.42, w: 8, h: 0.3, margin: 0,
    fontFace: BODY, fontSize: 9, color: GREY });
  s.addText(String(N).padStart(2, "0"), { x: W - 1.3, y: H - 0.42, w: 0.6, h: 0.3, margin: 0,
    fontFace: BODY, fontSize: 9, color: GREY, align: "right" });
  return s;
}

// ---------- dark slide scaffold (navy bg) ----------
function dark(kicker, titleText) {
  N++;
  const s = pres.addSlide();
  s.background = { color: INK };
  if (kicker) s.addText(kicker.toUpperCase(), { x: MX, y: 1.7, w: 11, h: 0.4, margin: 0,
    fontFace: HEAD, fontSize: 14, bold: true, color: "7FD9CD", charSpacing: 3 });
  if (titleText) s.addText(titleText, { x: MX, y: 2.1, w: 11.9, h: 1.8, margin: 0,
    fontFace: HEAD, fontSize: 44, bold: true, color: WHITE });
  return s;
}

// code line runs: lines beginning with // rendered grey
function codeRuns(lines) {
  return lines.map((ln, i) => {
    const isComment = ln.trimStart().startsWith("//");
    return { text: ln === "" ? " " : ln, options: {
      breakLine: true, fontFace: MONO, fontSize: 11.5,
      color: isComment ? GREY : "1D2733", italic: isComment } };
  });
}
function codePanel(s, x, y, w, h, lines, caption) {
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y, w, h, rectRadius: 0.06,
    fill: { color: CODEBG }, line: { color: CODEBORDER, width: 1 }, shadow: shadow() });
  if (caption) s.addText(caption, { x: x + 0.18, y: y + 0.1, w: w - 0.36, h: 0.28, margin: 0,
    fontFace: MONO, fontSize: 10.5, bold: true, color: ACCENT });
  s.addText(codeRuns(lines), { x: x + 0.22, y: y + (caption ? 0.42 : 0.16), w: w - 0.44, h: h - (caption ? 0.55 : 0.3),
    margin: 0, valign: "top", lineSpacingMultiple: 1.06 });
}
// labeled card
function card(s, x, y, w, h, title, lines, opt = {}) {
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y, w, h, rectRadius: 0.06,
    fill: { color: opt.fill || LIGHT }, line: { color: opt.border || CODEBORDER, width: 1 }, shadow: shadow() });
  let yy = y + 0.16;
  if (title) { s.addText(title, { x: x + 0.2, y: yy, w: w - 0.4, h: 0.34, margin: 0,
    fontFace: HEAD, fontSize: 14, bold: true, color: opt.titleColor || INK }); yy += 0.46; }
  if (lines && lines.length) {
    const runs = lines.map((t, i) => ({ text: t, options: { breakLine: true, bullet: opt.bullet ? { code: "2022", indent: 14 } : false,
      fontFace: BODY, fontSize: opt.fontSize || 12, color: opt.color || "1D2733", paraSpaceAfter: 4 } }));
    s.addText(runs, { x: x + 0.2, y: yy, w: w - 0.4, h: y + h - yy - 0.12, margin: 0, valign: "top" });
  }
}
// teal numbered chip
function chip(s, x, y, n) {
  s.addShape(pres.shapes.OVAL, { x, y, w: 0.42, h: 0.42, fill: { color: ACCENT } });
  s.addText(String(n), { x, y, w: 0.42, h: 0.42, margin: 0, align: "center", valign: "middle",
    fontFace: HEAD, fontSize: 15, bold: true, color: WHITE });
}

// ============================================================ SLIDE 1 — TITLE
{
  const s = dark("Workshop · Session 1", "OAuth 2.0 — the real flow, in Go");
  s.addText("Learn OAuth by running a real go-oidc demo end to end.", { x: MX, y: 3.85, w: 11, h: 0.5,
    margin: 0, fontFace: BODY, fontSize: 18, color: INKSOFT });
  // run box
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: 4.5, w: 6.1, h: 1.5, rectRadius: 0.06,
    fill: { color: "16304E" }, line: { color: "32517A", width: 1 } });
  s.addText("RUN THE LAB", { x: MX + 0.22, y: 4.62, w: 5, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 11, bold: true, color: "7FD9CD", charSpacing: 2 });
  s.addText([
    { text: "cd resource/labs/day1-go-oauth-demo", options: { breakLine: true } },
    { text: "go run .", options: { breakLine: true } },
    { text: "open http://localhost:8082", options: {} },
  ], { x: MX + 0.22, y: 4.95, w: 5.7, h: 0.95, margin: 0, fontFace: MONO, fontSize: 13, color: "D7E6F5", lineSpacingMultiple: 1.1 });
  // flow chips on the right
  const flows = ["Auth Code + PKCE", "Refresh Token", "Client Credentials", "Device Authorization"];
  flows.forEach((f, i) => {
    const yy = 4.55 + i * 0.4;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: 7.4, y: yy, w: 5.2, h: 0.34, rectRadius: 0.17,
      fill: { color: "16304E" }, line: { color: "32517A", width: 1 } });
    s.addText(f, { x: 7.4, y: yy, w: 5.2, h: 0.34, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 12, color: "BFE9E1" });
  });
}

// ============================================================ SLIDE 2 — OUTCOME
{
  const s = content("Session 1 outcome", "Explain the flow without hiding behind JWT");
  const items = [
    "Map the four OAuth actors and the endpoints onto the running app.",
    "Draw Authorization Code + PKCE from memory.",
    "Explain browser / front-channel vs back-channel requests.",
    "Explain state, redirect_uri, code_challenge, and code_verifier.",
    "Choose the right flow for common application shapes.",
  ];
  items.forEach((t, i) => {
    const yy = 1.85 + i * 1.0;
    chip(s, MX, yy, i + 1);
    s.addText(t, { x: MX + 0.65, y: yy - 0.04, w: 11.3, h: 0.55, margin: 0, valign: "middle", fontFace: BODY, fontSize: 16, color: "1D2733" });
  });
}

// ============================================================ SLIDE 3 — THE DEMO APP (3 servers)
{
  const s = content("The demo app", "One Go program, three separate servers");
  s.addText("Real network separation — each OAuth role is its own HTTP server (see main.go).",
    { x: MX, y: 1.55, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, color: GREY });
  const cols = [
    { t: "Authorization Server", p: ":8080  ·  op.go", b: ["go-oidc OpenID Provider", "/authorize  /token", "/introspect", "/device_authorization"] },
    { t: "Resource Server API", p: ":8081  ·  api.go", b: ["Protected endpoints", "/api/profile", "/api/fraud-report", "/api/deploy"] },
    { t: "Client App", p: ":8082  ·  client.go", b: ["RepoBoard + service + CLI", "/start  /callback", "/refresh  /call-api", "START HERE → :8082"] },
  ];
  const cw = 3.85, gap = 0.3; let x = MX;
  cols.forEach((c, i) => {
    const hl = i === 2;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y: 2.1, w: cw, h: 3.6, rectRadius: 0.06,
      fill: { color: hl ? "E9F5F2" : LIGHT }, line: { color: hl ? "A9D8CF" : CODEBORDER, width: hl ? 1.5 : 1 }, shadow: shadow() });
    s.addText(c.t, { x: x + 0.22, y: 2.32, w: cw - 0.44, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 15, bold: true, color: INK });
    s.addText(c.p, { x: x + 0.22, y: 2.74, w: cw - 0.44, h: 0.3, margin: 0, fontFace: MONO, fontSize: 11.5, bold: true, color: ACCENT });
    s.addText(c.b.map((t, j) => ({ text: t, options: { breakLine: true, fontFace: (t.includes("/") || t.includes(":8082")) ? MONO : BODY,
      fontSize: 12, color: "1D2733", paraSpaceAfter: 5 } })),
      { x: x + 0.22, y: 3.2, w: cw - 0.44, h: 2.4, margin: 0, valign: "top" });
    x += cw + gap;
  });
}

// ============================================================ SLIDE 4 — HOW WE TEACH
{
  const s = content("Teaching rhythm", "How we teach each flow");
  const rows = [
    ["Concept", "What problem are we solving?"],
    ["Diagram", "Who talks to whom, and on which channel?"],
    ["Go code", "Which function in client.go / op.go / api.go?"],
    ["Browser flow", "What URL or token actually appears?"],
    ["Checkpoint", "Can participants explain it back?"],
  ];
  rows.forEach((r, i) => {
    const yy = 1.8 + i * 0.92;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: yy, w: 2.7, h: 0.7, rectRadius: 0.06, fill: { color: INK } });
    s.addText(r[0], { x: MX, y: yy, w: 2.7, h: 0.7, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 15, bold: true, color: WHITE });
    s.addText(r[1], { x: MX + 3.0, y: yy, w: 9.0, h: 0.7, margin: 0, valign: "middle", fontFace: BODY, fontSize: 15, color: "1D2733" });
  });
  s.addText("Participants run and read the demo — they do not write code in Session 1.",
    { x: MX, y: 6.5, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13, italic: true, color: GREY });
}

// ============================================================ SLIDE 5 — OAUTH PROBLEM
{
  const s = content("The problem", "Limited access without the password");
  s.addText("RepoBoard needs to read the user's GitHub-like profile — without ever getting their password.",
    { x: MX, y: 1.55, w: 12, h: 0.5, margin: 0, fontFace: BODY, fontSize: 16, color: "1D2733" });
  card(s, MX, 2.3, 5.85, 3.2, "Password sharing — wrong shape", [
    "User hands RepoBoard a password",
    "App can do everything the user can",
    "No way to scope to one resource",
    "Revoking means changing the password",
  ], { fill: "FBECEA", border: "E6B6AC", titleColor: BAD, bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 2.3, 5.85, 3.2, "OAuth delegation — right shape", [
    "User approves limited API access",
    "Client receives a scoped token, not a password",
    "Access can be revoked on its own",
    "Token expires automatically",
  ], { fill: "E9F5F2", border: "A9D8CF", titleColor: ACCENT, bullet: true, fontSize: 13 });
  s.addText("The client receives permission, not the user's password.",
    { x: MX, y: 5.75, w: 12, h: 0.5, margin: 0, fontFace: BODY, fontSize: 15, bold: true, color: INK });
}

// ============================================================ SLIDE 6 — ACTOR MAP
{
  const s = content("Actors in the Go app", "Map the four actors onto the code");
  const rows = [
    ["Resource Owner", "Mock user “Alya Developer” (Subject = alya)"],
    ["Client", ":8082 · client.go — RepoBoard (client_id repoboard-web)"],
    ["Authorization Server", ":8080 · op.go — go-oidc provider (login, consent, tokens)"],
    ["Resource Server", ":8081 · api.go — protected APIs, validates tokens"],
  ];
  rows.forEach((r, i) => {
    const yy = 1.75 + i * 0.92;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: yy, w: 3.5, h: 0.74, rectRadius: 0.06, fill: { color: LIGHT }, line: { color: CODEBORDER, width: 1 } });
    s.addText(r[0], { x: MX + 0.15, y: yy, w: 3.2, h: 0.74, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: INK });
    s.addText(r[1], { x: MX + 3.8, y: yy, w: 8.3, h: 0.74, margin: 0, valign: "middle", fontFace: BODY, fontSize: 13.5, color: "1D2733" });
  });
  s.addText("Ask the room: which actor issues the token, and which actor validates it?",
    { x: MX, y: 5.7, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, italic: true, color: GREY });
}

// ============================================================ SLIDE 7 — ENDPOINT MAP
{
  const s = content("Endpoint map", "/authorize and /token are not interchangeable");
  card(s, MX, 1.8, 5.85, 2.9, "/authorize   (:8080, front-channel)", [
    "Browser redirect — the user sees it",
    "Login + consent happen here",
    "Returns: code + state in the redirect",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.8, 5.85, 2.9, "/token   (:8080, back-channel)", [
    "Server-to-server POST — no redirect",
    "Exchanges code (+ verifier) for tokens",
    "Returns: JSON token response",
  ], { bullet: true, fontSize: 13 });
  s.addText([
    { text: "Both endpoints are served by go-oidc ", options: {} },
    { text: "(op.Handler())", options: { fontFace: MONO, color: ACCENT } },
    { text: " — you do not implement them. Your client calls them from client.go.", options: {} },
  ], { x: MX, y: 5.0, w: 12, h: 0.7, margin: 0, fontFace: BODY, fontSize: 14, color: "1D2733" });
}

// ============================================================ SLIDE 8 — AUTH CODE + PKCE OVERVIEW
{
  const s = content("Authorization Code + PKCE", "The browser gets a code; the client gets tokens");
  const steps = [
    ["Client", "/start  (:8082)"],
    ["Authorize", "/authorize  (:8080)"],
    ["Callback", "/callback?code&state  (:8082)"],
    ["Token", "POST /token  (:8080)"],
    ["API", "GET /api/profile  (:8081)"],
  ];
  steps.forEach((st, i) => {
    const yy = 1.95 + i * 0.78;
    chip(s, MX, yy, i + 1);
    s.addText(st[0], { x: MX + 0.6, y: yy - 0.02, w: 2.4, h: 0.46, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 15, bold: true, color: INK });
    s.addText(st[1], { x: MX + 3.0, y: yy - 0.02, w: 5.5, h: 0.46, margin: 0, valign: "middle", fontFace: MONO, fontSize: 13, color: "1D2733" });
  });
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: 8.9, y: 1.95, w: 3.7, h: 3.6, rectRadius: 0.06, fill: { color: "E9F5F2" }, line: { color: "A9D8CF", width: 1 }, shadow: shadow() });
  s.addText("Demo anchor", { x: 9.1, y: 2.12, w: 3.3, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 13, bold: true, color: ACCENT });
  s.addText("Click “Start Auth Code + PKCE flow” on :8082, then ask participants what travels through the browser vs the back-channel.",
    { x: 9.1, y: 2.5, w: 3.35, h: 2.8, margin: 0, valign: "top", fontFace: BODY, fontSize: 13, color: "1D2733" });
}

// ============================================================ SLIDE 9 — START AUTHORIZATION (real Go)
{
  const s = content("Start authorization", "Client builds the request before redirecting");
  codePanel(s, MX, 1.7, 7.4, 4.6, [
    "func (c *ClientApp) StartAuthCode(w, r) {",
    "  state         := randomString(18)",
    "  codeVerifier  := randomString(48)",
    "  codeChallenge := pkceChallenge(codeVerifier)",
    "  // save state + verifier in the session",
    "",
    "  q := url.Values{}",
    "  q.Set(\"response_type\", \"code\")",
    "  q.Set(\"client_id\", \"repoboard-web\")",
    "  q.Set(\"redirect_uri\", clientBaseURL+\"/callback\")",
    "  q.Set(\"scope\", \"openid profile profile.read offline_access\")",
    "  q.Set(\"state\", state)",
    "  q.Set(\"code_challenge\", codeChallenge)",
    "  q.Set(\"code_challenge_method\", \"S256\")",
    "  http.Redirect(w, r, opBaseURL+\"/authorize?\"+q.Encode(), 302)",
    "}",
  ], "client.go");
  card(s, 8.35, 1.7, 4.25, 4.6, "What to point at", [
    "state + code_verifier are created on the client, before any redirect.",
    "Only the code_challenge (a hash) goes through the browser.",
    "The verifier stays in the session — it is never sent to /authorize.",
  ], { bullet: true, fontSize: 13 });
}

// ============================================================ SLIDE 10 — AUTH REQUEST PARAMS
{
  const s = content("Authorization request", "Every parameter in the redirect, explained");
  const rows = [
    ["response_type=code", "Ask for an authorization code"],
    ["client_id", "Identify the client app (repoboard-web)"],
    ["redirect_uri", "Where the browser returns (:8082/callback)"],
    ["scope", "Permissions requested"],
    ["state", "Callback binding — CSRF protection"],
    ["code_challenge (+ _method=S256)", "PKCE challenge for the code exchange"],
  ];
  rows.forEach((r, i) => {
    const yy = 1.75 + i * 0.8;
    s.addText(r[0], { x: MX, y: yy, w: 5.2, h: 0.6, margin: 0, valign: "middle", fontFace: MONO, fontSize: 13, bold: true, color: INK });
    s.addText(r[1], { x: MX + 5.4, y: yy, w: 6.6, h: 0.6, margin: 0, valign: "middle", fontFace: BODY, fontSize: 14, color: "1D2733" });
  });
}

// ============================================================ SLIDE 11 — USER APPROVAL (go-oidc policy)
{
  const s = content("User approval", "Login & consent is a go-oidc policy");
  codePanel(s, MX, 1.7, 7.4, 4.6, [
    "loginPolicy := goidc.NewPolicy(\"simple_login\",",
    "  func(...) bool { return true },   // applies to all",
    "  func(w, r, as, _) (goidc.Status, error) {",
    "    if r.Method == GET {",
    "      // render the consent screen",
    "      return goidc.StatusPending, nil",
    "    }",
    "    if r.FormValue(\"approve\") == \"true\" {",
    "      as.Subject       = \"alya\"",
    "      as.GrantedScopes = as.Scopes",
    "      return goidc.StatusSuccess, nil",
    "    }",
    "    return goidc.StatusFailure, nil",
    "  })",
  ], "op.go");
  card(s, 8.35, 1.7, 4.25, 4.6, "What to point at", [
    "This is the hook where a real app would check username / password / MFA.",
    "go-oidc creates and stores the authorization code internally on success.",
    "The code is short-lived and is not an access token.",
  ], { bullet: true, fontSize: 13 });
}

// ============================================================ SLIDE 12 — CALLBACK (real Go)
{
  const s = content("Callback", "Verify state, then exchange the code");
  codePanel(s, MX, 1.7, 7.4, 4.6, [
    "func (c *ClientApp) Callback(w, r) {",
    "  code  := r.URL.Query().Get(\"code\")",
    "  state := r.URL.Query().Get(\"state\")",
    "  if state == \"\" || state != expectedState {",
    "    http.Error(w, \"state mismatch — possible CSRF\", 400)",
    "    return",
    "  }",
    "  // back-channel POST to /token",
    "  postForm(opBaseURL+\"/token\", url.Values{",
    "    \"grant_type\":    {\"authorization_code\"},",
    "    \"client_id\":     {\"repoboard-web\"},",
    "    \"code\":          {code},",
    "    \"redirect_uri\":  {clientBaseURL + \"/callback\"},",
    "    \"code_verifier\": {codeVerifier},",
    "  })",
    "}",
  ], "client.go");
  card(s, 8.35, 1.7, 4.25, 4.6, "What to point at", [
    "The callback arrives through the browser (front-channel).",
    "state is checked first — reject anything the client did not start.",
    "The code-for-token exchange is back-channel, and carries the code_verifier.",
  ], { bullet: true, fontSize: 13 });
}

// ============================================================ SLIDE 13 — TOKEN EXCHANGE (go-oidc verifies)
{
  const s = content("Token exchange", "go-oidc verifies the exchange for you");
  card(s, MX, 1.75, 5.85, 3.5, "At /token, go-oidc checks:", [
    "The code exists and has not been used",
    "client_id and redirect_uri match the original request",
    "SHA-256(code_verifier) equals the stored code_challenge",
    "Then it issues access + refresh tokens",
  ], { bullet: true, fontSize: 13 });
  codePanel(s, MX + 6.2, 1.75, 5.85, 3.5, [
    "// enabled by provider options in op.go",
    "provider.WithAuthCodeGrant(store, ResponseTypeCode),",
    "provider.WithPKCE(CodeChallengeMethodSHA256),",
    "provider.WithRefreshTokenGrant(store),",
    "provider.WithRefreshTokenRotation(),",
  ], "op.go");
  s.addText("Your client never implements PKCE verification — it is a configuration choice on the server.",
    { x: MX, y: 5.5, w: 12, h: 0.5, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}

// ============================================================ SLIDE 14 — WHY redirect_uri AT /token
{
  const s = content("Why redirect_uri goes to /token", "It binds the exchange to the original request");
  const rows = [
    ["At /authorize", "go-oidc stores client_id, redirect_uri, scope, state, code_challenge."],
    ["At /callback", "The browser returns code + state to the client."],
    ["At /token", "The client resends redirect_uri + code_verifier; go-oidc checks they match."],
  ];
  rows.forEach((r, i) => {
    const yy = 1.95 + i * 1.1;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: yy, w: 2.9, h: 0.8, rectRadius: 0.06, fill: { color: LIGHT }, line: { color: CODEBORDER, width: 1 } });
    s.addText(r[0], { x: MX, y: yy, w: 2.9, h: 0.8, margin: 0, align: "center", valign: "middle", fontFace: MONO, fontSize: 13, bold: true, color: ACCENT });
    s.addText(r[1], { x: MX + 3.2, y: yy, w: 8.9, h: 0.8, margin: 0, valign: "middle", fontFace: BODY, fontSize: 14.5, color: "1D2733" });
  });
  s.addText("The token endpoint does not redirect — it verifies the exchange belongs to the same request.",
    { x: MX, y: 5.6, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, italic: true, color: GREY });
}

// ============================================================ SLIDE 15 — PKCE MENTAL MODEL
{
  const s = content("PKCE mental model", "Challenge first, verifier later");
  card(s, MX, 1.85, 5.85, 2.7, "Before redirect (client)", [
    "Create code_verifier (random secret)",
    "Send code_challenge = hash of verifier",
  ], { bullet: true, fontSize: 14 });
  card(s, MX + 6.2, 1.85, 5.85, 2.7, "At /token (client → server)", [
    "Client sends the code_verifier",
    "Server hashes it and compares",
  ], { bullet: true, fontSize: 14 });
  codePanel(s, MX, 4.85, 11.9, 1.1, [
    "code_challenge = BASE64URL( SHA256(code_verifier) )      // pkceChallenge() in client.go",
  ]);
  s.addText("SHA-256 is one-way: the server repeats the hash and compares — it never reverses it.",
    { x: MX, y: 6.15, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}

// ============================================================ SLIDE 16 — ACCESS TOKEN AT API
{
  const s = content("Access token at the API", "The resource server validates by introspection");
  codePanel(s, MX, 1.7, 7.4, 4.6, [
    "func apiProfile(w, r) {",
    "  info, ok := introspect(r, \"profile.read\")",
    "  if !ok {",
    "    writeJSON(w, 401, ...)   // invalid token or scope",
    "    return",
    "  }",
    "  // return profile JSON",
    "}",
    "",
    "// introspect(): POST /introspect (Basic auth),",
    "// require active == true AND the needed scope",
  ], "api.go");
  card(s, 8.35, 1.7, 4.25, 4.6, "What to point at", [
    "The API (:8081) does not share memory with the AS (:8080).",
    "It validates the bearer token by calling /introspect — the real decoupled pattern.",
    "A token for the wrong API or missing scope is rejected.",
  ], { bullet: true, fontSize: 13 });
}

// ============================================================ SLIDE 17 — REFRESH
{
  const s = content("Refresh token flow", "Renew access without re-authorizing");
  card(s, MX, 1.85, 5.85, 3.4, "How the demo does it", [
    "client.go Refresh() posts grant_type=refresh_token",
    "Rotation is on: a new refresh token replaces the old one",
    "Refresh tokens are issued only when offline_access scope is granted",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.85, 5.85, 3.4, "Why it matters", [
    "Access tokens should be short-lived",
    "Refresh tokens are more sensitive — they mint new access tokens",
    "Rotation makes refresh-token reuse detectable",
  ], { bullet: true, fontSize: 13, fill: "E9F5F2", border: "A9D8CF", titleColor: ACCENT });
}

// ============================================================ SLIDE 18 — CLIENT CREDENTIALS
{
  const s = content("Client Credentials flow", "Service-to-service — no user, no browser");
  card(s, MX, 1.85, 5.85, 3.4, "How the demo does it", [
    "order-service authenticates with Basic auth (id : secret)",
    "grant_type=client_credentials, scope fraud.check",
    "Then calls /api/fraud-report with the token",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.85, 5.85, 3.4, "Key idea", [
    "No browser redirect, no consent screen",
    "The token represents the service itself, not a user",
    "A confidential client can safely hold a secret",
  ], { bullet: true, fontSize: 13, fill: "E9F5F2", border: "A9D8CF", titleColor: ACCENT });
}

// ============================================================ SLIDE 19 — DEVICE
{
  const s = content("Device Authorization flow", "Approve on a second device while the CLI polls");
  card(s, MX, 1.85, 5.85, 3.4, "How the demo does it", [
    "DeviceStart() → POST /device_authorization (deploy-cli, deploy.write)",
    "User opens the verification page and approves",
    "DevicePoll() polls /token with the device_code",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.85, 5.85, 3.4, "What to point at", [
    "Good for CLI, TV, and limited-input devices",
    "Before approval: authorization_pending",
    "After approval: token response, then calls /api/deploy",
  ], { bullet: true, fontSize: 13, fill: "E9F5F2", border: "A9D8CF", titleColor: ACCENT });
}

// ============================================================ SLIDE 20 — CHOOSE THE FLOW
{
  const s = content("Choose the flow", "Start from the application shape");
  const rows = [
    ["Web / SPA / mobile with a user", "Authorization Code + PKCE"],
    ["Backend service, no user", "Client Credentials"],
    ["Renew an access token", "Refresh Token"],
    ["CLI / TV / limited input", "Device Authorization"],
    ["Need the user's identity / login", "OpenID Connect"],
  ];
  rows.forEach((r, i) => {
    const yy = 1.8 + i * 0.92;
    s.addText(r[0], { x: MX, y: yy, w: 6.4, h: 0.7, margin: 0, valign: "middle", fontFace: BODY, fontSize: 15, color: "1D2733" });
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX + 6.7, y: yy, w: 5.4, h: 0.62, rectRadius: 0.31, fill: { color: LIGHT }, line: { color: "A9D8CF", width: 1 } });
    s.addText(r[1], { x: MX + 6.7, y: yy, w: 5.4, h: 0.62, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: INK });
  });
}

// ============================================================ SLIDE 21 — OAUTH vs OIDC
{
  const s = content("OAuth vs OpenID Connect", "Two questions, two tokens");
  card(s, MX, 1.8, 5.85, 3.0, "OAuth 2.0", [
    "Delegated API access",
    "Access token",
    "Scopes; the resource API validates",
  ], { bullet: true, fontSize: 14 });
  card(s, MX + 6.2, 1.8, 5.85, 3.0, "OpenID Connect", [
    "Login identity",
    "ID token",
    "User claims; the client validates identity",
  ], { bullet: true, fontSize: 14, fill: "E9F5F2", border: "A9D8CF", titleColor: ACCENT });
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: 5.05, w: 11.9, h: 1.15, rectRadius: 0.06, fill: { color: "FFF6E9" }, line: { color: "E8C887", width: 1 } });
  s.addText([
    { text: "Heads-up for Session 2:  ", options: { bold: true, color: "9A6B00" } },
    { text: "this demo requests the ", options: {} },
    { text: "openid", options: { fontFace: MONO, color: "9A6B00" } },
    { text: " scope, so go-oidc also issues an ", options: {} },
    { text: "id_token", options: { fontFace: MONO, color: "9A6B00" } },
    { text: " — that is OpenID Connect, covered in Session 2.", options: {} },
  ], { x: MX + 0.25, y: 5.2, w: 11.4, h: 0.85, margin: 0, valign: "middle", fontFace: BODY, fontSize: 13.5, color: "5A4A1A" });
}

// ============================================================ SLIDE 22 — COMMON MISTAKES
{
  const s = content("Common Session 1 mistakes", "Most confusion mixes tokens, actors, and endpoints");
  const rows = [
    ["Calling OAuth “login”", "Use OIDC language when identity / login is the goal."],
    ["Access token as identity", "Access token is for APIs; ID token carries login claims."],
    ["Thinking /token redirects", "/token returns JSON data to the client; it never redirects."],
    ["Confusing state and PKCE", "state protects the callback; PKCE protects the code exchange."],
  ];
  rows.forEach((r, i) => {
    const yy = 1.8 + i * 1.05;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: yy, w: 4.4, h: 0.85, rectRadius: 0.06, fill: { color: "FBECEA" }, line: { color: "E6B6AC", width: 1 } });
    s.addText(r[0], { x: MX + 0.18, y: yy, w: 4.1, h: 0.85, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: BAD });
    s.addText(r[1], { x: MX + 4.7, y: yy, w: 7.4, h: 0.85, margin: 0, valign: "middle", fontFace: BODY, fontSize: 14, color: "1D2733" });
  });
}

// ============================================================ SLIDE 23 — INSTRUCTOR CHECKLIST
{
  const s = content("Instructor demo checklist", "Run it in this exact order");
  const steps = [
    "go run .   (starts :8080, :8081, :8082)",
    "Open http://localhost:8082",
    "“Start Auth Code + PKCE flow” → Approve for Alya Developer",
    "“Call resource API with current access token”",
    "“Use refresh token”",
    "“Run Client Credentials demo”",
    "“Start Device Authorization demo” → open verification page → approve → poll",
    "“Reset All Data” to run a clean pass again",
  ];
  steps.forEach((t, i) => {
    const col = i < 4 ? 0 : 1;
    const row = i % 4;
    const x = MX + col * 6.1;
    const yy = 1.95 + row * 1.05;
    chip(s, x, yy, i + 1);
    const mono = t.startsWith("go run") || t.startsWith("Open");
    s.addText(t, { x: x + 0.58, y: yy - 0.1, w: 5.3, h: 0.85, margin: 0, valign: "middle", fontFace: mono ? MONO : BODY, fontSize: 12.5, color: "1D2733" });
  });
}

// ============================================================ SLIDE 24 — CHECKPOINT (dark)
{
  const s = dark("Session 1 checkpoint", "Explain this chain from memory");
  const chain = ["browser redirect", "authorization code", "token endpoint", "access token", "resource API"];
  let x = MX;
  chain.forEach((c, i) => {
    const w = 2.15;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y: 4.2, w, h: 0.85, rectRadius: 0.1, fill: { color: "16304E" }, line: { color: "32517A", width: 1 } });
    s.addText(c, { x, y: 4.2, w, h: 0.85, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 12.5, color: "BFE9E1" });
    x += w;
    if (i < chain.length - 1) { s.addText("→", { x: x - 0.05, y: 4.2, w: 0.4, h: 0.85, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 20, color: "7FD9CD" }); x += 0.35; }
  });
  s.addText("Final task: say who sends each request, where the token appears, and which party validates it.",
    { x: MX, y: 5.5, w: 11.5, h: 0.6, margin: 0, fontFace: BODY, fontSize: 16, color: INKSOFT });
}

pres.writeFile({ fileName: "oauth_session1_go_demo_workshop.pptx" })
  .then(f => console.log("WROTE", f, "slides:", N))
  .catch(e => { console.error("ERR", e); process.exit(1); });
