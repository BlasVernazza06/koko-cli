package processors

import (
	"fmt"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

// ProcessEnvVariables genera automáticamente archivos .env y .env.example
// con las configuraciones y credenciales listas según el stack seleccionado.
func ProcessEnvVariables(v *vfs.VFS, cfg ProcessConfig) error {
	var envLines []string
	envLines = append(envLines, "# ----------------------------------------------------")
	envLines = append(envLines, fmt.Sprintf("# Environment variables for %s", cfg.ProjectName))
	envLines = append(envLines, "# ----------------------------------------------------\n")

	// 1. URLs de Frontend y Backend
	if cfg.Backend != "" && cfg.Backend != "none" {
		envLines = append(envLines, "PORT=4000")
		envLines = append(envLines, "NODE_ENV=development")
		if cfg.Frontend != "" && cfg.Frontend != "none" {
			if cfg.Frontend == "nextjs" {
				envLines = append(envLines, "NEXT_PUBLIC_API_URL=http://localhost:4000")
			} else {
				envLines = append(envLines, "VITE_API_URL=http://localhost:4000")
			}
		}
		envLines = append(envLines, "")
	}

	// 2. Base de datos
	db := strings.ToLower(cfg.Database)
	if db != "" && db != "none" {
		envLines = append(envLines, "# Database Connection")
		switch db {
		case "postgres":
			envLines = append(envLines, fmt.Sprintf("DATABASE_URL=\"postgresql://postgres:password@localhost:5432/%s?schema=public\"", cfg.ProjectName))
		case "mysql":
			envLines = append(envLines, fmt.Sprintf("DATABASE_URL=\"mysql://root:password@localhost:3306/%s\"", cfg.ProjectName))
		case "mongodb":
			envLines = append(envLines, fmt.Sprintf("DATABASE_URL=\"mongodb://localhost:27017/%s\"", cfg.ProjectName))
		case "sqlite":
			envLines = append(envLines, "DATABASE_URL=\"file:./local.db\"")
		}
		envLines = append(envLines, "")
	}

	// 3. Autenticación
	auth := strings.ToLower(cfg.Auth)
	if auth != "" && auth != "none" {
		envLines = append(envLines, "# Authentication")
		if strings.Contains(auth, "better") {
			envLines = append(envLines, "BETTER_AUTH_SECRET=\"supersecret-auth-key-change-in-production\"")
			envLines = append(envLines, "BETTER_AUTH_URL=\"http://localhost:4000\"")
		} else if strings.Contains(auth, "clerk") {
			envLines = append(envLines, "NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=\"pk_test_...\"")
			envLines = append(envLines, "CLERK_SECRET_KEY=\"sk_test_...\"")
		} else if strings.Contains(auth, "nextauth") {
			envLines = append(envLines, "NEXTAUTH_SECRET=\"supersecret-nextauth-key\"")
			envLines = append(envLines, "NEXTAUTH_URL=\"http://localhost:3000\"")
		}
		envLines = append(envLines, "")
	}

	content := strings.Join(envLines, "\n")
	// Si no existía un .env previo o si solo tenía placeholders, guardarlo
	if !v.Exists(".env") {
		v.WriteString(".env", content)
	}
	if !v.Exists(".env.example") {
		v.WriteString(".env.example", content)
	}

	return nil
}
