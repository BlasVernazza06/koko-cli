package skills

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInjectSkillToProject_Success(t *testing.T) {
	tmpDir := t.TempDir()

	skill := &SkillDefinition{
		ID:          "astro",
		Name:        "Astro Web Framework",
		Category:    "Astro",
		Description: "Astro guidelines",
		Content:     "Follow Astro strict conventions.",
		Files: map[string]string{
			"references/routing.md": "Astro routing rules",
			"core/island.md":        "Astro island architecture",
		},
		Source: "embedded",
	}

	if err := InjectSkillToProject(tmpDir, skill); err != nil {
		t.Fatalf("failed to inject skill: %v", err)
	}

	// 1. Verificar regla Cursor en .cursor/rules/astro.mdc
	targetFile := filepath.Join(tmpDir, ".cursor", "rules", "astro.mdc")
	data, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("expected rule file at %s: %v", targetFile, err)
	}

	content := string(data)
	if !strings.Contains(content, "description: Astro guidelines") {
		t.Errorf("expected frontmatter description, got: %s", content)
	}
	if !strings.Contains(content, "Follow Astro strict conventions.") {
		t.Errorf("expected skill body in file")
	}

	// 2. Verificar estructura modular en .cursor/skills/astro/SKILL.md
	skillFile := filepath.Join(tmpDir, ".cursor", "skills", "astro", "SKILL.md")
	skillData, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("expected modular skill file at %s: %v", skillFile, err)
	}
	if !strings.Contains(string(skillData), "Follow Astro strict conventions.") {
		t.Errorf("expected skill content in modular SKILL.md")
	}

	// 3. Verificar referencias en .cursor/skills/astro/references/routing.md
	refFile := filepath.Join(tmpDir, ".cursor", "skills", "astro", "references", "routing.md")
	refData, err := os.ReadFile(refFile)
	if err != nil {
		t.Fatalf("expected modular reference file at %s: %v", refFile, err)
	}
	if !strings.Contains(string(refData), "Astro routing rules") {
		t.Errorf("expected reference content in routing.md")
	}

	// 4. Verificar subcarpeta adicional en .cursor/skills/astro/core/island.md
	islandFile := filepath.Join(tmpDir, ".cursor", "skills", "astro", "core", "island.md")
	islandData, err := os.ReadFile(islandFile)
	if err != nil {
		t.Fatalf("expected modular core file at %s: %v", islandFile, err)
	}
	if !strings.Contains(string(islandData), "Astro island architecture") {
		t.Errorf("expected core content in island.md")
	}

	if !CheckIfSkillInstalled(tmpDir, "astro") {
		t.Errorf("CheckIfSkillInstalled should return true")
	}
}

func TestInjectSkillToProject_PathTraversalSecurity(t *testing.T) {
	tmpDir := t.TempDir()

	maliciousSkills := []*SkillDefinition{
		{ID: "../../evil"},
		{ID: "/etc/passwd"},
		{ID: "..\\..\\evil"},
		{ID: "."},
	}

	for _, ms := range maliciousSkills {
		err := InjectSkillToProject(tmpDir, ms)
		if err == nil {
			t.Errorf("expected path traversal security violation for ID %q, got nil", ms.ID)
		}
	}

	// Probar que archivos maliciosos dentro de skill.Files no escapen
	skillWithMaliciousFiles := &SkillDefinition{
		ID:          "safe-skill",
		Name:        "Safe Skill",
		Category:    "Testing",
		Description: "Testing path traversal inside files map",
		Content:     "Safe content",
		Files: map[string]string{
			"../../escaped.txt":        "should not be written",
			"/etc/evil.txt":            "should not be written",
			"sub/../../escaped2.txt":   "should not be written",
			"..\\..\\escaped3.txt":     "should not be written",
			"valid/sub/deep/file.txt":  "valid content",
		},
		Source: "embedded",
	}

	if err := InjectSkillToProject(tmpDir, skillWithMaliciousFiles); err != nil {
		t.Fatalf("expected successful injection with filtering of malicious files: %v", err)
	}

	escaped1 := filepath.Join(tmpDir, ".cursor", "escaped.txt")
	if _, err := os.Stat(escaped1); err == nil {
		t.Errorf("security violation: escaped file %s was created", escaped1)
	}

	validFile := filepath.Join(tmpDir, ".cursor", "skills", "safe-skill", "valid", "sub", "deep", "file.txt")
	if _, err := os.Stat(validFile); os.IsNotExist(err) {
		t.Errorf("expected valid nested file %s to be created", validFile)
	}
}

