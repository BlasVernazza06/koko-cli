package skills

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed embedded/*
var embeddedSkillsFS embed.FS

// SkillDefinition representa una habilidad de IA con metadatos y contenido.
type SkillDefinition struct {
	ID          string            // Identificador normalizado (ej: "frontend-design", "astro", "turborepo")
	Name        string            // Nombre descriptivo (ej: "Go Backend Architecture")
	Category    string            // Categoría para agrupación (ej: "Go", "Web Frameworks", "Databases & ORMs")
	Description string            // Descripción de la skill
	Content     string            // Contenido markdown principal de la skill (SKILL.md)
	References  map[string]string // Retrocompatibilidad para referencias directas
	Files       map[string]string // Árbol recursivo de archivos de soporte (ruta relativa slash -> contenido)
	Recommended bool              // Si está recomendada para el proyecto detectado
	Installed   bool              // Si ya existe instalada en .cursor/rules/
	Source      string            // Origen: "embedded", "cache", "network"
	IsCombo     bool              // Si es una habilidad que integra múltiples tecnologías
	ComboName   string            // Nombre comercial del combo (ej: "React + shadcn/ui")
}

// GetEmbeddedSkill obtiene una skill precompilada en el binario si existe.
func GetEmbeddedSkill(id string) (*SkillDefinition, bool) {
	normID := strings.ToLower(strings.TrimSpace(id))
	// Aliases conocidos
	if normID == "supabase-postgres-best-practices" {
		normID = "supabase"
	}

	mainFile := filepath.ToSlash(filepath.Join("embedded", normID, "SKILL.md"))
	data, err := embeddedSkillsFS.ReadFile(mainFile)
	if err != nil {
		// Compatibilidad hacia atrás si se intentara leer archivo plano
		flatFile := filepath.ToSlash(filepath.Join("embedded", normID+".md"))
		data, err = embeddedSkillsFS.ReadFile(flatFile)
		if err != nil {
			return nil, false
		}
	}

	content := strings.TrimPrefix(string(data), "\ufeff")

	// Escaneo recursivo de todo el subárbol de soporte de la skill
	files := make(map[string]string)
	refs := make(map[string]string)
	skillRoot := filepath.ToSlash(filepath.Join("embedded", normID))

	_ = fs.WalkDir(embeddedSkillsFS, skillRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}

		// Relativa al directorio raíz de la skill (ej: "references/conn-pooling.md", "core/page-object.md")
		relPath := strings.TrimPrefix(path, skillRoot+"/")
		if relPath == "" || relPath == "SKILL.md" {
			return nil
		}

		fileData, err := embeddedSkillsFS.ReadFile(path)
		if err != nil {
			return nil
		}

		fileContent := strings.TrimPrefix(string(fileData), "\ufeff")
		files[relPath] = fileContent

		// Poblar References con el nombre base para compatibilidad hacia atrás si está en references/
		if strings.HasPrefix(relPath, "references/") {
			refs[filepath.Base(relPath)] = fileContent
		}

		return nil
	})

	return &SkillDefinition{
		ID:          normID,
		Name:        FormatSkillTitle(normID),
		Category:    categorizeSkill(normID),
		Description: extractDescriptionFromFrontmatter(content),
		Content:     content,
		References:  refs,
		Files:       files,
		Source:      "embedded",
	}, true
}

// ListEmbeddedSkills retorna todas las skills embebidas en el binario.
func ListEmbeddedSkills() []*SkillDefinition {
	entries, err := embeddedSkillsFS.ReadDir("embedded")
	if err != nil {
		return nil
	}

	var result []*SkillDefinition
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		if skill, ok := GetEmbeddedSkill(id); ok {
			result = append(result, skill)
		}
	}
	return result
}

// FormatSkillTitle devuelve un título legible para un ID de skill.
func FormatSkillTitle(id string) string {
	switch id {
	// Go / Backend Core
	case "go-backend":
		return "Go Backend Architecture"
	case "go-cli-architecture":
		return "Go CLI & Cobra Architecture"
	case "go-testing":
		return "Go Testing & Benchmarking"

	// Frameworks Web & API & Auth
	case "nextjs":
		return "Next.js App Router Patterns"
	case "react":
		return "Modern React & Hooks"
	case "astro":
		return "Astro Web Framework"
	case "svelte":
		return "Svelte 5 & Runes"
	case "vue":
		return "Vue 3 Composition API & Pinia"
	case "hono":
		return "Hono Multi-Runtime Framework"
	case "fastify":
		return "Fastify High-Performance Server"
	case "express":
		return "Express & TypeScript REST APIs"
	case "nodejs":
		return "Node.js Server Patterns"
	case "nestjs-best-practices":
		return "NestJS Architecture & Best Practices"
	case "better-auth":
		return "Better Auth Modern Framework"
	case "clerk":
		return "Clerk User Management & Auth"
	case "clerk-expo":
		return "Clerk + Expo Authentication"
	case "clerk-nextjs-patterns":
		return "Clerk + Next.js App Router Patterns"
	case "clerk-react-patterns":
		return "Clerk + React Component Patterns"

	// Estilos & UI
	case "tailwindcss":
		return "Tailwind CSS Architecture"
	case "tailwind-v4-shadcn":
		return "Tailwind CSS v4 + shadcn/ui"
	case "shadcn-ui":
		return "shadcn/ui Component System"
	case "frontend-design":
		return "Frontend Design & Interface"
	case "accessibility":
		return "Web Accessibility (A11y)"
	case "seo":
		return "SEO & Metadata"
	case "expo-native-ui":
		return "Expo & React Native UI Patterns"
	case "design-mobile-apps":
		return "Mobile App UX & Ergonomics"

	// Bases de Datos & ORMs
	case "prisma":
		return "Prisma ORM & Schema Modeling"
	case "drizzle":
		return "Drizzle ORM Type-Safe SQL"
	case "postgresql":
		return "PostgreSQL Performance & Tuning"
	case "mongodb":
		return "MongoDB & Mongoose Modeling"
	case "supabase":
		return "Supabase Postgres & BaaS"

	// Infraestructura & Cloud & Payments
	case "cloudflare":
		return "Cloudflare Edge Workers & D1"
	case "docker":
		return "Docker Containerization & Compose"
	case "github-actions":
		return "GitHub Actions CI/CD Pipelines"
	case "turborepo":
		return "Turborepo Monorepo Architecture"
	case "stripe-best-practices":
		return "Stripe Billing & Payments Architecture"

	// Testing
	case "vitest":
		return "Vitest Unit & Integration Testing"
	case "playwright":
		return "Playwright End-to-End Testing"
	case "jest":
		return "Jest Unit & Integration Testing"

	// Lenguajes & Tooling
	case "typescript":
		return "TypeScript Strict Patterns"
	case "zod":
		return "Zod Type-Safe Schema Validation"
	case "python-fastapi":
		return "FastAPI & Modern Python"

	default:
		parts := strings.Split(id, "-")
		for i, p := range parts {
			if len(p) > 0 {
				parts[i] = strings.ToUpper(p[:1]) + p[1:]
			}
		}
		return strings.Join(parts, " ")
	}
}

// categorizeSkill determina la categoría visual para agrupar en la TUI.
func categorizeSkill(id string) string {
	switch id {
	case "go-backend", "go-cli-architecture", "go-testing":
		return "Go"
	case "nextjs", "react", "astro", "svelte", "vue", "hono", "fastify", "express", "nodejs", "nestjs-best-practices", "better-auth", "clerk", "clerk-expo", "clerk-nextjs-patterns", "clerk-react-patterns":
		return "Web Frameworks"
	case "tailwindcss", "tailwind-v4-shadcn", "shadcn-ui", "frontend-design", "accessibility", "seo", "expo-native-ui", "design-mobile-apps":
		return "Styling & UI"
	case "prisma", "drizzle", "postgresql", "mongodb", "supabase":
		return "Databases & ORMs"
	case "cloudflare", "docker", "github-actions", "turborepo", "stripe-best-practices":
		return "Infrastructure & Cloud"
	case "vitest", "playwright", "jest":
		return "Testing"
	case "typescript", "zod", "python-fastapi":
		return "Languages & Tooling"
	default:
		return "General"
	}
}

// extractDescriptionFromFrontmatter extrae la descripción del frontmatter YAML básico.
func extractDescriptionFromFrontmatter(content string) string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "description:") {
			desc := strings.TrimPrefix(trimmed, "description:")
			return strings.Trim(strings.TrimSpace(desc), "\"'")
		}
	}
	return "Curated AI engineering skill and contextual prompt rules."
}
