package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/example/edge-api/internal/auth"
	"github.com/example/edge-api/internal/config"
	"github.com/example/edge-api/internal/docs"
	"github.com/example/edge-api/internal/grpcclient"
	"github.com/example/edge-api/internal/handlers"
	"github.com/example/edge-api/internal/middleware"
)

type App struct {
	HTTPServer    *http.Server
	MetricsServer *http.Server
	grpcClients   *grpcclient.ClientSet
	logger        *slog.Logger
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	clients, err := grpcclient.New(ctx, grpcclient.DialConfig{
		Identity:     cfg.IdentityServiceURL,
		EventCommand: cfg.EventCommandServiceURL,
		EventQuery:   cfg.EventQueryServiceURL,
		Report:       cfg.ReportServiceURL,
	})
	if err != nil {
		return nil, err
	}

	h := handlers.New(clients, auth.CookieConfig{Secure: cfg.AuthCookieSecure})

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.InjectRequestID())
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.StructuredLogger(logger))
	r.Use(cors.Handler(cors.Options{AllowedOrigins: cfg.FrontendOrigins, AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, AllowCredentials: true, MaxAge: 300}))
	r.Use(middleware.ParseAuthCookie())
	r.Use(middleware.CSRFMiddleware())

	r.Get("/health", h.Health)
	r.Get("/healthz", h.Health)

	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := docs.OpenAPIJSON(docs.SpecConfig{ServerURL: cfg.OpenAPIServerURL})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	})
	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(docs.OpenAPIYAML(docs.SpecConfig{ServerURL: cfg.OpenAPIServerURL}))
	})
	r.Get("/docs", docs.SwaggerUI())

	r.Route("/api", func(api chi.Router) {
		api.Get("/auth/csrf", h.CSRF)
		api.Post("/auth/register", h.Register)
		api.Post("/auth/login", h.Login)
		api.Post("/auth/logout", h.Logout)
		api.Get("/auth/me", h.Me)

		api.Get("/categories", h.Categories)
		api.Get("/events", h.Events)
		api.Post("/events", h.Events)
		api.Get("/calendar", h.Calendar)
		api.Get("/dashboard", h.Dashboard)
		api.Get("/reports/summary", h.ReportSummary)
		api.Get("/reports/by-category", h.ReportByCategory)
		api.Get("/settings", h.Settings)
		api.Put("/settings", h.Settings)
		api.Get("/exports", h.Exports)
	})

	metrics := chi.NewRouter()
	metrics.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	return &App{
		HTTPServer: &http.Server{
			Addr:              fmt.Sprintf(":%s", cfg.Port),
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		MetricsServer: &http.Server{
			Addr:              fmt.Sprintf(":%s", cfg.MetricsPort),
			Handler:           metrics,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		grpcClients: clients,
		logger:      logger,
	}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 2)
	go func() { errCh <- a.HTTPServer.ListenAndServe() }()
	go func() { errCh <- a.MetricsServer.ListenAndServe() }()
	return <-errCh
}

func (a *App) Shutdown(ctx context.Context) error {
	a.grpcClients.Close()
	_ = a.MetricsServer.Shutdown(ctx)
	return a.HTTPServer.Shutdown(ctx)
}
