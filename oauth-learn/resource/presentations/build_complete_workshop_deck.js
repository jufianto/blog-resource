// Regenerates the complete-workshop deck (Lessons 1-8).
// Usage:
//   npm install pptxgenjs
//   node build_complete_workshop_deck.js   # writes oauth_complete_workshop_lessons_1_8.pptx in cwd
//   mv oauth_complete_workshop_lessons_1_8.pptx ../../export/presentations/

const P = require("pptxgenjs");
const pres = new P();
pres.layout = "LAYOUT_WIDE";            // 13.33 x 7.5
pres.author = "OAuth/JWT/FAPI Workshop";
pres.title = "OAuth 2.0, JWT/JOSE & FAPI 2.0 — Complete Workshop";

const W = 13.33, H = 7.5;
const INK = "1F3A5F", INK2 = "2F5A8F", ACCENT = "0F8A7E", GREY = "5A6B7B";
const LIGHT = "EEF3F8", CODEBG = "F5F7FB", CODEBORDER = "D8DEEA", WHITE = "FFFFFF";
const BAD = "C0492F", INKSOFT = "DCE7F2", NAVYCARD = "16304E", NAVYLINE = "32517A";
const TEALSOFT = "E9F5F2", TEALBORDER = "A9D8CF";
const HEAD = "Arial", BODY = "Arial", MONO = "Courier New";
const MX = 0.7;

let N = 0;
const shadow = () => ({ type: "outer", color: "1F3A5F", blur: 7, offset: 3, angle: 135, opacity: 0.13 });

