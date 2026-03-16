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
		_, _ = w.Write([]byte(`<!doctype html><html><head><title>Edge API Docs</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>window.ui = SwaggerUIBundle({url:'/openapi.json',dom_id:'#swagger-ui'});</script></body></html>`))
	}
}

func op(summary string) map[string]any {
	return map[string]any{"summary": summary, "responses": map[string]any{"200": map[string]any{"description": "OK"}}}
}

func getPath(summary string) map[string]any  { return map[string]any{"get": op(summary)} }
func postPath(summary string) map[string]any { return map[string]any{"post": op(summary)} }
