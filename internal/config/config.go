package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold/handlers"
	"github.com/BlasVernazza06/koko-cli/internal/vfs"
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
	API      *APIInfo      `json:"api,omitempty"`
}

type APIInfo struct {
	Layer string `json:"layer"` // "trpc", "orpc", "none"
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
	Payments       *PaymentInfo    `json:"payments,omitempty"`
	Email          *EmailInfo      `json:"email,omitempty"`
}

type PaymentInfo struct {
	Provider string `json:"provider"` // "stripe", "polar"
}

type EmailInfo struct {
	Provider string `json:"provider"` // "resend", "brevo"
}

type AuthInfo struct {
	Provider string `json:"provider"` // "better-auth", "next-auth", "firebase"
	Status   string `json:"status"`   // "installed", "pending"
}

type Infrastructure struct {
	DockerCompose bool   `json:"dockerCompose"`
	CICD          string `json:"ciCd"` // "github-actions", "gitlab-ci", "none"
}

// BuildKokoConfig construye la estructura de configuración tipada a partir de ScaffoldConfig.
func BuildKokoConfig(scaffoldCfg scaffold.ScaffoldConfig) KokoConfig {
	normRecipe := handlers.NormalizeRecipe(scaffoldCfg.Recipe)

	pm := scaffoldCfg.PackageManager
	if pm == "" {
		pm = determinePackageManager(normRecipe)
	}

	layout := determineLayout(scaffoldCfg)

	config := KokoConfig{
		Schema: "https://koko-cli.dev/schema.json",
		Project: ProjectInfo{
			Name:       scaffoldCfg.ProjectName,
			CLIVersion: CLIVersion,
			CreatedAt:  time.Now().UTC(),
		},
		Architecture: ArchitectureInfo{
			Layout:         layout,
			PackageManager: pm,
		},
	}

	if scaffoldCfg.Recipe != "" {
		switch normRecipe {
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
			config.Features.Payments = &PaymentInfo{
				Provider: "stripe",
			}
			config.Features.Email = &EmailInfo{
				Provider: "resend",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "github-actions",
			}

		case "java_spring":
			config.Stack.Frontend = &FrontendInfo{
				Framework: "react",
				Language:  "typescript",
				Styling:   "tailwindcss",
				Icons:     "lucide",
			}
			config.Stack.Backend = &BackendInfo{
				Framework: "spring",
				Language:  "java",
			}
			config.Stack.Database = &DatabaseInfo{
				Provider: "postgres",
				ORM:      "jpa",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "github-actions",
			}

		case "enterprise_nestjs":
			config.Stack.Frontend = &FrontendInfo{
				Framework: "next",
				Language:  "typescript",
				Styling:   "tailwindcss",
				Icons:     "lucide",
			}
			config.Stack.Backend = &BackendInfo{
				Framework: "nestjs",
				Language:  "typescript",
			}
			config.Stack.Database = &DatabaseInfo{
				Provider: "postgres",
				ORM:      "prisma",
			}
			config.Features.Auth = &AuthInfo{
				Provider: "better-auth",
				Status:   "installed",
			}
			config.Features.Payments = &PaymentInfo{
				Provider: "stripe",
			}
			config.Features.Email = &EmailInfo{
				Provider: "resend",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "github-actions",
			}

		case "mern":
			config.Stack.Frontend = &FrontendInfo{
				Framework: "react",
				Language:  "typescript",
				Styling:   "tailwindcss",
				Icons:     "lucide",
			}
			config.Stack.Backend = &BackendInfo{
				Framework: "express",
				Language:  "typescript",
			}
			config.Stack.Database = &DatabaseInfo{
				Provider: "mongodb",
				ORM:      "mongoose",
			}
			config.Features.Auth = &AuthInfo{
				Provider: "jwt",
				Status:   "installed",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "none",
			}

		case "pern":
			config.Stack.Frontend = &FrontendInfo{
				Framework: "react",
				Language:  "typescript",
				Styling:   "tailwindcss",
				Icons:     "lucide",
			}
			config.Stack.Backend = &BackendInfo{
				Framework: "express",
				Language:  "typescript",
			}
			config.Stack.Database = &DatabaseInfo{
				Provider: "postgres",
				ORM:      "prisma",
			}
			config.Features.Auth = &AuthInfo{
				Provider: "jwt",
				Status:   "installed",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "none",
			}

		case "fastapi_react":
			config.Stack.Frontend = &FrontendInfo{
				Framework: "react",
				Language:  "typescript",
				Styling:   "tailwindcss",
				Icons:     "lucide",
			}
			config.Stack.Backend = &BackendInfo{
				Framework: "fastapi",
				Language:  "python",
			}
			config.Stack.Database = &DatabaseInfo{
				Provider: "postgres",
				ORM:      "sqlalchemy",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "none",
			}

		case "mobile_expo":
			config.Stack.Frontend = &FrontendInfo{
				Framework: "expo",
				Language:  "typescript",
				Styling:   "react-native",
				Icons:     "lucide-react-native",
			}
			config.Stack.Backend = &BackendInfo{
				Framework: "express",
				Language:  "typescript",
			}
			config.Stack.Database = &DatabaseInfo{
				Provider: "postgres",
				ORM:      "prisma",
			}
			config.Features.Auth = &AuthInfo{
				Provider: "jwt",
				Status:   "installed",
			}
			config.Features.Infrastructure = &Infrastructure{
				DockerCompose: true,
				CICD:          "github-actions",
			}
		}
	} else {
		// Manual Configuration mapping
		if scaffoldCfg.Frontend != "" && scaffoldCfg.Frontend != "none" {
			config.Stack.Frontend = &FrontendInfo{
				Framework: scaffoldCfg.Frontend,
				Language:  "typescript",
				Styling:   "tailwindcss",
			}
		}

		if scaffoldCfg.Backend != "" && scaffoldCfg.Backend != "none" {
			lang := "typescript"
			framework := scaffoldCfg.Backend
			if scaffoldCfg.Backend == "fastapi" {
				lang = "python"
			} else if scaffoldCfg.Backend == "go_chi" {
				lang = "go"
			} else if scaffoldCfg.Backend == "spring_boot" || scaffoldCfg.Backend == "java_spring" || scaffoldCfg.Backend == "spring" {
				framework = "spring_boot"
				lang = "java"
			} else if scaffoldCfg.Backend == "self" {
				framework = scaffoldCfg.Frontend
				lang = "typescript"
			}
			config.Stack.Backend = &BackendInfo{
				Framework: framework,
				Language:  lang,
			}
		}

		if scaffoldCfg.API != "" && scaffoldCfg.API != "none" {
			config.Stack.API = &APIInfo{
				Layer: scaffoldCfg.API,
			}
		}

		if scaffoldCfg.Database != "" && scaffoldCfg.Database != "none" {
			config.Stack.Database = &DatabaseInfo{
				Provider: scaffoldCfg.Database,
				ORM:      scaffoldCfg.ORM,
			}
		}

		if scaffoldCfg.Auth != "" && scaffoldCfg.Auth != "none" {
			config.Features.Auth = &AuthInfo{
				Provider: scaffoldCfg.Auth,
				Status:   "installed",
			}
		}

		addons := strings.ToLower(scaffoldCfg.Addons)
		docker := strings.Contains(addons, "docker")
		ciCd := "none"
		if strings.Contains(addons, "github_actions") || strings.Contains(addons, "cicd") {
			ciCd = "github-actions"
		}

		if strings.Contains(addons, "stripe") {
			config.Features.Payments = &PaymentInfo{Provider: "stripe"}
		} else if strings.Contains(addons, "polar") {
			config.Features.Payments = &PaymentInfo{Provider: "polar"}
		}

		if strings.Contains(addons, "resend") {
			config.Features.Email = &EmailInfo{Provider: "resend"}
		} else if strings.Contains(addons, "brevo") {
			config.Features.Email = &EmailInfo{Provider: "brevo"}
		}

		if config.Stack.Frontend != nil {
			if strings.Contains(addons, "shadcn") {
				config.Stack.Frontend.UILibrary = "shadcn"
			}
			if strings.Contains(addons, "lucide") {
				config.Stack.Frontend.Icons = "lucide"
			}
		}

		config.Features.Infrastructure = &Infrastructure{
			DockerCompose: docker,
			CICD:          ciCd,
		}
	}

	return config
}

