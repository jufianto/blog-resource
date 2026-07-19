// Regenerates the full Session 1 deck (Lessons 1-4 + Session 2 teaser).
//   npm install pptxgenjs && node build_session1_deck.js

const P = require("pptxgenjs");
const pres = new P();
pres.layout = "LAYOUT_WIDE";
pres.author = "OAuth/JWT/FAPI Workshop";
pres.title = "OAuth 2.0 & JWT — Session 1 (Lessons 1–4)";

const W = 13.33, H = 7.5;
const INK = "1F3A5F", ACCENT = "0F8A7E", GREY = "5A6B7B";
const LIGHT = "EEF3F8", CODEBG = "F5F7FB", CODEBORDER = "D8DEEA", WHITE = "FFFFFF";
const BAD = "C0492F", INKSOFT = "DCE7F2", NAVYCARD = "16304E", NAVYLINE = "32517A";
const TEALSOFT = "E9F5F2", TEALBORDER = "A9D8CF";
const HEAD = "Arial", BODY = "Arial", MONO = "Courier New";
const MX = 0.7;

let N = 0;
const shadow = () => ({ type: "outer", color: "1F3A5F", blur: 7, offset: 3, angle: 135, opacity: 0.13 });
const FOOT = "OAuth 2.0  ·  Session 1  ·  Lessons 1–4";

function footer(s) {
  s.addText(FOOT, { x: MX, y: H - 0.42, w: 9, h: 0.3, margin: 0, fontFace: BODY, fontSize: 9, color: GREY });
  s.addText(String(N).padStart(2, "0"), { x: W - 1.3, y: H - 0.42, w: 0.6, h: 0.3, margin: 0, fontFace: BODY, fontSize: 9, color: GREY, align: "right" });
}
function content(kicker, titleText) {
  N++;
  const s = pres.addSlide(); s.background = { color: WHITE };
  if (kicker) s.addText(kicker.toUpperCase(), { x: MX, y: 0.4, w: 12, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 12, bold: true, color: ACCENT, charSpacing: 2 });
  s.addText(titleText, { x: MX, y: 0.7, w: 12, h: 0.85, margin: 0, fontFace: HEAD, fontSize: 26, bold: true, color: INK });
  footer(s); return s;
}
function dark(kicker, titleText, sub) {
  N++;
  const s = pres.addSlide(); s.background = { color: INK };
  if (kicker) s.addText(kicker.toUpperCase(), { x: MX, y: 1.7, w: 11, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 14, bold: true, color: "7FD9CD", charSpacing: 3 });
  if (titleText) s.addText(titleText, { x: MX, y: 2.1, w: 11.9, h: 1.7, margin: 0, fontFace: HEAD, fontSize: 40, bold: true, color: WHITE });
  if (sub) s.addText(sub, { x: MX, y: 3.95, w: 11.5, h: 0.7, margin: 0, fontFace: BODY, fontSize: 17, color: INKSOFT });
  return s;
}
function divider(numStr, titleText, sub) {
  N++;
  const s = pres.addSlide(); s.background = { color: INK };
  s.addText(numStr, { x: 8.7, y: 0.2, w: 4.4, h: 4.0, margin: 0, fontFace: HEAD, fontSize: 240, bold: true, color: "26456B", align: "right" });
  s.addText("LESSON " + numStr, { x: MX, y: 2.55, w: 8, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 14, bold: true, color: "7FD9CD", charSpacing: 3 });
  s.addText(titleText, { x: MX, y: 2.95, w: 8.2, h: 1.6, margin: 0, fontFace: HEAD, fontSize: 38, bold: true, color: WHITE });
  if (sub) s.addText(sub, { x: MX, y: 4.55, w: 8.0, h: 0.9, margin: 0, fontFace: BODY, fontSize: 15, color: INKSOFT });
  return s;
}
function codeRuns(lines) {
  return lines.map((ln) => { const c = ln.trimStart().startsWith("//"); return { text: ln === "" ? " " : ln, options: { breakLine: true, fontFace: MONO, fontSize: 11.5, color: c ? GREY : "1D2733", italic: c } }; });
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
function chip(s, x, y, n, sz) { sz = sz || 0.42; s.addShape(pres.shapes.OVAL, { x, y, w: sz, h: sz, fill: { color: ACCENT } }); s.addText(String(n), { x, y, w: sz, h: sz, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: WHITE }); }
function table(s, x, y, w, rows, colW, opt = {}) {
  const data = rows.map((r, ri) => r.map((c) => ri === 0
    ? { text: c, options: { fill: { color: INK }, color: WHITE, bold: true, fontSize: opt.headSize || 12.5, align: "left", valign: "middle" } }
    : { text: c, options: { fill: { color: ri % 2 ? WHITE : LIGHT }, color: "1D2733", fontSize: opt.fontSize || 12, align: "left", valign: "middle" } }));
  s.addTable(data, { x, y, w, colW, autoPage: false, border: { type: "solid", pt: 0.5, color: CODEBORDER }, fontFace: BODY, rowH: opt.rowH || 0.3, margin: [3, 5, 3, 5] });
}

// ===================================================== TITLE
{
  const s = dark("Workshop · Session 1 of 2", "OAuth 2.0 & JWT — the foundations", "Lessons 1–4: the mental model, the flows, and how tokens are validated.");
  const tags = ["Actors & flows", "Auth Code + PKCE", "Choosing a flow", "JWT validation"];
  let x = MX;
  tags.forEach((t) => { const tw = 0.5 + t.length * 0.125; s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y: 5.0, w: tw, h: 0.42, rectRadius: 0.21, fill: { color: NAVYCARD }, line: { color: NAVYLINE, width: 1 } }); s.addText(t, { x, y: 5.0, w: tw, h: 0.42, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 12, color: "BFE9E1" }); x += tw + 0.2; });
}

