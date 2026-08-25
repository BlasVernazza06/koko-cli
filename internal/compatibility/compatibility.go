package compatibility

import (
	"fmt"
	"strings"

	"github.com/BlasVernazza06/koko-cli/cmd/views"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

// Step indices
const (
	StepFrontend = iota
	StepBackend
	StepPackageManager
	StepDatabase
	StepORM
	StepAddons
	StepGit
)

// BaseOptions returns default template options for a given step.
func BaseOptions(stepIdx int) []views.SelectOption {
	switch stepIdx {
	case StepFrontend:
		return []views.SelectOption{
			{Value: "nextjs", Label: "Next.js", Hint: "React framework with SSR & Server Components"},
			{Value: "react", Label: "React + Vite", Hint: "Ultra-fast Single Page Application"},
			{Value: "nuxt", Label: "Nuxt", Hint: "Vue full-stack framework"},
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
			{Value: "none", Label: "None", Hint: "No dedicated backend (Server Actions or BaaS)"},
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
		if backend == "none" && (frontend == "react" || frontend == "svelte") {
			for i := range options {
				if options[i].Value == "postgres" || options[i].Value == "mongodb" || options[i].Value == "mysql" || options[i].Value == "sqlite" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Client-side SPAs cannot connect directly to DBs without a backend"
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

		isNodeOrTS := isNodeEcosystem(frontend, backend)
		isPython := backend == "fastapi"
		isGo := backend == "go_chi"

		for i := range options {
			val := options[i].Value
			switch val {
			case "moongose":
				if db != "mongodb" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Mongoose is exclusively for MongoDB"
				} else if !isNodeOrTS {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Mongoose is a Node.js/TS package"
				}

			case "drizzle":
				if db == "mongodb" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Drizzle is a SQL ORM (incompatible with MongoDB)"
				} else if !isNodeOrTS {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Drizzle is a TypeScript ORM"
				}

			case "prisma":
				if !isNodeOrTS {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Prisma is a Node.js/TS ORM"
				}

			case "sqlalchemy":
				if db == "mongodb" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: SQLAlchemy is a SQL ORM (incompatible with MongoDB)"
				} else if !isPython {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: SQLAlchemy is a Python-only ORM"
				}

			case "gorm":
				if db == "mongodb" {
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

func isNodeEcosystem(frontend, backend string) bool {
	if backend == "express" || backend == "nestjs" || backend == "hono" {
		return true
	}
	if backend == "none" && (frontend == "nextjs" || frontend == "nuxt") {
		return true
	}
	return false
}

// ValidateConfig validates the complete scaffold configuration against all compatibility rules.
func ValidateConfig(cfg scaffold.ScaffoldConfig) error {
	frontend := strings.ToLower(cfg.Frontend)
	backend := strings.ToLower(cfg.Backend)
	db := strings.ToLower(cfg.Database)
	orm := strings.ToLower(cfg.ORM)

	// 1. Both frontend and backend empty/none
	if (frontend == "none" || frontend == "") && (backend == "none" || backend == "") {
		return fmt.Errorf("invalid configuration: project cannot have both Frontend and Backend as 'none'")
	}

	// 2. Client-side SPA with no backend connecting to database
	if (backend == "none" || backend == "") && (frontend == "react" || frontend == "svelte") {
		if db != "none" && db != "" {
			return fmt.Errorf("incompatible configuration: client-side SPA (%s) with no backend cannot connect directly to database '%s'", frontend, db)
		}
	}

	// 3. Database vs ORM
	if (db == "none" || db == "") && orm != "none" && orm != "" {
		return fmt.Errorf("incompatible configuration: ORM '%s' selected but database is 'none'", orm)
	}

	if db == "mongodb" {
		if orm == "drizzle" || orm == "sqlalchemy" || orm == "gorm" {
			return fmt.Errorf("incompatible configuration: '%s' is a SQL ORM and cannot be used with MongoDB", orm)
		}
	}

	if db == "postgres" || db == "mysql" || db == "sqlite" {
		if orm == "moongose" || orm == "mongoose" {
			return fmt.Errorf("incompatible configuration: Mongoose is a MongoDB ODM and cannot be used with SQL database '%s'", db)
		}
	}

	// 4. Runtime / Ecosystem vs ORM
	isNode := isNodeEcosystem(frontend, backend)
	isPython := backend == "fastapi"
	isGo := backend == "go_chi"

	if isPython {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "gorm" {
			return fmt.Errorf("incompatible configuration: ORM '%s' is not supported in Python / FastAPI backend", orm)
		}
	}

	if isGo {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "sqlalchemy" {
			return fmt.Errorf("incompatible configuration: ORM '%s' is not supported in Go backend", orm)
		}
	}

	if isNode {
		if orm == "sqlalchemy" || orm == "gorm" {
			return fmt.Errorf("incompatible configuration: ORM '%s' is not supported in Node.js / TypeScript backend", orm)
		}
	}

	return nil
}
