package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type SpecConfig struct {
	ServerURL string
}

func OpenAPIJSON(cfg SpecConfig) ([]byte, error) {
	return json.MarshalIndent(spec(cfg), "", "  ")
}

func OpenAPIYAML(cfg SpecConfig) []byte {
	var buf bytes.Buffer
	writeYAML(&buf, spec(cfg), 0)
	return buf.Bytes()
}

func UIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerLikeHTML))
	}
}

func spec(cfg SpecConfig) map[string]any {
	serverURL := cfg.ServerURL
	if strings.TrimSpace(serverURL) == "" {
		serverURL = "http://localhost:8080"
	}

	return map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "NeverNet Edge API",
			"version":     "0.1.0",
			"description": "HTTP gateway / BFF for Rust gRPC backend",
		},
		"servers": []any{
			map[string]any{
				"url":         serverURL,
				"description": "Configured edge-api server",
			},
		},
		"tags": []any{
			tag("health", "Health, readiness, and browser-friendly bootstrap endpoints."),
			tag("auth", "Cookie-based authentication and CSRF bootstrap flow."),
			tag("events", "Event listing, creation, and calendar views."),
			tag("reports", "Reporting and dashboard-oriented aggregations."),
			tag("settings", "User settings read/write endpoints."),
			tag("misc", "Miscellaneous convenience endpoints exposed by the edge API."),
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"cookieAuth": map[string]any{
					"type":        "apiKey",
					"in":          "cookie",
					"name":        "edge_auth",
					"description": "Session token stored in the edge_auth cookie.",
				},
				"csrfHeader": map[string]any{
					"type":        "apiKey",
					"in":          "header",
					"name":        "X-CSRF-Token",
					"description": "Required for mutating requests and must match the edge_csrf cookie.",
				},
			},
			"schemas": schemas(),
		},
		"paths": paths(),
	}
}

