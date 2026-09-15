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
	StepAPI
	StepPackageManager
	StepDatabase
	StepORM
	StepAuth
	StepAddons
	StepGit
)

// Classification Predicates

// IsFullstackFrontend returns true if the frontend framework has native server capabilities (SSR / Server Actions / Endpoints).
func IsFullstackFrontend(frontend string) bool {
	f := strings.ToLower(frontend)
	return f == "nextjs" || f == "nuxt" || f == "svelte" || f == "astro"
}

// IsReactFrontend returns true if the frontend belongs to the React ecosystem (Next.js, React + Vite, React Native).
func IsReactFrontend(frontend string) bool {
	f := strings.ToLower(frontend)
	return f == "nextjs" || f == "react" || f == "native" || f == "react_native"
}

// IsClientSPA returns true if the frontend is a pure client-side Single Page Application or mobile app.
func IsClientSPA(frontend string) bool {
	f := strings.ToLower(frontend)
	return f == "react" || f == "native" || f == "react_native"
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

// IsJavaBackend returns true if the backend is Java-based.
func IsJavaBackend(backend string) bool {
	b := strings.ToLower(backend)
	return b == "spring_boot" || b == "java_spring" || b == "spring"
}

// IsNodeEcosystem returns true if the project contains a Node.js/TS server runtime.
func IsNodeEcosystem(frontend, backend string) bool {
	if IsNodeBackend(backend) {
		return true
	}
	if backend == "self" && IsFullstackFrontend(frontend) {
		return true
	}
	if (backend == "none" || backend == "") && (IsFullstackFrontend(frontend) || IsClientSPA(frontend)) {
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
			{Value: "astro", Label: "Astro", Hint: "Content-driven web apps with islands architecture"},
			{Value: "native", Label: "React Native / Expo", Hint: "Universal cross-platform mobile & web apps"},
			{Value: "none", Label: "None", Hint: "Backend / REST API only"},
		}
	case StepBackend:
		return []views.SelectOption{
			{Value: "express", Label: "Node.js / Express", Hint: "Lightweight REST API with TypeScript"},
			{Value: "hono", Label: "Hono", Hint: "Ultrafast multi-runtime web framework"},
			{Value: "fastapi", Label: "Python / FastAPI", Hint: "Async framework with Pydantic v2 validation"},
			{Value: "go_chi", Label: "Go / Chi Router", Hint: "High performance with strict types"},
			{Value: "spring_boot", Label: "Java / Spring Boot", Hint: "Enterprise robust backend with Spring Boot 3 & Maven"},
			{Value: "nestjs", Label: "NestJS", Hint: "Enterprise modular architecture with TypeScript"},
			{Value: "self", Label: "Self (Fullstack Framework)", Hint: "Uses frontend server capabilities (SSR / Server Actions / API Routes)"},
			{Value: "none", Label: "None", Hint: "No dedicated backend (Client-only or BaaS)"},
		}
	case StepAPI:
		return []views.SelectOption{
			{Value: "trpc", Label: "tRPC", Hint: "End-to-end typesafe APIs for TypeScript fullstack & monorepos"},
			{Value: "orpc", Label: "oRPC", Hint: "OpenAPI + RPC typesafe contract & client generator"},
			{Value: "none", Label: "None / REST API", Hint: "Standard HTTP endpoints without an RPC abstraction"},
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
			{Value: "jpa", Label: "Spring Data JPA / Hibernate", Hint: "Standard persistence layer for Java / Spring"},
			{Value: "none", Label: "None / Raw SQL", Hint: "Direct driver connection without ORM"},
		}
	case StepAuth:
		return []views.SelectOption{
			{Value: "better-auth", Label: "Better Auth", Hint: "Comprehensive TypeScript auth framework (Recommended)"},
			{Value: "clerk", Label: "Clerk", Hint: "Complete user management & authentication platform"},
			{Value: "next-auth", Label: "NextAuth / Auth.js", Hint: "Authentication solution for Next.js apps"},
			{Value: "none", Label: "None", Hint: "No authentication layer"},
		}
	case StepAddons:
		return []views.SelectOption{
			{Label: "── UI & Components ──", IsHeader: true},
			{Value: "shadcn", Label: "shadcn/ui", Hint: "Re-usable component library built on Radix UI & Tailwind in packages/ui"},

			{Label: "── Icons ──", IsHeader: true},
			{Value: "lucide", Label: "Lucide Icons", Hint: "Clean and consistent open-source icons pack"},
			{Value: "svgl", Label: "SVGL", Hint: "Beautiful SVG logos and developer brand library"},

			{Label: "── Animations ──", IsHeader: true},
			{Value: "motion", Label: "Framer Motion", Hint: "Production-ready declarative animations library"},

			{Label: "── Pagos ──", IsHeader: true},
			{Value: "stripe", Label: "Stripe", Hint: "Payments infrastructure and subscription billing"},
			{Value: "polar", Label: "Polar (polar.sh)", Hint: "Developer-first monetization platform and billing engine"},

			{Label: "── Servicio de Correo ──", IsHeader: true},
			{Value: "resend", Label: "Resend", Hint: "Email API for developers with React email components"},
			{Value: "brevo", Label: "Brevo", Hint: "Transactional email delivery and marketing automation"},

			{Label: "── Validation & Typing ──", IsHeader: true},
			{Value: "zod", Label: "Zod", Hint: "TypeScript-first schema declaration and validation"},

			{Label: "── DevOps & Tooling ──", IsHeader: true},
			{Value: "docker", Label: "Docker Compose", Hint: "Local containerized services and database"},
			{Value: "github_actions", Label: "GitHub Actions CI", Hint: "Automated linting and test workflows"},
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
		// Rule 1: If Frontend is "none", Backend cannot be "none" or "self"
		if frontend == "none" {
			for i := range options {
				if options[i].Value == "none" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Project cannot have both Frontend and Backend as None"
				}
				if options[i].Value == "self" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: 'Self' backend requires a fullstack frontend to be selected"
				}
			}
		} else if !IsFullstackFrontend(frontend) {
			// Rule: If frontend is not fullstack (e.g. React SPA, Native), "self" is disabled
			for i := range options {
				if options[i].Value == "self" {
					options[i].Disabled = true
					options[i].DisabledReason = fmt.Sprintf("Incompatible: 'Self' backend requires a fullstack frontend (Next.js, Nuxt, Svelte, or Astro), but '%s' was selected", frontend)
				}
			}
		}

	case StepAPI:
		isNode := IsNodeEcosystem(frontend, backend)
		for i := range options {
			val := options[i].Value
			switch val {
			case "trpc":
				if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: tRPC requires TypeScript on frontend/backend"
				} else if frontend != "" && frontend != "none" && !IsReactFrontend(frontend) {
					options[i].Disabled = true
					options[i].DisabledReason = fmt.Sprintf("Incompatible: tRPC requires a React frontend (Next.js, React + Vite, React Native), but '%s' was selected. Use oRPC or None", frontend)
				}
			case "orpc":
				if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: oRPC requires TypeScript on frontend/backend"
				}
			}
		}

	case StepPackageManager:
		// If pure Go, Python, or Java project without JS/TS frontend
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
			if IsJavaBackend(backend) {
				return []views.SelectOption{
					{Value: "mvn", Label: "Maven", Hint: "Standard Java build & dependency tool"},
					{Value: "gradle", Label: "Gradle", Hint: "High performance build automation tool"},
				}
			}
		}

	case StepDatabase:
		// Rule 2: Client-side SPAs (React + Vite, Native) without a backend cannot directly connect to DB servers
		if (backend == "none" || backend == "") && IsClientSPA(frontend) {
			for i := range options {
				if options[i].Value != "none" {
					options[i].Disabled = true
					options[i].DisabledReason = fmt.Sprintf("Incompatible: Client-side app (%s) cannot connect directly to DBs without a backend", frontend)
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
		isJava := IsJavaBackend(backend)

		for i := range options {
			val := options[i].Value
			switch val {
			case "jpa":
				if !isJava {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Spring Data JPA is exclusively for Java / Spring Boot"
				}
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

			if isJava && val != "jpa" && val != "none" {
				options[i].Disabled = true
				options[i].DisabledReason = "Incompatible: Java Spring Boot uses JPA / Hibernate or raw SQL"
			}
		}

	case StepAuth:
		isNode := IsNodeEcosystem(frontend, backend)
		for i := range options {
			val := options[i].Value
			switch val {
			case "next-auth":
				if frontend != "nextjs" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: NextAuth is exclusively for Next.js frontend"
				}
			case "better-auth":
				if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Better Auth requires a TypeScript/Node ecosystem"
				}
			case "clerk":
				if !isNode {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Clerk requires a TypeScript/Node ecosystem"
				} else if frontend != "" && frontend != "none" && !IsReactFrontend(frontend) {
					options[i].Disabled = true
					options[i].DisabledReason = fmt.Sprintf("Incompatible: Clerk requires a React frontend (Next.js, React + Vite, React Native), but '%s' was selected", frontend)
				} else if backend == "self" && frontend != "nextjs" {
					options[i].Disabled = true
					options[i].DisabledReason = "Incompatible: Clerk with 'self' fullstack backend is only supported on Next.js"
				}
			}
		}
	}

	return options
}

