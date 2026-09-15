package doctor

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/BlasVernazza06/koko-cli/internal/config"
)

// PackageJSON representa la estructura simplificada de un archivo package.json.
type PackageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// fileExists verifica si un archivo o directorio existe en el disco.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// readPackageJSON lee y deserializa un archivo package.json.
func readPackageJSON(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	if pkg.Dependencies == nil {
		pkg.Dependencies = make(map[string]string)
	}
	if pkg.DevDependencies == nil {
		pkg.DevDependencies = make(map[string]string)
	}

	return &pkg, nil
}

// hasDep revisa si una dependencia está presente en dependencies o devDependencies.
func hasDep(pkg *PackageJSON, name string) bool {
	if pkg == nil {
		return false
	}
	if _, ok := pkg.Dependencies[name]; ok {
		return true
	}
	if _, ok := pkg.DevDependencies[name]; ok {
		return true
	}
	return false
}

// AnalyzeProject escanea el sistema de archivos del proyecto y construye el KokoConfig detectado.
func AnalyzeProject(projectDir string) config.KokoConfig {
	// 1. Cargar manifiestos package.json de ubicaciones conocidas
	rootPkg, _ := readPackageJSON(filepath.Join(projectDir, "package.json"))
	webPkg, _ := readPackageJSON(filepath.Join(projectDir, "apps", "web", "package.json"))
	apiPkg, _ := readPackageJSON(filepath.Join(projectDir, "apps", "api", "package.json"))
	dbPkg, _ := readPackageJSON(filepath.Join(projectDir, "packages", "db", "package.json"))

	pkgs := []*PackageJSON{rootPkg, webPkg, apiPkg, dbPkg}

	// 2. Detectar Arquitectura y Package Manager
	arch := detectArchitecture(projectDir)

	detected := config.KokoConfig{
		Schema:       "https://koko-cli.dev/schema.json",
		Architecture: arch,
		Stack:        config.StackInfo{},
		Features:     config.FeaturesInfo{},
	}

	// 3. Escanear dependencias con TechCatalog
	frontendPkg := webPkg
	if frontendPkg == nil {
		frontendPkg = rootPkg
	}

	backendPkg := apiPkg
	if backendPkg == nil && arch.Layout == "standalone" {
		backendPkg = rootPkg
	}

	// Detectar Frontend
	detected.Stack.Frontend = detectFrontend(projectDir, frontendPkg)

	// Detectar Backend
	detected.Stack.Backend = detectBackend(projectDir, backendPkg)

	// Detectar Database
	detected.Stack.Database = detectDatabase(projectDir, pkgs...)

	// Detectar API Layer
	detected.Stack.API = detectAPI(projectDir, pkgs...)

	// Detectar Features (Auth, Infra, Payments, Email)
	detected.Features = detectFeatures(projectDir, pkgs...)

	return detected
}

func detectArchitecture(root string) config.ArchitectureInfo {
	layout := "standalone"
	if fileExists(filepath.Join(root, "apps")) ||
		fileExists(filepath.Join(root, "packages")) ||
		fileExists(filepath.Join(root, "pnpm-workspace.yaml")) ||
		fileExists(filepath.Join(root, "turbo.json")) {
		layout = "monorepo"
	}

	pm := "pnpm"
	if fileExists(filepath.Join(root, "bun.lockb")) || fileExists(filepath.Join(root, "bun.lock")) {
		pm = "bun"
	} else if fileExists(filepath.Join(root, "pnpm-lock.yaml")) {
		pm = "pnpm"
	} else if fileExists(filepath.Join(root, "yarn.lock")) {
		pm = "yarn"
	} else if fileExists(filepath.Join(root, "package-lock.json")) {
		pm = "npm"
	}

	return config.ArchitectureInfo{
		Layout:         layout,
		PackageManager: pm,
	}
}