// ===================================================== AGENDA
{
  const s = content("Session 1 roadmap", "What we cover today — Lessons 1 to 4");
  const rows = [
    ["1", "OAuth Problem, Actors & Endpoints", "Why OAuth exists; the four actors; tokens, scopes, consent"],
    ["2", "Authorization Code + PKCE", "The main user flow; front vs back channel; PKCE"],
    ["3", "Client Credentials, Refresh & Device", "Choose the right flow for the application shape"],
    ["4", "Token Validation & JWT Basics", "Claims; decoding ≠ validation; the validation checklist"],
  ];
  rows.forEach((r, i) => { const yy = 1.7 + i * 0.78; chip(s, MX, yy, r[0], 0.42); s.addText(r[1], { x: MX + 0.6, y: yy - 0.06, w: 5.1, h: 0.55, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: INK }); s.addText(r[2], { x: MX + 5.9, y: yy - 0.06, w: 6.2, h: 0.55, margin: 0, valign: "middle", fontFace: BODY, fontSize: 12.5, color: "1D2733" }); });
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: 5.05, w: 11.9, h: 1.0, rectRadius: 0.06, fill: { color: TEALSOFT }, line: { color: TEALBORDER, width: 1 } });
  s.addText([
    { text: "Session 2 (date to be announced):  ", options: { bold: true, color: ACCENT } },
    { text: "JWS/JWE/JWKS, OAuth security controls, advanced controls (PAR, JAR, JARM, DPoP, mTLS), and FAPI 2.0.", options: { color: "1D2733" } },
  ], { x: MX + 0.25, y: 5.2, w: 11.4, h: 0.7, margin: 0, valign: "middle", fontFace: BODY, fontSize: 13.5 });
}