function footer(s) {
  s.addText("OAuth 2.0  ·  JWT/JOSE  ·  FAPI 2.0  —  Complete Workshop", { x: MX, y: H - 0.42, w: 9, h: 0.3, margin: 0, fontFace: BODY, fontSize: 9, color: GREY });
  s.addText(String(N).padStart(2, "0"), { x: W - 1.3, y: H - 0.42, w: 0.6, h: 0.3, margin: 0, fontFace: BODY, fontSize: 9, color: GREY, align: "right" });
}
function content(kicker, titleText) {
  N++;
  const s = pres.addSlide();
  s.background = { color: WHITE };
  if (kicker) s.addText(kicker.toUpperCase(), { x: MX, y: 0.4, w: 12, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 12, bold: true, color: ACCENT, charSpacing: 2 });
  s.addText(titleText, { x: MX, y: 0.7, w: 12, h: 0.85, margin: 0, fontFace: HEAD, fontSize: 26, bold: true, color: INK });
  footer(s);
  return s;
}
function dark(kicker, titleText, sub) {
  N++;
  const s = pres.addSlide();
  s.background = { color: INK };
  if (kicker) s.addText(kicker.toUpperCase(), { x: MX, y: 1.7, w: 11, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 14, bold: true, color: "7FD9CD", charSpacing: 3 });
  if (titleText) s.addText(titleText, { x: MX, y: 2.1, w: 11.9, h: 1.7, margin: 0, fontFace: HEAD, fontSize: 40, bold: true, color: WHITE });
  if (sub) s.addText(sub, { x: MX, y: 3.95, w: 11.5, h: 0.7, margin: 0, fontFace: BODY, fontSize: 17, color: INKSOFT });
  return s;
}
function divider(numStr, titleText, sub) {
  N++;
  const s = pres.addSlide();
  s.background = { color: INK };
  s.addText(numStr, { x: 8.7, y: 0.2, w: 4.4, h: 4.0, margin: 0, fontFace: HEAD, fontSize: 240, bold: true, color: "26456B", align: "right" });
  s.addText("PART " + numStr, { x: MX, y: 2.55, w: 8, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 14, bold: true, color: "7FD9CD", charSpacing: 3 });
  s.addText(titleText, { x: MX, y: 2.95, w: 8.2, h: 1.6, margin: 0, fontFace: HEAD, fontSize: 38, bold: true, color: WHITE });
  if (sub) s.addText(sub, { x: MX, y: 4.55, w: 8.0, h: 0.9, margin: 0, fontFace: BODY, fontSize: 15, color: INKSOFT });
  return s;
}
function codeRuns(lines) {
  return lines.map((ln) => {
    const isC = ln.trimStart().startsWith("//");
    return { text: ln === "" ? " " : ln, options: { breakLine: true, fontFace: MONO, fontSize: 11.5, color: isC ? GREY : "1D2733", italic: isC } };
  });
}
function codePanel(s, x, y, w, h, lines, caption) {
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y, w, h, rectRadius: 0.06, fill: { color: CODEBG }, line: { color: CODEBORDER, width: 1 }, shadow: shadow() });
  if (caption) s.addText(caption, { x: x + 0.18, y: y + 0.1, w: w - 0.36, h: 0.28, margin: 0, fontFace: MONO, fontSize: 10.5, bold: true, color: ACCENT });
  s.addText(codeRuns(lines), { x: x + 0.22, y: y + (caption ? 0.42 : 0.16), w: w - 0.44, h: h - (caption ? 0.55 : 0.3), margin: 0, valign: "top", lineSpacingMultiple: 1.06 });
}
function card(s, x, y, w, h, title, lines, opt = {}) {
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y, w, h, rectRadius: 0.06, fill: { color: opt.fill || LIGHT }, line: { color: opt.border || CODEBORDER, width: 1 }, shadow: shadow() });
  let yy = y + 0.16;
  if (title) { s.addText(title, { x: x + 0.2, y: yy, w: w - 0.4, h: 0.34, margin: 0, fontFace: HEAD, fontSize: opt.titleSize || 14, bold: true, color: opt.titleColor || INK }); yy += 0.46; }
  if (lines && lines.length) {
    const runs = lines.map((t) => ({ text: t, options: { breakLine: true, bullet: opt.bullet ? { code: "2022", indent: 14 } : false, fontFace: BODY, fontSize: opt.fontSize || 12, color: opt.color || "1D2733", paraSpaceAfter: 4 } }));
    s.addText(runs, { x: x + 0.2, y: yy, w: w - 0.4, h: y + h - yy - 0.12, margin: 0, valign: "top" });
  }
}
function chip(s, x, y, n, sz) {
  sz = sz || 0.42;
  s.addShape(pres.shapes.OVAL, { x, y, w: sz, h: sz, fill: { color: ACCENT } });
  s.addText(String(n), { x, y, w: sz, h: sz, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: WHITE });
}
// styled table; rows = array of arrays of strings (first row = header)
function table(s, x, y, w, rows, colW, opt = {}) {
  const data = rows.map((r, ri) => r.map((c) => {
    if (ri === 0) return { text: c, options: { fill: { color: INK }, color: WHITE, bold: true, fontSize: opt.headSize || 12.5, align: "left", valign: "middle" } };
    return { text: c, options: { fill: { color: ri % 2 ? WHITE : LIGHT }, color: "1D2733", fontSize: opt.fontSize || 12, align: "left", valign: "middle" } };
  }));
  s.addTable(data, { x, y, w, colW, autoPage: false, border: { type: "solid", pt: 0.5, color: CODEBORDER }, fontFace: BODY, rowH: opt.rowH || 0.3, margin: [3, 5, 3, 5] });
}

// ===================================================== TITLE
{
  const s = dark("Complete workshop · Lessons 1–8", "OAuth 2.0, JWT/JOSE & FAPI 2.0", "From the core mental model to high-security FAPI — one connected story.");
  const tags = ["Flows", "Tokens", "JWT / JOSE", "Security controls", "FAPI 2.0"];
  let x = MX;
  tags.forEach((t) => {
    const tw = 0.5 + t.length * 0.13;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y: 5.0, w: tw, h: 0.42, rectRadius: 0.21, fill: { color: NAVYCARD }, line: { color: NAVYLINE, width: 1 } });
    s.addText(t, { x, y: 5.0, w: tw, h: 0.42, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 12, color: "BFE9E1" });
    x += tw + 0.2;
  });
}

