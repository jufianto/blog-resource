// Regenerates the complete-workshop deck (Lessons 1-8).
// Usage:
//   npm install pptxgenjs
//   node build_complete_workshop_deck.js   # writes oauth_complete_workshop_lessons_1_8.pptx in cwd
//   mv oauth_complete_workshop_lessons_1_8.pptx ../../export/presentations/
//   # then, on a machine with LibreOffice, convert to PDF:
//   soffice --headless --convert-to pdf oauth_complete_workshop_lessons_1_8.pptx
//
// Design note: presenter scripts in resource/presentations/oauth_deck_*.md are keyed to
// these slides by page number AND title. Keep the slide COUNT, ORDER, and TITLES stable
// when editing — trim/rephrase body text freely, but do not add/remove/reorder slides
// without re-keying the scripts.

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
  s.addText(titleText, { x: MX, y: 0.7, w: 12, h: 0.9, margin: 0, fontFace: HEAD, fontSize: 30, bold: true, color: INK });
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
    return { text: ln === "" ? " " : ln, options: { breakLine: true, fontFace: MONO, fontSize: 12.5, color: isC ? GREY : "1D2733", italic: isC } };
  });
}
function codePanel(s, x, y, w, h, lines, caption) {
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y, w, h, rectRadius: 0.06, fill: { color: CODEBG }, line: { color: CODEBORDER, width: 1 }, shadow: shadow() });
  if (caption) s.addText(caption, { x: x + 0.18, y: y + 0.1, w: w - 0.36, h: 0.28, margin: 0, fontFace: MONO, fontSize: 10.5, bold: true, color: ACCENT });
  s.addText(codeRuns(lines), { x: x + 0.22, y: y + (caption ? 0.42 : 0.16), w: w - 0.44, h: h - (caption ? 0.55 : 0.3), margin: 0, valign: "top", lineSpacingMultiple: 1.1 });
}
// A card holds a title and a FEW short bullets. Keep bullets punchy (<= ~6 words).
function card(s, x, y, w, h, title, lines, opt = {}) {
  s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x, y, w, h, rectRadius: 0.06, fill: { color: opt.fill || LIGHT }, line: { color: opt.border || CODEBORDER, width: 1 }, shadow: shadow() });
  let yy = y + 0.18;
  if (title) { s.addText(title, { x: x + 0.22, y: yy, w: w - 0.44, h: 0.38, margin: 0, fontFace: HEAD, fontSize: opt.titleSize || 16, bold: true, color: opt.titleColor || INK }); yy += 0.52; }
  if (lines && lines.length) {
    const runs = lines.map((t) => ({ text: t, options: { breakLine: true, bullet: opt.bullet ? { code: "2022", indent: 16 } : false, fontFace: BODY, fontSize: opt.fontSize || 14.5, color: opt.color || "1D2733", paraSpaceAfter: 8 } }));
    s.addText(runs, { x: x + 0.22, y: yy, w: w - 0.44, h: y + h - yy - 0.14, margin: 0, valign: "top" });
  }
}
// A bold one-line takeaway ("land it"). Kept large and dark.
function takeaway(s, y, text, size) {
  s.addText(text, { x: MX, y, w: 12, h: 0.5, margin: 0, fontFace: BODY, fontSize: size || 16, bold: true, color: INK });
}
function chip(s, x, y, n, sz) {
  sz = sz || 0.42;
  s.addShape(pres.shapes.OVAL, { x, y, w: sz, h: sz, fill: { color: ACCENT } });
  s.addText(String(n), { x, y, w: sz, h: sz, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: WHITE });
}
// styled table; rows = array of arrays of strings (first row = header)
function table(s, x, y, w, rows, colW, opt = {}) {
  const data = rows.map((r, ri) => r.map((c) => {
    if (ri === 0) return { text: c, options: { fill: { color: INK }, color: WHITE, bold: true, fontSize: opt.headSize || 13, align: "left", valign: "middle" } };
    return { text: c, options: { fill: { color: ri % 2 ? WHITE : LIGHT }, color: "1D2733", fontSize: opt.fontSize || 12.5, align: "left", valign: "middle" } };
  }));
  s.addTable(data, { x, y, w, colW, autoPage: false, border: { type: "solid", pt: 0.5, color: CODEBORDER }, fontFace: BODY, rowH: opt.rowH || 0.34, margin: [3, 5, 3, 5] });
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
    ["1", "OAuth Problem, Actors & Endpoints", "Why OAuth exists; the four actors; tokens & consent"],
    ["2", "Authorization Code + PKCE", "The main user flow; front vs back channel; PKCE"],
    ["3", "Client Credentials, Refresh & Device", "Pick the right flow for the app shape"],
    ["4", "Token Validation & JWT Basics", "Decoding ≠ validation; the checklist"],
    ["5", "JWS, JWE, JWK / JWKS", "Signing vs encryption; publishing & rotating keys"],
    ["6", "OAuth Security Controls", "Attacks mapped to mitigations"],
    ["7", "Advanced: PAR, JAR, JARM, DPoP, mTLS", "Harden request, response, token"],
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
  card(s, MX, 1.8, 5.85, 3.3, "OAuth 2.0 IS", [
    "An authorization framework (RFC 6749)",
    "Limited, scoped access via tokens",
    "The answer to: what may this client access?",
  ], { bullet: true, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  card(s, MX + 6.2, 1.8, 5.85, 3.3, "OAuth 2.0 is NOT", [
    "Password sharing",
    "One endpoint or token format",
    "A login protocol (that's OpenID Connect)",
    "Automatically secure",
  ], { bullet: true, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  takeaway(s, 5.4, "Framework = the model · RFCs = the spec · Flows = ways to use it.");
}
{
  const s = content("Lesson 1 · the problem", "Password sharing is the wrong shape");
  card(s, MX, 1.8, 5.85, 3.1, "Give away the password", [
    "Does everything you can",
    "No way to scope",
    "Revoke = change your password",
    "Breach exposes the whole account",
  ], { bullet: true, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.8, 5.85, 3.1, "Give a scoped token", [
    "Limited to what's needed",
    "Password stays with the trusted service",
    "Revoke just this app, anytime",
    "Expires on its own",
  ], { bullet: true, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  takeaway(s, 5.2, "The client receives permission, not the user's password.");
}
{
  const s = content("Lesson 1 · the cast", "The four actors");
  const a = [
    ["Resource Owner", "Owns the data — the user"],
    ["Client", "Wants access — the app"],
    ["Authorization Server", "Authenticates, gets consent, issues tokens"],
    ["Resource Server", "Hosts the API, validates tokens"],
  ];
  a.forEach((r, i) => {
    const col = i % 2, row = (i / 2) | 0;
    const x = MX + col * 6.1, yy = 1.8 + row * 1.7;
    card(s, x, yy, 5.85, 1.45, r[0], [r[1]]);
  });
  s.addText("Ask: who owns the data, who wants access, who grants it, who serves the API?", { x: MX, y: 5.35, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, italic: true, color: GREY });
}
{
  const s = content("Lesson 1 · endpoints & channels", "Two endpoints, two channels");
  table(s, MX, 1.8, 11.9, [
    ["Endpoint", "Job", "Channel"],
    ["/authorize", "User logs in & approves", "Front-channel (browser) — exposed"],
    ["/token", "Exchange a grant for a token", "Back-channel (server) — protected"],
  ], [2.6, 5.1, 4.2], { rowH: 0.6 });
  takeaway(s, 4.0, "Front-channel vs back-channel is the idea behind every later control.");
  card(s, MX, 4.6, 11.9, 1.4, "The access token", [
    "Permission, not a password · scoped · expiring · revocable · checked on every call.",
  ]);
}
{
  const s = content("Lesson 1 · the boundary", "OAuth vs OpenID Connect");
  table(s, MX, 1.8, 11.9, [
    ["", "OAuth 2.0", "OpenID Connect"],
    ["Purpose", "Authorization", "Authentication"],
    ["Question", "What can this client access?", "Who is the user?"],
    ["Main token", "Access token", "ID token"],
    ["Example", "Connect GitHub / Drive", "“Log in with Google”"],
  ], [2.2, 4.85, 4.85], { rowH: 0.55 });
  takeaway(s, 5.4, "Call an API → OAuth. Know who you are → OpenID Connect.");
}

// ##################################################### LESSON 2
divider("02", "Authorization Code + PKCE", "The browser carries only a short-lived code; tokens come from the back-channel. PKCE protects the exchange.");
{
  const s = content("Lesson 2 · the flow", "Authorization Code, step by step");
  const steps = [
    "Client redirects the browser to /authorize",
    "User logs in and approves",
    "Server redirects back with code + state",
    "Client posts the code to /token (back-channel)",
    "Client receives the access token (+ refresh)",
    "Client calls the API with the access token",
  ];
  steps.forEach((t, i) => {
    const yy = 1.7 + i * 0.78;
    chip(s, MX, yy, i + 1, 0.4);
    s.addText(t, { x: MX + 0.58, y: yy - 0.06, w: 11.4, h: 0.5, margin: 0, valign: "middle", fontFace: BODY, fontSize: 16, color: "1D2733" });
  });
}
{
  const s = content("Lesson 2 · the request", "Authorization request parameters");
  table(s, MX, 1.8, 11.9, [
    ["Parameter", "Purpose"],
    ["response_type=code", "Ask for an authorization code"],
    ["client_id", "Identify the client"],
    ["redirect_uri", "Where the browser returns (must match)"],
    ["scope", "Permissions requested"],
    ["state", "Bind the callback to this client (CSRF)"],
    ["code_challenge (S256)", "PKCE hash for the code exchange"],
  ], [4.3, 7.6], { rowH: 0.5 });
}
{
  const s = content("Lesson 2 · why PKCE", "A stolen code should not be enough");
  card(s, MX, 1.8, 5.85, 3.1, "Without PKCE", [
    "Attacker intercepts the code",
    "Sends it to /token",
    "Public client → may get tokens",
  ], { bullet: true, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.8, 5.85, 3.1, "With PKCE", [
    "Client sends code_challenge (a hash)",
    "Proves it with code_verifier at /token",
    "No verifier → attack fails",
  ], { bullet: true, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  codePanel(s, MX, 5.15, 11.9, 0.9, ["code_challenge = BASE64URL( SHA256(code_verifier) )      // S256, one-way"]);
}
{
  const s = content("Lesson 2 · not interchangeable", "state vs redirect_uri vs PKCE");
  table(s, MX, 1.8, 11.9, [
    ["Check", "Checked by", "Protects against"],
    ["state", "Client", "Callbacks the client didn't start (CSRF)"],
    ["redirect_uri", "Authorization server", "Binds the code to the original target"],
    ["code_verifier", "Token endpoint", "Proves the redeemer started the flow"],
  ], [3.4, 3.6, 4.9], { rowH: 0.6 });
  takeaway(s, 4.1, "Different problems — a modern browser flow uses all three.");
}

// ##################################################### LESSON 3
divider("03", "Client Credentials, Refresh & Device", "OAuth is not one flow — different trust situations need different flows.");
{
  const s = content("Lesson 3 · pick by situation", "Three more flows, three problems");
  const a = [
    ["Client Credentials", "Machine-to-machine. No user, no browser. The token represents the service."],
    ["Refresh Token", "Renew a short-lived access token without sending the user through the flow again."],
    ["Device Authorization", "CLI / TV: the user approves on a second device while the client polls."],
  ];
  a.forEach((r, i) => { card(s, MX, 1.8 + i * 1.35, 11.9, 1.2, r[0], [r[1]]); });
}
{
  const s = content("Lesson 3 · refresh tokens", "Rotation makes reuse detectable");
  table(s, MX, 1.8, 11.9, [
    ["Step", "What happens"],
    ["Issue", "Client gets access token AT1 + refresh token RT1"],
    ["Refresh", "AT1 expires; client sends RT1 to /token"],
    ["Rotate", "Server returns AT2 + RT2, marks RT1 used"],
    ["Reuse seen", "RT1 appears again → theft; revoke the family"],
  ], [2.4, 9.5], { rowH: 0.55 });
  takeaway(s, 5.0, "An access token opens the door once; a refresh token keeps making keys.");
}
{
  const s = content("Lesson 3 · cheat sheet", "Choose the flow");
  table(s, MX, 1.8, 11.9, [
    ["Flow", "User?", "Redirect?", "Main use", "Today?"],
    ["Auth Code + PKCE", "Yes", "Yes", "Web, SPA, mobile with a user", "Recommended"],
    ["Client Credentials", "No", "No", "Service-to-service APIs", "Recommended"],
    ["Refresh Token", "Sometimes", "No", "Renew an access token", "With controls"],
    ["Device Authorization", "Yes", "Partial", "CLI / TV / limited input", "When needed"],
    ["Implicit / Password", "Yes", "—", "Legacy", "Avoid"],
  ], [3.2, 1.2, 1.5, 4.0, 2.0], { rowH: 0.5, fontSize: 12, headSize: 12 });
}

// ##################################################### LESSON 4
divider("04", "Token Validation & JWT Basics", "A JWT is only a format. Security comes from validation — decoding is not validation.");
{
  const s = content("Lesson 4 · the format", "JWT: header . payload . signature");
  codePanel(s, MX, 1.8, 6.0, 3.5, [
    "{",
    "  \"iss\": \"https://auth.example.com\",",
    "  \"sub\": \"user_123\",",
    "  \"aud\": \"https://api.example.com\",",
    "  \"exp\": 1760000000,",
    "  \"scope\": \"repo.read repo.write\"",
    "}",
  ], "decoded payload");
  card(s, MX + 6.3, 1.8, 5.6, 3.5, "Claims that drive validation", [
    "iss — trust this issuer?",
    "aud — meant for this API?",
    "exp / nbf — valid right now?",
    "scope — enough permission?",
  ], { bullet: true });
  s.addText("A signed JWT is only base64url-encoded — anyone holding it can read the payload.", { x: MX, y: 5.5, w: 12, h: 0.4, margin: 0, fontFace: BODY, fontSize: 13.5, italic: true, color: GREY });
}
{
  const s = content("Lesson 4 · decoding ≠ validation", "What the resource server must check");
  card(s, MX, 1.8, 5.85, 3.5, "Validation checklist", [
    "Signature verifies with a trusted key",
    "alg expected; kid is trusted",
    "iss trusted · aud is this API",
    "exp / nbf valid in time",
    "scope sufficient",
  ], { bullet: true });
  card(s, MX + 6.2, 1.8, 5.85, 3.5, "Reject for the right reason", [
    "401 — missing / expired / bad token",
    "403 — valid token, not enough scope",
    "Opaque (non-JWT) token → introspect",
  ], { bullet: true, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
}

// ##################################################### LESSON 5
divider("05", "JWS, JWE, JWK / JWKS", "JWS signs, JWE encrypts, JWK/JWKS publish the keys. Signed does not mean private.");
{
  const s = content("Lesson 5 · the JOSE family", "Five pieces that fit together");
  table(s, MX, 1.8, 11.9, [
    ["Term", "Purpose"],
    ["JWT", "Compact claims format (the container)"],
    ["JWS", "Signature — integrity / authenticity"],
    ["JWE", "Encryption — confidentiality"],
    ["JWK", "JSON representation of one key"],
    ["JWKS", "A published set of keys (selected by kid)"],
  ], [2.2, 9.7], { rowH: 0.5 });
  takeaway(s, 5.0, "JWS proves the content wasn't modified. JWE hides the content.");
}
{
  const s = content("Lesson 5 · keys & rotation", "Where the verifier gets the key");
  card(s, MX, 1.8, 5.85, 3.4, "How verification works", [
    "AS signs with its private key",
    "AS publishes the public key in JWKS",
    "RS fetches JWKS from a trusted uri",
    "Selects by kid, then verifies",
  ], { bullet: true });
  card(s, MX + 6.2, 1.8, 5.85, 3.4, "Common alg / kid mistakes", [
    "Accepting alg=none",
    "Trusting the token's algorithm",
    "Keys from a token-controlled URL",
    "Dropping an old key too early",
  ], { bullet: true, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  takeaway(s, 5.4, "Allowed algorithms and key sources come from configuration — never the token.");
}

// ##################################################### LESSON 6
divider("06", "OAuth Security Controls", "Connect common attacks to the controls that stop them.");
{
  const s = content("Lesson 6 · legacy flows", "Avoid Implicit and Password grant");
  card(s, MX, 1.8, 5.85, 3.1, "Implicit flow", [
    "Tokens returned in the browser redirect",
    "Front-channel exposure (history, logs)",
    "Replaced by Auth Code + PKCE",
  ], { bullet: true, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  card(s, MX + 6.2, 1.8, 5.85, 3.1, "Password credentials", [
    "Client collects the user's real password",
    "Breaks the core OAuth idea",
    "Blocks MFA; a breach exposes creds",
  ], { bullet: true, fill: "FBECEA", border: "E6B6AC", titleColor: BAD });
  takeaway(s, 5.2, "Auth Code + PKCE for users · Client Credentials for services.");
}
{
  const s = content("Lesson 6 · threat model", "Attacks mapped to mitigations");
  table(s, MX, 1.8, 11.9, [
    ["Attack / mistake", "Mitigation"],
    ["Missing state", "Generate, store, and verify state"],
    ["Weak redirect URI matching", "Exact, pre-registered redirect URIs"],
    ["Authorization code interception", "PKCE + short-lived single-use codes"],
    ["Long-lived / stolen access token", "Short token lifetime"],
    ["Refresh token theft", "Rotation + reuse detection"],
    ["Bearer token replay", "Sender-constrained tokens (DPoP / mTLS)"],
  ], [5.6, 6.3], { rowH: 0.48 });
}

// ##################################################### LESSON 7
divider("07", "Advanced: PAR, JAR, JARM, DPoP, mTLS", "Harden the request, the response, and proof of token possession.");
{
  const s = content("Lesson 7 · request & response", "PAR, JAR, JARM");
  table(s, MX, 1.8, 11.9, [
    ["Control", "What it does", "Protects"],
    ["PAR", "Push request details to a back-channel endpoint", "Request exposure in the browser"],
    ["JAR", "Sign the authorization request (a JWT)", "Request-parameter integrity"],
    ["JARM", "Sign the authorization response (a JWT)", "Response integrity, mix-up"],
  ], [1.7, 6.0, 4.2], { rowH: 0.6 });
  takeaway(s, 4.1, "PAR pushes the request · JAR signs the request · JARM signs the response.");
}
{
  const s = content("Lesson 7 · stop token replay", "DPoP vs mTLS");
  card(s, MX, 1.8, 5.85, 3.3, "DPoP", [
    "App-layer signed proof per request",
    "Token bound to a public key",
    "Fits public clients & APIs",
  ], { bullet: true });
  card(s, MX + 6.2, 1.8, 5.85, 3.3, "mTLS", [
    "TLS-layer client certificate",
    "Certificate-bound access tokens",
    "Fits confidential / high-security clients",
  ], { bullet: true, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  takeaway(s, 5.4, "Both sender-constrain the token: a stolen bearer token alone is not enough.");
}

// ##################################################### LESSON 8
divider("08", "FAPI 2.0", "A high-security profile — stricter OAuth for high-value APIs, not a replacement protocol.");
{
  const s = content("Lesson 8 · what FAPI is", "Stricter OAuth, combined");
  card(s, MX, 1.8, 5.85, 3.4, "FAPI 2.0 combines", [
    "Strict OAuth configuration",
    "PAR + sender-constrained tokens",
    "Strong client authentication",
    "Strict token validation",
  ], { bullet: true });
  card(s, MX + 6.2, 1.8, 5.85, 3.4, "Where it's used", [
    "Open banking & payments",
    "Healthcare & sensitive identity",
    "Regulated, high-impact APIs",
  ], { bullet: true, fill: TEALSOFT, border: TEALBORDER, titleColor: ACCENT });
  takeaway(s, 5.4, "Normal OAuth gives options; FAPI chooses the strict ones.");
}
{
  const s = content("Lesson 8 · controls → attacks", "FAPI controls mapped to risk");
  table(s, MX, 1.8, 11.9, [
    ["Risk", "FAPI-oriented control"],
    ["Authorization request tampering", "PAR + signed request objects (JAR)"],
    ["Bearer token replay", "Sender-constrained tokens (DPoP / mTLS)"],
    ["Weak client authentication", "Private-key or certificate-based auth"],
    ["Token accepted by the wrong API", "Strict audience validation"],
    ["Payment / consent dispute", "Message signing + audit evidence"],
  ], [5.9, 6.0], { rowH: 0.52 });
}

// ===================================================== CLOSING
{
  const s = dark("You can teach this when…", "Ready to teach", null);
  const items = [
    "Draw Auth Code + PKCE and Client Credentials from memory",
    "Explain code vs access vs refresh vs ID token",
    "List the JWT validation checks (decoding isn't validation)",
    "Connect attacks to controls: PKCE, state, redirect, DPoP/mTLS",
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
