package docs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SpecConfig struct {
	ServerURL string
}

func OpenAPIJSON(cfg SpecConfig) ([]byte, error) {
	spec := map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":   "Edge API",
			"version": "0.1.0",
		},
		"servers": []map[string]string{{"url": cfg.ServerURL}},
		"paths": map[string]any{
			"/health":                  getPath("Healthcheck"),
			"/healthz":                 getPath("Healthcheck"),
			"/api/auth/csrf":           getPath("Get CSRF token"),
			"/api/auth/register":       postPath("Register"),
			"/api/auth/login":          postPath("Login"),
			"/api/auth/logout":         postPath("Logout"),
			"/api/auth/me":             getPath("Current user"),
			"/api/categories":          getPath("Categories"),
			"/api/events":              map[string]any{"get": op("Events"), "post": op("Create event")},
			"/api/calendar":            getPath("Calendar"),
			"/api/dashboard":           getPath("Dashboard"),
			"/api/reports/summary":     getPath("Reports summary"),
			"/api/reports/by-category": getPath("Reports by category"),
			"/api/settings":            map[string]any{"get": op("Get settings"), "put": op("Update settings")},
			"/api/exports":             getPath("Exports"),
		},
	}
	return json.MarshalIndent(spec, "", "  ")
}

func OpenAPIYAML(cfg SpecConfig) []byte {
	return []byte(fmt.Sprintf("openapi: 3.1.0\ninfo:\n  title: Edge API\n  version: 0.1.0\nservers:\n  - url: %s\n", cfg.ServerURL))
}

func SwaggerUI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Edge API Docs</title>
  <style>
    body{font-family:system-ui,-apple-system,sans-serif;max-width:1100px;margin:24px auto;padding:0 16px;line-height:1.5}
    h1{margin-bottom:8px}
    .links{margin-bottom:16px}
    a{margin-right:16px}
    pre{background:#0f172a;color:#e2e8f0;padding:16px;border-radius:8px;overflow:auto}
  </style>
</head>
<body>
  <h1>Edge API Docs</h1>
  <div class="links">
    <a href="/openapi.json" target="_blank" rel="noopener">openapi.json</a>
    <a href="/openapi.yaml" target="_blank" rel="noopener">openapi.yaml</a>
  </div>
  <pre id="spec">Loading /openapi.json ...</pre>
  <script>
    fetch('/openapi.json')
      .then(function(res){ if(!res.ok) throw new Error('HTTP '+res.status); return res.json(); })
      .then(function(spec){ document.getElementById('spec').textContent = JSON.stringify(spec, null, 2); })
      .catch(function(err){ document.getElementById('spec').textContent = 'Failed to load /openapi.json: ' + err.message; });
  </script>
</body>
</html>`))
	}
}

func op(summary string) map[string]any {
	return map[string]any{"summary": summary, "responses": map[string]any{"200": map[string]any{"description": "OK"}}}
}

func getPath(summary string) map[string]any  { return map[string]any{"get": op(summary)} }
func postPath(summary string) map[string]any { return map[string]any{"post": op(summary)} }
