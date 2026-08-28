package compatibility

import (
	"fmt"
	"strings"

	"github.com/BlasVernazza06/koko-cli/cmd/views"
	"github.com/BlasVernazza06/koko-cli/internal/errors"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

// Step indices for interactive TUI wizard
const (
	StepFrontend = iota
	StepBackend
	StepPackageManager
	StepDatabase
	StepORM
	StepAddons
	StepGit
)

// Classification Predicates

// IsFullstackFrontend returns true if the frontend framework has native server capabilities (SSR / Server Actions).
func IsFullstackFrontend(frontend string) bool {
	f := strings.ToLower(frontend)
	return f == "nextjs" || f == "nuxt"
}

// IsClientSPA returns true if the frontend is a pure client-side Single Page Application.
func IsClientSPA(frontend string) bool {
	f := strings.ToLower(frontend)
	return f == "react" || f == "svelte"
}

// IsNodeBackend returns true if the backend framework belongs to Node.js / TypeScript.
func IsNodeBackend(backend string) bool {
	b := strings.ToLower(backend)
	return b == "express" || b == "nestjs" || b == "hono"
}

// IsPythonBackend returns true if the backend is Python-based.
func IsPythonBackend(backend string) bool {
	return strings.ToLower(backend) == "fastapi"
}

// IsGoBackend returns true if the backend is Go-based.
func IsGoBackend(backend string) bool {
	return strings.ToLower(backend) == "go_chi"
}

// IsNodeEcosystem returns true if the project contains a Node.js/TS server runtime.
func IsNodeEcosystem(frontend, backend string) bool {
	if IsNodeBackend(backend) {
		return true
	}
	if (backend == "none" || backend == "") && IsFullstackFrontend(frontend) {
		return true
	}
	return false
}

// IsSQLDatabase returns true if the database is a relational SQL engine.
func IsSQLDatabase(db string) bool {
	d := strings.ToLower(db)
	return d == "postgres" || d == "mysql" || d == "sqlite"
}

// IsNoSQLDatabase returns true if the database is document-based / NoSQL.
func IsNoSQLDatabase(db string) bool {
	return strings.ToLower(db) == "mongodb"
}

// BaseOptions returns default template options for a given step.
func BaseOptions(stepIdx int) []views.SelectOption {
	switch stepIdx {
	case StepFrontend:
		return []views.SelectOption{
			{Value: "nextjs", Label: "Next.js", Hint: "React framework with SSR & Server Components"},
			{Value: "react", Label: "React + Vite", Hint: "Ultra-fast Single Page Application"},
			{Value: "nuxt", Label: "Nuxt", Hint: "Vue full-stack framework with Nitro engine"},
			{Value: "svelte", Label: "Svelte", Hint: "Cybernetically enhanced web apps"},
			{Value: "none", Label: "None", Hint: "Backend / REST API only"},
		}
	case StepBackend:
		return []views.SelectOption{
			{Value: "express", Label: "Node.js / Express", Hint: "Lightweight REST API with TypeScript"},
			{Value: "fastapi", Label: "Python / FastAPI", Hint: "Async framework with Pydantic v2 validation"},
			{Value: "go_chi", Label: "Go / Chi Router", Hint: "High performance with strict types"},
			{Value: "nestjs", Label: "NestJS", Hint: "Enterprise modular architecture with TypeScript"},
			{Value: "hono", Label: "Hono", Hint: "Ultrafast web framework"},
			{Value: "none", Label: "None", Hint: "No dedicated backend (Fullstack Server Actions or BaaS)"},
		}
	case StepPackageManager:
		return []views.SelectOption{
			{Value: "pnpm", Label: "PNPM", Hint: "Fast and disk space efficient (Recommended)"},
			{Value: "npm", Label: "NPM", Hint: "Standard Node package manager"},
			{Value: "bun", Label: "Bun", Hint: "All-in-one JavaScript runtime & package manager"},
		}
	case StepDatabase:
		return []views.SelectOption{
			{Value: "postgres", Label: "PostgreSQL", Hint: "Standard relational database with Docker"},
			{Value: "mongodb", Label: "MongoDB", Hint: "NoSQL document database"},
			{Value: "mysql", Label: "MySQL / MariaDB", Hint: "Traditional SQL database"},
			{Value: "sqlite", Label: "SQLite", Hint: "Embedded lightweight database"},
			{Value: "none", Label: "None", Hint: "No database persistence"},
		}
	case StepORM:
		return []views.SelectOption{
			{Value: "drizzle", Label: "Drizzle ORM", Hint: "Lightweight, type-safe with native SQL support"},
			{Value: "prisma", Label: "Prisma", Hint: "Next-gen ORM with auto type generation"},
			{Value: "moongose", Label: "Mongoose", Hint: "Elegant object modeling tool for MongoDB"},
			{Value: "sqlalchemy", Label: "SQLAlchemy / SQLModel", Hint: "Standard ORM for Python"},
			{Value: "gorm", Label: "GORM", Hint: "Feature-rich ORM for Go"},
			{Value: "none", Label: "None / Raw SQL", Hint: "Direct driver connection without ORM"},
		}
	case StepAddons:
		return []views.SelectOption{
			{Value: "docker_cicd", Label: "Docker Compose + GitHub Actions", Hint: "Full containerization & CI/CD workflow"},
			{Value: "docker", Label: "Docker Compose", Hint: "Local containerized services"},
			{Value: "github_actions", Label: "GitHub Actions CI", Hint: "Automated linting and test workflows"},
			{Value: "none", Label: "None", Hint: "No extra tooling"},
		}
	case StepGit:
		return []views.SelectOption{
			{Value: "yes", Label: "Yes", Hint: "Initialize a new Git repository (git init)"},
			{Value: "no", Label: "No", Hint: "Skip Git repository initialization"},
		}
	default:
		return nil
	}
}

// GetStepOptions returns options for a given step evaluated dynamically against current selections.
func GetStepOptions(stepIdx int, currentSelections []views.SelectOption) []views.SelectOption {
	base := BaseOptions(stepIdx)
	options := make([]views.SelectOption, len(base))
	copy(options, base)

	var frontend, backend, db string
	if len(currentSelections) > StepFrontend {
		frontend = currentSelections[StepFrontend].Value
	}
	if len(currentSelections) > StepBackend {
		backend = currentSelections[StepBackend].Value
	}
	if len(currentSelections) > StepDatabase {
		db = currentSelections[StepDatabase].Value
	}

	switch stepIdx {
	case StepBackend:
		// Rule 1: If Frontend is "none", Backend cannot be "none"
		if frontend == "none" {
			for i := range options {
				if options[i].Value == "none" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Project cannot have both Frontend and Backend as None"
				}
			}
		}

	case StepPackageManager:
		// If pure Go or pure Python project without JS/TS frontend
		if frontend == "none" {
			if backend == "go_chi" {
				return []views.SelectOption{
					{Value: "go_mod", Label: "Go Modules", Hint: "Native Go dependency management"},
				}
			}
			if backend == "fastapi" {
				return []views.SelectOption{
					{Value: "pip", Label: "pip / requirements.txt", Hint: "Standard Python package installer"},
					{Value: "uv", Label: "uv", Hint: "Extremely fast Python package installer"},
				}
			}
		}

	case StepDatabase:
		// Rule 2: Client-side SPAs (React + Vite, Svelte) without a backend cannot directly connect to DB servers
		if backend == "none" && IsClientSPA(frontend) {
			for i := range options {
				if options[i].Value != "none" {
					options[i].Disabled = true
					options[i].DisabledReason = fmt.Sprintf("Incompatible: Client-side SPA (%s) cannot connect directly to DBs without a backend", frontend)
				}
			}
		}

	case StepORM:
		// Rule 3: Database vs ORM
		if db == "none" {
			for i := range options {
				if options[i].Value != "none" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Requires a database to be selected"
				}
			}
			return options
		}

		isNode := IsNodeEcosystem(frontend, backend)
		isPython := IsPythonBackend(backend)
		isGo := IsGoBackend(backend)

		for i := range options {
			val := options[i].Value
			switch val {
			case "moongose":
				if !IsNoSQLDatabase(db) {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Mongoose is exclusively for MongoDB"
				} else if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Mongoose is a Node.js/TS package"
				}

			case "drizzle":
				if IsNoSQLDatabase(db) {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Drizzle is a SQL ORM (incompatible with MongoDB)"
				} else if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Drizzle is a TypeScript ORM"
				}

			case "prisma":
				if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Prisma is a Node.js/TS ORM"
				}

			case "sqlalchemy":
				if IsNoSQLDatabase(db) {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: SQLAlchemy is a SQL ORM (incompatible with MongoDB)"
				} else if !isPython {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: SQLAlchemy is a Python-only ORM"
				}

			case "gorm":
				if IsNoSQLDatabase(db) {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: GORM is a SQL ORM (incompatible with MongoDB)"
				} else if !isGo {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: GORM is a Go-only ORM"
				}

			case "none":
				// Always available
			}
		}
	}

	return options
}