// GenerateConfigToVFS genera koko.config.json directamente en el Virtual File System.
func GenerateConfigToVFS(v *vfs.VFS, scaffoldCfg scaffold.ScaffoldConfig) error {
	config := BuildKokoConfig(scaffoldCfg)
	return v.WriteJSON("koko.config.json", config, "  ")
}

// GenerateConfig genera el archivo directamente en disco (legacy / standalone).
func GenerateConfig(targetDir string, scaffoldCfg scaffold.ScaffoldConfig) error {
	config := BuildKokoConfig(scaffoldCfg)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(targetDir, "koko.config.json")
	return os.WriteFile(configFilePath, data, 0644)
}

func determineLayout(scaffoldCfg scaffold.ScaffoldConfig) string {
	if scaffoldCfg.Recipe != "" {
		return "monorepo"
	}
	if scaffoldCfg.Frontend != "" && scaffoldCfg.Frontend != "none" && scaffoldCfg.Backend != "" && scaffoldCfg.Backend != "none" && scaffoldCfg.Backend != "self" {
		return "monorepo"
	}
	return "standalone"
}

func determinePackageManager(recipe string) string {
	norm := handlers.NormalizeRecipe(recipe)
	switch norm {
	case "mern":
		return "npm"
	case "saas", "pern", "enterprise_nestjs", "mobile_expo", "java_spring", "fastapi_react":
		return "pnpm"
	default:
		return "pnpm"
	}
}
