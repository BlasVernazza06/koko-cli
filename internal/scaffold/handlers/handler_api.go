package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// EvaluateAPI evalúa y mapea los archivos de la capa de API (tRPC u oRPC).
func EvaluateAPI(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.API == "" || config.API == "none" {
		return "", false
	}

	apiPrefix := "manual/api/" + config.API + "/"
	if !strings.HasPrefix(rel, apiPrefix) {
		return "", false
	}

	trimmed := strings.TrimPrefix(rel, apiPrefix)

	// 1. Paquete compartido packages/api/
	if strings.HasPrefix(trimmed, "packages/api/") {
		return trimmed, true
	}

	isFullstack := config.Backend == "" || config.Backend == "none" || config.Backend == "self"

	// 2. Integración Server en Frontend Fullstack (ej. Next.js App Router)
	if isFullstack {
		fullstackPrefix := "fullstack/" + config.Frontend + "/"
		if strings.HasPrefix(trimmed, fullstackPrefix) {
			dest := "apps/web/" + strings.TrimPrefix(trimmed, fullstackPrefix)
			return dest, true
		}
	}

	// 3. Cliente Web (React, Next.js, etc.)
	webPrefix := "web/" + config.Frontend + "/"
	if strings.HasPrefix(trimmed, webPrefix) {
		dest := "apps/web/" + strings.TrimPrefix(trimmed, webPrefix)
		return dest, true
	}

	// 4. Integración Server en Backend Dedicado (Express, Hono, etc.)
	if !isFullstack {
		serverPrefix := "server/" + config.Backend + "/"
		if strings.HasPrefix(trimmed, serverPrefix) {
			dest := "apps/api/" + strings.TrimPrefix(trimmed, serverPrefix)
			return dest, true
		}
	}

	return "", false
}
