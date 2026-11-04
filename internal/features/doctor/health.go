package doctor

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BlasVernazza06/koko-cli/internal/core/config"
)

const (
	tcpDialTimeout    = 150 * time.Millisecond
	processExecTimeout = 1500 * time.Millisecond
)

// HealthChecker orquesta los checks de salud del sistema, variables, puertos y skills.
type HealthChecker struct {
	lookPath   func(string) (string, error)
	runCommand func(ctx context.Context, name string, args ...string) (string, error)
	tcpDial    func(network, address string, timeout time.Duration) (net.Conn, error)
}

// NewHealthChecker crea una instancia con los ejecutores reales del SO.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		lookPath: exec.LookPath,
		runCommand: func(ctx context.Context, name string, args ...string) (string, error) {
			cmd := exec.CommandContext(ctx, name, args...)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out
			err := cmd.Run()
			return strings.TrimSpace(out.String()), err
		},
		tcpDial: net.DialTimeout,
	}
}

// CheckEnvironment verifica Node.js, gestor de paquetes, Git y Docker.
func (h *HealthChecker) CheckEnvironment(ctx context.Context, projectDir string, detected config.KokoConfig) []HealthItem {
	var items []HealthItem

	// 1. Node.js check
	nodeItem := HealthItem{Name: "Node.js"}
	ctxNode, cancelNode := context.WithTimeout(ctx, processExecTimeout)
	defer cancelNode()

	if _, err := h.lookPath("node"); err != nil {
		nodeItem.Status = StatusFail
		nodeItem.Message = "Node.js not found in PATH"
		nodeItem.Detail = "Required for JavaScript/TypeScript development (>= 18.0.0)"
		nodeItem.FixCommand = "Install Node.js (https://nodejs.org/) or use nvm"
	} else {
		out, err := h.runCommand(ctxNode, "node", "-v")
		if err != nil {
			nodeItem.Status = StatusWarn
			nodeItem.Message = "Installed but failed to query version"
		} else {
			cleanVer := strings.TrimPrefix(out, "v")
			parts := strings.Split(cleanVer, ".")
			major := 0
			if len(parts) > 0 {
				major, _ = strconv.Atoi(parts[0])
			}
			if major >= 18 {
				nodeItem.Status = StatusOk
				nodeItem.Message = fmt.Sprintf("v%s (satisfies >= 18.0.0)", cleanVer)
			} else {
				nodeItem.Status = StatusWarn
				nodeItem.Message = fmt.Sprintf("v%s (outdated, recommend >= 18.0.0)", cleanVer)
				nodeItem.FixCommand = "Upgrade Node.js to >= 18.0.0"
			}
		}
	}
	items = append(items, nodeItem)

	// 2. Package Manager
	pm := detected.Architecture.PackageManager
	if pm == "" {
		pm = "pnpm"
	}
	pmItem := HealthItem{Name: fmt.Sprintf("Package Manager (%s)", pm)}
	ctxPM, cancelPM := context.WithTimeout(ctx, processExecTimeout)
	defer cancelPM()

	if _, err := h.lookPath(pm); err != nil {
		pmItem.Status = StatusFail
		pmItem.Message = fmt.Sprintf("'%s' is not installed or not in PATH", pm)
		pmItem.FixCommand = fmt.Sprintf("npm install -g %s", pm)
	} else {
		out, err := h.runCommand(ctxPM, pm, "-v")
		if err != nil {
			pmItem.Status = StatusOk
			pmItem.Message = "Installed"
		} else {
			pmItem.Status = StatusOk
			pmItem.Message = fmt.Sprintf("v%s installed", strings.TrimSpace(out))
		}
	}
	items = append(items, pmItem)

	// 3. Git repository
	gitItem := HealthItem{Name: "Git repository"}
	dotGit := filepath.Join(projectDir, ".git")
	if fileExists(dotGit) {
		gitItem.Status = StatusOk
		gitItem.Message = "Initialized (.git found)"
	} else {
		gitItem.Status = StatusWarn
		gitItem.Message = "Repository not initialized"
		gitItem.Detail = "Tracking project changes is recommended"
		gitItem.FixCommand = "git init"
	}
	items = append(items, gitItem)

	// 4. Docker (si aplica)
	hasDocker := (detected.Features.Infrastructure != nil && detected.Features.Infrastructure.DockerCompose) ||
		fileExists(filepath.Join(projectDir, "docker-compose.yml")) ||
		fileExists(filepath.Join(projectDir, "docker-compose.yaml")) ||
		fileExists(filepath.Join(projectDir, "compose.yaml"))

	if hasDocker {
		dockerItem := HealthItem{Name: "Docker daemon"}
		if _, err := h.lookPath("docker"); err != nil {
			dockerItem.Status = StatusWarn
			dockerItem.Message = "Docker CLI not found in PATH"
			dockerItem.Detail = "Required to run compose services locally"
			dockerItem.FixCommand = "Install Docker Desktop"
		} else {
			ctxDoc, cancelDoc := context.WithTimeout(ctx, processExecTimeout)
			defer cancelDoc()
			_, err := h.runCommand(ctxDoc, "docker", "info")
			if err != nil {
				dockerItem.Status = StatusWarn
				dockerItem.Message = "Docker daemon is not running"
				dockerItem.Detail = "Local database/compose containers will fail to start"
				dockerItem.FixCommand = "Start Docker Desktop or service"
			} else {
				dockerItem.Status = StatusOk
				dockerItem.Message = "Daemon is active and running"
			}
		}
		items = append(items, dockerItem)
	}

	return items
}

