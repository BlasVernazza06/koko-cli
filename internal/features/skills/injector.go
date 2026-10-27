package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InstalledSkillMetadata representa una skill registrada en koko.config.json.
type InstalledSkillMetadata struct {
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Source      string    `json:"source"`
	InstalledAt time.Time `json:"installedAt"`
}

// InjectSkillToProject inyecta una skill en `.cursor/rules/<skill>.mdc` de forma segura.
func InjectSkillToProject(projectDir string, skill *SkillDefinition) error {
	cleanProjectDir := filepath.Clean(projectDir)

	// Prevenir Path Traversal en el identificador de skill
	rawID := skill.ID
	if strings.Contains(rawID, "..") || strings.Contains(rawID, "/") || strings.Contains(rawID, "\\") {
		return fmt.Errorf("invalid skill ID '%s': path traversal detected", rawID)
	}

	cleanID := filepath.Base(filepath.Clean(rawID))
	if cleanID == "." || cleanID == "/" || cleanID == "\\" || cleanID == "" {
		return fmt.Errorf("invalid skill ID '%s': path traversal detected", skill.ID)
	}

	rulesDir := filepath.Join(cleanProjectDir, ".cursor", "rules")
	// Validar que rulesDir siga perteneciendo al proyecto
	if !strings.HasPrefix(rulesDir, cleanProjectDir) {
		return fmt.Errorf("security violation: target directory outside project boundary")
	}

	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create .cursor/rules directory: %w", err)
	}

	targetFile := filepath.Join(rulesDir, cleanID+".mdc")
	if !strings.HasPrefix(targetFile, rulesDir) {
		return fmt.Errorf("security violation: target file escapes rules directory")
	}

	// Formatear contenido como .mdc con frontmatter
	formattedContent := formatMDCContent(skill)
	if err := os.WriteFile(targetFile, []byte(formattedContent), 0644); err != nil {
		return fmt.Errorf("failed to write skill rule to '%s': %w", targetFile, err)
	}

	// Sincronizar en koko.config.json solo si existe previamente
	_ = syncKokoConfigSkills(cleanProjectDir, skill)

	return nil
}

// CheckIfSkillInstalled verifica si la regla ya existe físicamente en el proyecto.
func CheckIfSkillInstalled(projectDir string, skillID string) bool {
	cleanProjectDir := filepath.Clean(projectDir)
	cleanID := filepath.Base(filepath.Clean(skillID))

	mdcPath := filepath.Join(cleanProjectDir, ".cursor", "rules", cleanID+".mdc")
	if _, err := os.Stat(mdcPath); err == nil {
		return true
	}

	mdPath := filepath.Join(cleanProjectDir, ".cursor", "rules", cleanID+".md")
	if _, err := os.Stat(mdPath); err == nil {
		return true
	}

	return false
}

func formatMDCContent(skill *SkillDefinition) string {
	content := strings.TrimSpace(skill.Content)
	if strings.HasPrefix(content, "---") {
		return content + "\n"
	}

	// Agregar frontmatter estándar si no lo tiene
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("description: %s\n", skill.Description))
	sb.WriteString("globs: [\"*.*\"]\n")
	sb.WriteString("---\n\n")
	sb.WriteString(content)
	sb.WriteString("\n")
	return sb.String()
}

func syncKokoConfigSkills(projectDir string, skill *SkillDefinition) error {
	kokoConfigPath := filepath.Join(projectDir, "koko.config.json")
	data, err := os.ReadFile(kokoConfigPath)
	if err != nil {
		// Si no existe, no crearlo a la fuerza (Regla del PM)
		return nil
	}

	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		return err
	}

	rawSkills, ok := root["skills"].([]any)
	if !ok {
		rawSkills = []any{}
	}

	newSkillEntry := map[string]any{
		"name":        skill.ID,
		"category":    skill.Category,
		"source":      skill.Source,
		"installedAt": time.Now().UTC().Format(time.RFC3339),
	}

	// Evitar duplicados reemplazando la entrada existente
	found := false
	for i, item := range rawSkills {
		if m, ok := item.(map[string]any); ok {
			if name, ok := m["name"].(string); ok && strings.EqualFold(name, skill.ID) {
				rawSkills[i] = newSkillEntry
				found = true
				break
			}
		}
	}
	if !found {
		rawSkills = append(rawSkills, newSkillEntry)
	}

	root["skills"] = rawSkills

	updated, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	updated = append(updated, '\n')

	return os.WriteFile(kokoConfigPath, updated, 0644)
}