func schemas() map[string]any {
	return map[string]any{
		"Envelope": map[string]any{
			"type":        "object",
			"description": "Standard success envelope used by every successful JSON response.",
			"required":    []any{"data"},
			"properties": map[string]any{
				"data": map[string]any{
					"description": "Endpoint-specific payload.",
				},
			},
			"example": map[string]any{"data": map[string]any{"status": "ok"}},
		},
		"ErrorBody": map[string]any{
			"type":        "object",
			"description": "Machine-readable error payload.",
			"required":    []any{"code", "message"},
			"properties": map[string]any{
				"code":       strSchema("Stable error code.", "csrf_mismatch"),
				"message":    strSchema("Human-readable description.", "invalid csrf token"),
				"request_id": strSchema("Correlates logs and upstream calls.", "4fbbf3c24c3d4d41"),
			},
		},
		"ErrorEnvelope": map[string]any{
			"type":        "object",
			"description": "Standard error envelope used by every failed JSON response.",
			"required":    []any{"error"},
			"properties": map[string]any{
				"error": map[string]any{"$ref": "#/components/schemas/ErrorBody"},
			},
			"example": map[string]any{
				"error": map[string]any{
					"code":       "upstream_error",
					"message":    "upstream service error",
					"request_id": "4fbbf3c24c3d4d41",
				},
			},
		},
		"LoginRequest": map[string]any{
			"type":        "object",
			"required":    []any{"email", "password"},
			"description": "Credentials used to create an authenticated session.",
			"properties": map[string]any{
				"email":    strSchema("User e-mail address.", "rei@nerv.mil"),
				"password": strSchema("Plain-text password submitted by the client.", "correct-horse-battery-staple"),
			},
			"example": map[string]any{"email": "rei@nerv.mil", "password": "correct-horse-battery-staple"},
		},
		"RegisterRequest": map[string]any{
			"type":        "object",
			"required":    []any{"email", "password", "display_name"},
			"description": "Registration payload forwarded to the identity service.",
			"properties": map[string]any{
				"email":        strSchema("User e-mail address.", "shinji@nerv.mil"),
				"password":     strSchema("Requested account password.", "eva-unit-01"),
				"display_name": strSchema("Display name shown in the UI.", "Shinji Ikari"),
			},
			"example": map[string]any{"email": "shinji@nerv.mil", "password": "eva-unit-01", "display_name": "Shinji Ikari"},
		},
		"CreateEventRequest": map[string]any{
			"type":        "object",
			"description": "Loose event payload accepted by the current mock CreateEvent handler.",
			"required":    []any{"title", "starts_at", "category"},
			"properties": map[string]any{
				"title":       strSchema("Human-readable event title.", "Training sortie"),
				"description": strSchema("Optional event notes.", "Full sync rehearsal before deployment."),
				"category":    strSchema("Category slug or name.", "operations"),
				"starts_at":   map[string]any{"type": "string", "format": "date-time", "description": "Event start timestamp in RFC3339 format.", "example": "2026-03-19T09:00:00Z"},
				"ends_at":     map[string]any{"type": "string", "format": "date-time", "description": "Optional event end timestamp in RFC3339 format.", "example": "2026-03-19T10:00:00Z"},
			},
			"example": map[string]any{"title": "Training sortie", "description": "Full sync rehearsal before deployment.", "category": "operations", "starts_at": "2026-03-19T09:00:00Z", "ends_at": "2026-03-19T10:00:00Z"},
		},
		"UpdateSettingsRequest": map[string]any{
			"type":        "object",
			"description": "Hackathon-friendly settings payload accepted by the stub update handler.",
			"properties": map[string]any{
				"timezone":       strSchema("Preferred IANA timezone.", "UTC"),
				"theme":          strSchema("UI theme hint.", "blue"),
				"weekly_reports": map[string]any{"type": "boolean", "description": "Whether weekly report delivery is enabled.", "example": true},
			},
			"example": map[string]any{"timezone": "UTC", "theme": "blue", "weekly_reports": true},
		},
		"HealthData": map[string]any{
			"type":       "object",
			"required":   []any{"status"},
			"properties": map[string]any{"status": strSchema("Service health status.", "ok")},
		},
		"CSRFData": map[string]any{
			"type":       "object",
			"required":   []any{"csrf_token"},
			"properties": map[string]any{"csrf_token": strSchema("Token mirrored into the edge_csrf cookie.", "2hWf0B3BsvS3x1Q6uS4WjA")},
		},
		"Session": map[string]any{
			"type":       "object",
			"required":   []any{"token"},
			"properties": map[string]any{"token": strSchema("Opaque session token returned by the mock auth flow.", "todo-login-token")},
		},
		"AuthSessionData": map[string]any{
			"type":       "object",
			"properties": map[string]any{"session": map[string]any{"$ref": "#/components/schemas/Session"}, "todo": strSchema("Current stub note from the edge handler.", "wire login")},
		},
		"LogoutData": map[string]any{
			"type":       "object",
			"required":   []any{"logged_out"},
			"properties": map[string]any{"logged_out": map[string]any{"type": "boolean", "example": true}},
		},
		"MessageData": map[string]any{
			"type":                 "object",
			"description":          "Generic object wrapper for current stubbed read endpoints.",
			"additionalProperties": true,
			"example":              map[string]any{"todo": "wire dashboard RPC"},
		},
		"CollectionData": map[string]any{
			"type":        "object",
			"description": "Generic collection wrapper used by current stubbed list endpoints.",
			"properties": map[string]any{
				"items": map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": true}, "example": []any{}},
				"todo":  strSchema("Stub note returned by the current handler.", "wire Events RPC"),
			},
		},
		"CreateEventData": map[string]any{
			"type":       "object",
			"properties": map[string]any{"created": map[string]any{"type": "boolean", "example": true}, "payload": map[string]any{"$ref": "#/components/schemas/CreateEventRequest"}, "todo": strSchema("Current stub note from the edge handler.", "wire CreateEvent RPC")},
		},
		"UpdateSettingsData": map[string]any{
			"type":       "object",
			"properties": map[string]any{"updated": map[string]any{"type": "boolean", "example": true}, "todo": strSchema("Current stub note from the edge handler.", "wire UpdateSettings RPC")},
		},
	}
}

