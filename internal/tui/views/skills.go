package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SkillItem representa un item para la vista de selección de skills.
type SkillItem struct {
	ID          string
	Name        string
	Category    string
	Description string
	Selected    bool
	Installed   bool
	Recommended bool
	IsCombo     bool
	ComboLabel  string
}

// RenderSkillsMenu renderiza la interfaz TUI enriquecida para selección de skills.
func RenderSkillsMenu(skills []SkillItem, cursor int, categories []string) string {
	var b strings.Builder

	selectedCount := 0
	for _, s := range skills {
		if s.Selected {
			selectedCount++
		}
	}

	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Koko Skills · AI Agent Knowledge Manager")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render("Select AI agent skills to configure:")))

	// Agrupar items por categoría
	catMap := make(map[string][]int)
	for i, s := range skills {
		catMap[s.Category] = append(catMap[s.Category], i)
	}

	for _, cat := range categories {
		indices := catMap[cat]
		if len(indices) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf("%s\n%s  %s\n", BarSymbol, BarSymbol, StyleCategoryHeader.Render("── "+cat+" ──")))

		for _, idx := range indices {
			item := skills[idx]
			isActive := (idx == cursor)

			var checkIndicator string
			if item.Selected {
				if isActive {
					checkIndicator = CheckboxCheckedActive
				} else {
					checkIndicator = CheckboxChecked
				}
			} else {
				if isActive {
					checkIndicator = CheckboxUncheckedActive
				} else {
					checkIndicator = CheckboxUnchecked
				}
			}

			var nameStyle lipgloss.Style
			if isActive {
				nameStyle = StyleActiveItem
			} else {
				nameStyle = StyleInactiveItem
			}

			var statusBadges []string
			if item.IsCombo {
				comboText := "[Combo]"
				if item.ComboLabel != "" {
					comboText = fmt.Sprintf("[Combo: %s]", item.ComboLabel)
				}
				statusBadges = append(statusBadges, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Bold(true).Render(comboText))
			}

			if item.Installed {
				statusBadges = append(statusBadges, lipgloss.NewStyle().Foreground(ColorLightViolet).Render("[installed]"))
			} else if item.Recommended {
				statusBadges = append(statusBadges, lipgloss.NewStyle().Foreground(ColorGreenCheck).Render("[recommended]"))
			}

			line := fmt.Sprintf("%s  %s %s", BarSymbol, checkIndicator, nameStyle.Render(item.Name))
			if len(statusBadges) > 0 {
				line += " " + strings.Join(statusBadges, " ")
			}

			if item.Description != "" {
				line += " " + StyleHint.Render("· "+item.Description)
			}

			b.WriteString(line + "\n")
		}
	}

	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	footerNav := fmt.Sprintf("[↑↓] move • [space] toggle • [a] all • [n] new • [i] installed • [enter] confirm (%d/%d)", selectedCount, len(skills))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render(footerNav)))

	return b.String()
}
