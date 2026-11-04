package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/core/config"
	"github.com/BlasVernazza06/koko-cli/internal/features/doctor"
)

type rawPackageJSON struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
}

func readScriptsFromPackage(path string) map[string]string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var pkg rawPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	if pkg.Scripts == nil {
		return make(map[string]string)
	}
	return pkg.Scripts
}

// DetectProjectDBContext analiza el proyecto ubicado en targetDir y construye su contexto de DB.
func DetectProjectDBContext(targetDir string) (*ProjectDBContext, error) {
	cleanDir := filepath.Clean(targetDir)
	info, err := os.Stat(cleanDir)
	if err != nil {
		return nil, fmt.Errorf("target directory '%s' not accessible: %w", cleanDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("target path '%s' is not a directory", cleanDir)
	}

	ctx := &ProjectDBContext{
		ProjectDir: cleanDir,
	}

	// 1. Verificar si existe koko.config.json
	kokoConfigPath := filepath.Join(cleanDir, "koko.config.json")
	if configData, err := os.ReadFile(kokoConfigPath); err == nil {
		configData = bytes.TrimPrefix(configData, []byte("\xef\xbb\xbf"))
		var cfg config.KokoConfig
		if err := json.Unmarshal(configData, &cfg); err == nil {
			ctx.ConfigExists = true
			ctx.Layout = cfg.Architecture.Layout
			ctx.PackageManager = cfg.Architecture.PackageManager
			if cfg.Stack.Database != nil {
				ctx.Provider = cfg.Stack.Database.Provider
				ctx.ORM = cfg.Stack.Database.ORM
			}
		}
	}

	// 2. Fallback / Enriquecimiento mediante análisis estático (doctor)
	if !ctx.ConfigExists || ctx.ORM == "" || ctx.PackageManager == "" {
		detected := doctor.AnalyzeProject(cleanDir)
		if ctx.Layout == "" {
			ctx.Layout = detected.Architecture.Layout
		}
		if ctx.PackageManager == "" {
			ctx.PackageManager = detected.Architecture.PackageManager
		}
		if ctx.Provider == "" && detected.Stack.Database != nil {
			ctx.Provider = detected.Stack.Database.Provider
		}
		if ctx.ORM == "" && detected.Stack.Database != nil {
			ctx.ORM = detected.Stack.Database.ORM
		}
	}

	// Normalizaciones por defecto
	if ctx.PackageManager == "" {
		ctx.PackageManager = "pnpm"
	}
	ctx.PackageManager = strings.ToLower(strings.TrimSpace(ctx.PackageManager))
	ctx.ORM = strings.ToLower(strings.TrimSpace(ctx.ORM))
	ctx.Provider = strings.ToLower(strings.TrimSpace(ctx.Provider))
	if ctx.Layout == "" {
		ctx.Layout = "standalone"
	}

	// 3. Cargar scripts de package.json
	ctx.RootScripts = readScriptsFromPackage(filepath.Join(cleanDir, "package.json"))
	dbPkgPath := filepath.Join(cleanDir, "packages", "db", "package.json")
	if _, err := os.Stat(dbPkgPath); err == nil {
		ctx.HasDBPackage = true
		ctx.DBScripts = readScriptsFromPackage(dbPkgPath)
	}

	return ctx, nil
}

// ResolveCommand traduce una acción de base de datos en un ResolvedCommand ejecutable.
func ResolveCommand(pCtx *ProjectDBContext, action Action, extraArgs []string) (*ResolvedCommand, error) {
	if pCtx == nil {
		return nil, fmt.Errorf("project db context cannot be nil")
	}

	// Normalizar acción
	actionStr := string(action)

	// =========================================================================
	// PRIORIDAD 1: Script explícito en root package.json
	// =========================================================================
	// Nombres convencionales según la acción
	scriptCandidates := getScriptCandidatesForAction(action)
	for _, scriptName := range scriptCandidates {
		if _, exists := pCtx.RootScripts[scriptName]; exists {
			exe, args := buildPackageManagerRunCmd(pCtx.PackageManager, scriptName, extraArgs)
			return &ResolvedCommand{
				Executable: exe,
				Args:       args,
				WorkDir:    pCtx.ProjectDir,
				SourceDesc: fmt.Sprintf("package.json script '%s' via %s", scriptName, pCtx.PackageManager),
			}, nil
		}
	}

	// =========================================================================
	// PRIORIDAD 1.5: Monorepo con packages/db pero sin script directo en root
	// =========================================================================
	if pCtx.HasDBPackage && pCtx.Layout == "monorepo" {
		for _, scriptName := range scriptCandidates {
			if _, exists := pCtx.DBScripts[scriptName]; exists {
				exe, args := buildMonorepoFilterCmd(pCtx.PackageManager, scriptName, extraArgs)
				return &ResolvedCommand{
					Executable: exe,
					Args:       args,
					WorkDir:    pCtx.ProjectDir,
					SourceDesc: fmt.Sprintf("packages/db script '%s' via %s filter", scriptName, pCtx.PackageManager),
				}, nil
			}
		}
	}

	// =========================================================================
	// PRIORIDAD 2: Resolución nativa según ORM detectado
	// =========================================================================
	switch pCtx.ORM {
	case "drizzle":
		return resolveDrizzleCommand(pCtx, action, extraArgs)
	case "prisma":
		return resolvePrismaCommand(pCtx, action, extraArgs)
	case "mongoose":
		return resolveMongooseCommand(pCtx, action, extraArgs)
	default:
		return nil, fmt.Errorf("no database ORM configured or detected (detected: '%s'). Unable to resolve '%s'", pCtx.ORM, actionStr)
	}
}

func getScriptCandidatesForAction(action Action) []string {
	switch action {
	case ActionMigrate:
		return []string{"db:migrate", "migrate"}
	case ActionPush:
		return []string{"db:push", "push"}
	case ActionGenerate:
		return []string{"db:generate", "generate"}
	case ActionStudio:
		return []string{"db:studio", "studio"}
	case ActionSeed:
		return []string{"db:seed", "seed"}
	case ActionReset:
		return []string{"db:reset", "reset"}
	default:
		return []string{string(action)}
	}
}

func buildPackageManagerRunCmd(pm string, script string, extraArgs []string) (string, []string) {
	switch pm {
	case "bun":
		args := []string{"run", script}
		if len(extraArgs) > 0 {
			args = append(args, "--")
			args = append(args, extraArgs...)
		}
		return "bun", args
	case "npm":
		args := []string{"run", script}
		if len(extraArgs) > 0 {
			args = append(args, "--")
			args = append(args, extraArgs...)
		}
		return "npm", args
	case "pnpm":
		fallthrough
	default:
		args := []string{"run", script}
		if len(extraArgs) > 0 {
			args = append(args, "--")
			args = append(args, extraArgs...)
		}
		return "pnpm", args
	}
}

func buildMonorepoFilterCmd(pm string, script string, extraArgs []string) (string, []string) {
	switch pm {
	case "bun":
		args := []string{"run", "--filter", "@*/db", script}
		if len(extraArgs) > 0 {
			args = append(args, "--")
			args = append(args, extraArgs...)
		}
		return "bun", args
	case "npm":
		args := []string{"run", script, "--workspace", "packages/db"}
		if len(extraArgs) > 0 {
			args = append(args, "--")
			args = append(args, extraArgs...)
		}
		return "npm", args
	case "pnpm":
		fallthrough
	default:
		args := []string{"--filter", "@*/db", "run", script}
		if len(extraArgs) > 0 {
			args = append(args, "--")
			args = append(args, extraArgs...)
		}
		return "pnpm", args
	}
}

func resolveDrizzleCommand(pCtx *ProjectDBContext, action Action, extraArgs []string) (*ResolvedCommand, error) {
	var subArgs []string
	switch action {
	case ActionMigrate:
		subArgs = []string{"migrate"}
	case ActionPush:
		subArgs = []string{"push"}
	case ActionGenerate:
		subArgs = []string{"generate"}
	case ActionStudio:
		subArgs = []string{"studio"}
	case ActionSeed:
		// Drizzle no tiene comando built-in para seed; sugerir o verificar script
		return nil, fmt.Errorf("drizzle does not provide a native 'seed' command. Please define a 'db:seed' script in your package.json")
	case ActionReset:
		return nil, fmt.Errorf("drizzle does not provide a native 'reset' command. Use 'push --force' or define a 'db:reset' script in package.json")
	default:
		return nil, fmt.Errorf("unsupported action '%s' for drizzle", action)
	}

	execBin, execArgs := getORMRunner(pCtx.PackageManager, "drizzle-kit", subArgs, extraArgs)
	workDir := pCtx.ProjectDir
	if pCtx.HasDBPackage && pCtx.Layout == "monorepo" {
		workDir = filepath.Join(pCtx.ProjectDir, "packages", "db")
	}

	return &ResolvedCommand{
		Executable: execBin,
		Args:       execArgs,
		WorkDir:    workDir,
		SourceDesc: fmt.Sprintf("native drizzle-kit %s via %s", strings.Join(subArgs, " "), execBin),
	}, nil
}

func resolvePrismaCommand(pCtx *ProjectDBContext, action Action, extraArgs []string) (*ResolvedCommand, error) {
	var subArgs []string
	switch action {
	case ActionMigrate:
		subArgs = []string{"migrate", "dev"}
	case ActionPush:
		subArgs = []string{"db", "push"}
	case ActionGenerate:
		subArgs = []string{"generate"}
	case ActionStudio:
		subArgs = []string{"studio"}
	case ActionSeed:
		subArgs = []string{"db", "seed"}
	case ActionReset:
		subArgs = []string{"migrate", "reset"}
	default:
		return nil, fmt.Errorf("unsupported action '%s' for prisma", action)
	}

	execBin, execArgs := getORMRunner(pCtx.PackageManager, "prisma", subArgs, extraArgs)
	workDir := pCtx.ProjectDir
	if pCtx.HasDBPackage && pCtx.Layout == "monorepo" {
		workDir = filepath.Join(pCtx.ProjectDir, "packages", "db")
	}

	return &ResolvedCommand{
		Executable: execBin,
		Args:       execArgs,
		WorkDir:    workDir,
		SourceDesc: fmt.Sprintf("native prisma %s via %s", strings.Join(subArgs, " "), execBin),
	}, nil
}

func resolveMongooseCommand(pCtx *ProjectDBContext, action Action, extraArgs []string) (*ResolvedCommand, error) {
	// Mongoose no cuenta con CLI oficial de schema migration; depende enteramente de scripts de usuario
	return nil, fmt.Errorf("mongoose does not provide a native CLI for '%s'. Please define a 'db:%s' script in package.json", action, action)
}

func getORMRunner(pm string, tool string, toolArgs []string, extraArgs []string) (string, []string) {
	var exe string
	var args []string

	switch pm {
	case "bun":
		exe = "bunx"
		args = append([]string{tool}, toolArgs...)
	case "npm":
		exe = "npx"
		args = append([]string{tool}, toolArgs...)
	case "pnpm":
		fallthrough
	default:
		exe = "pnpm"
		args = append([]string{"exec", tool}, toolArgs...)
	}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	return exe, args
}