// ValidateConfig validates the complete scaffold configuration against all compatibility rules.
func ValidateConfig(cfg scaffold.ScaffoldConfig) error {
	frontend := strings.ToLower(cfg.Frontend)
	backend := strings.ToLower(cfg.Backend)
	pm := strings.ToLower(cfg.PackageManager)
	db := strings.ToLower(cfg.Database)
	orm := strings.ToLower(cfg.ORM)
	auth := strings.ToLower(cfg.Auth)

	// 1. Both frontend and backend empty/none
	if (frontend == "none" || frontend == "") && (backend == "none" || backend == "") {
		return errors.NewValidationError(
			"No se puede crear un proyecto con Frontend y Backend configurados como 'none'",
			"Selecciona al menos un Frontend o un Backend para tu proyecto",
		)
	}

	// 2. Client-side SPA with no backend connecting to database
	if (backend == "none" || backend == "") && IsClientSPA(frontend) {
		if db != "none" && db != "" {
			return errors.NewValidationError(
				fmt.Sprintf("Una aplicación SPA cliente (%s) sin backend no puede conectarse directamente a la base de datos '%s'", frontend, db),
				"Agrega un backend (como Express, Hono, FastAPI o Go Chi) o usa un framework fullstack (como Next.js o Nuxt)",
			)
		}
		if orm != "none" && orm != "" {
			return errors.NewValidationError(
				fmt.Sprintf("Una aplicación SPA cliente (%s) sin backend no puede utilizar el ORM '%s'", frontend, orm),
				"Agrega un backend para gestionar la persistencia y ORM",
			)
		}
	}

	// 3. Database vs ORM consistency
	if (db == "none" || db == "") && orm != "none" && orm != "" {
		return errors.NewValidationError(
			fmt.Sprintf("Se seleccionó el ORM '%s' pero la base de datos está marcada como 'none'", orm),
			"Selecciona una base de datos compatible o marca el ORM como 'none'",
		)
	}

	// 4. SQL vs NoSQL ORMs
	if IsNoSQLDatabase(db) {
		if orm == "drizzle" || orm == "sqlalchemy" || orm == "gorm" {
			return errors.NewValidationError(
				fmt.Sprintf("'%s' es un ORM relacional SQL y no es compatible con MongoDB", orm),
				"Para MongoDB utiliza 'Mongoose', 'Prisma' o 'none'",
			)
		}
	}

	if IsSQLDatabase(db) {
		if orm == "moongose" || orm == "mongoose" {
			return errors.NewValidationError(
				fmt.Sprintf("Mongoose es exclusivo para MongoDB y no es compatible con '%s'", db),
				"Para bases de datos relacionales SQL utiliza 'Drizzle', 'Prisma', 'SQLAlchemy' (Python) o 'GORM' (Go)",
			)
		}
	}

	// 5. Runtime & Ecosystem vs ORM
	isNode := IsNodeEcosystem(frontend, backend)
	isPython := IsPythonBackend(backend)
	isGo := IsGoBackend(backend)

	if isPython {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "gorm" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el backend Python / FastAPI", orm),
				"Para Python utiliza 'SQLAlchemy' o 'none'",
			)
		}
	}

	if isGo {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "sqlalchemy" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el backend Go", orm),
				"Para Go utiliza 'GORM' o 'none'",
			)
		}
	}

	if isNode {
		if orm == "sqlalchemy" || orm == "gorm" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el ecosistema Node.js / TypeScript", orm),
				"Para TypeScript utiliza 'Drizzle', 'Prisma' o 'Mongoose'",
			)
		}
	}

	// 6. Package Manager compatibility
	if frontend == "none" {
		if isGo && (pm == "pnpm" || pm == "npm" || pm == "bun") {
			return errors.NewValidationError(
				fmt.Sprintf("Un proyecto exclusivo de Go no utiliza el gestor de paquetes '%s'", pm),
				"Utiliza 'go_mod' para proyectos en Go",
			)
		}
		if isPython && (pm == "pnpm" || pm == "npm" || pm == "bun") {
			return errors.NewValidationError(
				fmt.Sprintf("Un proyecto exclusivo de Python no utiliza el gestor de paquetes '%s'", pm),
				"Utiliza 'pip' o 'uv' para proyectos en Python",
			)
		}
	}

	// 7. Auth Provider compatibility
	if auth != "" && auth != "none" {
		if auth == "better-auth" {
			if !isNode {
				return errors.NewValidationError(
					"Better Auth requiere un entorno Node.js / TypeScript (Next.js, Express, Hono)",
					"Agrega un frontend o backend compatible con TypeScript o cambia el proveedor de autenticación",
				)
			}
		} else if auth == "next-auth" {
			if frontend != "nextjs" {
				return errors.NewValidationError(
					"NextAuth.js es exclusivo para proyectos con frontend Next.js",
					"Selecciona Next.js como frontend o utiliza Better Auth",
				)
			}
		}
	}

	return nil
}
