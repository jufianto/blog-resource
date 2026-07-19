// Regenerates the 1-page Session 1 outline (A4 portrait).
//   npm install pptxgenjs && node build_session1_outline.js

const P = require("pptxgenjs");
const pres = new P();
pres.defineLayout({ name: "A4P", width: 8.27, height: 11.69 });
pres.layout = "A4P";

const W = 8.27, H = 11.69;
const INK = "1F3A5F", ACCENT = "0F8A7E", GREY = "5A6B7B", WHITE = "FFFFFF";
const LIGHT = "EEF3F8", CODEBORDER = "D8DEEA", INKSOFT = "DCE7F2", NAVYCARD = "16304E", NAVYLINE = "32517A";
const TEALSOFT = "E9F5F2", TEALBORDER = "A9D8CF";
const HEAD = "Arial", BODY = "Arial";
const MX = 0.62;

const s = pres.addSlide();
s.background = { color: WHITE };

// header
s.addShape(pres.shapes.RECTANGLE, { x: 0, y: 0, w: W, h: 2.95, fill: { color: INK } });
s.addText("OAUTH 2.0 & JWT WORKSHOP  ·  SESSION 1 OF 2", { x: MX, y: 0.55, w: 7.1, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 12, bold: true, color: "7FD9CD", charSpacing: 2 });
s.addText("Session 1 — the foundations", { x: MX, y: 0.92, w: 7.1, h: 0.8, margin: 0, fontFace: HEAD, fontSize: 32, bold: true, color: WHITE });
s.addText("Lessons 1 to 4: the mental model, the flows, and how tokens are validated.", { x: MX, y: 1.78, w: 7.0, h: 0.5, margin: 0, fontFace: BODY, fontSize: 13.5, color: INKSOFT });
const tags = ["Actors & flows", "Auth Code + PKCE", "Choosing a flow", "JWT validation"];
let tx = MX;
tags.forEach((t) => { const tw = 0.4 + t.length * 0.092; s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: tx, y: 2.38, w: tw, h: 0.34, rectRadius: 0.17, fill: { color: NAVYCARD }, line: { color: NAVYLINE, width: 1 } }); s.addText(t, { x: tx, y: 2.38, w: tw, h: 0.34, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 10, color: "BFE9E1" }); tx += tw + 0.14; });

function heading(y, text, color) { s.addText(text.toUpperCase(), { x: MX, y, w: 7.1, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 12, bold: true, color: color || ACCENT, charSpacing: 2 }); }

// session 1 lessons
heading(3.2, "What Session 1 covers");
const lessons = [
  ["1", "OAuth Problem, Actors & Endpoints", "Why OAuth beats password sharing · the four actors · /authorize vs /token (front- vs back-channel) · tokens, scopes, consent · OAuth vs OpenID Connect"],
  ["2", "Authorization Code + PKCE", "The main user flow step by step · request parameters · why PKCE protects the code exchange · state vs redirect_uri vs PKCE"],
  ["3", "Client Credentials, Refresh & Device", "Machine-to-machine (no user) · refresh tokens & rotation · device flow for CLI/TV · choosing the right flow"],
  ["4", "Token Validation & JWT Basics", "JWT claims (iss, aud, exp, scope) · decoding ≠ validation · the validation checklist · 401 vs 403 · opaque tokens & introspection"],
];
lessons.forEach((l, i) => {
  const y = 3.6 + i * 1.18;
  s.addShape(pres.shapes.OVAL, { x: MX, y: y, w: 0.4, h: 0.4, fill: { color: ACCENT } });
  s.addText(l[0], { x: MX, y: y, w: 0.4, h: 0.4, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 14, bold: true, color: WHITE });
  s.addText(l[1], { x: MX + 0.55, y: y - 0.04, w: 6.9, h: 0.34, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 13.5, bold: true, color: INK });
  s.addText(l[2], { x: MX + 0.55, y: y + 0.34, w: 7.0, h: 0.7, margin: 0, valign: "top", fontFace: BODY, fontSize: 11, color: "1D2733", lineSpacingMultiple: 1.0 });
});

// after session 1 strip
const sy = 8.35;
s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: sy, w: 7.03, h: 0.78, rectRadius: 0.06, fill: { color: TEALSOFT }, line: { color: TEALBORDER, width: 1 } });
s.addText([
  { text: "After Session 1 you can:  ", options: { bold: true, color: ACCENT } },
  { text: "explain OAuth without JWT, draw Auth Code + PKCE from memory, pick the right flow, and validate a JWT for the right reason.", options: { color: "1D2733" } },
], { x: MX + 0.2, y: sy, w: 6.63, h: 0.78, margin: 0, valign: "middle", fontFace: BODY, fontSize: 11.5, lineSpacingMultiple: 1.0 });

// session 2 preview
heading(9.45, "Session 2 — date to be announced");
const s2 = [
  ["Lesson 5", "JWS, JWE, JWK / JWKS — signing vs encryption, key rotation"],
  ["Lesson 6", "OAuth security controls — attacks mapped to mitigations"],
  ["Lesson 7", "Advanced: PAR, JAR, JARM, DPoP, mTLS"],
  ["Lesson 8", "FAPI 2.0 — strict OAuth for high-value APIs"],
];
s2.forEach((r, i) => {
  const y = 9.85 + i * 0.4;
  s.addText(r[0], { x: MX, y, w: 1.2, h: 0.34, margin: 0, valign: "middle", fontFace: HEAD, fontSize: 11, bold: true, color: ACCENT });
  s.addText(r[1], { x: MX + 1.2, y, w: 6.0, h: 0.34, margin: 0, valign: "middle", fontFace: BODY, fontSize: 11, color: "1D2733" });
});

// footer logistics
s.addText([
  { text: "Session 1 date: ", options: { bold: true, color: INK } }, { text: "________", options: { color: GREY } },
  { text: "      Format: ", options: { bold: true, color: INK } }, { text: "in-person / online", options: { color: GREY } },
  { text: "      Register: ", options: { bold: true, color: INK } }, { text: "________", options: { color: GREY } },
], { x: MX, y: 11.32, w: 7.03, h: 0.3, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 10.5 });

pres.writeFile({ fileName: "oauth_session1_outline.pptx" })
  .then(f => console.log("WROTE", f)).catch(e => { console.error("ERR", e); process.exit(1); });
