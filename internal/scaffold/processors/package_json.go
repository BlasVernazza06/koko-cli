package processors

import (
	"fmt"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/catalog"
	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

type PackageJSON struct {
	Name            string                 `json:"name,omitempty"`
	Private         bool                   `json:"private,omitempty"`
	Version         string                 `json:"version,omitempty"`
	Workspaces      []string               `json:"workspaces,omitempty"`
	PackageManager  string                 `json:"packageManager,omitempty"`
	Scripts         map[string]string      `json:"scripts,omitempty"`
	Dependencies    map[string]interface{} `json:"dependencies,omitempty"`
	DevDependencies map[string]interface{} `json:"devDependencies,omitempty"`
}

type ProcessConfig struct {
	ProjectName    string
	PackageManager string
	Frontend       string
	Backend        string
	Database       string
	ORM            string
	Auth           string
	Addons         string
}

// AddDependency agrega o actualiza una dependencia en un objeto package.json en memoria,
// consultando la versión oficial registrada en el catálogo maestro.
func AddDependency(pkg map[string]interface{}, pkgName string, isDev bool) {
	field := "dependencies"
	if isDev {
		field = "devDependencies"
	}

	deps, ok := pkg[field].(map[string]interface{})
	if !ok {
		deps = make(map[string]interface{})
	}

	// Consultamos el catálogo maestro
	deps[pkgName] = catalog.GetVersion(pkgName)
	pkg[field] = deps
}

// ProcessPackageJSONs inspecciona y muta los package.json en el VFS en memoria
// para inyectar los scripts exactos y dependencias dinámicas según el stack.
func ProcessPackageJSONs(v *vfs.VFS, cfg ProcessConfig) error {
	updateRootPackageJSON(v, cfg)
	updateWorkspacePackageJSONs(v, cfg)
	return nil
}

func updateRootPackageJSON(v *vfs.VFS, cfg ProcessConfig) {
	if !v.Exists("package.json") {
		return
	}

	var pkg map[string]interface{}
	if err := v.ReadJSON("package.json", &pkg); err != nil {
		return
	}

	pkg["name"] = cfg.ProjectName

	scripts, ok := pkg["scripts"].(map[string]interface{})
	if !ok {
		scripts = make(map[string]interface{})
	}

	pm := strings.ToLower(cfg.PackageManager)
	if pm == "" {
		pm = "pnpm"
	}

	// 1. Scripts de desarrollo y build según Package Manager
	switch pm {
	case "bun":
		scripts["dev"] = "bun run --filter '*' dev"
		scripts["build"] = "bun run --filter '*' build"
		scripts["check-types"] = "bun run --filter '*' check-types"
	case "npm":
		scripts["dev"] = "npm run dev --workspaces --if-present"
		scripts["build"] = "npm run build --workspaces --if-present"
		scripts["check-types"] = "npm run check-types --workspaces --if-present"
	case "pnpm":
		fallthrough
	default:
		scripts["dev"] = "pnpm -r dev"
		scripts["build"] = "pnpm -r build"
		scripts["check-types"] = "pnpm -r check-types"
	}

	// 2. Scripts de base de datos a nivel raíz si existe DB
	orm := strings.ToLower(cfg.ORM)
	dbPkgName := fmt.Sprintf("@%s/db", cfg.ProjectName)

	if cfg.Database != "" && cfg.Database != "none" {
		filterCmd := func(script string) string {
			switch pm {
			case "bun":
				return fmt.Sprintf("bun run --filter %s %s", dbPkgName, script)
			case "npm":
				return fmt.Sprintf("npm run %s --workspace %s", script, dbPkgName)
			case "pnpm":
				fallthrough
			default:
				return fmt.Sprintf("pnpm --filter %s %s", dbPkgName, script)
			}
		}

		if orm == "drizzle" {
			scripts["db:generate"] = filterCmd("db:generate")
			scripts["db:push"] = filterCmd("db:push")
			scripts["db:studio"] = filterCmd("db:studio")
			scripts["db:migrate"] = filterCmd("db:migrate")
		} else if orm == "prisma" {
			scripts["db:generate"] = filterCmd("db:generate")
			scripts["db:migrate"] = filterCmd("db:migrate")
			scripts["db:studio"] = filterCmd("db:studio")
			scripts["db:push"] = filterCmd("db:push")
		}
	}

	// 3. Scripts de Docker si está activo
	if strings.Contains(cfg.Addons, "docker") {
		scripts["docker:up"] = "docker compose up -d"
		scripts["docker:down"] = "docker compose down"
		scripts["docker:logs"] = "docker compose logs -f"
	}

	pkg["scripts"] = scripts
	_ = v.WriteJSON("package.json", pkg, "  ")
}

func updateWorkspacePackageJSONs(v *vfs.VFS, cfg ProcessConfig) {
	// Actualizar nombres en apps/ y packages/
	workspacePaths := []struct {
		file   string
		suffix string
	}{
		{"apps/web/package.json", "web"},
		{"apps/api/package.json", "api"},
		{"packages/db/package.json", "db"},
		{"packages/config/package.json", "config"},
		{"packages/ui/package.json", "ui"},
		{"packages/eslint-config/package.json", "eslint-config"},
		{"packages/typescript-config/package.json", "typescript-config"},
	}

	for _, wp := range workspacePaths {
		if !v.Exists(wp.file) {
			continue
		}
		var pkg map[string]interface{}
		if err := v.ReadJSON(wp.file, &pkg); err != nil {
			continue
		}

		pkg["name"] = fmt.Sprintf("@%s/%s", cfg.ProjectName, wp.suffix)

		// Inyección dinámica de dependencias opcionales
		if wp.suffix == "db" && cfg.ORM == "drizzle" {
			if cfg.Database == "postgres" {
				AddDependency(pkg, "postgres", false)
			} else if cfg.Database == "mysql" {
				AddDependency(pkg, "mysql2", false)
			}
		}

		if wp.suffix == "api" && strings.Contains(cfg.Auth, "better") {
			AddDependency(pkg, "better-auth", false)
		}

		if wp.suffix == "web" && strings.Contains(cfg.Auth, "clerk") {
			AddDependency(pkg, "@clerk/nextjs", false)
		}

		_ = v.WriteJSON(wp.file, pkg, "  ")
	}
}