func detectFrontend(projectDir string, pkg *PackageJSON) *config.FrontendInfo {
	if pkg == nil {
		return nil
	}

	framework := ""
	styling := "none"
	uiLib := ""
	icons := ""
	lang := "javascript"

	// Comprobar presencia de TypeScript
	if fileExists(filepath.Join(projectDir, "tsconfig.json")) ||
		fileExists(filepath.Join(projectDir, "apps", "web", "tsconfig.json")) ||
		hasDep(pkg, "typescript") {
		lang = "typescript"
	}

	// Comprobar shadcn UI por archivo components.json o packages/ui
	if fileExists(filepath.Join(projectDir, "components.json")) ||
		fileExists(filepath.Join(projectDir, "apps", "web", "components.json")) ||
		fileExists(filepath.Join(projectDir, "packages", "ui")) {
		uiLib = "shadcn"
	}

	// Consultar catálogo sobre dependencias del package
	allDeps := getCombinedDeps(pkg)
	for dep := range allDeps {
		if rule, ok := TechCatalog[dep]; ok {
			switch rule.Section {
			case SectionFrontend:
				if framework == "" || framework == "react" {
					framework = rule.Value
				}
				if rule.Lang != "" {
					lang = rule.Lang
				}
			case SectionFrontendStyle:
				styling = rule.Value
			case SectionFrontendUI:
				if uiLib == "" {
					uiLib = rule.Value
				}
			case SectionFrontendIcons:
				icons = rule.Value
			}
		}
	}

	// Heurística de archivos de configuración
	if framework == "" {
		if fileExists(filepath.Join(projectDir, "next.config.js")) ||
			fileExists(filepath.Join(projectDir, "next.config.mjs")) ||
			fileExists(filepath.Join(projectDir, "apps", "web", "next.config.mjs")) ||
			fileExists(filepath.Join(projectDir, "apps", "web", "next.config.js")) {
			framework = "next"
			lang = "typescript"
		}
	}

	if styling == "none" && (fileExists(filepath.Join(projectDir, "tailwind.config.ts")) ||
		fileExists(filepath.Join(projectDir, "tailwind.config.js")) ||
		fileExists(filepath.Join(projectDir, "apps", "web", "tailwind.config.ts"))) {
		styling = "tailwindcss"
	}

	if framework == "" {
		return nil
	}

	return &config.FrontendInfo{
		Framework: framework,
		Language:  lang,
		Styling:   styling,
		UILibrary: uiLib,
		Icons:     icons,
	}
}

func detectBackend(projectDir string, pkg *PackageJSON) *config.BackendInfo {
	// A. Backend en Node/TS mediante catálogo
	if pkg != nil {
		allDeps := getCombinedDeps(pkg)
		var backendDeps []string
		framework := ""
		lang := "typescript"

		for dep := range allDeps {
			if rule, ok := TechCatalog[dep]; ok {
				if rule.Section == SectionBackend {
					framework = rule.Value
					if rule.Lang != "" {
						lang = rule.Lang
					}
				} else if rule.Section == SectionBackendDep {
					backendDeps = append(backendDeps, rule.Value)
				}
			}
		}

		if framework != "" {
			return &config.BackendInfo{
				Framework:    framework,
				Language:     lang,
				Dependencies: backendDeps,
			}
		}
	}

	// B. Backend en Python (FastAPI)
	if fileExists(filepath.Join(projectDir, "apps", "api", "requirements.txt")) ||
		fileExists(filepath.Join(projectDir, "requirements.txt")) ||
		fileExists(filepath.Join(projectDir, "pyproject.toml")) {
		return &config.BackendInfo{Framework: "fastapi", Language: "python"}
	}

	// C. Backend en Go (Chi / Fiber)
	if fileExists(filepath.Join(projectDir, "apps", "api", "go.mod")) ||
		fileExists(filepath.Join(projectDir, "go.mod")) {
		return &config.BackendInfo{Framework: "go_chi", Language: "go"}
	}

	// D. Backend en Java (Spring Boot)
	if fileExists(filepath.Join(projectDir, "apps", "api", "pom.xml")) ||
		fileExists(filepath.Join(projectDir, "pom.xml")) ||
		fileExists(filepath.Join(projectDir, "build.gradle")) {
		return &config.BackendInfo{Framework: "spring_boot", Language: "java"}
	}

	return nil
}