// ===================================================== AGENDA
{
  const s = content("Roadmap", "Eight lessons, one mental model");
  const rows = [
    ["1", "OAuth Problem, Actors & Endpoints", "Why OAuth exists; the four actors; tokens, scopes, consent"],
    ["2", "Authorization Code + PKCE", "The main user flow; front vs back channel; PKCE"],
    ["3", "Client Credentials, Refresh & Device", "Choose the right flow for the application shape"],
    ["4", "Token Validation & JWT Basics", "Claims, decoding ≠ validation, the validation checklist"],
    ["5", "JWS, JWE, JWK / JWKS", "Signing vs encryption; publishing and rotating keys"],
    ["6", "OAuth Security Controls", "Deprecated flows; attacks mapped to mitigations"],
    ["7", "Advanced: PAR, JAR, JARM, DPoP, mTLS", "Harden requests, responses, and token possession"],
    ["8", "FAPI 2.0", "Stricter OAuth for high-value APIs"],
  ];
  rows.forEach((r, i) => {
    const yy = 1.62 + i * 0.66;
    chip(s, MX, yy, r[0], 0.4);
    s.addText(r[1], { x: MX + 0.58, y: yy - 0.08, w: 5.0, h: 0.55, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 13.5, bold: true, color: INK });
    s.addText(r[2], { x: MX + 5.7, y: yy - 0.08, w: 6.4, h: 0.55, margin: 0, valign: "middle", fontFace: BODY, fontSize: 12.5, color: "1D2733" });
  });
}

