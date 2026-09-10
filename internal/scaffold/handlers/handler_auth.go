package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// EvaluateAuth evalúa y mapea los archivos de configuración del proveedor de autenticación.
func EvaluateAuth(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.Auth == "" || config.Auth == "none" {
		return "", false
	}

	authPrefix := "manual/auth/" + config.Auth + "/"
	if !strings.HasPrefix(rel, authPrefix) {
		return "", false
	}

	trimmed := strings.TrimPrefix(rel, authPrefix)

	// 1. Paquete compartido packages/auth/ (para arquitecturas monorepo)
	if strings.HasPrefix(trimmed, "packages/auth/") {
		return trimmed, true
	}

	// 2. Archivos de configuración específicos para el frontend fullstack seleccionado
	fullstackPrefix := "fullstack/" + config.Frontend + "/"
	if strings.HasPrefix(trimmed, fullstackPrefix) {
		dest := "apps/web/" + strings.TrimPrefix(trimmed, fullstackPrefix)
		return dest, true
	}

	// 3. Cliente web de autenticación (authClient) para el frontend seleccionado
	webPrefix := "web/" + config.Frontend + "/"
	if strings.HasPrefix(trimmed, webPrefix) {
		dest := "apps/web/" + strings.TrimPrefix(trimmed, webPrefix)
		return dest, true
	}

	return "", false
}

