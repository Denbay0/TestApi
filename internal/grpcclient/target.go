package grpcclient

import (
	"net/url"
	"strings"
)

func NormalizeTarget(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return value
	}
	if strings.Contains(value, "://") {
		u, err := url.Parse(value)
		if err == nil && u.Host != "" {
			return u.Host
		}
	}
	return value
}
