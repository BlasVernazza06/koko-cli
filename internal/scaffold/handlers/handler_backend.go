package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateBackend evalúa y mapea los archivos del servidor Backend seleccionado hacia apps/api/.
func EvaluateBackend(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.Backend == "" || config.Backend == "none" {
		return "", false
	}

	backendPrefix := "manual/backend/" + config.Backend + "/"
	if strings.HasPrefix(rel, backendPrefix) {
		dest := "apps/api/" + strings.TrimPrefix(rel, backendPrefix)
		return dest, true
	}

	return "", false
}
