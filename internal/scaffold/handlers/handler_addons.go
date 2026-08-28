package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateAddons evalúa y mapea herramientas complementarias como Docker Compose y GitHub Actions CI.
func EvaluateAddons(rel string, config types.ScaffoldConfig) (string, bool) {
	addons := strings.ToLower(config.Addons)

	// Addon: Docker Compose
	if rel == "manual/docker/docker-compose.yml" {
		if strings.Contains(addons, "docker") || (config.Database != "" && config.Database != "none" && config.Database != "sqlite") {
			return "docker-compose.yml", true
		}
	}

	// Addon: GitHub Actions CI
	if rel == "manual/github/ci.yml" {
		if strings.Contains(addons, "github_actions") || strings.Contains(addons, "cicd") {
			return ".github/workflows/ci.yml", true
		}
	}

	return "", false
}
