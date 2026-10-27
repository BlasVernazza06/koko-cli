package skills

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ComboSkillDefinition define una habilidad combo resultante de la confluencia de tecnologías.
type ComboSkillDefinition struct {
	ID        string // Skill ID a sugerir
	ComboName string // Nombre visible del combo (ej: "React + shadcn/ui")
	Reason    string // Explicación de la sinergia
}

// DetectionResult consolida los hallazgos tras escanear un directorio de proyecto.
type DetectionResult struct {
	HasPackageJSON   bool
	Dependencies     map[string]bool
	HasKokoConfig    bool
	KokoConfigStack  []string

	// Go detection
	HasGoMod         bool
	HasGoFiles       bool

	// Frontend files
	HasHTMLFiles     bool
	HasCSSFiles      bool
	HasAssetsDir     bool
	HasStylesDir     bool

	// Web frameworks configs
	HasAstroConfig   bool
	HasNextConfig    bool
	HasSvelteConfig  bool
	HasViteConfig    bool
	HasTailwindConf  bool
	HasTSConfig      bool

	// Infra, Monorepo & DevOps
	HasTurboConfig   bool
	HasPnpmWorkspace bool
	HasWranglerConf  bool
	HasDockerFiles   bool
	HasGithubActions bool

	// Testing configs
	HasVitestConf     bool
	HasPlaywrightConf bool
	HasJestConf       bool

	SuggestedSkills   []string
	CombosDetected    map[string]ComboSkillDefinition // skillID -> Combo metadata
}

// DetectProjectSkills escanea archivos de configuración y el árbol físico del directorio.
func DetectProjectSkills(projectDir string) (*DetectionResult, error) {
	cleanDir := filepath.Clean(projectDir)
	res := &DetectionResult{
		Dependencies:   make(map[string]bool),
		CombosDetected: make(map[string]ComboSkillDefinition),
	}

	// 1. Escanear package.json si existe
	pkgPath := filepath.Join(cleanDir, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		res.HasPackageJSON = true
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			for dep := range pkg.Dependencies {
				res.Dependencies[strings.ToLower(dep)] = true
			}
			for dep := range pkg.DevDependencies {
				res.Dependencies[strings.ToLower(dep)] = true
			}
		}
	}

	// 2. Escanear koko.config.json si existe
	kokoPath := filepath.Join(cleanDir, "koko.config.json")
	if data, err := os.ReadFile(kokoPath); err == nil {
		res.HasKokoConfig = true
		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err == nil {
			extractKokoStackTokens(raw, &res.KokoConfigStack)
		}
	}

	// 3. Inspeccionar árbol físico de archivos (hasta 4 niveles de profundidad para cold start veloz)
	_ = filepath.WalkDir(cleanDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, relErr := filepath.Rel(cleanDir, path)
		if relErr != nil {
			return nil
		}
		// Ignorar node_modules, .git, dist, build, .cursor
		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || name == ".git" || name == "dist" || name == "build" || name == ".next" || name == ".turbo" {
				return filepath.SkipDir
			}
			lowerName := strings.ToLower(name)
			if lowerName == "assets" || lowerName == "public" || lowerName == "static" {
				res.HasAssetsDir = true
			}
			if lowerName == "styles" || lowerName == "css" {
				res.HasStylesDir = true
			}
			if lowerName == "workflows" && strings.Contains(filepath.ToSlash(rel), ".github") {
				res.HasGithubActions = true
			}
			// Limitar profundidad de recorrido para evitar demoras
			if strings.Count(filepath.ToSlash(rel), "/") > 4 {
				return filepath.SkipDir
			}
			return nil
		}

		lower := strings.ToLower(d.Name())

		// Monorepo & Workspaces
		if lower == "turbo.json" {
			res.HasTurboConfig = true
		}
		if lower == "pnpm-workspace.yaml" || lower == "pnpm-workspace.yml" {
			res.HasPnpmWorkspace = true
		}

		// Go detection
		if lower == "go.mod" {
			res.HasGoMod = true
		}
		if strings.HasSuffix(lower, ".go") {
			res.HasGoFiles = true
		}

		// Frontend static files
		if strings.HasSuffix(lower, ".html") {
			res.HasHTMLFiles = true
		}
		if strings.HasSuffix(lower, ".css") || strings.HasSuffix(lower, ".scss") || strings.HasSuffix(lower, ".sass") {
			res.HasCSSFiles = true
		}

		// Framework configs
		if strings.HasPrefix(lower, "astro.config") {
			res.HasAstroConfig = true
		}
		if strings.HasPrefix(lower, "next.config") {
			res.HasNextConfig = true
		}
		if strings.HasPrefix(lower, "svelte.config") {
			res.HasSvelteConfig = true
		}
		if strings.HasPrefix(lower, "vite.config") {
			res.HasViteConfig = true
		}
		if strings.HasPrefix(lower, "tailwind.config") || strings.Contains(lower, "tailwind") {
			res.HasTailwindConf = true
		}
		if lower == "tsconfig.json" || strings.HasSuffix(lower, ".ts") || strings.HasSuffix(lower, ".tsx") {
			res.HasTSConfig = true
		}

		// Cloud & Infra
		if lower == "wrangler.toml" || lower == "wrangler.json" {
			res.HasWranglerConf = true
		}
		if lower == "dockerfile" || strings.HasPrefix(lower, "dockerfile.") || strings.HasPrefix(lower, "docker-compose") {
			res.HasDockerFiles = true
		}

		// Testing configs & files
		if strings.Contains(lower, "vitest.config") || strings.Contains(lower, ".test.ts") || strings.Contains(lower, ".spec.ts") {
			res.HasVitestConf = true
		}
		if strings.Contains(lower, "playwright.config") {
			res.HasPlaywrightConf = true
		}
		if strings.Contains(lower, "jest.config") {
			res.HasJestConf = true
		}

		return nil
	})

	// 4. Evaluar Combos
	evaluateCombos(res)

	// 5. Determinar skills sugeridas de forma determinista y exhaustiva
	res.SuggestedSkills = evaluateSuggestions(res)
	return res, nil
}

