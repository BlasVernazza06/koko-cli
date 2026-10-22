package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BlasVernazza06/koko-cli/internal/compatibility"
	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/doctor"
	"github.com/BlasVernazza06/koko-cli/internal/injector"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/BlasVernazza06/koko-cli/internal/types"
	"github.com/BlasVernazza06/koko-cli/internal/validator"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) {
	// 1. Tool: koko_get_options
	s.AddTool(
		mcp.NewTool("koko_get_options",
			mcp.WithDescription("Lists all supported recipes, frontend frameworks, backend frameworks, databases, ORMs, auth providers, and addons in Koko CLI"),
		),
		handleGetOptions,
	)

	// 2. Tool: koko_validate_stack
	s.AddTool(
		mcp.NewTool("koko_validate_stack",
			mcp.WithDescription("Validates compatibility between selected frontend, backend, orm, database, and auth before scaffolding"),
			mcp.WithString("frontend", mcp.Description("Frontend framework (next, react, nuxt, svelte, astro, native)")),
			mcp.WithString("backend", mcp.Description("Backend framework (hono, express, nestjs, fastapi, go-chi, fiber, spring)")),
			mcp.WithString("database", mcp.Description("Database engine (postgres, sqlite, mysql, mongodb, none)")),
			mcp.WithString("orm", mcp.Description("ORM tool (drizzle, prisma, sqlx, none)")),
			mcp.WithString("auth", mcp.Description("Auth provider (better-auth, next-auth, firebase, none)")),
			mcp.WithString("api", mcp.Description("API layer (trpc, orpc, none)")),
		),
		handleValidateStack,
	)

	// 3. Tool: koko_init_project
	s.AddTool(
		mcp.NewTool("koko_init_project",
			mcp.WithDescription("Scaffolds a new project or monorepo with production-ready structure and configurations"),
			mcp.WithString("projectName", mcp.Required(), mcp.Description("Name of the project to create")),
			mcp.WithString("recipe", mcp.Description("Predefined recipe: saas-next, enterprise-fullstack, modern-monorepo, api-service, minimal")),
			mcp.WithString("frontend", mcp.Description("Frontend framework (next, react, nuxt, svelte, astro, native)")),
			mcp.WithString("backend", mcp.Description("Backend framework (hono, express, nestjs, fastapi, go-chi, fiber, spring)")),
			mcp.WithString("database", mcp.Description("Database (postgres, sqlite, mysql, mongodb, none)")),
			mcp.WithString("orm", mcp.Description("ORM (drizzle, prisma, sqlx, none)")),
			mcp.WithString("auth", mcp.Description("Auth provider (better-auth, next-auth, firebase, none)")),
			mcp.WithString("packageManager", mcp.Description("Package manager (pnpm, npm, bun)")),
			mcp.WithString("addons", mcp.Description("Comma-separated addons: stripe, resend, docker, biome, shadcn")),
			mcp.WithBoolean("initGit", mcp.Description("Initialize git repository (default false)")),
			mcp.WithString("targetDirectory", mcp.Description("Target parent directory (default current working directory)")),
		),
		handleInitProject,
	)

	// 4. Tool: koko_add_service
	s.AddTool(
		mcp.NewTool("koko_add_service",
			mcp.WithDescription("Injects an addon or service into an existing Koko project (e.g., stripe, resend, polar, docker, biome, shadcn)"),
			mcp.WithString("addon", mcp.Required(), mcp.Description("The addon to inject: stripe, resend, polar, brevo, docker, biome, shadcn")),
			mcp.WithString("projectDirectory", mcp.Description("Directory of the existing Koko project (default current working directory)")),
		),
		handleAddService,
	)

	// 5. Tool: koko_doctor
	s.AddTool(
		mcp.NewTool("koko_doctor",
			mcp.WithDescription("Diagnoses project health, verifies dependencies, detects drift, and optionally auto-fixes koko.config.json"),
			mcp.WithString("projectDirectory", mcp.Description("Directory of the Koko project to diagnose (default current working directory)")),
			mcp.WithBoolean("autoFix", mcp.Description("If true, auto-repairs koko.config.json to match detected files")),
		),
		handleDoctor,
	)

	// 6. Tool: koko_set_default_directory
	s.AddTool(
		mcp.NewTool("koko_set_default_directory",
			mcp.WithDescription("Configures the persistent default target directory where new Koko projects will be created if not explicitly specified"),
			mcp.WithString("directory", mcp.Required(), mcp.Description("The absolute path to the directory (e.g. C:\\Projects or /home/user/projects)")),
		),
		handleSetDefaultDirectory,
	)

	// 7. Tool: koko_get_default_directory
	s.AddTool(
		mcp.NewTool("koko_get_default_directory",
			mcp.WithDescription("Returns the currently configured default target directory for project scaffolding"),
		),
		handleGetDefaultDirectory,
	)
}

