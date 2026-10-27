package skills

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetectProjectSkills_PackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	pkgJSON := `{
		"dependencies": {
			"react": "^18.0.0",
			"tailwindcss": "^3.4.0"
		},
		"devDependencies": {
			"typescript": "^5.0.0",
			"vitest": "^1.0.0"
		}
	}`

	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("failed to write package.json: %v", err)
	}

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.HasPackageJSON {
		t.Errorf("expected HasPackageJSON to be true")
	}

	expectedDeps := []string{"react", "tailwindcss", "typescript", "vitest"}
	for _, dep := range expectedDeps {
		if !res.Dependencies[dep] {
			t.Errorf("expected dependency %q to be detected", dep)
		}
	}

	// Verify suggestions include tailwindcss, typescript, vitest, and frontend trio
	suggestedSet := make(map[string]bool)
	for _, s := range res.SuggestedSkills {
		suggestedSet[s] = true
	}

	expectedSuggestions := []string{"tailwindcss", "typescript", "vitest", "frontend-design", "accessibility", "seo"}
	for _, s := range expectedSuggestions {
		if !suggestedSet[s] {
			t.Errorf("expected suggestion %q to be present, got: %v", s, res.SuggestedSkills)
		}
	}
}

func TestDetectProjectSkills_FileTree(t *testing.T) {
	tmpDir := t.TempDir()

	// Create astro.config.mjs, wrangler.toml, assets directory, tsconfig.json
	if err := os.WriteFile(filepath.Join(tmpDir, "astro.config.mjs"), []byte("// astro"), 0644); err != nil {
		t.Fatalf("failed to create astro.config.mjs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "wrangler.toml"), []byte("name = 'worker'"), 0644); err != nil {
		t.Fatalf("failed to create wrangler.toml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create tsconfig.json: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "src", "assets"), 0755); err != nil {
		t.Fatalf("failed to create assets dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "src", "styles.css"), []byte("body {}"), 0644); err != nil {
		t.Fatalf("failed to create styles.css: %v", err)
	}

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.HasAstroConfig {
		t.Errorf("expected HasAstroConfig to be true")
	}
	if !res.HasWranglerConf {
		t.Errorf("expected HasWranglerConf to be true")
	}
	if !res.HasTSConfig {
		t.Errorf("expected HasTSConfig to be true")
	}
	if !res.HasAssetsDir {
		t.Errorf("expected HasAssetsDir to be true")
	}
	if !res.HasCSSFiles {
		t.Errorf("expected HasCSSFiles to be true")
	}

	suggestedSet := make(map[string]bool)
	for _, s := range res.SuggestedSkills {
		suggestedSet[s] = true
	}

	for _, expected := range []string{"astro", "cloudflare", "typescript", "frontend-design", "accessibility", "seo"} {
		if !suggestedSet[expected] {
			t.Errorf("expected suggestion %q, got %v", expected, res.SuggestedSkills)
		}
	}
}

func TestDetectProjectSkills_KokoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	kokoConfig := `{
		"stack": {
			"framework": "astro",
			"backend": "node-express"
		}
	}`

	if err := os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(kokoConfig), 0644); err != nil {
		t.Fatalf("failed to create koko.config.json: %v", err)
	}

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.HasKokoConfig {
		t.Errorf("expected HasKokoConfig to be true")
	}

	suggestedSet := make(map[string]bool)
	for _, s := range res.SuggestedSkills {
		suggestedSet[s] = true
	}

	if !suggestedSet["astro"] {
		t.Errorf("expected astro to be suggested from koko.config.json")
	}
	if !suggestedSet["nodejs"] {
		t.Errorf("expected nodejs to be suggested from koko.config.json")
	}
}

func TestDetectProjectSkills_EmptyDirFallback(t *testing.T) {
	tmpDir := t.TempDir()

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFallback := []string{"frontend-design", "accessibility", "seo"}
	if !reflect.DeepEqual(res.SuggestedSkills, expectedFallback) {
		t.Errorf("expected fallback suggestions %v, got %v", expectedFallback, res.SuggestedSkills)
	}
}

func TestDetectProjectSkills_GoProject(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module github.com/test/my-go-app\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644); err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.HasGoMod {
		t.Errorf("expected HasGoMod to be true")
	}
	if !res.HasGoFiles {
		t.Errorf("expected HasGoFiles to be true")
	}

	suggestedSet := make(map[string]bool)
	for _, s := range res.SuggestedSkills {
		suggestedSet[s] = true
	}

	for _, expected := range []string{"go-backend", "go-cli-architecture", "go-testing"} {
		if !suggestedSet[expected] {
			t.Errorf("expected Go suggestion %q, got %v", expected, res.SuggestedSkills)
		}
	}
}

func TestDetectProjectSkills_RichStack(t *testing.T) {
	tmpDir := t.TempDir()

	pkgContent := `{
		"name": "fullstack-app",
		"dependencies": {
			"next": "^15.0.0",
			"@prisma/client": "^5.0.0",
			"tailwindcss": "^4.0.0",
			"pg": "^8.11.0"
		},
		"devDependencies": {
			"@playwright/test": "^1.40.0"
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgContent), 0644); err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "Dockerfile"), []byte("FROM alpine\n"), 0644); err != nil {
		t.Fatalf("failed to create Dockerfile: %v", err)
	}

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	suggestedSet := make(map[string]bool)
	for _, s := range res.SuggestedSkills {
		suggestedSet[s] = true
	}

	expectedSkills := []string{
		"nextjs",
		"tailwindcss",
		"prisma",
		"postgresql",
		"docker",
		"playwright",
		"frontend-design",
		"accessibility",
		"seo",
	}

	for _, s := range expectedSkills {
		if !suggestedSet[s] {
			t.Errorf("expected rich stack suggestion %q, got %v", s, res.SuggestedSkills)
		}
	}
}

func TestDetectProjectSkills_Combos(t *testing.T) {
	tmpDir := t.TempDir()

	pkgContent := `{
		"name": "combo-app",
		"dependencies": {
			"react": "^19.0.0",
			"clsx": "^2.1.0",
			"class-variance-authority": "^0.7.0",
			"@radix-ui/react-slot": "^1.0.0",
			"react-hook-form": "^7.50.0",
			"zod": "^3.22.0",
			"@supabase/supabase-js": "^2.39.0",
			"better-auth": "^1.0.0",
			"drizzle-orm": "^0.30.0"
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgContent), 0644); err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "turbo.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create turbo.json: %v", err)
	}

	res, err := DetectProjectSkills(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCombos := []string{
		"shadcn-ui",
		"zod",
		"turborepo",
		"better-auth",
	}

	for _, comboID := range expectedCombos {
		combo, ok := res.CombosDetected[comboID]
		if !ok {
			t.Errorf("expected combo %q to be detected, combos: %v", comboID, res.CombosDetected)
		}
		if combo.ComboName == "" {
			t.Errorf("expected combo %q to have non-empty ComboName", comboID)
		}
	}
}
