package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                   string
	MetricsPort            string
	RedisURL               string
	IdentityServiceURL     string
	EventCommandServiceURL string
	EventQueryServiceURL   string
	ReportServiceURL       string
	FrontendOrigins        []string
	AuthCookieSecure       bool
	OpenAPIServerURL       string
}

func Load() Config {
	return Config{
		Port:                   getEnv("PORT", "8080"),
		MetricsPort:            getEnv("METRICS_PORT", "9100"),
		RedisURL:               getEnv("REDIS_URL", "redis://redis:6379"),
		IdentityServiceURL:     getEnv("IDENTITY_SERVICE_URL", "http://identity-svc:50051"),
		EventCommandServiceURL: getEnv("EVENT_COMMAND_SERVICE_URL", "http://event-command-svc:50052"),
		EventQueryServiceURL:   getEnv("EVENT_QUERY_SERVICE_URL", "http://event-query-svc:50053"),
		ReportServiceURL:       getEnv("REPORT_SERVICE_URL", "http://report-svc:50054"),
		FrontendOrigins:        splitCSV(getEnv("FRONTEND_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		AuthCookieSecure:       getBool("AUTH_COOKIE_SECURE", false),
		OpenAPIServerURL:       getEnv("OPENAPI_SERVER_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}