func handleGetOptions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	options := map[string]interface{}{
		"recipes": []map[string]string{
			{"id": "saas-next", "name": "SaaS Next.js", "description": "Fullstack Next.js 15, Drizzle ORM, PostgreSQL, Better-Auth, Stripe"},
			{"id": "enterprise-fullstack", "name": "Enterprise Fullstack", "description": "Next.js + NestJS + PostgreSQL + Prisma + Docker"},
			{"id": "modern-monorepo", "name": "Modern Monorepo", "description": "Next.js + Hono RPC + Drizzle + Turborepo"},
			{"id": "api-service", "name": "API Service", "description": "Go Chi or FastAPI backend microservice with Docker and DB"},
			{"id": "minimal", "name": "Minimal", "description": "Clean modular starting point"},
		},
		"frontends":       []string{"next", "react", "nuxt", "svelte", "astro", "native", "none"},
		"backends":        []string{"hono", "express", "nestjs", "fastapi", "go-chi", "fiber", "spring", "none"},
		"databases":       []string{"postgres", "sqlite", "mysql", "mongodb", "none"},
		"orms":            []string{"drizzle", "prisma", "sqlx", "none"},
		"auth":            []string{"better-auth", "next-auth", "firebase", "none"},
		"packageManagers": []string{"pnpm", "npm", "bun"},
		"addons":          []string{"stripe", "polar", "resend", "brevo", "docker", "biome", "shadcn"},
	}

	bytes, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode options: %v", err)), nil
	}

	return mcp.NewToolResultText(string(bytes)), nil
}

func handleValidateStack(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	frontend, _ := request.RequireString("frontend")
	backend, _ := request.RequireString("backend")
	database, _ := request.RequireString("database")
	orm, _ := request.RequireString("orm")
	auth, _ := request.RequireString("auth")
	api, _ := request.RequireString("api")

	cfg := scaffold.ScaffoldConfig{
		Frontend: frontend,
		Backend:  backend,
		Database: database,
		ORM:      orm,
		Auth:     auth,
		API:      api,
	}

	err := compatibility.ValidateConfig(cfg)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Stack validation error: %v", err)), nil
	}

	return mcp.NewToolResultText("Configuration is valid and fully compatible!"), nil
}

func handleInitProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectName, err := request.RequireString("projectName")
	if err != nil || projectName == "" {
		return mcp.NewToolResultError("projectName is required"), nil
	}

	// 1. Validar nombre de proyecto con reglas estrictas de Koko
	if err := validator.ValidateProjectName(projectName); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid project name '%s': %v", projectName, err)), nil
	}

	rawTargetDir, _ := request.RequireString("targetDirectory")
	targetDir := kokoConfig.ResolveDefaultTargetDirectory(rawTargetDir)

	recipe, _ := request.RequireString("recipe")
	frontend, _ := request.RequireString("frontend")
	backend, _ := request.RequireString("backend")
	database, _ := request.RequireString("database")
	orm, _ := request.RequireString("orm")
	auth, _ := request.RequireString("auth")
	api, _ := request.RequireString("api")
	pm, _ := request.RequireString("packageManager")
	addons, _ := request.RequireString("addons")
	initGit, _ := request.RequireBool("initGit")

	scaffoldCfg := types.ScaffoldConfig{
		ProjectName:    projectName,
		Recipe:         recipe,
		Frontend:       frontend,
		Backend:        backend,
		API:            api,
		Database:       database,
		ORM:            orm,
		Auth:           auth,
		PackageManager: pm,
		Addons:         addons,
		InitGit:        initGit,
	}

	// 2. Validar compatibilidad tecnológica cuando no se use una receta predefinida completa
	if recipe == "" {
		if err := compatibility.ValidateConfig(scaffoldCfg); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Incompatible stack configuration: %v", err)), nil
		}
	} else if frontend != "" || backend != "" || api != "" || database != "" || orm != "" || auth != "" {
		// Si se especifica receta pero además se sobreescriben flags personalizados
		if err := compatibility.ValidateConfig(scaffoldCfg); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Incompatible stack configuration: %v", err)), nil
		}
	}

	fullPath := filepath.Join(targetDir, projectName)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return mcp.NewToolResultError(fmt.Sprintf("Target directory '%s' already exists", fullPath)), nil
	}

	if err := scaffold.RunScaffold(fullPath, scaffoldCfg); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to scaffold project: %v", err)), nil
	}

	// Generar el manifiesto koko.config.json del proyecto
	if err := kokoConfig.GenerateConfig(fullPath, scaffoldCfg); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to generate koko.config.json: %v", err)), nil
	}

	resp := map[string]interface{}{
		"success":     true,
		"projectName": projectName,
		"projectPath": fullPath,
		"recipe":      recipe,
		"message":     fmt.Sprintf("Project '%s' successfully generated at '%s'", projectName, fullPath),
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")

	return mcp.NewToolResultText(string(bytes)), nil
}

func handleAddService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	addon, err := request.RequireString("addon")
	if err != nil || addon == "" {
		return mcp.NewToolResultError("addon name is required"), nil
	}

	projectDir, _ := request.RequireString("projectDirectory")
	if projectDir == "" {
		projectDir = "."
	}

	if err := injector.AddAddon(projectDir, addon); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to inject addon '%s': %v", addon, err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully injected addon '%s' into project at '%s'", addon, projectDir)), nil
}

func handleDoctor(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectDir, _ := request.RequireString("projectDirectory")
	if projectDir == "" {
		projectDir = "."
	}

	autoFix, _ := request.RequireBool("autoFix")

	result, err := doctor.RunDiagnostics(projectDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run diagnostics: %v", err)), nil
	}

	diffCount := len(result.Differences)
	fixed := false

	if autoFix && diffCount > 0 {
		if err := doctor.ApplyFixes(projectDir, result.DetectedConfig, result.ExistingConfig); err == nil {
			fixed = true
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Diagnostics found differences but auto-fix failed: %v", err)), nil
		}
	}

	summary := map[string]interface{}{
		"configExists":     result.ConfigExists,
		"differencesFound": diffCount,
		"differences":      result.Differences,
		"autoFixed":        fixed,
		"status":           "healthy",
	}

	if diffCount > 0 && !fixed {
		summary["status"] = "inconsistent"
	}

	bytes, _ := json.MarshalIndent(summary, "", "  ")
	return mcp.NewToolResultText(string(bytes)), nil
}

func handleSetDefaultDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir, err := request.RequireString("directory")
	if err != nil || dir == "" {
		return mcp.NewToolResultError("directory path is required"), nil
	}

	cleanDir := filepath.Clean(dir)
	if err := os.MkdirAll(cleanDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create or access directory '%s': %v", cleanDir, err)), nil
	}

	settings, err := kokoConfig.LoadGlobalSettings()
	if err != nil {
		settings = &kokoConfig.GlobalSettings{}
	}
	settings.DefaultProjectsDirectory = cleanDir

	if err := kokoConfig.SaveGlobalSettings(settings); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to save default directory settings: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Default projects directory successfully set to: %s", cleanDir)), nil
}

func handleGetDefaultDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resolved := kokoConfig.ResolveDefaultTargetDirectory("")
	settings, _ := kokoConfig.LoadGlobalSettings()

	resp := map[string]interface{}{
		"configuredDirectory": settings.DefaultProjectsDirectory,
		"activeTargetDirectory": resolved,
	}
	bytes, _ := json.MarshalIndent(resp, "", "  ")
	return mcp.NewToolResultText(string(bytes)), nil
}

