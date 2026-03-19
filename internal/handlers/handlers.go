package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/edge-api/internal/auth"
	"github.com/example/edge-api/internal/grpcclient"
	"github.com/example/edge-api/internal/middleware"
	"github.com/example/edge-api/internal/response"
)

type Handler struct {
	GRPC      *grpcclient.ClientSet
	CookieCfg auth.CookieConfig
}

func New(grpc *grpcclient.ClientSet, cookieCfg auth.CookieConfig) *Handler {
	return &Handler{GRPC: grpc, CookieCfg: cookieCfg}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>NeverNet Edge API</title>
  <style>
    :root{color-scheme:dark;--bg:#081226;--panel:#0d1b33;--line:rgba(139,233,255,.18);--text:#eaf8ff;--muted:#9ab6cf;--blue:#4fb7ff;}
    *{box-sizing:border-box} body{margin:0;font-family:Inter,system-ui,-apple-system,sans-serif;background:radial-gradient(circle at top,#123a73 0%,#081226 45%,#050a15 100%);color:var(--text)}
    .wrap{width:min(920px,calc(100% - 32px));margin:0 auto;padding:48px 0} .panel{background:rgba(13,27,51,.94);border:1px solid var(--line);border-radius:24px;padding:28px;box-shadow:0 18px 60px rgba(0,0,0,.28)}
    h1{margin:0 0 12px;font-size:clamp(28px,5vw,46px)} p{color:var(--muted);line-height:1.6} .links{display:grid;gap:12px;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));margin-top:22px}
    a{display:block;padding:16px 18px;border-radius:18px;border:1px solid rgba(139,233,255,.18);background:rgba(255,255,255,.03);color:#8be9ff;text-decoration:none;font-weight:700} a small{display:block;margin-top:6px;color:var(--muted);font-weight:500}
    .hint{margin-top:18px;color:var(--muted);font-size:14px} code{font-family:ui-monospace,SFMono-Regular,Menlo,monospace;color:#8be9ff}
  </style>
</head>
<body>
  <div class="wrap">
    <section class="panel">
      <h1>NeverNet Edge API</h1>
      <p>HTTP gateway / BFF for Rust gRPC backend. Start here for docs, machine-readable OpenAPI files, and health endpoints.</p>
      <div class="links">
        <a href="/docs">/docs<small>Interactive local API docs with grouped operations and Try it out.</small></a>
        <a href="/openapi.json">/openapi.json<small>Full OpenAPI specification in JSON.</small></a>
        <a href="/openapi.yaml">/openapi.yaml<small>Full OpenAPI specification in YAML.</small></a>
        <a href="/health">/health<small>Primary health endpoint.</small></a>
      </div>
      <div class="hint">No default 404 here, and <code>/favicon.ico</code> returns <code>204 No Content</code> to keep browser logs clean.</div>
    </section>
  </div>
</body>
</html>`)
}

func (h *Handler) Favicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) CSRF(w http.ResponseWriter, r *http.Request) {
	token, err := auth.NewCSRFToken()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "csrf_generation_failed", "failed to generate csrf token", middleware.RequestIDFromContext(r.Context()))
		return
	}
	auth.SetCSRFCookie(w, token, h.CookieCfg)
	response.JSON(w, http.StatusOK, map[string]string{"csrf_token": token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	var input grpcclient.RegisterParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_json", "invalid json body", requestID)
		return
	}
	session, err := h.GRPC.Identity.Register(r.Context(), requestID, input)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	if session != nil {
		auth.SetAuthCookie(w, session.Token, h.CookieCfg)
	}
	response.JSON(w, http.StatusCreated, map[string]any{"session": session, "todo": "wire register"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	var input grpcclient.LoginParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_json", "invalid json body", requestID)
		return
	}
	session, err := h.GRPC.Identity.Login(r.Context(), requestID, input)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	if session != nil {
		auth.SetAuthCookie(w, session.Token, h.CookieCfg)
	}
	response.JSON(w, http.StatusOK, map[string]any{"session": session, "todo": "wire login"})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearAuthCookie(w, h.CookieCfg)
	auth.ClearCSRFCookie(w, h.CookieCfg)
	response.JSON(w, http.StatusOK, map[string]any{"logged_out": true})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.Identity.Me(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.EventQuery.Categories(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	if r.Method == http.MethodGet {
		data, err := h.GRPC.EventQuery.Events(r.Context(), requestID)
		if err != nil {
			response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
			return
		}
		response.JSON(w, http.StatusOK, data)
		return
	}
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_json", "invalid json body", requestID)
		return
	}
	data, err := h.GRPC.EventCommand.CreateEvent(r.Context(), requestID, body)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusCreated, data)
}

func (h *Handler) Calendar(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.EventQuery.Calendar(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.EventQuery.Dashboard(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) ReportSummary(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.Report.Summary(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) ReportByCategory(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.Report.ByCategory(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	if r.Method == http.MethodGet {
		data, err := h.GRPC.EventQuery.Settings(r.Context(), requestID)
		if err != nil {
			response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
			return
		}
		response.JSON(w, http.StatusOK, data)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"updated": true, "todo": "wire UpdateSettings RPC"})
}

func (h *Handler) Exports(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.RequestIDFromContext(r.Context())
	data, err := h.GRPC.EventQuery.Exports(r.Context(), requestID)
	if err != nil {
		response.Error(w, http.StatusBadGateway, "upstream_error", "upstream service error", requestID)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
