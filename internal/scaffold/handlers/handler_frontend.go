package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateFrontend evalúa y mapea los archivos del framework Frontend seleccionado hacia apps/web/.
func EvaluateFrontend(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.Frontend == "" || config.Frontend == "none" {
		return "", false
	}

	frontendPrefix := "manual/frontend/" + config.Frontend + "/"
	if strings.HasPrefix(rel, frontendPrefix) {
		dest := "apps/web/" + strings.TrimPrefix(rel, frontendPrefix)
		return dest, true
	}

	return "", false
}