// CheckEnvFiles escanea variables requeridas en .env.example contra .env.
func (h *HealthChecker) CheckEnvFiles(projectDir string) []HealthItem {
	var items []HealthItem

	// Buscar pares (.env.example vs .env) en root y en apps/web, apps/api
	searchDirs := []string{
		projectDir,
		filepath.Join(projectDir, "apps", "web"),
		filepath.Join(projectDir, "apps", "api"),
		filepath.Join(projectDir, "packages", "db"),
	}

	for _, dir := range searchDirs {
		examplePath := filepath.Join(dir, ".env.example")
		if !fileExists(examplePath) {
			continue
		}

		relDir, _ := filepath.Rel(projectDir, dir)
		if relDir == "." || relDir == "" {
			relDir = "root"
		}

		envPath := filepath.Join(dir, ".env")
		envLocalPath := filepath.Join(dir, ".env.local")

		if !fileExists(envPath) && !fileExists(envLocalPath) {
			items = append(items, HealthItem{
				Name:       fmt.Sprintf(".env (%s)", relDir),
				Status:     StatusFail,
				Message:    "Missing .env file (found .env.example)",
				Detail:     "Environment configuration is required before starting development",
				FixCommand: fmt.Sprintf("cp %s %s", filepath.Join(dir, ".env.example"), filepath.Join(dir, ".env")),
			})
			continue
		}

		// Leer variables de .env.example
		requiredKeys := parseEnvKeys(examplePath)
		if len(requiredKeys) == 0 {
			continue
		}

		// Leer variables activas
		activeVars := make(map[string]string)
		if fileExists(envPath) {
			for k, v := range parseEnvVariables(envPath) {
				activeVars[k] = v
			}
		}
		if fileExists(envLocalPath) {
			for k, v := range parseEnvVariables(envLocalPath) {
				activeVars[k] = v
			}
		}

		var missingKeys []string
		var placeholderKeys []string

		for _, key := range requiredKeys {
			val, found := activeVars[key]
			if !found {
				missingKeys = append(missingKeys, key)
			} else if isPlaceholderValue(val) {
				placeholderKeys = append(placeholderKeys, key)
			}
		}

		if len(missingKeys) > 0 {
			items = append(items, HealthItem{
				Name:       fmt.Sprintf("Keys in %s/.env", relDir),
				Status:     StatusFail,
				Message:    fmt.Sprintf("%d missing keys defined in .env.example", len(missingKeys)),
				Detail:     strings.Join(missingKeys, ", "),
				FixCommand: fmt.Sprintf("Add missing keys to %s/.env", relDir),
			})
		}

		if len(placeholderKeys) > 0 {
			items = append(items, HealthItem{
				Name:       fmt.Sprintf("Placeholders in %s/.env", relDir),
				Status:     StatusWarn,
				Message:    fmt.Sprintf("%d keys contain placeholder values", len(placeholderKeys)),
				Detail:     strings.Join(placeholderKeys, ", "),
				FixCommand: fmt.Sprintf("Update values in %s/.env with real credentials", relDir),
			})
		}

		if len(missingKeys) == 0 && len(placeholderKeys) == 0 {
			items = append(items, HealthItem{
				Name:    fmt.Sprintf(".env (%s)", relDir),
				Status:  StatusOk,
				Message: fmt.Sprintf("All %d keys configured and in sync with .env.example", len(requiredKeys)),
			})
		}
	}

	if len(items) == 0 {
		// No se encontraron archivos de ejemplo
		items = append(items, HealthItem{
			Name:    ".env files",
			Status:  StatusInfo,
			Message: "No .env.example found; skipping key audit",
		})
	}

	return items
}

