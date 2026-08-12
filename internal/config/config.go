package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

const CLIVersion = "v1.0"

type KokoConfig struct {
	Schema       string           `json:"$schema"`
	Project      ProjectInfo      `json:"project"`
	Architecture ArchitectureInfo `json:"architecture"`
	Stack        StackInfo        `json:"stack"`
	Features     FeaturesInfo     `json:"features"`
}

type ProjectInfo struct {
	Name       string    `json:"name"`
	CLIVersion string    `json:"cliVersion"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ArchitectureInfo struct {
	Layout         string `json:"layout"`         // "monorepo" o "standalone"
	PackageManager string `json:"packageManager"` // "pnpm", "npm", "bun", etc.
}

type StackInfo struct {
	Frontend *FrontendInfo `json:"frontend,omitempty"`
	Backend  *BackendInfo  `json:"backend,omitempty"`
	Database *DatabaseInfo `json:"database,omitempty"`
}

type FrontendInfo struct {
	Framework string   `json:"framework"`          // "next", "react" (Vite), "vue"
	Language  string   `json:"language"`           // "typescript", "javascript"
	Styling   string   `json:"styling"`            // "tailwindcss", "css-modules", "none"
	UILibrary string   `json:"uiLibrary,omitempty"` // "shadcn", "radix", etc. (Futura expansión)
	Icons     string   `json:"icons,omitempty"`     // "lucide", "react-icons", etc. (Futura expansión)
}

type BackendInfo struct {
	Framework    string   `json:"framework"`              // "express", "fiber", "hono", "fastapi"
	Language     string   `json:"language"`               // "typescript", "go", "python"
	Dependencies []string `json:"dependencies,omitempty"`  // ["zod", "cors", "dotenv"] (Para inyecciones)
}

type DatabaseInfo struct {
	Provider string `json:"provider"` // "postgres", "mysql", "mongodb", "none"
	ORM      string `json:"orm"`      // "prisma", "sqlx", "drizzle", "none"
}

type FeaturesInfo struct {
	Auth           *AuthInfo       `json:"auth,omitempty"`
	Infrastructure *Infrastructure `json:"infrastructure,omitempty"`
}

type AuthInfo struct {
	Provider string `json:"provider"` // "better-auth", "next-auth", "firebase"
	Status   string `json:"status"`   // "installed", "pending"
}

type Infrastructure struct {
	DockerCompose bool   `json:"dockerCompose"`
	CICD          string `json:"ciCd"` // "github-actions", "gitlab-ci", "none"
}

func GenerateConfig(targetDir string, scaffoldCfg scaffold.ScaffoldConfig) error {
	config := KokoConfig{
		Schema: "https://koko-cli.dev/schema.json",
		Project: ProjectInfo{
			Name:       scaffoldCfg.ProjectName,
			CLIVersion: CLIVersion,
			CreatedAt:  time.Now().UTC(),
		},
		Architecture: ArchitectureInfo{
			Layout:         determineLayout(scaffoldCfg.Recipe),
			PackageManager: determinePackageManager(scaffoldCfg.Recipe),
		},
	}

	switch scaffoldCfg.Recipe {
	case "saas":
		config.Stack.Frontend = &FrontendInfo{
			Framework: "next",
			Language:  "typescript",
			Styling:   "tailwindcss",
			UILibrary: "shadcn",
			Icons:     "lucide",
		}
		config.Stack.Backend = &BackendInfo{
			Framework: "next",
			Language:  "typescript",
		}
		config.Stack.Database = &DatabaseInfo{
			Provider: "postgres",
			ORM:      "drizzle",
		}
		config.Features.Auth = &AuthInfo{
			Provider: "better-auth",
			Status:   "installed",
		}

	case "pern":
		config.Stack.Frontend = &FrontendInfo{
			Framework: "react",
			Language:  "typescript",
			Styling:   "tailwindcss",
		}
		config.Stack.Backend = &BackendInfo{
			Framework: "express",
			Language:  "typescript",
		}
		config.Stack.Database = &DatabaseInfo{
			Provider: "postgres",
			ORM:      "prisma",
		}

	case "mern":
		config.Stack.Frontend = &FrontendInfo{
			Framework: "react",
			Language:  "typescript",
			Styling:   "tailwindcss",
		}
		config.Stack.Backend = &BackendInfo{
			Framework: "express",
			Language:  "typescript",
		}
		config.Stack.Database = &DatabaseInfo{
			Provider: "mongodb",
			ORM:      "mongoose",
		}

	case "fastapi_react":
		config.Stack.Frontend = &FrontendInfo{
			Framework: "react",
			Language:  "typescript",
			Styling:   "tailwindcss",
		}
		config.Stack.Backend = &BackendInfo{
			Framework: "fastapi",
			Language:  "python",
		}
	}

	var docker bool
	var ciCd string = "none"
	if scaffoldCfg.Recipe == "saas" {
		docker = true
		ciCd = "github-actions"
	} else if scaffoldCfg.Recipe == "pern" || scaffoldCfg.Recipe == "mern" {
		docker = true
	}

	config.Features.Infrastructure = &Infrastructure{
		DockerCompose: docker,
		CICD:          ciCd,
	}

	// Marshal JSON with indents
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(targetDir, "koko.config.json")
	return os.WriteFile(configFilePath, data, 0644)
}

func determineLayout(recipe string) string {
	if recipe == "saas" || recipe == "pern" || recipe == "mern" {
		return "monorepo"
	}
	return "standalone"
}

func determinePackageManager(recipe string) string {
	if recipe == "saas" || recipe == "pern" || recipe == "mern" {
		return "pnpm"
	}
	return "npm"
}