func paths() map[string]any {
	return map[string]any{
		"/": map[string]any{
			"get": map[string]any{
				"tags":        []any{"misc"},
				"summary":     "HTML index",
				"description": "Browser-friendly landing page with quick links to docs, OpenAPI artifacts, and health checks.",
				"responses":   map[string]any{"200": htmlResponse("Small landing page for humans opening the service root in a browser.")},
			},
		},
		"/favicon.ico": map[string]any{
			"get": map[string]any{
				"tags":        []any{"misc"},
				"summary":     "Favicon placeholder",
				"description": "Returns 204 No Content so browsers do not generate noisy 404 logs.",
				"responses":   map[string]any{"204": map[string]any{"description": "No favicon content."}},
			},
		},
		"/health":  healthOp("Primary health endpoint."),
		"/healthz": healthOp("Secondary health endpoint for readiness or liveness probes."),
		"/api/auth/csrf": map[string]any{"get": map[string]any{
			"tags":        []any{"auth"},
			"summary":     "Fetch CSRF token",
			"description": "Generates a CSRF token, sets the edge_csrf cookie, and returns the same token in the standard success envelope.",
			"responses": map[string]any{
				"200": jsonResponse("CSRF token generated.", envelopeSchemaRef("#/components/schemas/CSRFData"), map[string]any{"data": map[string]any{"csrf_token": "2hWf0B3BsvS3x1Q6uS4WjA"}}),
				"500": errorResponse("Token generation failed.", "csrf_generation_failed"),
			},
		}},
		"/api/auth/register": map[string]any{"post": map[string]any{
			"tags":        []any{"auth"},
			"summary":     "Register a new user",
			"description": "Creates a mock identity session and sets the auth cookie when registration succeeds.",
			"security":    []any{map[string]any{"csrfHeader": []any{}}},
			"requestBody": jsonRequestBody("Registration payload forwarded to the identity service.", "#/components/schemas/RegisterRequest", true),
			"responses": map[string]any{
				"201": jsonResponse("User registered and session created.", envelopeSchemaRef("#/components/schemas/AuthSessionData"), map[string]any{"data": map[string]any{"session": map[string]any{"token": "todo-register-token"}, "todo": "wire register"}}),
				"400": errorResponse("Malformed JSON payload.", "invalid_json"),
				"403": errorResponse("Missing or invalid CSRF token.", "csrf_mismatch"),
				"502": errorResponse("Identity service returned an error.", "upstream_error"),
			},
		}},
		"/api/auth/login": map[string]any{"post": map[string]any{
			"tags":        []any{"auth"},
			"summary":     "Login",
			"description": "Authenticates a user, sets the edge_auth cookie, and returns the created session envelope.",
			"security":    []any{map[string]any{"csrfHeader": []any{}}},
			"requestBody": jsonRequestBody("Credentials used to create a session.", "#/components/schemas/LoginRequest", true),
			"responses": map[string]any{
				"200": jsonResponse("Authenticated session created.", envelopeSchemaRef("#/components/schemas/AuthSessionData"), map[string]any{"data": map[string]any{"session": map[string]any{"token": "todo-login-token"}, "todo": "wire login"}}),
				"400": errorResponse("Malformed JSON payload.", "invalid_json"),
				"403": errorResponse("Missing or invalid CSRF token.", "csrf_mismatch"),
				"502": errorResponse("Identity service returned an error.", "upstream_error"),
			},
		}},
		"/api/auth/logout": map[string]any{"post": map[string]any{
			"tags":        []any{"auth"},
			"summary":     "Logout",
			"description": "Clears both auth and CSRF cookies and returns a simple success envelope.",
			"security":    []any{map[string]any{"cookieAuth": []any{}, "csrfHeader": []any{}}},
			"responses": map[string]any{
				"200": jsonResponse("Session cookies cleared.", envelopeSchemaRef("#/components/schemas/LogoutData"), map[string]any{"data": map[string]any{"logged_out": true}}),
				"403": errorResponse("Missing or invalid CSRF token.", "csrf_mismatch"),
			},
		}},
		"/api/auth/me": map[string]any{"get": map[string]any{
			"tags":        []any{"auth"},
			"summary":     "Current user",
			"description": "Returns the current mock identity payload resolved from the cookie-backed session.",
			"security":    []any{map[string]any{"cookieAuth": []any{}}},
			"responses": map[string]any{
				"200": jsonResponse("Current user payload.", envelopeSchemaRef("#/components/schemas/MessageData"), map[string]any{"data": map[string]any{"status": "TODO: wire identity Me rpc"}}),
				"502": errorResponse("Identity service returned an error.", "upstream_error"),
			},
		}},
		"/api/categories": collectionGet("events", "List categories", "Returns categories for event creation filters and selectors."),
		"/api/events": map[string]any{
			"get": map[string]any{
				"tags":        []any{"events"},
				"summary":     "List events",
				"description": "Returns the current stubbed events list envelope.",
				"responses": map[string]any{
					"200": jsonResponse("Event collection returned.", envelopeSchemaRef("#/components/schemas/CollectionData"), map[string]any{"data": map[string]any{"items": []any{}, "todo": "wire Events RPC"}}),
					"502": errorResponse("Event query service returned an error.", "upstream_error"),
				},
			},
			"post": map[string]any{
				"tags":        []any{"events"},
				"summary":     "Create event",
				"description": "Passes a loose JSON event payload to the command service and returns the created envelope.",
				"security":    []any{map[string]any{"cookieAuth": []any{}, "csrfHeader": []any{}}},
				"requestBody": jsonRequestBody("Event payload accepted by the current mock CreateEvent handler.", "#/components/schemas/CreateEventRequest", true),
				"responses": map[string]any{
					"201": jsonResponse("Event created.", envelopeSchemaRef("#/components/schemas/CreateEventData"), map[string]any{"data": map[string]any{"created": true, "payload": map[string]any{"title": "Training sortie", "description": "Full sync rehearsal before deployment.", "category": "operations", "starts_at": "2026-03-19T09:00:00Z", "ends_at": "2026-03-19T10:00:00Z"}, "todo": "wire CreateEvent RPC"}}),
					"400": errorResponse("Malformed JSON payload.", "invalid_json"),
					"403": errorResponse("Missing or invalid CSRF token.", "csrf_mismatch"),
					"502": errorResponse("Event command service returned an error.", "upstream_error"),
				},
			},
		},
		"/api/calendar":            collectionGet("events", "Calendar view", "Returns calendar-oriented event items for the active user."),
		"/api/dashboard":           messageGet("reports", "Dashboard snapshot", "Returns the current dashboard summary payload."),
		"/api/reports/summary":     messageGet("reports", "Reports summary", "Returns the top-level reporting summary used by the dashboard."),
		"/api/reports/by-category": collectionGet("reports", "Reports by category", "Returns grouped report entries by category."),
		"/api/settings": map[string]any{
			"get": map[string]any{
				"tags":        []any{"settings"},
				"summary":     "Get settings",
				"description": "Returns the current settings payload from the query service.",
				"security":    []any{map[string]any{"cookieAuth": []any{}}},
				"responses": map[string]any{
					"200": jsonResponse("Settings payload returned.", envelopeSchemaRef("#/components/schemas/MessageData"), map[string]any{"data": map[string]any{"todo": "wire Settings RPC"}}),
					"502": errorResponse("Event query service returned an error.", "upstream_error"),
				},
			},
			"put": map[string]any{
				"tags":        []any{"settings"},
				"summary":     "Update settings",
				"description": "Accepts a lightweight settings payload and returns the current stub update response envelope.",
				"security":    []any{map[string]any{"cookieAuth": []any{}, "csrfHeader": []any{}}},
				"requestBody": jsonRequestBody("Settings payload accepted by the current stub handler.", "#/components/schemas/UpdateSettingsRequest", true),
				"responses": map[string]any{
					"200": jsonResponse("Settings update acknowledged.", envelopeSchemaRef("#/components/schemas/UpdateSettingsData"), map[string]any{"data": map[string]any{"updated": true, "todo": "wire UpdateSettings RPC"}}),
					"403": errorResponse("Missing or invalid CSRF token.", "csrf_mismatch"),
				},
			},
		},
		"/api/exports": collectionGet("misc", "Exports list", "Returns available export artifacts exposed by the query service."),
	}
}