// ##################################################### LESSON 1
divider("01", "OAuth Problem, Actors & Endpoints", "OAuth lets a client get limited API access on a user's behalf — without ever receiving the password.");
{
  const s = content("Lesson 1 · what OAuth is", "A framework, not a login or a token format");
  card(s, MX, 1.7, 5.85, 3.6, "OAuth 2.0 IS", [
    "An authorization framework (RFC 6749)",
    "A way to get limited, scoped access via tokens",
    "Answering: what may this client access?",
  ], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  card(s, MX + 6.2, 1.7, 5.85, 3.6, "OAuth 2.0 is NOT", [
    "Password sharing",
    "A single endpoint or token format",
    "A complete login protocol (that's OpenID Connect)",
    "Automatically secure — implementation matters",
  ], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  s.addText("Framework = the model · RFCs = the spec · Flows = ways to use it.", { x: MX, y: 5.55, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}
{
  const s = content("Lesson 1 · the problem", "Password sharing is the wrong shape");
  card(s, MX, 1.7, 5.85, 3.3, "Give away the password", [
    "App can do everything you can",
    "No way to scope to one resource",
    "Revoking means changing the password",
    "App breach exposes the whole account",
  ], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.7, 5.85, 3.3, "Give a scoped token", [
    "Limited to exactly what's needed",
    "Password never leaves the trusted service",
    "Revoke just this app, anytime",
    "Token expires on its own",
  ], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  s.addText("The client receives permission, not the user's password.", { x: MX, y: 5.3, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 15, bold: true, color: INK });
}
{
  const s = content("Lesson 1 · the cast", "The four actors");
  const a = [
    ["Resource Owner", "Owns the data — the user (e.g. you on GitHub)"],
    ["Client", "Wants access — the app (e.g. a deployment dashboard)"],
    ["Authorization Server", "Authenticates the user, gets consent, issues tokens"],
    ["Resource Server", "Hosts the API and validates access tokens"],
  ];
  a.forEach((r, i) => {
    const col = i % 2, row = (i / 2) | 0;
    const x = MX + col * 6.1, yy = 1.8 + row * 1.7;
    card(s, x, yy, 5.85, 1.5, r[0], [r[1]], { fontSize: 13 });
  });
  s.addText("Identify them by asking: who owns the data, who wants access, who grants it, who serves the API?", { x: MX, y: 5.4, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, italic: true, color: GREY });
}
{
  const s = content("Lesson 1 · endpoints & channels", "Two endpoints, two channels");
  table(s, MX, 1.8, 11.9, [
    ["Endpoint", "Job", "Channel"],
    ["/authorize", "User logs in & approves access", "Front-channel (browser) — more exposed"],
    ["/token", "Exchange a grant for an access token", "Back-channel (server-to-server) — protected"],
  ], [2.6, 5.1, 4.2], { rowH: 0.55 });
  s.addText("Front-channel vs back-channel is the idea behind every later security control.", { x: MX, y: 3.9, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
  card(s, MX, 4.5, 11.9, 1.5, "Also: the access token", [
    "Represents permission, not the user's password · scoped, expiring, revocable · checked by the resource server on every call.",
  ], { fontSize: 13 });
}
{
  const s = content("Lesson 1 · the boundary", "OAuth vs OpenID Connect");
  table(s, MX, 1.8, 11.9, [
    ["", "OAuth 2.0", "OpenID Connect"],
    ["Purpose", "Authorization", "Authentication"],
    ["Question", "What can this client access?", "Who is the user?"],
    ["Main token", "Access token", "ID token"],
    ["Example", "Connect GitHub repos / Drive files", "“Log in with Google”"],
  ], [2.2, 4.85, 4.85], { rowH: 0.5 });
  s.addText("Wants to call an API → OAuth. Wants to know who you are → OpenID Connect.", { x: MX, y: 5.3, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14.5, bold: true, color: INK });
}

// ##################################################### LESSON 2
divider("02", "Authorization Code + PKCE", "The browser carries only a short-lived code; tokens come from the back-channel. PKCE protects the exchange.");
{
  const s = content("Lesson 2 · the flow", "Authorization Code, step by step");
  const steps = [
    "Client redirects the browser to /authorize",
    "User logs in and approves at the authorization server",
    "Server redirects back with a short-lived code + state",
    "Client posts the code to /token (back-channel)",
    "Client receives the access token (+ refresh token)",
    "Client calls the API with the access token",
  ];
  steps.forEach((t, i) => {
    const yy = 1.7 + i * 0.78;
    chip(s, MX, yy, i + 1, 0.4);
    s.addText(t, { x: MX + 0.58, y: yy - 0.06, w: 11.4, h: 0.5, margin: 0, valign: "middle", fontFace: BODY, fontSize: 15, color: "1D2733" });
  });
}
{
  const s = content("Lesson 2 · the request", "Authorization request parameters");
  table(s, MX, 1.7, 11.9, [
    ["Parameter", "Purpose"],
    ["response_type=code", "Ask for an authorization code"],
    ["client_id", "Identify the client"],
    ["redirect_uri", "Where the browser returns (must match registration)"],
    ["scope", "Permissions requested"],
    ["state", "Bind the callback to a request the client started (CSRF)"],
    ["code_challenge (+ _method=S256)", "PKCE challenge for the code exchange"],
  ], [4.3, 7.6], { rowH: 0.46 });
}
{
  const s = content("Lesson 2 · why PKCE", "A stolen code should not be enough");
  card(s, MX, 1.7, 5.85, 3.4, "Without PKCE", [
    "Attacker intercepts the authorization code",
    "Sends it to /token",
    "If the client is public, may get tokens",
  ], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.7, 5.85, 3.4, "With PKCE", [
    "Client first sends code_challenge (a hash)",
    "Later proves it with the code_verifier at /token",
    "Attacker has the code but not the verifier → fails",
  ], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  codePanel(s, MX, 5.25, 11.9, 0.9, ["code_challenge = BASE64URL( SHA256(code_verifier) )      // S256, one-way"]);
}
{
  const s = content("Lesson 2 · not interchangeable", "state vs redirect_uri vs PKCE");
  table(s, MX, 1.8, 11.9, [
    ["Check", "Checked by", "Protects"],
    ["state", "Client", "Callbacks the client did not start (CSRF)"],
    ["redirect_uri", "Authorization server", "Binds the code to the original request target"],
    ["code_verifier / challenge", "Token endpoint", "Proves the redeemer knows the original secret"],
  ], [3.4, 3.6, 4.9], { rowH: 0.55 });
  s.addText("They solve different problems — a modern browser flow uses all three.", { x: MX, y: 4.0, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}

// ##################################################### LESSON 3
divider("03", "Client Credentials, Refresh & Device", "OAuth is not one flow — different trust situations need different flows.");
{
  const s = content("Lesson 3 · pick by situation", "Three more flows, three problems");
  const a = [
    ["Client Credentials", "Machine-to-machine — no user, no browser, no consent. Token represents the service."],
    ["Refresh Token", "Renew a short-lived access token without sending the user through the flow again."],
    ["Device Authorization", "Limited-input devices (CLI, TV): user approves on a second device while the client polls."],
  ];
  a.forEach((r, i) => { card(s, MX, 1.7 + i * 1.35, 11.9, 1.2, r[0], [r[1]], { fontSize: 13 }); });
}
{
  const s = content("Lesson 3 · refresh tokens", "Rotation makes reuse detectable");
  table(s, MX, 1.7, 11.9, [
    ["Step", "What happens"],
    ["Issue", "Client gets access token AT1 + refresh token RT1"],
    ["Refresh", "AT1 expires; client sends RT1 to /token"],
    ["Rotate", "Server returns AT2 + RT2, marks RT1 used"],
    ["Reuse seen", "If RT1 appears again → treat as theft, revoke the family"],
  ], [2.4, 9.5], { rowH: 0.5 });
  s.addText("An access token opens the door once; a refresh token can keep making new keys — store it well.", { x: MX, y: 5.0, w: 12, h: 0.5, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}
{
  const s = content("Lesson 3 · cheat sheet", "Choose the flow");
  table(s, MX, 1.7, 11.9, [
    ["Flow", "User?", "Redirect?", "Main use", "Today?"],
    ["Auth Code + PKCE", "Yes", "Yes", "Web, SPA, mobile with a user", "Recommended"],
    ["Client Credentials", "No", "No", "Service-to-service APIs", "Recommended"],
    ["Refresh Token", "Sometimes", "No", "Renew an access token", "Recommended w/ controls"],
    ["Device Authorization", "Yes", "Partial", "CLI / TV / limited input", "When appropriate"],
    ["Implicit / Password", "Yes", "—", "Legacy", "Avoid"],
  ], [3.2, 1.2, 1.5, 4.0, 2.0], { rowH: 0.46, fontSize: 11.5, headSize: 11.5 });
}

// ##################################################### LESSON 4
divider("04", "Token Validation & JWT Basics", "A JWT is only a format. Security comes from validation — decoding is not validation.");
{
  const s = content("Lesson 4 · the format", "JWT: header . payload . signature");
  codePanel(s, MX, 1.7, 6.0, 3.7, [
    "{",
    "  \"iss\": \"https://auth.example.com\",",
    "  \"sub\": \"user_123\",",
    "  \"aud\": \"https://api.example.com\",",
    "  \"exp\": 1760000000,",
    "  \"scope\": \"repo.read repo.write\"",
    "}",
  ], "decoded payload");
  card(s, MX + 6.3, 1.7, 5.6, 3.7, "Claims that drive validation", [
    "iss — do we trust this issuer?",
    "aud — is this token meant for this API?",
    "exp / nbf — is it valid right now?",
    "scope — enough permission for this action?",
  ], { bullet: true, fontSize: 13 });
  s.addText("A signed JWT is only base64url-encoded — anyone holding it can read the payload.", { x: MX, y: 5.6, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, italic: true, color: GREY });
}
{
  const s = content("Lesson 4 · decoding ≠ validation", "What the resource server must check");
  card(s, MX, 1.7, 5.85, 3.7, "Validation checklist", [
    "Signature verifies with a trusted key",
    "alg is expected; kid is a trusted key",
    "iss trusted · aud is this API",
    "exp / nbf valid in time",
    "scope sufficient · correct token type",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.7, 5.85, 3.7, "Reject for the right reason", [
    "401 Unauthorized — missing, expired, malformed, invalid token",
    "403 Forbidden — valid token, but not enough scope",
    "Opaque (non-JWT) token? Call the introspection endpoint",
  ], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
}

// ##################################################### LESSON 5
divider("05", "JWS, JWE, JWK / JWKS", "JWS signs, JWE encrypts, JWK/JWKS publish the keys. Signed does not mean private.");
{
  const s = content("Lesson 5 · the JOSE family", "Five pieces that fit together");
  table(s, MX, 1.7, 11.9, [
    ["Term", "Purpose"],
    ["JWT", "Compact claims format (the container)"],
    ["JWS", "Signature — integrity / authenticity"],
    ["JWE", "Encryption — confidentiality"],
    ["JWK", "JSON representation of one key"],
    ["JWKS", "A published set of keys (selected by kid)"],
  ], [2.2, 9.7], { rowH: 0.46 });
  s.addText("JWS proves the content was not modified. JWE hides the content.", { x: MX, y: 4.9, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14.5, bold: true, color: INK });
}
{
  const s = content("Lesson 5 · keys & rotation", "Where the verifier gets the key");
  card(s, MX, 1.7, 5.85, 3.6, "How verification works", [
    "AS signs with its private key",
    "AS publishes the public key in JWKS",
    "Resource server fetches JWKS from a trusted jwks_uri",
    "Selects the key by kid, then verifies",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.7, 5.85, 3.6, "Common alg / kid mistakes", [
    "Accepting alg=none (unsigned tokens)",
    "Accepting any algorithm from the token header",
    "Fetching keys from a URL the token controls",
    "Removing an old key before old tokens expire",
  ], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  s.addText("Allowed algorithms and key sources come from configuration — never from the token.", { x: MX, y: 5.5, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, bold: true, color: INK });
}

// ##################################################### LESSON 6
divider("06", "OAuth Security Controls", "Connect common attacks to the controls that stop them.");
{
  const s = content("Lesson 6 · legacy flows", "Avoid Implicit and Password grant");
  card(s, MX, 1.7, 5.85, 3.3, "Implicit flow", [
    "Returned tokens directly in the browser redirect",
    "Front-channel token exposure (history, logs, scripts)",
    "Replaced by Auth Code + PKCE",
  ], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.7, 5.85, 3.3, "Password credentials", [
    "Client collects the user's real password",
    "Breaks the core OAuth idea",
    "Blocks MFA; a breached client exposes credentials",
  ], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  s.addText("Use Auth Code + PKCE for users, Client Credentials for services.", { x: MX, y: 5.25, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}
{
  const s = content("Lesson 6 · threat model", "Attacks mapped to mitigations");
  table(s, MX, 1.7, 11.9, [
    ["Attack / mistake", "Mitigation"],
    ["Missing state", "Generate, store, and verify state"],
    ["Weak redirect URI matching", "Exact, pre-registered redirect URIs"],
    ["Authorization code interception", "PKCE + short-lived single-use codes"],
    ["Long-lived / stolen access token", "Short token lifetime"],
    ["Refresh token theft", "Rotation + reuse detection"],
    ["Bearer token replay", "Sender-constrained tokens (DPoP / mTLS)"],
  ], [5.6, 6.3], { rowH: 0.44 });
}

// ##################################################### LESSON 7
divider("07", "Advanced: PAR, JAR, JARM, DPoP, mTLS", "Harden the request, the response, and proof of token possession.");
{
  const s = content("Lesson 7 · request & response", "PAR, JAR, JARM");
  table(s, MX, 1.7, 11.9, [
    ["Control", "What it does", "Protects"],
    ["PAR", "Push request details to a back-channel endpoint", "Request exposure in the browser"],
    ["JAR", "Sign the authorization request (a JWT)", "Request-parameter integrity"],
    ["JARM", "Sign the authorization response (a JWT)", "Response integrity, mix-up"],
  ], [1.7, 6.0, 4.2], { rowH: 0.55 });
  s.addText("PAR pushes the request back-channel · JAR signs the request · JARM signs the response.", { x: MX, y: 4.1, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, bold: true, color: INK });
}
{
  const s = content("Lesson 7 · stop token replay", "DPoP vs mTLS");
  card(s, MX, 1.7, 5.85, 3.5, "DPoP", [
    "Application-layer signed proof per request",
    "Token bound to a public key",
    "Fits public clients & APIs",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.7, 5.85, 3.5, "mTLS", [
    "TLS-layer client certificate",
    "Certificate-bound access tokens",
    "Fits confidential / high-security clients",
  ], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  s.addText("Both sender-constrain the token: a stolen bearer token alone is no longer enough.", { x: MX, y: 5.4, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}

// ##################################################### LESSON 8
divider("08", "FAPI 2.0", "A high-security profile — stricter OAuth for high-value APIs, not a replacement protocol.");
{
  const s = content("Lesson 8 · what FAPI is", "Stricter OAuth, combined");
  card(s, MX, 1.7, 5.85, 3.6, "FAPI 2.0 combines", [
    "Strict OAuth configuration",
    "PAR + sender-constrained tokens",
    "Strong client authentication",
    "Strict token validation",
    "Optional message signing (non-repudiation)",
  ], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.7, 5.85, 3.6, "Where it's used", [
    "Open banking & payments",
    "Healthcare & sensitive identity APIs",
    "Regulated / high-impact data APIs",
  ], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  s.addText("Normal OAuth gives options; FAPI chooses the strict ones for high-value APIs.", { x: MX, y: 5.5, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}
{
  const s = content("Lesson 8 · controls → attacks", "FAPI controls mapped to risk");
  table(s, MX, 1.7, 11.9, [
    ["Risk", "FAPI-oriented control"],
    ["Authorization request tampering", "PAR + signed request objects (JAR)"],
    ["Bearer token replay", "Sender-constrained tokens (DPoP / mTLS)"],
    ["Weak client authentication", "Private-key or certificate-based auth"],
    ["Token accepted by the wrong API", "Strict audience validation"],
    ["Payment / consent dispute", "Message signing + audit evidence"],
  ], [5.9, 6.0], { rowH: 0.48 });
}

// ===================================================== CLOSING
{
  const s = dark("You can teach this when…", "Ready to teach", null);
  const items = [
    "Draw Auth Code + PKCE and Client Credentials from memory",
    "Explain code vs access vs refresh vs ID token",
    "List the JWT validation checks and why decoding isn't validation",
    "Connect attacks to mitigations: PKCE, state, redirect matching, DPoP/mTLS",
    "Explain FAPI as strict OAuth for high-value APIs",
  ];
  items.forEach((t, i) => {
    const yy = 4.0 + i * 0.55;
    s.addText("✓", { x: MX, y: yy, w: 0.4, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 16, bold: true, color: "7FD9CD" });
    s.addText(t, { x: MX + 0.45, y: yy, w: 11.5, h: 0.4, margin: 0, valign: "middle", fontFace: BODY, fontSize: 14, color: INKSOFT });
  });
}

pres.writeFile({ fileName: "oauth_complete_workshop_lessons_1_8.pptx" })
  .then(f => console.log("WROTE", f, "slides:", N))
  .catch(e => { console.error("ERR", e); process.exit(1); });