func extractKokoStackTokens(node any, tokens *[]string) {
	switch v := node.(type) {
	case map[string]any:
		for _, val := range v {
			extractKokoStackTokens(val, tokens)
		}
	case []any:
		for _, val := range v {
			extractKokoStackTokens(val, tokens)
		}
	case string:
		*tokens = append(*tokens, strings.ToLower(v))
	}
}

func evaluateCombos(res *DetectionResult) {
	hasReact := res.Dependencies["react"] || res.HasNextConfig || res.Dependencies["next"] || res.HasViteConfig
	hasShadcn := res.Dependencies["@radix-ui/react-slot"] || res.Dependencies["class-variance-authority"] || res.Dependencies["clsx"]
	hasTailwind := res.HasTailwindConf || res.Dependencies["tailwindcss"] || res.Dependencies["@tailwindcss/vite"]
	hasClerk := res.Dependencies["@clerk/clerk-react"] || res.Dependencies["@clerk/nextjs"] || res.Dependencies["@clerk/clerk-expo"]
	hasExpo := res.Dependencies["expo"] || res.Dependencies["react-native"]

	// Combo 1: React + shadcn/ui
	if hasReact && hasShadcn {
		res.CombosDetected["shadcn-ui"] = ComboSkillDefinition{
			ID:        "shadcn-ui",
			ComboName: "React + shadcn/ui",
			Reason:    "Detected React framework combined with Radix UI / Tailwind component primitives",
		}
	}

	// Combo 2: Tailwind v4 + shadcn/ui
	if hasTailwind && hasShadcn {
		res.CombosDetected["tailwind-v4-shadcn"] = ComboSkillDefinition{
			ID:        "tailwind-v4-shadcn",
			ComboName: "Tailwind CSS v4 + shadcn/ui",
			Reason:    "Detected Tailwind utility styling coupled with shadcn component architecture",
		}
	}

	// Combo 3: Clerk + Expo
	if hasClerk && hasExpo {
		res.CombosDetected["clerk-expo"] = ComboSkillDefinition{
			ID:        "clerk-expo",
			ComboName: "Clerk + Expo Auth",
			Reason:    "Detected Clerk authentication integrated into React Native / Expo environment",
		}
	}

	// Combo 4: Clerk + Next.js
	if hasClerk && (res.HasNextConfig || res.Dependencies["next"] || res.Dependencies["@clerk/nextjs"]) {
		res.CombosDetected["clerk-nextjs-patterns"] = ComboSkillDefinition{
			ID:        "clerk-nextjs-patterns",
			ComboName: "Clerk + Next.js App Router",
			Reason:    "Detected Clerk authentication middleware and Server Components in Next.js",
		}
	}

	// Combo 5: Clerk + React Patterns
	if hasClerk && hasReact && !res.HasNextConfig && !hasExpo {
		res.CombosDetected["clerk-react-patterns"] = ComboSkillDefinition{
			ID:        "clerk-react-patterns",
			ComboName: "Clerk + React Components",
			Reason:    "Detected Clerk provider with client-side React SPA",
		}
	}

	// Combo 6: React Hook Form + Zod
	if res.Dependencies["react-hook-form"] && (res.Dependencies["zod"] || res.Dependencies["@hookform/resolvers"]) {
		res.CombosDetected["zod"] = ComboSkillDefinition{
			ID:        "zod",
			ComboName: "React Hook Form + Zod",
			Reason:    "Detected form state management paired with schema validation resolver",
		}
	}

	// Combo 7: Next.js + Supabase
	hasSupabase := res.Dependencies["@supabase/supabase-js"] || res.Dependencies["@supabase/ssr"]
	if (res.HasNextConfig || res.Dependencies["next"]) && hasSupabase {
		res.CombosDetected["supabase"] = ComboSkillDefinition{
			ID:        "supabase",
			ComboName: "Next.js + Supabase",
			Reason:    "Detected Next.js Server Components paired with Supabase SSR cookie auth & Postgres",
		}
	}

	// Combo 8: Turborepo Monorepo Architecture
	if res.HasTurboConfig || res.HasPnpmWorkspace {
		res.CombosDetected["turborepo"] = ComboSkillDefinition{
			ID:        "turborepo",
			ComboName: "Turborepo Monorepo Architecture",
			Reason:    "Detected monorepo configuration with task pipeline and workspace boundaries",
		}
	}

	// Combo 9: Better Auth + Database Integration
	hasBetterAuth := res.Dependencies["better-auth"]
	hasDBORM := res.Dependencies["drizzle-orm"] || res.Dependencies["@prisma/client"] || res.Dependencies["prisma"] || res.Dependencies["pg"]
	if hasBetterAuth && hasDBORM {
		res.CombosDetected["better-auth"] = ComboSkillDefinition{
			ID:        "better-auth",
			ComboName: "Better Auth + DB Adapter",
			Reason:    "Detected Better Auth coupled with an ORM database adapter (Drizzle/Prisma)",
		}
	}

	// Combo 10: Expo Native UI & Mobile Design
	if hasExpo {
		res.CombosDetected["expo-native-ui"] = ComboSkillDefinition{
			ID:        "expo-native-ui",
			ComboName: "Expo + React Native UI",
			Reason:    "Detected Expo application with cross-platform native UI primitives",
		}
		res.CombosDetected["design-mobile-apps"] = ComboSkillDefinition{
			ID:        "design-mobile-apps",
			ComboName: "Mobile App UX Ergonomics",
			Reason:    "Detected mobile touch interface standards and ergonomic requirements",
		}
	}
}