// ##################################################### LESSON 1
divider("01", "OAuth Problem, Actors & Endpoints", "OAuth lets a client get limited API access on a user's behalf — without ever receiving the password.");
{
  const s = content("Lesson 1 · what OAuth is", "A framework, not a login or a token format");
  card(s, MX, 1.7, 5.85, 3.6, "OAuth 2.0 IS", ["An authorization framework (RFC 6749)", "A way to get limited, scoped access via tokens", "Answering: what may this client access?"], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  card(s, MX + 6.2, 1.7, 5.85, 3.6, "OAuth 2.0 is NOT", ["Password sharing", "A single endpoint or token format", "A complete login protocol (that's OpenID Connect)", "Automatically secure — implementation matters"], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  s.addText("Framework = the model · RFCs = the spec · Flows = ways to use it.", { x: MX, y: 5.55, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}
{
  const s = content("Lesson 1 · the problem", "Password sharing is the wrong shape");
  card(s, MX, 1.7, 5.85, 3.3, "Give away the password", ["App can do everything you can", "No way to scope to one resource", "Revoking means changing the password", "App breach exposes the whole account"], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.7, 5.85, 3.3, "Give a scoped token", ["Limited to exactly what's needed", "Password never leaves the trusted service", "Revoke just this app, anytime", "Token expires on its own"], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  s.addText("The client receives permission, not the user's password.", { x: MX, y: 5.3, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 15, bold: true, color: INK });
}
{
  const s = content("Lesson 1 · the cast", "The four actors");
  const a = [["Resource Owner", "Owns the data — the user (e.g. you on GitHub)"], ["Client", "Wants access — the app (e.g. a deployment dashboard)"], ["Authorization Server", "Authenticates the user, gets consent, issues tokens"], ["Resource Server", "Hosts the API and validates access tokens"]];
  a.forEach((r, i) => { const col = i % 2, row = (i / 2) | 0; card(s, MX + col * 6.1, 1.8 + row * 1.7, 5.85, 1.5, r[0], [r[1]], { fontSize: 13 }); });
  s.addText("Ask: who owns the data, who wants access, who grants it, who serves the API?", { x: MX, y: 5.4, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, italic: true, color: GREY });
}
{
  const s = content("Lesson 1 · endpoints & channels", "Two endpoints, two channels");
  table(s, MX, 1.8, 11.9, [["Endpoint", "Job", "Channel"], ["/authorize", "User logs in & approves access", "Front-channel (browser) — more exposed"], ["/token", "Exchange a grant for an access token", "Back-channel (server-to-server) — protected"]], [2.6, 5.1, 4.2], { rowH: 0.55 });
  s.addText("Front-channel vs back-channel is the idea behind every later security control.", { x: MX, y: 3.9, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
  card(s, MX, 4.5, 11.9, 1.5, "Also: the access token", ["Represents permission, not the user's password · scoped, expiring, revocable · checked by the resource server on every call."], { fontSize: 13 });
}
{
  const s = content("Lesson 1 · the boundary", "OAuth vs OpenID Connect");
  table(s, MX, 1.8, 11.9, [["", "OAuth 2.0", "OpenID Connect"], ["Purpose", "Authorization", "Authentication"], ["Question", "What can this client access?", "Who is the user?"], ["Main token", "Access token", "ID token"], ["Example", "Connect GitHub repos / Drive files", "“Log in with Google”"]], [2.2, 4.85, 4.85], { rowH: 0.5 });
  s.addText("Wants to call an API → OAuth. Wants to know who you are → OpenID Connect.", { x: MX, y: 5.3, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14.5, bold: true, color: INK });
}

// ##################################################### LESSON 2
divider("02", "Authorization Code + PKCE", "The browser carries only a short-lived code; tokens come from the back-channel. PKCE protects the exchange.");
{
  const s = content("Lesson 2 · the flow", "Authorization Code, step by step");
  ["Client redirects the browser to /authorize", "User logs in and approves at the authorization server", "Server redirects back with a short-lived code + state", "Client posts the code to /token (back-channel)", "Client receives the access token (+ refresh token)", "Client calls the API with the access token"].forEach((t, i) => { const yy = 1.7 + i * 0.78; chip(s, MX, yy, i + 1, 0.4); s.addText(t, { x: MX + 0.58, y: yy - 0.06, w: 11.4, h: 0.5, margin: 0, valign: "middle", fontFace: BODY, fontSize: 15, color: "1D2733" }); });
}
{
  const s = content("Lesson 2 · the request", "Authorization request parameters");
  table(s, MX, 1.7, 11.9, [["Parameter", "Purpose"], ["response_type=code", "Ask for an authorization code"], ["client_id", "Identify the client"], ["redirect_uri", "Where the browser returns (must match registration)"], ["scope", "Permissions requested"], ["state", "Bind the callback to a request the client started (CSRF)"], ["code_challenge (+ _method=S256)", "PKCE challenge for the code exchange"]], [4.3, 7.6], { rowH: 0.46 });
}
{
  const s = content("Lesson 2 · why PKCE", "A stolen code should not be enough");
  card(s, MX, 1.7, 5.85, 3.4, "Without PKCE", ["Attacker intercepts the authorization code", "Sends it to /token", "If the client is public, may get tokens"], { bullet: true, fontSize: 13, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.7, 5.85, 3.4, "With PKCE", ["Client first sends code_challenge (a hash)", "Later proves it with the code_verifier at /token", "Attacker has the code but not the verifier → fails"], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  codePanel(s, MX, 5.25, 11.9, 0.9, ["code_challenge = BASE64URL( SHA256(code_verifier) )      // S256, one-way"]);
}
{
  const s = content("Lesson 2 · not interchangeable", "state vs redirect_uri vs PKCE");
  table(s, MX, 1.8, 11.9, [["Check", "Checked by", "Protects"], ["state", "Client", "Callbacks the client did not start (CSRF)"], ["redirect_uri", "Authorization server", "Binds the code to the original request target"], ["code_verifier / challenge", "Token endpoint", "Proves the redeemer knows the original secret"]], [3.4, 3.6, 4.9], { rowH: 0.55 });
  s.addText("They solve different problems — a modern browser flow uses all three.", { x: MX, y: 4.0, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}

// ##################################################### LESSON 3
divider("03", "Client Credentials, Refresh & Device", "OAuth is not one flow — different trust situations need different flows.");
{
  const s = content("Lesson 3 · pick by situation", "Three more flows, three problems");
  [["Client Credentials", "Machine-to-machine — no user, no browser, no consent. Token represents the service."], ["Refresh Token", "Renew a short-lived access token without sending the user through the flow again."], ["Device Authorization", "Limited-input devices (CLI, TV): user approves on a second device while the client polls."]].forEach((r, i) => card(s, MX, 1.7 + i * 1.35, 11.9, 1.2, r[0], [r[1]], { fontSize: 13 }));
}
{
  const s = content("Lesson 3 · refresh tokens", "Rotation makes reuse detectable");
  table(s, MX, 1.7, 11.9, [["Step", "What happens"], ["Issue", "Client gets access token AT1 + refresh token RT1"], ["Refresh", "AT1 expires; client sends RT1 to /token"], ["Rotate", "Server returns AT2 + RT2, marks RT1 used"], ["Reuse seen", "If RT1 appears again → treat as theft, revoke the family"]], [2.4, 9.5], { rowH: 0.5 });
  s.addText("An access token opens the door once; a refresh token can keep making new keys — store it well.", { x: MX, y: 5.0, w: 12, h: 0.5, margin: 0, fontFace: BODY, fontSize: 14, bold: true, color: INK });
}
{
  const s = content("Lesson 3 · cheat sheet", "Choose the flow");
  table(s, MX, 1.7, 11.9, [["Flow", "User?", "Redirect?", "Main use", "Today?"], ["Auth Code + PKCE", "Yes", "Yes", "Web, SPA, mobile with a user", "Recommended"], ["Client Credentials", "No", "No", "Service-to-service APIs", "Recommended"], ["Refresh Token", "Sometimes", "No", "Renew an access token", "Recommended w/ controls"], ["Device Authorization", "Yes", "Partial", "CLI / TV / limited input", "When appropriate"], ["Implicit / Password", "Yes", "—", "Legacy", "Avoid"]], [3.2, 1.2, 1.5, 4.0, 2.0], { rowH: 0.46, fontSize: 11.5, headSize: 11.5 });
}

// ##################################################### LESSON 4
divider("04", "Token Validation & JWT Basics", "A JWT is only a format. Security comes from validation — decoding is not validation.");
{
  const s = content("Lesson 4 · the format", "JWT: header . payload . signature");
  codePanel(s, MX, 1.7, 6.0, 3.7, ["{", "  \"iss\": \"https://auth.example.com\",", "  \"sub\": \"user_123\",", "  \"aud\": \"https://api.example.com\",", "  \"exp\": 1760000000,", "  \"scope\": \"repo.read repo.write\"", "}"], "decoded payload");
  card(s, MX + 6.3, 1.7, 5.6, 3.7, "Claims that drive validation", ["iss — do we trust this issuer?", "aud — is this token meant for this API?", "exp / nbf — is it valid right now?", "scope — enough permission for this action?"], { bullet: true, fontSize: 13 });
  s.addText("A signed JWT is only base64url-encoded — anyone holding it can read the payload.", { x: MX, y: 5.6, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, italic: true, color: GREY });
}
{
  const s = content("Lesson 4 · decoding ≠ validation", "What the resource server must check");
  card(s, MX, 1.7, 5.85, 3.7, "Validation checklist", ["Signature verifies with a trusted key", "alg is expected; kid is a trusted key", "iss trusted · aud is this API", "exp / nbf valid in time", "scope sufficient · correct token type"], { bullet: true, fontSize: 13 });
  card(s, MX + 6.2, 1.7, 5.85, 3.7, "Reject for the right reason", ["401 Unauthorized — missing, expired, malformed, invalid token", "403 Forbidden — valid token, but not enough scope", "Opaque (non-JWT) token? Call the introspection endpoint"], { bullet: true, fontSize: 13, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
}

// ===================================================== SESSION 2 TEASER
{
  const s = dark("Coming in Session 2 · date to be announced", "Where Session 2 goes next", null);
  const items = [
    ["Lesson 5", "JWS, JWE, JWK / JWKS — signing vs encryption, publishing & rotating keys"],
    ["Lesson 6", "OAuth security controls — deprecated flows; attacks mapped to mitigations"],
    ["Lesson 7", "Advanced: PAR, JAR, JARM, DPoP, mTLS — hardening requests & token possession"],
    ["Lesson 8", "FAPI 2.0 — strict OAuth for high-value APIs (open banking, payments, health)"],
  ];
  items.forEach((r, i) => {
    const yy = 4.0 + i * 0.62;
    s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: yy, w: 1.9, h: 0.46, rectRadius: 0.06, fill: { color: NAVYCARD }, line: { color: NAVYLINE, width: 1 } });
    s.addText(r[0], { x: MX, y: yy, w: 1.9, h: 0.46, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 13, bold: true, color: "7FD9CD" });
    s.addText(r[1], { x: MX + 2.15, y: yy, w: 10.0, h: 0.46, margin: 0, valign: "middle", fontFace: BODY, fontSize: 14, color: INKSOFT });
  });
}

// ===================================================== CLOSING
{
  const s = dark("Session 1 · you're ready when…", "Ready for Session 2", null);
  ["Explain OAuth in two minutes without mentioning JWT", "Draw Authorization Code + PKCE from memory", "Choose the right flow for a given application shape", "List the JWT validation checks and why decoding isn't validation"].forEach((t, i) => {
    const yy = 4.1 + i * 0.55;
    s.addText("✓", { x: MX, y: yy, w: 0.4, h: 0.4, margin: 0, fontFace: HEAD, fontSize: 16, bold: true, color: "7FD9CD" });
    s.addText(t, { x: MX + 0.45, y: yy, w: 11.5, h: 0.4, margin: 0, valign: "middle", fontFace: BODY, fontSize: 14, color: INKSOFT });
  });
}

pres.writeFile({ fileName: "oauth_session1_workshop_lessons_1_4.pptx" })
  .then(f => console.log("WROTE", f, "slides:", N)).catch(e => { console.error("ERR", e); process.exit(1); });
