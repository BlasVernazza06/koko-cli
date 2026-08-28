package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateRecipe evalúa las rutas de plantillas cuando el usuario selecciona una receta preconfigurada
// (ej: saas, pern, mern, fastapi_react).
func EvaluateRecipe(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.Recipe == "" {
		return "", false
	}

	recipePrefix := "templates/recipes/" + config.Recipe + "/"
	if strings.HasPrefix(rel, recipePrefix) {
		dest := strings.TrimPrefix(rel, recipePrefix)
		return dest, true
	}

	// Archivos compartidos de Docker para recetas
	if rel == "manual/docker/docker-compose.yml" || rel == "templates/docker/docker-compose.yml" {
		if config.Recipe == "saas" || config.Recipe == "pern" || config.Recipe == "mern" {
			return "docker-compose.yml", true
		}
	}

	// Archivos de CI/CD para recetas
	if rel == "manual/github/ci.yml" || rel == "templates/github/ci.yml" {
		if config.Recipe == "saas" {
			return ".github/workflows/ci.yml", true
		}
	}

	return "", false
}
