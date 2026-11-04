package skills

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/tui/views"
)

// Orchestrator gestiona el flujo integral de detección, TUI, provisión híbrida e inyección.
type Orchestrator struct {
	Provider *Provider
	RunTUI   func(items []views.SkillItem, categories []string) ([]string, error)
}

// NewOrchestrator crea un nuevo orquestador de skills nativo.
func NewOrchestrator() (*Orchestrator, error) {
	prov, err := NewProvider()
	if err != nil {
		return nil, err
	}
	return &Orchestrator{
		Provider: prov,
		RunTUI:   RunSkillsTUI,
	}, nil
}

// ExecuteNativeSkills ejecuta el flujo nativo completo de Koko Skills.
func (o *Orchestrator) ExecuteNativeSkills(ctx context.Context, opts Options) error {
	targetDir := opts.Dir
	if targetDir == "" {
		targetDir = "."
	}
	cleanDir := filepath.Clean(targetDir)

	info, err := os.Stat(cleanDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrDirectoryNotFound, cleanDir)
		}
		return fmt.Errorf("failed to check target directory '%s': %w", cleanDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrDirectoryNotADir, cleanDir)
	}

	// 1. Escaneo estático del proyecto y árbol físico
	detection, err := DetectProjectSkills(cleanDir)
	if err != nil {
		return fmt.Errorf("failed to detect project skills: %w", err)
	}

	recommendedSet := make(map[string]bool)
	for _, id := range detection.SuggestedSkills {
		recommendedSet[id] = true
	}

	// 2. Compilar catálogo completo de skills disponibles (embebidas + detectadas)
	allEmbedded := ListEmbeddedSkills()
	skillsMap := make(map[string]*SkillDefinition)
	for _, s := range allEmbedded {
		skillsMap[s.ID] = s
	}

	// Asegurar que todas las sugeridas estén en el mapa
	for _, id := range detection.SuggestedSkills {
		if _, exists := skillsMap[id]; !exists {
			if skill, err := o.Provider.FetchSkill(ctx, id); err == nil {
				skillsMap[id] = skill
			}
		}
	}

	// Orden canónico de categorías para la TUI
	categoryOrder := []string{
		"Go",
		"Web Frameworks",
		"Styling & UI",
		"Databases & ORMs",
		"Infrastructure & Cloud",
		"Testing",
		"Languages & Tooling",
	}

	// Construir items para la TUI con soporte de Combo Skills
	var items []views.SkillItem
	for id, s := range skillsMap {
		installed := CheckIfSkillInstalled(cleanDir, id)
		rec := recommendedSet[id]
		selected := rec && !installed // Recomendadas no instaladas seleccionadas por defecto

		comboMeta, isCombo := detection.CombosDetected[id]
		comboLabel := ""
		if isCombo {
			comboLabel = comboMeta.ComboName
		}

		items = append(items, views.SkillItem{
			ID:          id,
			Name:        s.Name,
			Category:    s.Category,
			Description: s.Description,
			Selected:    selected,
			Installed:   installed,
			Recommended: rec,
			IsCombo:     isCombo,
			ComboLabel:  comboLabel,
		})
	}

	// Ordenar items agrupados por categoría canónica y luego nombre alfabético
	catRank := func(cat string) int {
		for i, c := range categoryOrder {
			if strings.EqualFold(c, cat) {
				return i
			}
		}
		return len(categoryOrder)
	}

	sort.Slice(items, func(i, j int) bool {
		rankI := catRank(items[i].Category)
		rankJ := catRank(items[j].Category)
		if rankI != rankJ {
			return rankI < rankJ
		}
		// Si una es combo y la otra no, priorizar el combo
		if items[i].IsCombo != items[j].IsCombo {
			return items[i].IsCombo
		}
		return items[i].Name < items[j].Name
	})

	// 3. Manejo de flags: Dry-Run
	if opts.DryRun {
		fmt.Printf("\n\033[90m┌\033[0m  \033[1mKoko Skills · Dry Run Analysis\033[0m\n")
		fmt.Printf("\033[90m│\033[0m  Target Directory: \033[38;2;167;139;250m%s\033[0m\n", cleanDir)
		fmt.Printf("\033[90m│\033[0m\n")
		fmt.Printf("\033[90m│\033[0m  \033[1mDetected / Recommended Skills (%d):\033[0m\n", len(detection.SuggestedSkills))
		for _, id := range detection.SuggestedSkills {
			title := FormatSkillTitle(id)
			installedStr := ""
			if CheckIfSkillInstalled(cleanDir, id) {
				installedStr = "\033[90m (already installed)\033[0m"
			}
			comboBadge := ""
			if c, ok := detection.CombosDetected[id]; ok {
				comboBadge = fmt.Sprintf(" \033[38;2;255;184;108m[Combo: %s]\033[0m", c.ComboName)
			}
			fmt.Printf("\033[90m│\033[0m    \033[32m•\033[0m \033[1m%-24s\033[0m \033[90m[%s]\033[0m%s%s\n", title, id, comboBadge, installedStr)
		}
		fmt.Printf("\033[90m│\033[0m\n")
		fmt.Printf("\033[90m└\033[0m  \033[33mDry run completed: no files were written to disk.\033[0m\n\n")
		return nil
	}

	// 4. Determinar lista de skills a instalar
	var selectedIDs []string

	if opts.Yes {
		// Modo automático no interactivo: seleccionar recomendadas no instaladas (o todas las recomendadas)
		for _, it := range items {
			if it.Recommended {
				selectedIDs = append(selectedIDs, it.ID)
			}
		}
	} else {
		// Abrir la TUI interactiva con Bubble Tea
		tuiFn := o.RunTUI
		if tuiFn == nil {
			tuiFn = RunSkillsTUI
		}
		chosen, err := tuiFn(items, categoryOrder)
		if err != nil {
			return err
		}
		if len(chosen) == 0 {
			fmt.Printf("\n\033[90m┌\033[0m  \033[1mKoko Skills\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[33mNo skills selected. Exiting without making changes.\033[0m\n\n")
			return nil
		}
		selectedIDs = chosen
	}

	// 5. Inyección de skills seleccionadas
	fmt.Printf("\n\033[90m┌\033[0m  \033[1mKoko Skills · Provisioning %d skill(s)\033[0m\n", len(selectedIDs))
	fmt.Printf("\033[90m│\033[0m\n")

	for _, id := range selectedIDs {
		skill, err := o.Provider.FetchSkill(ctx, id)
		if err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗\033[0m  Failed to fetch skill '%s': %v\n", id, err)
			continue
		}

		if err := InjectSkillToProject(cleanDir, skill); err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗\033[0m  Failed to inject '%s': %v\n", id, err)
			continue
		}

		sourceLabel := skill.Source
		if sourceLabel == "" {
			sourceLabel = "embedded"
		}
		comboNotice := ""
		if c, ok := detection.CombosDetected[id]; ok {
			comboNotice = fmt.Sprintf(" (%s)", c.ComboName)
		}
		fmt.Printf("\033[90m│\033[0m  \033[32m✓\033[0m  Installed \033[1m%-24s\033[0m \033[90m(.cursor/rules/%s.mdc, via %s)%s\033[0m\n", skill.Name, skill.ID, sourceLabel, comboNotice)
	}

	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m└\033[0m  \033[32m\033[1mAll requested AI agent skills were successfully configured!\033[0m\n\n")

	return nil
}
