package scaffold

import (
	"path/filepath"

	"github.com/BlasVernazza06/koko-cli/internal/scaffold/handlers"
	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluatePath es el enrutador central que orquesta los handlers especializados
// para determinar si un archivo de plantilla debe incluirse y cuál es su destino final.
func evaluatePath(path string, config types.ScaffoldConfig) (string, bool) {
	rel := filepath.ToSlash(path)

	// 1. Modo Recetas predefinidas (SaaS, PERN, MERN, FastAPI+React)
	if config.Recipe != "" {
		return handlers.EvaluateRecipe(rel, config)
	}

	// 2. Modo Configuración Manual:
	// A. Archivos raíz y paquetes compartidos del monorepo
	if dest, ok := handlers.EvaluateRoot(rel, config); ok {
		return dest, true
	}

	// B. Framework Frontend (Next.js, React, Nuxt, Svelte)
	if dest, ok := handlers.EvaluateFrontend(rel, config); ok {
		return dest, true
	}

	// C. Servidor Backend (Express, Hono, NestJS, FastAPI, Go Chi)
	if dest, ok := handlers.EvaluateBackend(rel, config); ok {
		return dest, true
	}

	// D. Base de datos y ORM (Drizzle, Prisma, Mongoose, SQLAlchemy, GORM)
	if dest, ok := handlers.EvaluateDatabase(rel, config); ok {
		return dest, true
	}

	// E. Addons (Docker Compose, GitHub Actions)
	if dest, ok := handlers.EvaluateAddons(rel, config); ok {
		return dest, true
	}

	return "", false
}
