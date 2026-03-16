package handlers

import (
	"encoding/json"
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