func tag(name, desc string) map[string]any {
	return map[string]any{"name": name, "description": desc}
}

func strSchema(description, example string) map[string]any {
	return map[string]any{"type": "string", "description": description, "example": example}
}

func healthOp(description string) map[string]any {
	return map[string]any{"get": map[string]any{
		"tags":        []any{"health"},
		"summary":     "Health check",
		"description": description,
		"responses": map[string]any{
			"200": jsonResponse("Service is healthy.", envelopeSchemaRef("#/components/schemas/HealthData"), map[string]any{"data": map[string]any{"status": "ok"}}),
		},
	}}
}

func collectionGet(tagName, summary, description string) map[string]any {
	return map[string]any{"get": map[string]any{
		"tags":        []any{tagName},
		"summary":     summary,
		"description": description,
		"responses": map[string]any{
			"200": jsonResponse("Collection returned.", envelopeSchemaRef("#/components/schemas/CollectionData"), map[string]any{"data": map[string]any{"items": []any{}, "todo": "wire collection RPC"}}),
			"502": errorResponse("Upstream service returned an error.", "upstream_error"),
		},
	}}
}

func messageGet(tagName, summary, description string) map[string]any {
	return map[string]any{"get": map[string]any{
		"tags":        []any{tagName},
		"summary":     summary,
		"description": description,
		"responses": map[string]any{
			"200": jsonResponse("Payload returned.", envelopeSchemaRef("#/components/schemas/MessageData"), map[string]any{"data": map[string]any{"todo": "wire RPC"}}),
			"502": errorResponse("Upstream service returned an error.", "upstream_error"),
		},
	}}
}

func envelopeSchemaRef(dataRef string) map[string]any {
	return map[string]any{
		"allOf": []any{
			map[string]any{"$ref": "#/components/schemas/Envelope"},
			map[string]any{
				"type": "object",
				"properties": map[string]any{
					"data": map[string]any{"$ref": dataRef},
				},
			},
		},
	}
}

func jsonRequestBody(description, schemaRef string, required bool) map[string]any {
	return map[string]any{
		"required":    required,
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{"$ref": schemaRef},
			},
		},
	}
}

func jsonResponse(description string, schema map[string]any, example map[string]any) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema":  schema,
				"example": example,
			},
		},
	}
}

func errorResponse(description, code string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{"$ref": "#/components/schemas/ErrorEnvelope"},
				"example": map[string]any{
					"error": map[string]any{
						"code":       code,
						"message":    description,
						"request_id": "4fbbf3c24c3d4d41",
					},
				},
			},
		},
	}
}

