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
		Source:      "embedded",
	}

	if err := InjectSkillToProject(tmpDir, skill); err != nil {
		t.Fatalf("failed to inject skill: %v", err)
	}

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
