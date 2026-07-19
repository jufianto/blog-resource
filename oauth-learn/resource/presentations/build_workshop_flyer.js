// Regenerates the 1-page workshop flyer (A4 portrait, Session 1 / Session 2 framing).
//   npm install pptxgenjs && node build_workshop_flyer.js

const P = require("pptxgenjs");
const pres = new P();
pres.defineLayout({ name: "A4P", width: 8.27, height: 11.69 });
pres.layout = "A4P";
pres.author = "OAuth/JWT/FAPI Workshop";

const W = 8.27, H = 11.69;
const INK = "1F3A5F", ACCENT = "0F8A7E", GREY = "5A6B7B", WHITE = "FFFFFF";
const LIGHT = "EEF3F8", CODEBORDER = "D8DEEA", INKSOFT = "DCE7F2", NAVYCARD = "16304E", NAVYLINE = "32517A";
const TEALSOFT = "E9F5F2", TEALBORDER = "A9D8CF";
const HEAD = "Arial", BODY = "Arial";
const MX = 0.62;

const s = pres.addSlide();
s.background = { color: WHITE };

// ---------- header ----------
s.addShape(pres.shapes.RECTANGLE, { x: 0, y: 0, w: W, h: 3.95, fill: { color: INK } });
s.addText("HANDS-ON WORKSHOP  ·  TWO SESSIONS", { x: MX, y: 0.6, w: 7, h: 0.32, margin: 0, fontFace: HEAD, fontSize: 13, bold: true, color: "7FD9CD", charSpacing: 3 });
s.addText("OAuth 2.0, JWT/JOSE\n& FAPI 2.0", { x: MX, y: 0.98, w: 7.1, h: 1.5, margin: 0, fontFace: HEAD, fontSize: 38, bold: true, color: WHITE, lineSpacingMultiple: 0.98 });
s.addText("Build it, break it, and learn to teach it — not just recognise the diagrams.", { x: MX, y: 2.6, w: 7.0, h: 0.5, margin: 0, fontFace: BODY, fontSize: 14.5, color: INKSOFT });
const tags = ["Flows", "Tokens", "JWT / JOSE", "Security", "FAPI 2.0"];
let tx = MX;
tags.forEach((t) => { const tw = 0.42 + t.length * 0.105; s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: tx, y: 3.28, w: tw, h: 0.36, rectRadius: 0.18, fill: { color: NAVYCARD }, line: { color: NAVYLINE, width: 1 } }); s.addText(t, { x: tx, y: 3.28, w: tw, h: 0.36, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 10.5, color: "BFE9E1" }); tx += tw + 0.16; });

function heading(y, text, color) { s.addText(text.toUpperCase(), { x: MX, y, w: 7.1, h: 0.3, margin: 0, fontFace: HEAD, fontSize: 12.5, bold: true, color: color || ACCENT, charSpacing: 2 }); }

// ---------- outcomes ----------
heading(4.26, "What you'll walk away able to do");
[
  "Choose the right OAuth flow and draw Authorization Code + PKCE from memory",
  "Validate JWTs properly — and explain why decoding is not validation",
  "Map real-world attacks to the controls that stop them",
  "Explain FAPI 2.0 as strict OAuth for high-value APIs",
].forEach((t, i) => { const y = 4.64 + i * 0.44; s.addText("✓", { x: MX, y, w: 0.3, h: 0.32, margin: 0, fontFace: HEAD, fontSize: 14, bold: true, color: ACCENT }); s.addText(t, { x: MX + 0.36, y, w: 6.9, h: 0.32, margin: 0, valign: "middle", fontFace: BODY, fontSize: 12, color: "1D2733" }); });