func evaluateSuggestions(res *DetectionResult) []string {
	suggestions := make(map[string]bool)

	// Reglas Go
	if res.HasGoMod || res.HasGoFiles {
		suggestions["go-backend"] = true
		suggestions["go-cli-architecture"] = true
		suggestions["go-testing"] = true
	}

	// Regla Frontend Core: si hay HTML, CSS, assets, styles o librerías web
	if res.HasHTMLFiles || res.HasCSSFiles || res.HasAssetsDir || res.HasStylesDir ||
		res.Dependencies["react"] || res.Dependencies["next"] || res.Dependencies["vue"] ||
		res.Dependencies["svelte"] || res.Dependencies["astro"] {
		suggestions["frontend-design"] = true
		suggestions["accessibility"] = true
		suggestions["seo"] = true
	}

	// Reglas Frameworks Web
	if res.HasNextConfig || res.Dependencies["next"] {
		suggestions["nextjs"] = true
	}
	if res.HasViteConfig || res.Dependencies["react"] || res.Dependencies["react-dom"] {
		suggestions["react"] = true
	}
	if res.HasAstroConfig || res.Dependencies["astro"] {
		suggestions["astro"] = true
	}
	if res.HasSvelteConfig || res.Dependencies["svelte"] || res.Dependencies["@sveltejs/kit"] {
		suggestions["svelte"] = true
	}
	if res.Dependencies["vue"] || res.Dependencies["nuxt"] {
		suggestions["vue"] = true
	}

	// Reglas Backend API
	if res.Dependencies["hono"] {
		suggestions["hono"] = true
	}
	if res.Dependencies["fastify"] {
		suggestions["fastify"] = true
	}
	if res.Dependencies["express"] {
		suggestions["express"] = true
	}
	if res.Dependencies["@nestjs/core"] || res.Dependencies["@nestjs/common"] {
		suggestions["nestjs-best-practices"] = true
	}
	if res.HasPackageJSON || res.Dependencies["express"] || res.Dependencies["hono"] || res.Dependencies["dotenv"] {
		suggestions["nodejs"] = true
	}

	// Reglas Auth
	if res.Dependencies["better-auth"] {
		suggestions["better-auth"] = true
	}
	if res.Dependencies["@clerk/nextjs"] || res.Dependencies["@clerk/clerk-react"] || res.Dependencies["@clerk/clerk-expo"] {
		suggestions["clerk"] = true
	}

	// Reglas Estilos & UI
	if res.HasTailwindConf || res.Dependencies["tailwindcss"] || res.Dependencies["@tailwindcss/vite"] {
		suggestions["tailwindcss"] = true
	}
	if res.Dependencies["@radix-ui/react-slot"] || res.Dependencies["class-variance-authority"] || res.Dependencies["clsx"] {
		suggestions["shadcn-ui"] = true
	}

	// Reglas Mobile
	if res.Dependencies["expo"] || res.Dependencies["react-native"] {
		suggestions["expo-native-ui"] = true
		suggestions["design-mobile-apps"] = true
	}

	// Reglas Payments
	if res.Dependencies["stripe"] || res.Dependencies["@stripe/stripe-js"] {
		suggestions["stripe-best-practices"] = true
	}

	// Reglas Databases & ORMs
	if res.Dependencies["@prisma/client"] || res.Dependencies["prisma"] {
		suggestions["prisma"] = true
	}
	if res.Dependencies["drizzle-orm"] || res.Dependencies["drizzle-kit"] {
		suggestions["drizzle"] = true
	}
	if res.Dependencies["pg"] || res.Dependencies["postgres"] || res.Dependencies["@neondatabase/serverless"] {
		suggestions["postgresql"] = true
	}
	if res.Dependencies["mongoose"] || res.Dependencies["mongodb"] {
		suggestions["mongodb"] = true
	}
	if res.Dependencies["@supabase/supabase-js"] || res.Dependencies["@supabase/ssr"] {
		suggestions["supabase"] = true
	}

	// Reglas Infraestructura, Monorepo & Cloud
	if res.HasTurboConfig || res.HasPnpmWorkspace || res.Dependencies["turbo"] {
		suggestions["turborepo"] = true
	}
	if res.HasWranglerConf || res.Dependencies["wrangler"] || res.Dependencies["@cloudflare/workers-types"] {
		suggestions["cloudflare"] = true
	}
	if res.HasDockerFiles {
		suggestions["docker"] = true
	}
	if res.HasGithubActions {
		suggestions["github-actions"] = true
	}

	// Reglas Testing
	if res.HasVitestConf || res.Dependencies["vitest"] {
		suggestions["vitest"] = true
	}
	if res.HasPlaywrightConf || res.Dependencies["@playwright/test"] {
		suggestions["playwright"] = true
	}
	if res.HasJestConf || res.Dependencies["jest"] || res.Dependencies["ts-jest"] {
		suggestions["jest"] = true
	}

	// Reglas Lenguajes & Schemas
	if res.HasTSConfig || res.Dependencies["typescript"] {
		suggestions["typescript"] = true
	}
	if res.Dependencies["zod"] {
		suggestions["zod"] = true
	}

	// Inyectar Combos detectados directamente a sugerencias
	for comboID := range res.CombosDetected {
		suggestions[comboID] = true
	}

	// Inspección de tokens de koko.config.json para enriquecimiento
	for _, tok := range res.KokoConfigStack {
		if strings.Contains(tok, "go") {
			suggestions["go-backend"] = true
			suggestions["go-cli-architecture"] = true
			suggestions["go-testing"] = true
		}
		if strings.Contains(tok, "next") {
			suggestions["nextjs"] = true
		}
		if strings.Contains(tok, "react") {
			suggestions["react"] = true
		}
		if strings.Contains(tok, "astro") {
			suggestions["astro"] = true
		}
		if strings.Contains(tok, "svelte") {
			suggestions["svelte"] = true
		}
		if strings.Contains(tok, "vue") || strings.Contains(tok, "nuxt") {
			suggestions["vue"] = true
		}
		if strings.Contains(tok, "hono") {
			suggestions["hono"] = true
		}
		if strings.Contains(tok, "fastify") {
			suggestions["fastify"] = true
		}
		if strings.Contains(tok, "nestjs") {
			suggestions["nestjs-best-practices"] = true
		}
		if strings.Contains(tok, "node") || strings.Contains(tok, "express") || strings.Contains(tok, "hono") || strings.Contains(tok, "fastify") {
			suggestions["nodejs"] = true
		}
		if strings.Contains(tok, "better-auth") {
			suggestions["better-auth"] = true
		}
		if strings.Contains(tok, "clerk") {
			suggestions["clerk"] = true
		}
		if strings.Contains(tok, "tailwind") {
			suggestions["tailwindcss"] = true
		}
		if strings.Contains(tok, "shadcn") {
			suggestions["shadcn-ui"] = true
		}
		if strings.Contains(tok, "prisma") {
			suggestions["prisma"] = true
		}
		if strings.Contains(tok, "drizzle") {
			suggestions["drizzle"] = true
		}
		if strings.Contains(tok, "postgres") {
			suggestions["postgresql"] = true
		}
		if strings.Contains(tok, "mongo") {
			suggestions["mongodb"] = true
		}
		if strings.Contains(tok, "supabase") {
			suggestions["supabase"] = true
		}
		if strings.Contains(tok, "turbo") || strings.Contains(tok, "monorepo") {
			suggestions["turborepo"] = true
		}
		if strings.Contains(tok, "docker") {
			suggestions["docker"] = true
		}
		if strings.Contains(tok, "github") || strings.Contains(tok, "cicd") {
			suggestions["github-actions"] = true
		}
		if strings.Contains(tok, "stripe") {
			suggestions["stripe-best-practices"] = true
		}
		if strings.Contains(tok, "expo") || strings.Contains(tok, "native") {
			suggestions["expo-native-ui"] = true
			suggestions["design-mobile-apps"] = true
		}
		if strings.Contains(tok, "vitest") {
			suggestions["vitest"] = true
		}
		if strings.Contains(tok, "playwright") {
			suggestions["playwright"] = true
		}
		if strings.Contains(tok, "typescript") {
			suggestions["typescript"] = true
		}
		if strings.Contains(tok, "zod") {
			suggestions["zod"] = true
		}
		if strings.Contains(tok, "fastapi") || strings.Contains(tok, "python") {
			suggestions["python-fastapi"] = true
		}
	}

	// Orden canónico por categorías para presentación consistente
	canonicalOrder := []string{
		// Go
		"go-backend", "go-cli-architecture", "go-testing",
		// Web Frameworks
		"nextjs", "react", "astro", "svelte", "vue", "hono", "fastify", "express", "nodejs", "nestjs-best-practices", "better-auth", "clerk", "clerk-expo", "clerk-nextjs-patterns", "clerk-react-patterns",
		// Styling & UI
		"tailwindcss", "tailwind-v4-shadcn", "shadcn-ui", "frontend-design", "accessibility", "seo", "expo-native-ui", "design-mobile-apps",
		// Databases & ORMs
		"prisma", "drizzle", "postgresql", "mongodb", "supabase",
		// Infrastructure & Cloud
		"turborepo", "cloudflare", "docker", "github-actions", "stripe-best-practices",
		// Testing
		"vitest", "playwright", "jest",
		// Languages & Tooling
		"typescript", "zod", "python-fastapi",
	}

	var ordered []string
	for _, id := range canonicalOrder {
		if suggestions[id] {
			ordered = append(ordered, id)
		}
	}

	// Si no se detectó nada específico, sugerir el trío base de frontend
	if len(ordered) == 0 {
		ordered = []string{"frontend-design", "accessibility", "seo"}
	}

	return ordered
}
