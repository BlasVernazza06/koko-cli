package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/core/types"
)

// EvaluateTesting evalúa y mapea los archivos de configuración y tests del framework de pruebas seleccionado.
func EvaluateTesting(rel string, config types.ScaffoldConfig) (string, bool) {
	testingFw := strings.ToLower(strings.TrimSpace(config.Testing))
	if testingFw == "" || testingFw == "none" {
		return "", false
	}

	testingPrefix := "templates/testing/" + testingFw + "/"
	if !strings.HasPrefix(rel, testingPrefix) {
		return "", false
	}

	trimmed := strings.TrimPrefix(rel, testingPrefix)

	// 1. Pytest para FastAPI backend -> apps/api/
	if testingFw == "pytest" {
		fastapiPrefix := "fastapi/"
		if strings.HasPrefix(trimmed, fastapiPrefix) && strings.ToLower(config.Backend) == "fastapi" {
			dest := "apps/api/" + strings.TrimPrefix(trimmed, fastapiPrefix)
			return dest, true
		}
		return "", false
	}

	// 2. Tests de backend Node (Express, Hono, NestJS) para Jest / Vitest -> apps/api/
	if strings.HasPrefix(trimmed, "node_backend/") {
		b := strings.ToLower(config.Backend)
		if b == "express" || b == "hono" || b == "nestjs" {
			dest := "apps/api/" + strings.TrimPrefix(trimmed, "node_backend/")
			return dest, true
		}
		return "", false
	}

	// 3. Tests y configuraciones para Frontend -> apps/web/
	fe := strings.ToLower(config.Frontend)
	if fe != "" && fe != "none" {
		fePrefix := fe + "/"
		if strings.HasPrefix(trimmed, fePrefix) {
			dest := "apps/web/" + strings.TrimPrefix(trimmed, fePrefix)
			return dest, true
		}
		// Compatibilidad de Vue para Nuxt si aplica
		if fe == "nuxt" && strings.HasPrefix(trimmed, "vue/") {
			dest := "apps/web/" + strings.TrimPrefix(trimmed, "vue/")
			return dest, true
		}
	}

	return "", false
}