// CheckDatabaseAndPorts audita esquemas de DB, conectividad de puertos de DB y conflicto de puertos dev.
func (h *HealthChecker) CheckDatabaseAndPorts(projectDir string, detected config.KokoConfig) []HealthItem {
	var items []HealthItem

	// 1. Verificación de Esquema de Base de Datos
	if detected.Stack.Database != nil && detected.Stack.Database.ORM != "none" && detected.Stack.Database.ORM != "" {
		orm := detected.Stack.Database.ORM
		schemaItem := HealthItem{Name: fmt.Sprintf("Database Schema (%s)", orm)}

		switch orm {
		case "prisma":
			schema1 := filepath.Join(projectDir, "packages", "db", "prisma", "schema.prisma")
			schema2 := filepath.Join(projectDir, "prisma", "schema.prisma")
			if fileExists(schema1) || fileExists(schema2) {
				schemaItem.Status = StatusOk
				schemaItem.Message = "schema.prisma found"
			} else {
				schemaItem.Status = StatusWarn
				schemaItem.Message = "schema.prisma missing"
				schemaItem.FixCommand = "koko db generate"
			}
		case "drizzle":
			cfg1 := filepath.Join(projectDir, "packages", "db", "drizzle.config.ts")
			cfg2 := filepath.Join(projectDir, "drizzle.config.ts")
			if fileExists(cfg1) || fileExists(cfg2) {
				schemaItem.Status = StatusOk
				schemaItem.Message = "drizzle.config.ts found"
			} else {
				schemaItem.Status = StatusWarn
				schemaItem.Message = "drizzle.config.ts missing"
				schemaItem.FixCommand = "koko db generate"
			}
		default:
			schemaItem.Status = StatusOk
			schemaItem.Message = "Schema configured"
		}
		items = append(items, schemaItem)

		// 2. Conectividad local del servicio de Base de Datos (si es Postgres, Mongo, MySQL)
		dbPort := 0
		dbServiceName := ""
		switch detected.Stack.Database.Provider {
		case "postgres":
			dbPort = 5432
			dbServiceName = "PostgreSQL"
		case "mongodb":
			dbPort = 27017
			dbServiceName = "MongoDB"
		case "mysql":
			dbPort = 3306
			dbServiceName = "MySQL"
		}

		if dbPort > 0 {
			dbItem := HealthItem{Name: fmt.Sprintf("Local DB Service (%s :%d)", dbServiceName, dbPort)}
			conn, err := h.tcpDial("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(dbPort)), tcpDialTimeout)
			if err == nil {
				conn.Close()
				dbItem.Status = StatusOk
				dbItem.Message = "Listening on localhost"
			} else {
				dbItem.Status = StatusWarn
				dbItem.Message = "Service is not reachable on localhost"
				dbItem.Detail = "Start local service or docker compose up -d"
				dbItem.FixCommand = "docker compose up -d"
			}
			items = append(items, dbItem)
		}
	}

	// 3. Chequeo de puertos de desarrollo en conflicto
	devPorts := getDevPortsForStack(detected)
	for _, dp := range devPorts {
		portItem := HealthItem{Name: fmt.Sprintf("Port :%d (%s)", dp.Port, dp.Description)}
		conn, err := h.tcpDial("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(dp.Port)), tcpDialTimeout)
		if err == nil {
			conn.Close()
			portItem.Status = StatusWarn
			portItem.Message = "Port is already in use by another process"
			portItem.Detail = fmt.Sprintf("May conflict when starting %s", dp.Description)
			portItem.FixCommand = fmt.Sprintf("Stop the process using port %d before running dev", dp.Port)
		} else {
			portItem.Status = StatusOk
			portItem.Message = "Available"
		}
		items = append(items, portItem)
	}

	return items
}

type devPortDef struct {
	Port        int
	Description string
}

func getDevPortsForStack(detected config.KokoConfig) []devPortDef {
	var ports []devPortDef
	seen := make(map[int]bool)

	addPort := func(p int, desc string) {
		if !seen[p] && p > 0 {
			seen[p] = true
			ports = append(ports, devPortDef{Port: p, Description: desc})
		}
	}

	// Frontend ports
	if detected.Stack.Frontend != nil {
		switch detected.Stack.Frontend.Framework {
		case "next":
			addPort(3000, "Next.js dev server")
		case "react":
			addPort(5173, "Vite dev server")
		case "astro":
			addPort(4321, "Astro dev server")
		case "nuxt":
			addPort(3000, "Nuxt dev server")
		case "svelte":
			addPort(5173, "SvelteKit dev server")
		default:
			addPort(3000, "Frontend server")
		}
	}

	// Backend ports
	if detected.Stack.Backend != nil {
		switch detected.Stack.Backend.Framework {
		case "express", "hono":
			addPort(3001, "API server")
		case "fastapi":
			addPort(8000, "FastAPI server")
		case "spring_boot", "nestjs":
			addPort(8080, "Backend service")
		case "go_chi":
			addPort(8080, "Go API service")
		}
	}

	return ports
}

// CheckSkillsContext audita la presencia de AI Agent skills (.cursor/rules o .cursor/skills).
func (h *HealthChecker) CheckSkillsContext(projectDir string, detected config.KokoConfig) []HealthItem {
	var items []HealthItem

	rulesDir := filepath.Join(projectDir, ".cursor", "rules")
	skillsDir := filepath.Join(projectDir, ".cursor", "skills")

	hasRulesDir := fileExists(rulesDir) || fileExists(skillsDir)

	// Candidatos a tener skill instalada
	var candidates []string

	if detected.Stack.Frontend != nil {
		switch detected.Stack.Frontend.Framework {
		case "next":
			candidates = append(candidates, "nextjs")
		case "react":
			candidates = append(candidates, "react")
		case "astro":
			candidates = append(candidates, "astro")
		case "svelte":
			candidates = append(candidates, "svelte")
		case "vue", "nuxt":
			candidates = append(candidates, "vue")
		}
		if detected.Stack.Frontend.Styling == "tailwindcss" {
			candidates = append(candidates, "tailwindcss")
		}
		if detected.Stack.Frontend.UILibrary == "shadcn" {
			candidates = append(candidates, "shadcn-ui")
		}
	}

	if detected.Stack.Database != nil {
		if detected.Stack.Database.ORM == "drizzle" {
			candidates = append(candidates, "drizzle")
		} else if detected.Stack.Database.ORM == "prisma" {
			candidates = append(candidates, "prisma")
		}
	}

	if detected.Features.Payments != nil && detected.Features.Payments.Provider == "stripe" {
		candidates = append(candidates, "stripe-best-practices")
	}

	var missingSkills []string
	for _, c := range candidates {
		ruleMDC := filepath.Join(rulesDir, c+".mdc")
		ruleMD := filepath.Join(rulesDir, c+".md")
		skillMD := filepath.Join(skillsDir, c, "SKILL.md")

		if !fileExists(ruleMDC) && !fileExists(ruleMD) && !fileExists(skillMD) {
			missingSkills = append(missingSkills, c)
		}
	}

	if !hasRulesDir && len(candidates) > 0 {
		items = append(items, HealthItem{
			Name:       "AI Agent Skills",
			Status:     StatusWarn,
			Message:    "No .cursor/rules or skills installed",
			Detail:     fmt.Sprintf("Detected technologies: %s", strings.Join(candidates, ", ")),
			FixCommand: "koko skills",
		})
	} else if len(missingSkills) > 0 {
		items = append(items, HealthItem{
			Name:       "AI Agent Skills",
			Status:     StatusWarn,
			Message:    fmt.Sprintf("%d skills missing for detected stack", len(missingSkills)),
			Detail:     fmt.Sprintf("Missing: %s", strings.Join(missingSkills, ", ")),
			FixCommand: "koko skills",
		})
	} else if len(candidates) > 0 {
		items = append(items, HealthItem{
			Name:    "AI Agent Skills",
			Status:  StatusOk,
			Message: fmt.Sprintf("All %d stack skills are installed in .cursor", len(candidates)),
		})
	} else {
		items = append(items, HealthItem{
			Name:    "AI Agent Skills",
			Status:  StatusInfo,
			Message: "No specific skills required for current stack",
		})
	}

	return items
}

// ConsolidateRecommendedActions reúne todas las acciones correctivas sugeridas.
func ConsolidateRecommendedActions(environment, envVars, dbPorts, skills []HealthItem) []string {
	var actions []string
	seen := make(map[string]bool)

	addAction := func(cmd string) {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" && !seen[cmd] {
			seen[cmd] = true
			actions = append(actions, cmd)
		}
	}

	for _, list := range [][]HealthItem{environment, envVars, dbPorts, skills} {
		for _, it := range list {
			if it.Status == StatusFail || it.Status == StatusWarn {
				if it.FixCommand != "" {
					addAction(it.FixCommand)
				}
			}
		}
	}

	return actions
}

// RunFullHealthCheck ejecuta el diagnóstico completo en 5 módulos.
func RunFullHealthCheck(ctx context.Context, projectDir string, detected config.KokoConfig) *HealthCheckReport {
	checker := NewHealthChecker()

	envItems := checker.CheckEnvironment(ctx, projectDir, detected)
	envVarItems := checker.CheckEnvFiles(projectDir)
	dbPortItems := checker.CheckDatabaseAndPorts(projectDir, detected)
	skillItems := checker.CheckSkillsContext(projectDir, detected)

	actions := ConsolidateRecommendedActions(envItems, envVarItems, dbPortItems, skillItems)

	return &HealthCheckReport{
		Environment:        envItems,
		EnvVars:            envVarItems,
		DatabaseAndPorts:   dbPortItems,
		Skills:             skillItems,
		RecommendedActions: actions,
	}
}

// FixEnvFiles busca archivos .env.example que no tengan .env correspondiente y los copia.
func FixEnvFiles(projectDir string) (int, error) {
	copiedCount := 0
	searchDirs := []string{
		projectDir,
		filepath.Join(projectDir, "apps", "web"),
		filepath.Join(projectDir, "apps", "api"),
		filepath.Join(projectDir, "packages", "db"),
	}

	for _, dir := range searchDirs {
		examplePath := filepath.Join(dir, ".env.example")
		envPath := filepath.Join(dir, ".env")
		envLocalPath := filepath.Join(dir, ".env.local")

		if fileExists(examplePath) && !fileExists(envPath) && !fileExists(envLocalPath) {
			data, err := os.ReadFile(examplePath)
			if err != nil {
				return copiedCount, err
			}
			if err := os.WriteFile(envPath, data, 0644); err != nil {
				return copiedCount, err
			}
			copiedCount++
		}
	}

	return copiedCount, nil
}

func parseEnvKeys(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var keys []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) > 0 {
			key := strings.TrimSpace(parts[0])
			if key != "" {
				keys = append(keys, key)
			}
		}
	}
	return keys
}

func parseEnvVariables(path string) map[string]string {
	vars := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return vars
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			vars[key] = val
		}
	}
	return vars
}

var placeholderPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)your_.*_here`),
	regexp.MustCompile(`(?i)change_me`),
	regexp.MustCompile(`(?i)todo`),
	regexp.MustCompile(`(?i)xxx+`),
	regexp.MustCompile(`(?i)replace_me`),
}

func isPlaceholderValue(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return true
	}
	for _, p := range placeholderPatterns {
		if p.MatchString(v) {
			return true
		}
	}
	return false
}
