package base

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Render writes a minimal HTML page. Body is treated as trusted HTML the caller built;
// escape any untrusted values with Esc before interpolating them.
func Render(w http.ResponseWriter, title, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	const page = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; line-height: 1.5; max-width: 980px; margin: 40px auto; padding: 0 20px; color: #172033; }
    h1, h2 { line-height: 1.2; }
    pre { background: #f5f7fb; border: 1px solid #d8deea; padding: 14px; overflow-x: auto; border-radius: 6px; white-space: pre-wrap; word-break: break-all; }
    code { background: #f5f7fb; padding: 1px 4px; border-radius: 4px; }
    .button { display: inline-block; background: #1b5cff; color: white; border-radius: 6px; padding: 10px 14px; text-decoration: none; font-weight: 700; margin: 4px 0; }
    .secondary { background: #32415c; }
    .bad { background: #c0492f; }
    .ok { color: #0f8a5f; font-weight: 700; }
    .fail { color: #c0492f; font-weight: 700; }
  </style>
</head>
<body>
  <h1>{{.Title}}</h1>
  {{.Body}}
</body>
</html>`
	t := template.Must(template.New("p").Parse(page))
	_ = t.Execute(w, map[string]any{"Title": title, "Body": template.HTML(body)})
}

// Esc HTML-escapes a value for safe interpolation into a Render body.
func Esc(v string) string { return template.HTMLEscapeString(v) }

// RandomString returns n URL-safe random characters.
func RandomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

// PKCEChallenge returns the S256 code_challenge for a verifier.
func PKCEChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// PostForm posts a urlencoded form, optionally with HTTP Basic client auth, and
// returns the body, status code, and any transport error.
func PostForm(target string, form url.Values, basicUser, basicPass string) (string, int, error) {
	req, _ := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if basicUser != "" {
		req.SetBasicAuth(basicUser, basicPass)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b), resp.StatusCode, nil
}