func htmlResponse(description string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"text/html": map[string]any{
				"schema": map[string]any{"type": "string"},
			},
		},
	}
}

func writeYAML(buf *bytes.Buffer, v any, indent int) {
	switch value := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			buf.WriteString(strings.Repeat(" ", indent))
			buf.WriteString(key)
			buf.WriteString(":")
			if isScalar(value[key]) {
				buf.WriteString(" ")
				buf.WriteString(yamlScalar(value[key]))
				buf.WriteByte('\n')
				continue
			}
			buf.WriteByte('\n')
			writeYAML(buf, value[key], indent+2)
		}
	case []any:
		if len(value) == 0 {
			buf.WriteString(strings.Repeat(" ", indent))
			buf.WriteString("[]\n")
			return
		}
		for _, item := range value {
			buf.WriteString(strings.Repeat(" ", indent))
			buf.WriteString("-")
			if isScalar(item) {
				buf.WriteString(" ")
				buf.WriteString(yamlScalar(item))
				buf.WriteByte('\n')
				continue
			}
			buf.WriteByte('\n')
			writeYAML(buf, item, indent+2)
		}
	default:
		buf.WriteString(strings.Repeat(" ", indent))
		buf.WriteString(yamlScalar(value))
		buf.WriteByte('\n')
	}
}

func isScalar(v any) bool {
	switch v.(type) {
	case nil, string, bool, int, int64, float64, float32, json.Number:
		return true
	default:
		return false
	}
}