func detectDatabase(projectDir string, pkgs ...*PackageJSON) *config.DatabaseInfo {
	orm := ""
	provider := ""

	// 1. Detección por archivos de esquema
	if fileExists(filepath.Join(projectDir, "packages", "db", "prisma", "schema.prisma")) ||
		fileExists(filepath.Join(projectDir, "prisma", "schema.prisma")) {
		orm = "prisma"
	}
	if fileExists(filepath.Join(projectDir, "packages", "db", "drizzle.config.ts")) ||
		fileExists(filepath.Join(projectDir, "drizzle.config.ts")) {
		orm = "drizzle"
	}

	// 2. Detección por catálogo de dependencias
	for _, pkg := range pkgs {
		if pkg == nil {
			continue
		}
		for dep := range getCombinedDeps(pkg) {
			if rule, ok := TechCatalog[dep]; ok {
				if rule.Section == SectionDatabaseORM && orm == "" {
					orm = rule.Value
					if rule.Value == "mongoose" {
						provider = "mongodb"
					}
				} else if rule.Section == SectionDatabaseDriver && provider == "" {
					provider = rule.Value
				}
			}
		}
	}

	if orm == "" && provider == "" {
		return nil
	}

	if orm != "" && provider == "" {
		provider = "postgres"
	}
	if orm == "" && provider != "" {
		orm = "none"
	}

	return &config.DatabaseInfo{
		Provider: provider,
		ORM:      orm,
	}
}

func detectAPI(projectDir string, pkgs ...*PackageJSON) *config.APIInfo {
	if fileExists(filepath.Join(projectDir, "packages", "api")) {
		for _, pkg := range pkgs {
			if pkg != nil && hasDep(pkg, "@orpc/server") {
				return &config.APIInfo{Layer: "orpc"}
			}
		}
		return &config.APIInfo{Layer: "trpc"}
	}

	for _, pkg := range pkgs {
		if pkg == nil {
			continue
		}
		for dep := range getCombinedDeps(pkg) {
			if rule, ok := TechCatalog[dep]; ok && rule.Section == SectionAPI {
				return &config.APIInfo{Layer: rule.Value}
			}
		}
	}

	return nil
}

func detectFeatures(projectDir string, pkgs ...*PackageJSON) config.FeaturesInfo {
	features := config.FeaturesInfo{}

	for _, pkg := range pkgs {
		if pkg == nil {
			continue
		}
		for dep := range getCombinedDeps(pkg) {
			if rule, ok := TechCatalog[dep]; ok {
				switch rule.Section {
				case SectionAuth:
					if features.Auth == nil {
						features.Auth = &config.AuthInfo{Provider: rule.Value, Status: "installed"}
					}
				case SectionPayments:
					if features.Payments == nil {
						features.Payments = &config.PaymentInfo{Provider: rule.Value}
					}
				case SectionEmail:
					if features.Email == nil {
						features.Email = &config.EmailInfo{Provider: rule.Value}
					}
				}
			}
		}
	}

	// Infraestructura (Docker y CI/CD)
	hasDocker := fileExists(filepath.Join(projectDir, "docker-compose.yml")) ||
		fileExists(filepath.Join(projectDir, "docker-compose.yaml")) ||
		fileExists(filepath.Join(projectDir, "compose.yaml"))

	ciCd := "none"
	if fileExists(filepath.Join(projectDir, ".github", "workflows")) {
		ciCd = "github-actions"
	} else if fileExists(filepath.Join(projectDir, ".gitlab-ci.yml")) {
		ciCd = "gitlab-ci"
	}

	if hasDocker || ciCd != "none" {
		features.Infrastructure = &config.Infrastructure{
			DockerCompose: hasDocker,
			CICD:          ciCd,
		}
	}

	return features
}

func getCombinedDeps(pkg *PackageJSON) map[string]bool {
	combined := make(map[string]bool)
	if pkg == nil {
		return combined
	}
	for k := range pkg.Dependencies {
		combined[k] = true
	}
	for k := range pkg.DevDependencies {
		combined[k] = true
	}
	return combined
}