// ValidateConfig validates the complete scaffold configuration against all compatibility rules.
func ValidateConfig(cfg scaffold.ScaffoldConfig) error {
	frontend := strings.ToLower(cfg.Frontend)
	backend := strings.ToLower(cfg.Backend)
	api := strings.ToLower(cfg.API)
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

	// 1.1 Backend 'self' requires a fullstack frontend
	if backend == "self" {
		if !IsFullstackFrontend(frontend) {
			return errors.NewValidationError(
				fmt.Sprintf("La opción de backend 'self' requiere un framework frontend fullstack (Next.js, Nuxt, Svelte o Astro), pero se seleccionó '%s'", frontend),
				"Selecciona un frontend con soporte de servidor (Next.js, Nuxt, Svelte, Astro) o elige un backend dedicado (Express, Hono, FastAPI, Go Chi, Spring Boot)",
			)
		}
	}

	// 1.2 API Layer validation
	if api == "trpc" || api == "orpc" {
		if !IsNodeEcosystem(frontend, backend) {
			return errors.NewValidationError(
				fmt.Sprintf("'%s' requiere un entorno Node.js / TypeScript", api),
				"Selecciona un frontend o backend TypeScript (Next.js, Express, Hono, etc.) o elige 'none' para la capa de API",
			)
		}
		if api == "trpc" && frontend != "" && frontend != "none" && !IsReactFrontend(frontend) {
			return errors.NewValidationError(
				fmt.Sprintf("tRPC requiere un frontend basado en React (Next.js, React + Vite, React Native), pero se seleccionó '%s'", frontend),
				"Para Nuxt, Svelte o Astro utiliza 'oRPC' o marca la API como 'none'",
			)
		}
	}

	// 2. Client-side SPA with no backend connecting to database
	if (backend == "none" || backend == "") && IsClientSPA(frontend) {
		if db != "none" && db != "" {
			return errors.NewValidationError(
				fmt.Sprintf("Una aplicación cliente (%s) sin backend no puede conectarse directamente a la base de datos '%s'", frontend, db),
				"Agrega un backend (como Express, Hono, FastAPI, Go Chi o Spring Boot) o usa un framework fullstack (como Next.js o Nuxt)",
			)
		}
		if orm != "none" && orm != "" {
			return errors.NewValidationError(
				fmt.Sprintf("Una aplicación cliente (%s) sin backend no puede utilizar el ORM '%s'", frontend, orm),
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
		if orm == "drizzle" || orm == "sqlalchemy" || orm == "gorm" || orm == "jpa" {
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
				"Para bases de datos relacionales SQL utiliza 'Drizzle', 'Prisma', 'SQLAlchemy' (Python), 'GORM' (Go) o 'JPA' (Java)",
			)
		}
	}

	// 5. Runtime & Ecosystem vs ORM
	isNode := IsNodeEcosystem(frontend, backend)
	isPython := IsPythonBackend(backend)
	isGo := IsGoBackend(backend)
	isJava := IsJavaBackend(backend)

	if isPython {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "gorm" || orm == "jpa" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el backend Python / FastAPI", orm),
				"Para Python utiliza 'SQLAlchemy' o 'none'",
			)
		}
	}

	if isGo {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "sqlalchemy" || orm == "jpa" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el backend Go", orm),
				"Para Go utiliza 'GORM' o 'none'",
			)
		}
	}

	if isJava {
		if orm == "drizzle" || orm == "prisma" || orm == "moongose" || orm == "mongoose" || orm == "sqlalchemy" || orm == "gorm" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el backend Java / Spring Boot", orm),
				"Para Java Spring Boot utiliza 'JPA' o 'none'",
			)
		}
	}

	if isNode {
		if orm == "sqlalchemy" || orm == "gorm" || orm == "jpa" {
			return errors.NewValidationError(
				fmt.Sprintf("El ORM '%s' no es compatible con el ecosistema Node.js / TypeScript", orm),
				"Para TypeScript utiliza 'Drizzle', 'Prisma' o 'Mongoose'",
			)
		}
	}

	// 6. Package Manager compatibility
	if frontend == "none" {
		if isGo && (pm == "pnpm" || pm == "npm" || pm == "bun" || pm == "mvn" || pm == "gradle") {
			return errors.NewValidationError(
				fmt.Sprintf("Un proyecto exclusivo de Go no utiliza el gestor de paquetes '%s'", pm),
				"Utiliza 'go_mod' para proyectos en Go",
			)
		}
		if isPython && (pm == "pnpm" || pm == "npm" || pm == "bun" || pm == "mvn" || pm == "gradle") {
			return errors.NewValidationError(
				fmt.Sprintf("Un proyecto exclusivo de Python no utiliza el gestor de paquetes '%s'", pm),
				"Utiliza 'pip' o 'uv' para proyectos en Python",
			)
		}
		if isJava && (pm == "pnpm" || pm == "npm" || pm == "bun" || pm == "go_mod") {
			return errors.NewValidationError(
				fmt.Sprintf("Un proyecto exclusivo de Java no utiliza el gestor de paquetes '%s'", pm),
				"Utiliza 'mvn' o 'gradle' para proyectos en Java",
			)
		}
	}

	// 7. Auth Provider compatibility
	if auth != "" && auth != "none" {
		if auth == "better-auth" {
			if !isNode {
				return errors.NewValidationError(
					"Better Auth requiere un entorno Node.js / TypeScript (Next.js, Express, Hono, Astro, Svelte, Nuxt)",
					"Agrega un frontend o backend compatible con TypeScript o cambia el proveedor de autenticación",
				)
			}
		} else if auth == "clerk" {
			if !isNode {
				return errors.NewValidationError(
					"Clerk requiere un entorno Node.js / TypeScript",
					"Agrega un frontend o backend compatible con TypeScript o cambia el proveedor de autenticación",
				)
			}
			if frontend != "" && frontend != "none" && !IsReactFrontend(frontend) {
				return errors.NewValidationError(
					fmt.Sprintf("Clerk requiere un frontend basado en React (Next.js, React + Vite, React Native), pero se seleccionó '%s'", frontend),
					"Selecciona un frontend React o utiliza 'Better Auth'",
				)
			}
			if backend == "self" && frontend != "nextjs" {
				return errors.NewValidationError(
					"Clerk con backend fullstack 'self' solo es compatible con Next.js",
					"Selecciona Next.js como frontend, utiliza un backend dedicado (Express, Hono, etc.) o cambia el proveedor de autenticación",
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