func yamlScalar(v any) string {
	switch value := v.(type) {
	case nil:
		return "null"
	case string:
		return strconv.Quote(value)
	case bool:
		if value {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case json.Number:
		return value.String()
	default:
		return strconv.Quote(fmt.Sprint(value))
	}
}

const swaggerLikeHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>NeverNet Edge API Docs</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #081226;
      --panel: #0d1b33;
      --panel-soft: #10223f;
      --line: rgba(139,233,255,.18);
      --line-strong: rgba(139,233,255,.32);
      --text: #eaf8ff;
      --muted: #9ab6cf;
      --blue: #4fb7ff;
      --green: #27c281;
      --yellow: #f7c948;
      --red: #ff7d7d;
      --purple: #9b8cff;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: radial-gradient(circle at top, #123a73 0%, #081226 45%, #050a15 100%);
      color: var(--text);
    }
    a { color: #8be9ff; }
    .wrap { width: min(1240px, calc(100% - 32px)); margin: 0 auto; padding: 24px 0 48px; }
    .hero, .panel, .tag-group, .op-card, .response-card { background: rgba(13, 27, 51, 0.94); border: 1px solid var(--line); border-radius: 20px; box-shadow: 0 18px 60px rgba(0,0,0,.28); }
    .hero { padding: 28px; margin-bottom: 20px; }
    .hero h1 { margin: 0 0 10px; font-size: clamp(28px, 4vw, 44px); }
    .hero p { margin: 0 0 16px; color: var(--muted); max-width: 860px; }
    .hero-grid, .meta-grid { display: grid; gap: 14px; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); }
    .pill-row, .link-row { display: flex; flex-wrap: wrap; gap: 10px; }
    .pill, .method { display: inline-flex; align-items: center; border-radius: 999px; font-weight: 700; letter-spacing: .02em; }
    .pill { padding: 8px 12px; border: 1px solid var(--line-strong); color: var(--muted); background: rgba(255,255,255,.03); }
    .method { min-width: 72px; justify-content: center; padding: 8px 12px; color: #07111f; }
    .method.get { background: #5bc0ff; }
    .method.post { background: #3ddc97; }
    .method.put { background: #f7c948; }
    .method.delete { background: #ff7d7d; }
    .panel { padding: 18px 20px; margin-bottom: 18px; }
    .panel h2, .tag-group h2 { margin: 0 0 12px; font-size: 20px; }
    .meta-card { padding: 16px; border: 1px solid var(--line); border-radius: 16px; background: rgba(255,255,255,.02); }
    .meta-card strong { display: block; margin-bottom: 8px; }
    #operations { display: grid; gap: 18px; }
    .tag-group { padding: 18px; }
    .tag-group > p { color: var(--muted); margin: 0 0 14px; }
    .ops-list { display: grid; gap: 14px; }
    .op-card { padding: 18px; }
    .op-top { display: grid; gap: 12px; grid-template-columns: auto 1fr; align-items: start; }
    .op-path { font-size: 20px; font-weight: 700; word-break: break-all; }
    .op-summary { margin: 0 0 6px; font-size: 18px; }
    .op-description, .hint, .muted { color: var(--muted); }
    .grid-2 { display: grid; gap: 14px; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); }
    textarea, input, pre, code { font-family: "SFMono-Regular", ui-monospace, Menlo, Consolas, monospace; }
    textarea, input {
      width: 100%;
      border: 1px solid var(--line-strong);
      background: rgba(3, 10, 21, 0.9);
      color: var(--text);
      border-radius: 14px;
      padding: 12px 14px;
      outline: none;
    }
    textarea { min-height: 180px; resize: vertical; }
    pre {
      margin: 0;
      white-space: pre-wrap;
      word-break: break-word;
      background: rgba(3, 10, 21, 0.92);
      border: 1px solid rgba(255,255,255,.06);
      border-radius: 14px;
      padding: 14px;
      max-height: 420px;
      overflow: auto;
    }
    button {
      cursor: pointer;
      border: 0;
      border-radius: 12px;
      padding: 12px 16px;
      font-weight: 800;
      color: #06111f;
      background: linear-gradient(135deg, #8be9ff 0%, #4fb7ff 100%);
      box-shadow: 0 10px 24px rgba(79,183,255,.22);
    }
    button.secondary { background: rgba(255,255,255,.06); color: var(--text); box-shadow: none; border: 1px solid var(--line); }
    .actions { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; }
    .response-card { padding: 16px; margin-top: 14px; }
    .response-head { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; margin-bottom: 10px; }
    .status { display: inline-flex; padding: 6px 10px; border-radius: 999px; font-weight: 700; }
    .status.ok { background: rgba(39,194,129,.18); color: #80ffd0; }
    .status.err { background: rgba(255,125,125,.18); color: #ffb3b3; }
    .small { font-size: 13px; }
    .footer { margin-top: 22px; color: var(--muted); font-size: 14px; }
    .loading { padding: 40px 0; text-align: center; color: var(--muted); }
    @media (max-width: 700px) {
      .wrap { width: min(100% - 20px, 1240px); }
      .hero, .panel, .tag-group, .op-card, .response-card { border-radius: 16px; }
      .op-top { grid-template-columns: 1fr; }
      .op-path { font-size: 18px; }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <section class="hero">
      <div class="pill-row">
        <span class="pill">Swagger-style local docs</span>
        <span class="pill" id="specVersion">Loading spec…</span>
      </div>
      <h1 id="title">NeverNet Edge API</h1>
      <p id="description">Loading OpenAPI document…</p>
      <div class="link-row">
        <a href="/">Root</a>
        <a href="/openapi.json" target="_blank" rel="noopener">openapi.json</a>
        <a href="/openapi.yaml" target="_blank" rel="noopener">openapi.yaml</a>
        <a href="/health" target="_blank" rel="noopener">health</a>
      </div>
    </section>

    <section class="panel">
      <h2>Servers & security</h2>
      <div class="meta-grid" id="metaGrid"></div>
    </section>

    <section class="panel">
      <h2>Try it out notes</h2>
      <div class="grid-2">
        <div class="meta-card">
          <strong>Cookie auth</strong>
          <div class="muted">Use <code>GET /api/auth/csrf</code> first, then call login or register. Browser cookies are preserved because requests are sent with <code>credentials: include</code>.</div>
        </div>
        <div class="meta-card">
          <strong>CSRF header</strong>
          <div class="muted">For mutating operations, copy the returned token into the request header field below. It must match the <code>edge_csrf</code> cookie.</div>
          <label class="small" for="globalCsrf">Default X-CSRF-Token</label>
          <input id="globalCsrf" placeholder="Paste csrf token here for POST/PUT requests">
        </div>
      </div>
    </section>

    <section id="operations">
      <div class="loading">Loading operations…</div>
    </section>

    <div class="footer">Local documentation UI generated by the service itself. No external CDN required.</div>
  </div>
  <script>
    function resolveRef(spec, ref) {
      if (!ref || ref.indexOf('#/') !== 0) return null;
      var node = spec;
      ref.slice(2).split('/').forEach(function (part) {
        if (node) node = node[part];
      });
      return node || null;
    }

    function derefSchema(spec, schema) {
      if (!schema) return null;
      if (schema.$ref) return derefSchema(spec, resolveRef(spec, schema.$ref));
      if (schema.allOf) {
        var merged = { type: 'object', properties: {}, required: [] };
        schema.allOf.forEach(function (part) {
          var current = derefSchema(spec, part) || {};
          if (current.properties) Object.assign(merged.properties, current.properties);
          if (Array.isArray(current.required)) merged.required = merged.required.concat(current.required);
          Object.keys(current).forEach(function (key) {
            if (key !== 'properties' && key !== 'required') merged[key] = current[key];
          });
        });
        return merged;
      }
      return schema;
    }

    function schemaExample(spec, schema) {
      var actual = derefSchema(spec, schema) || {};
      if (actual.example !== undefined) return actual.example;
      if (actual.properties) {
        var obj = {};
        Object.keys(actual.properties).forEach(function (key) {
          var prop = derefSchema(spec, actual.properties[key]) || {};
          if (prop.example !== undefined) obj[key] = prop.example;
          else if (prop.type === 'object') obj[key] = schemaExample(spec, prop);
          else if (prop.type === 'array') obj[key] = [];
          else if (prop.type === 'boolean') obj[key] = false;
          else if (prop.type === 'number' || prop.type === 'integer') obj[key] = 0;
          else obj[key] = '';
        });
        return obj;
      }
      if (actual.type === 'array') return [];
      return {};
    }

    function pretty(value) {
      return JSON.stringify(value, null, 2);
    }

    function methodClass(method) {
      return method.toLowerCase();
    }

    function buildMeta(spec) {
      var metaGrid = document.getElementById('metaGrid');
      metaGrid.innerHTML = '';
      var serverCard = document.createElement('div');
      serverCard.className = 'meta-card';
      serverCard.innerHTML = '<strong>Servers</strong><div class="muted">' + (spec.servers || []).map(function (s) { return s.url + (s.description ? ' — ' + s.description : ''); }).join('<br>') + '</div>';
      metaGrid.appendChild(serverCard);

      var securityCard = document.createElement('div');
      securityCard.className = 'meta-card';
      var schemes = spec.components && spec.components.securitySchemes ? spec.components.securitySchemes : {};
      var rows = Object.keys(schemes).map(function (name) {
        var scheme = schemes[name];
        return '<div><strong>' + name + '</strong><div class="muted">' + (scheme.description || '') + '</div></div>';
      }).join('<div style="height:10px"></div>');
      securityCard.innerHTML = '<strong>Security schemes</strong>' + rows;
      metaGrid.appendChild(securityCard);
    }

    function opSecurityText(operation) {
      if (!operation.security || !operation.security.length) return 'None';
      return operation.security.map(function (req) {
        return Object.keys(req).join(' + ');
      }).join(' or ');
    }

    function renderOperations(spec) {
      var target = document.getElementById('operations');
      target.innerHTML = '';
      var tagDescriptions = {};
      (spec.tags || []).forEach(function (tag) { tagDescriptions[tag.name] = tag.description || ''; });
      var grouped = {};

      Object.keys(spec.paths || {}).forEach(function (path) {
        var pathItem = spec.paths[path];
        Object.keys(pathItem).forEach(function (method) {
          var op = pathItem[method];
          var tag = (op.tags && op.tags[0]) || 'misc';
          if (!grouped[tag]) grouped[tag] = [];
          grouped[tag].push({ path: path, method: method.toUpperCase(), operation: op });
        });
      });

      Object.keys(grouped).sort().forEach(function (tag) {
        var group = document.createElement('section');
        group.className = 'tag-group';
        group.innerHTML = '<h2>' + tag + '</h2><p>' + (tagDescriptions[tag] || 'Operations') + '</p>';
        var list = document.createElement('div');
        list.className = 'ops-list';

        grouped[tag].forEach(function (entry) {
          var op = entry.operation;
          var card = document.createElement('article');
          card.className = 'op-card';

          var requestBody = op.requestBody && op.requestBody.content && op.requestBody.content['application/json'];
          var requestSchema = requestBody ? requestBody.schema : null;
          var requestExample = requestBody && requestBody.example ? requestBody.example : (requestSchema ? schemaExample(spec, requestSchema) : null);
          var responses = op.responses || {};
          var preferredResponse = responses['200'] || responses['201'] || responses['204'] || responses[Object.keys(responses)[0]] || {};
          var responseJSON = preferredResponse.content && preferredResponse.content['application/json'];
          var responseSchema = responseJSON ? responseJSON.schema : null;
          var responseExample = responseJSON && responseJSON.example ? responseJSON.example : (responseSchema ? schemaExample(spec, responseSchema) : null);

          card.innerHTML = '' +
            '<div class="op-top">' +
              '<span class="method ' + methodClass(entry.method) + '">' + entry.method + '</span>' +
              '<div>' +
                '<div class="op-path">' + entry.path + '</div>' +
                '<h3 class="op-summary">' + (op.summary || '') + '</h3>' +
                '<div class="op-description">' + (op.description || '') + '</div>' +
                '<div class="pill-row" style="margin-top:10px">' +
                  '<span class="pill">Security: ' + opSecurityText(op) + '</span>' +
                  '<span class="pill">Responses: ' + Object.keys(responses).join(', ') + '</span>' +
                '</div>' +
              '</div>' +
            '</div>' +
            '<div class="grid-2" style="margin-top:16px">' +
              '<div class="meta-card"><strong>Request body</strong><pre>' + (requestSchema ? pretty(derefSchema(spec, requestSchema)) : 'No request body') + '</pre></div>' +
              '<div class="meta-card"><strong>Primary response</strong><pre>' + (responseSchema ? pretty(derefSchema(spec, responseSchema)) : pretty(responses)) + '</pre></div>' +
            '</div>';

          if (entry.method !== 'GET' && entry.method !== 'HEAD') {
            var tryPanel = document.createElement('div');
            tryPanel.className = 'response-card';
            tryPanel.innerHTML = '<div class="response-head"><strong>Try it out</strong><span class="small muted">Requests are sent against the current origin and include cookies.</span></div>' +
              '<div class="grid-2">' +
                '<div>' +
                  '<label class="small">X-CSRF-Token override</label>' +
                  '<input class="csrf-input" placeholder="Uses the global token above when left blank">' +
                  '<div style="height:12px"></div>' +
                  '<label class="small">Request JSON body</label>' +
                  '<textarea class="body-input"></textarea>' +
                '</div>' +
                '<div>' +
                  '<label class="small">Reference example</label>' +
                  '<pre>' + (requestExample ? pretty(requestExample) : 'No request body required') + '</pre>' +
                '</div>' +
              '</div>' +
              '<div class="actions" style="margin-top:14px">' +
                '<button type="button" class="send-btn">Send request</button>' +
                '<button type="button" class="secondary fill-btn">Fill example</button>' +
              '</div>' +
              '<div class="response-output" style="margin-top:14px"></div>';
            card.appendChild(tryPanel);

            var csrfInput = tryPanel.querySelector('.csrf-input');
            var bodyInput = tryPanel.querySelector('.body-input');
            var output = tryPanel.querySelector('.response-output');
            bodyInput.value = requestExample ? pretty(requestExample) : '{}';

            tryPanel.querySelector('.fill-btn').addEventListener('click', function () {
              bodyInput.value = requestExample ? pretty(requestExample) : '{}';
            });

            tryPanel.querySelector('.send-btn').addEventListener('click', function () {
              output.innerHTML = '<div class="muted">Sending…</div>';
              var headers = {};
              var token = csrfInput.value.trim() || document.getElementById('globalCsrf').value.trim();
              if (token) headers['X-CSRF-Token'] = token;
              var bodyText = bodyInput.value.trim();
              var options = { method: entry.method, credentials: 'include', headers: headers };
              if (bodyText && bodyText !== '{}') {
                headers['Content-Type'] = 'application/json';
                try {
                  options.body = JSON.stringify(JSON.parse(bodyText));
                } catch (err) {
                  output.innerHTML = '<span class="status err">Invalid JSON</span><pre>' + err.message + '</pre>';
                  return;
                }
              } else if (requestSchema) {
                headers['Content-Type'] = 'application/json';
                options.body = '{}';
              }

              fetch(entry.path, options).then(async function (res) {
                var text = await res.text();
                var statusClass = res.ok ? 'ok' : 'err';
                output.innerHTML = '<div class="response-head"><span class="status ' + statusClass + '">' + res.status + ' ' + res.statusText + '</span><span class="small muted">content-type: ' + (res.headers.get('content-type') || 'unknown') + '</span></div><pre>' + text.replace(/</g, '&lt;').replace(/>/g, '&gt;') + '</pre>';
              }).catch(function (err) {
                output.innerHTML = '<span class="status err">Request failed</span><pre>' + err.message + '</pre>';
              });
            });
          }

          var responsePanel = document.createElement('div');
          responsePanel.className = 'response-card';
          responsePanel.innerHTML = '<div class="response-head"><strong>Examples & responses</strong><span class="small muted">Primary example shown below.</span></div><pre>' + pretty(responseExample || responses) + '</pre>';
          card.appendChild(responsePanel);

          list.appendChild(card);
        });

        group.appendChild(list);
        target.appendChild(group);
      });
    }

    fetch('/openapi.json').then(function (res) {
      if (!res.ok) throw new Error('Failed to load spec: HTTP ' + res.status);
      return res.json();
    }).then(function (spec) {
      document.getElementById('title').textContent = spec.info && spec.info.title ? spec.info.title : 'API Docs';
      document.getElementById('description').textContent = spec.info && spec.info.description ? spec.info.description : 'OpenAPI explorer';
      document.getElementById('specVersion').textContent = 'OpenAPI ' + spec.openapi + ' • v' + (spec.info && spec.info.version ? spec.info.version : 'n/a');
      buildMeta(spec);
      renderOperations(spec);
    }).catch(function (err) {
      document.getElementById('operations').innerHTML = '<div class="panel"><strong>Failed to load /openapi.json</strong><pre>' + err.message + '</pre></div>';
    });
  </script>
</body>
</html>`
