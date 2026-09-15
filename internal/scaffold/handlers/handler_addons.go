package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// EvaluateAddons evalúa y mapea herramientas complementarias como Docker Compose y GitHub Actions CI.
func EvaluateAddons(rel string, config types.ScaffoldConfig) (string, bool) {
	addons := strings.ToLower(config.Addons)
	db := strings.ToLower(config.Database)
	orm := strings.ToLower(config.ORM)

	// Addon: Docker Compose
	hasDocker := strings.Contains(addons, "docker") || (db != "" && db != "none" && db != "sqlite")
	if hasDocker {
		if (db == "postgres" || db == "postgresql" || ((db == "" || db == "none") && strings.Contains(addons, "docker"))) && rel == "manual/docker/postgres/docker-compose.yml" {
			return "docker-compose.yml", true
		}
		if db == "mysql" && rel == "manual/docker/mysql/docker-compose.yml" {
			return "docker-compose.yml", true
		}
		if (db == "mongodb" || orm == "mongoose" || orm == "moongose") && rel == "manual/docker/mongoose/docker-compose.yml" {
			return "docker-compose.yml", true
		}
	}

	// Addon: GitHub Actions CI
	if rel == "manual/github/ci.yml" {
		if strings.Contains(addons, "github_actions") || strings.Contains(addons, "cicd") {
			return ".github/workflows/ci.yml", true
		}
	}

	// Addon: shadcn/ui -> packages/ui
	if strings.Contains(addons, "shadcn") {
		prefix := "manual/addons/shadcn/"
		if strings.HasPrefix(rel, prefix) {
			dest := "packages/ui/" + strings.TrimPrefix(rel, prefix)
			return dest, true
		}
	}

	return "", false
}
