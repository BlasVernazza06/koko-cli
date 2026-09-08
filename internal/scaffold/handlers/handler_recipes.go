package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// NormalizeRecipe normaliza el nombre o alias de la receta seleccionada
func NormalizeRecipe(recipe string) string {
	r := strings.ToLower(strings.TrimSpace(recipe))
	r = strings.ReplaceAll(r, "-", "_")
	switch r {
	case "saas", "saas_next", "next_saas":
		return "saas"
	case "java_spring", "spring", "spring_boot", "java":
		return "java_spring"
	case "enterprise_nestjs", "nestjs_next", "nestjs":
		return "enterprise_nestjs"
	case "mern":
		return "mern"
	case "pern":
		return "pern"
	case "fastapi_react", "python_fastapi", "fastapi":
		return "fastapi_react"
	case "mobile_expo", "expo_mobile", "expo":
		return "mobile_expo"
	default:
		return r
	}
}

// EvaluateRecipe evalúa las rutas de plantillas cuando el usuario selecciona una receta preconfigurada
func EvaluateRecipe(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.Recipe == "" {
		return "", false
	}

	normRecipe := NormalizeRecipe(config.Recipe)

	recipePrefix := "templates/recipes/" + normRecipe + "/"
	if strings.HasPrefix(rel, recipePrefix) {
		dest := strings.TrimPrefix(rel, recipePrefix)
		return dest, true
	}

	// Archivos compartidos de Docker para recetas (si la receta no incluye su propio docker-compose)
	if rel == "manual/docker/docker-compose.yml" || rel == "templates/docker/docker-compose.yml" {
		if normRecipe == "saas" || normRecipe == "pern" || normRecipe == "mern" || normRecipe == "fastapi_react" || normRecipe == "java_spring" || normRecipe == "enterprise_nestjs" || normRecipe == "mobile_expo" {
			return "docker-compose.yml", true
		}
	}

	// Archivos de CI/CD para recetas (si la receta no incluye su propio workflow)
	if rel == "manual/github/ci.yml" || rel == "templates/github/ci.yml" {
		if normRecipe == "saas" || normRecipe == "enterprise_nestjs" || normRecipe == "mobile_expo" || normRecipe == "java_spring" {
			return ".github/workflows/ci.yml", true
		}
	}

	return "", false
}
