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

// InjectSkillToProject inyecta una skill de forma modular en:
// 1. `.cursor/skills/<cleanID>/SKILL.md` (y `references/*.md` opcionales)
// 2. `.cursor/rules/<cleanID>.mdc` (regla activa para Cursor con frontmatter)
// Garantiza estricta validación contra Path Traversal.
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

	// 1. Inyección de estructura modular en .cursor/skills/<cleanID>/
	skillModularDir := filepath.Join(cleanProjectDir, ".cursor", "skills", cleanID)
	if !strings.HasPrefix(skillModularDir, cleanProjectDir) {
		return fmt.Errorf("security violation: skill directory outside project boundary")
	}

	if err := os.MkdirAll(skillModularDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory '%s': %w", skillModularDir, err)
	}

	skillMDPath := filepath.Join(skillModularDir, "SKILL.md")
	if !strings.HasPrefix(skillMDPath, skillModularDir) {
		return fmt.Errorf("security violation: skill file escapes skill directory")
	}

	if err := os.WriteFile(skillMDPath, []byte(strings.TrimSpace(skill.Content)+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write SKILL.md to '%s': %w", skillMDPath, err)
	}

	// Inyectar árbol recursivo completo de archivos de soporte (references/, rules/, core/, etc.)
	if len(skill.Files) > 0 {
		for relPath, fileContent := range skill.Files {
			// Prevenir path traversal en cada componente de la ruta
			cleanRel := filepath.Clean(filepath.ToSlash(relPath))
			if strings.HasPrefix(cleanRel, "..") || strings.Contains(cleanRel, "/../") || strings.HasPrefix(cleanRel, "/") || strings.HasPrefix(cleanRel, "\\") {
				continue
			}

			targetFilePath := filepath.Join(skillModularDir, filepath.FromSlash(cleanRel))
			if !strings.HasPrefix(targetFilePath, skillModularDir) {
				continue
			}

			parentDir := filepath.Dir(targetFilePath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				continue
			}

			_ = os.WriteFile(targetFilePath, []byte(strings.TrimSpace(fileContent)+"\n"), 0644)
		}
	} else if len(skill.References) > 0 {
		// Fallback para instancias construidas en tests que solo definan skill.References
		refDir := filepath.Join(skillModularDir, "references")
		if strings.HasPrefix(refDir, skillModularDir) {
			_ = os.MkdirAll(refDir, 0755)
			for refName, refContent := range skill.References {
				if strings.Contains(refName, "..") || strings.Contains(refName, "/") || strings.Contains(refName, "\\") {
					continue
				}
				cleanRefName := filepath.Base(filepath.Clean(refName))
				targetRefFile := filepath.Join(refDir, cleanRefName)
				if strings.HasPrefix(targetRefFile, refDir) {
					_ = os.WriteFile(targetRefFile, []byte(strings.TrimSpace(refContent)+"\n"), 0644)
				}
			}
		}
	}

	// 2. Inyección de regla Cursor en .cursor/rules/<cleanID>.mdc
	rulesDir := filepath.Join(cleanProjectDir, ".cursor", "rules")
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

	// Verifica regla Cursor o archivo de skill modular
	mdcPath := filepath.Join(cleanProjectDir, ".cursor", "rules", cleanID+".mdc")
	if _, err := os.Stat(mdcPath); err == nil {
		return true
	}

	modularSkillPath := filepath.Join(cleanProjectDir, ".cursor", "skills", cleanID, "SKILL.md")
	if _, err := os.Stat(modularSkillPath); err == nil {
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