// ---------- what we cover, split across two sessions ----------
heading(6.55, "What we cover — across two sessions");
// column labels
s.addText("SESSION 1  ·  AVAILABLE NOW", { x: MX, y: 6.92, w: 3.6, h: 0.28, margin: 0, fontFace: HEAD, fontSize: 10.5, bold: true, color: ACCENT, charSpacing: 1 });
s.addText("SESSION 2  ·  DATE TBA", { x: MX + 3.72, y: 6.92, w: 3.6, h: 0.28, margin: 0, fontFace: HEAD, fontSize: 10.5, bold: true, color: GREY, charSpacing: 1 });
const topics = [
  ["OAuth problem, actors & endpoints", 0], ["Authorization Code + PKCE", 0], ["Client Credentials, Refresh & Device", 0], ["Token validation & JWT basics", 0],
  ["JWS, JWE, JWK / JWKS", 1], ["OAuth security controls", 1], ["Advanced: PAR, JAR, JARM, DPoP, mTLS", 1], ["FAPI 2.0", 1],
];
topics.forEach((t, i) => {
  const col = t[1], row = i % 4;
  const x = MX + col * 3.72, y = 7.3 + row * 0.46;
  const n = i + 1;
  s.addShape(pres.shapes.OVAL, { x, y, w: 0.34, h: 0.34, fill: { color: col === 0 ? ACCENT : GREY } });
  s.addText(String(n), { x, y, w: 0.34, h: 0.34, margin: 0, align: "center", valign: "middle", fontFace: HEAD, fontSize: 12, bold: true, color: WHITE });
  s.addText(t[0], { x: x + 0.46, y: y - 0.03, w: 3.05, h: 0.4, margin: 0, valign: "middle", fontFace: BODY, fontSize: 11, color: col === 0 ? "1D2733" : GREY });
});

// ---------- how it works + who it's for ----------
const cardY = 9.35, cardH = 1.58, cw = 3.46;
s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: MX, y: cardY, w: cw, h: cardH, rectRadius: 0.06, fill: { color: TEALSOFT }, line: { color: TEALBORDER, width: 1 } });
s.addText("HOW IT WORKS", { x: MX + 0.2, y: cardY + 0.16, w: cw - 0.4, h: 0.28, margin: 0, fontFace: HEAD, fontSize: 11.5, bold: true, color: ACCENT, charSpacing: 1 });
s.addText([
  { text: "Two sessions, hands-on", options: { breakLine: true, bold: true } },
  { text: "Each topic: concept → diagram → run a real Go demo → explain it back.", options: {} },
], { x: MX + 0.2, y: cardY + 0.5, w: cw - 0.4, h: cardH - 0.6, margin: 0, valign: "top", fontFace: BODY, fontSize: 11.5, color: "1D2733", lineSpacingMultiple: 1.02, paraSpaceAfter: 4 });

const x2 = MX + cw + 0.34;
s.addShape(pres.shapes.ROUNDED_RECTANGLE, { x: x2, y: cardY, w: cw, h: cardH, rectRadius: 0.06, fill: { color: LIGHT }, line: { color: CODEBORDER, width: 1 } });
s.addText("WHO IT'S FOR", { x: x2 + 0.2, y: cardY + 0.16, w: cw - 0.4, h: 0.28, margin: 0, fontFace: HEAD, fontSize: 11.5, bold: true, color: INK, charSpacing: 1 });
s.addText([
  { text: "Backend, API & security engineers", options: { breakLine: true } },
  { text: "Anyone building APIs, adding “Login with…”, or heading toward open banking / FAPI.", options: {} },
], { x: x2 + 0.2, y: cardY + 0.5, w: cw - 0.4, h: cardH - 0.6, margin: 0, valign: "top", fontFace: BODY, fontSize: 11.5, color: "1D2733", lineSpacingMultiple: 1.02, paraSpaceAfter: 4 });

// ---------- logistics ----------
s.addText([
  { text: "Session 1: ", options: { bold: true, color: INK } }, { text: "________", options: { color: GREY } },
  { text: "    Session 2: ", options: { bold: true, color: INK } }, { text: "to be announced", options: { color: GREY } },
  { text: "    Register: ", options: { bold: true, color: INK } }, { text: "________", options: { color: GREY } },
], { x: MX, y: 11.2, w: 7.03, h: 0.32, margin: 0, align: "center", valign: "middle", fontFace: BODY, fontSize: 11 });

pres.writeFile({ fileName: "oauth_workshop_flyer.pptx" })
  .then(f => console.log("WROTE", f)).catch(e => { console.error("ERR", e); process.exit(1); });