func TestInjectSkillToProject_DeepTopology(t *testing.T) {
	tmpDir := t.TempDir()

	playwrightSkill, ok := GetEmbeddedSkill("playwright")
	if !ok {
		t.Fatalf("playwright skill not found in embedded catalog")
	}

	if err := InjectSkillToProject(tmpDir, playwrightSkill); err != nil {
		t.Fatalf("failed to inject playwright skill: %v", err)
	}

	// Verificar existencia de archivos en subcarpetas temáticas profundas
	expectedPaths := []string{
		filepath.Join(tmpDir, ".cursor", "rules", "playwright.mdc"),
		filepath.Join(tmpDir, ".cursor", "skills", "playwright", "SKILL.md"),
		filepath.Join(tmpDir, ".cursor", "skills", "playwright", "core", "page-object-model.md"),
		filepath.Join(tmpDir, ".cursor", "skills", "playwright", "frameworks", "nextjs.md"),
		filepath.Join(tmpDir, ".cursor", "skills", "playwright", "testing-patterns", "accessibility.md"),
	}

	for _, ep := range expectedPaths {
		if _, err := os.Stat(ep); os.IsNotExist(err) {
			t.Errorf("expected deep topology file not found: %s", ep)
		}
	}
}


func TestInjectSkillToProject_SyncKokoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	kokoConfigPath := filepath.Join(tmpDir, "koko.config.json")

	initialConfig := `{
		"name": "my-project",
		"skills": [
			{
				"name": "old-skill",
				"category": "General",
				"source": "network",
				"installedAt": "2026-01-01T00:00:00Z"
			}
		]
	}`

	if err := os.WriteFile(kokoConfigPath, []byte(initialConfig), 0644); err != nil {
		t.Fatalf("failed to write initial config: %v", err)
	}

	skill := &SkillDefinition{
		ID:          "tailwindcss",
		Name:        "Tailwind CSS Architecture",
		Category:    "Tailwind CSS",
		Description: "Tailwind best practices",
		Content:     "Tailwind content",
		Source:      "embedded",
	}

	if err := InjectSkillToProject(tmpDir, skill); err != nil {
		t.Fatalf("failed to inject skill: %v", err)
	}

	updatedData, err := os.ReadFile(kokoConfigPath)
	if err != nil {
		t.Fatalf("failed to read updated config: %v", err)
	}

	var parsed struct {
		Name   string                   `json:"name"`
		Skills []InstalledSkillMetadata `json:"skills"`
	}
	if err := json.Unmarshal(updatedData, &parsed); err != nil {
		t.Fatalf("failed to parse updated config: %v", err)
	}

	if len(parsed.Skills) != 2 {
		t.Fatalf("expected 2 skills in config, got %d", len(parsed.Skills))
	}

	found := false
	for _, s := range parsed.Skills {
		if s.Name == "tailwindcss" {
			found = true
			if s.Category != "Tailwind CSS" {
				t.Errorf("expected category 'Tailwind CSS', got %q", s.Category)
			}
		}
	}
	if !found {
		t.Errorf("tailwindcss skill not found in koko.config.json")
	}

	// Test idempotency / update without duplication
	if err := InjectSkillToProject(tmpDir, skill); err != nil {
		t.Fatalf("failed re-injecting skill: %v", err)
	}

	updatedData2, _ := os.ReadFile(kokoConfigPath)
	_ = json.Unmarshal(updatedData2, &parsed)
	if len(parsed.Skills) != 2 {
		t.Errorf("expected still 2 skills (no duplicate), got %d", len(parsed.Skills))
	}
}

func TestInjectSkillToProject_NoKokoConfigCreatedIfMissing(t *testing.T) {
	tmpDir := t.TempDir()

	skill := &SkillDefinition{
		ID:      "seo",
		Content: "SEO content",
	}

	if err := InjectSkillToProject(tmpDir, skill); err != nil {
		t.Fatalf("failed injecting skill: %v", err)
	}

	kokoConfigPath := filepath.Join(tmpDir, "koko.config.json")
	if _, err := os.Stat(kokoConfigPath); err == nil {
		t.Errorf("koko.config.json should NOT be created if it was missing initially")
	}
}
