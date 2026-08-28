package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateRoot evalúa y mapea los archivos raíz del monorepo (configuración global)
// y los paquetes compartidos en packages/.
func EvaluateRoot(rel string, config types.ScaffoldConfig) (string, bool) {
	// Archivos raíz del monorepo -> raíz del proyecto
	if strings.HasPrefix(rel, "manual/root/") {
		if rel == "manual/root/apps/.gitkeep" {
			hasApps := (config.Frontend != "" && config.Frontend != "none") || (config.Backend != "" && config.Backend != "none")
			if hasApps {
				return "", false
			}
		}
		dest := strings.TrimPrefix(rel, "manual/root/")
		return dest, true
	}

	// Paquetes compartidos -> packages/
	if strings.HasPrefix(rel, "manual/packages/") {
		dest := strings.TrimPrefix(rel, "manual/")
		return dest, true
	}

	return "", false
}
