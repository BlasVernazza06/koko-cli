package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateRoot evalúa y mapea los archivos raíz del monorepo (configuración global)
// y los paquetes compartidos en packages/.
func EvaluateRoot(rel string, config types.ScaffoldConfig) (string, bool) {
	// Ignorar siempre archivos .gitkeep
	if strings.HasSuffix(rel, ".gitkeep") {
		return "", false
	}

	// Archivos raíz del monorepo -> raíz del proyecto
	if strings.HasPrefix(rel, "templates/root/") {
		dest := strings.TrimPrefix(rel, "templates/root/")
		return dest, true
	}

	// Paquetes compartidos -> packages/
	if strings.HasPrefix(rel, "templates/packages/") {
		dest := strings.TrimPrefix(rel, "templates/")
		return dest, true
	}

	return "", false
}
