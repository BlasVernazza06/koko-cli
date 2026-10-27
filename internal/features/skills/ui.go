package skills

import (
	"fmt"

	"github.com/BlasVernazza06/koko-cli/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

// SkillsTUIModel gestiona el estado del modelo interactivo de Bubble Tea para skills.
type SkillsTUIModel struct {
	Items      []views.SkillItem
	Categories []string
	Cursor     int
	Confirmed  bool
	Quitting   bool
}

// NewSkillsTUIModel inicializa el modelo de Bubble Tea.
func NewSkillsTUIModel(items []views.SkillItem, categories []string) *SkillsTUIModel {
	return &SkillsTUIModel{
		Items:      items,
		Categories: categories,
		Cursor:     0,
		Confirmed:  false,
		Quitting:   false,
	}
}

func (m *SkillsTUIModel) Init() tea.Cmd {
	return nil
}

func (m *SkillsTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Items) - 1
			}

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			} else {
				m.Cursor = 0
			}

		case " ":
			if m.Cursor >= 0 && m.Cursor < len(m.Items) {
				m.Items[m.Cursor].Selected = !m.Items[m.Cursor].Selected
			}

		case "a":
			// Toggle all: si todos seleccionados, deseleccionar todos; si no, seleccionar todos
			allSelected := true
			for _, it := range m.Items {
				if !it.Selected {
					allSelected = false
					break
				}
			}
			for i := range m.Items {
				m.Items[i].Selected = !allSelected
			}

		case "n":
			// Seleccionar solo los "new" (no instalados)
			for i := range m.Items {
				m.Items[i].Selected = !m.Items[i].Installed
			}

		case "i":
			// Seleccionar solo los "installed"
			for i := range m.Items {
				m.Items[i].Selected = m.Items[i].Installed
			}

		case "enter":
			m.Confirmed = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *SkillsTUIModel) View() string {
	if m.Quitting && !m.Confirmed {
		return "\n  Operation cancelled.\n\n"
	}
	return views.RenderSkillsMenu(m.Items, m.Cursor, m.Categories)
}

// RunSkillsTUI ejecuta el programa Bubble Tea y retorna los IDs de las skills seleccionadas.
func RunSkillsTUI(items []views.SkillItem, categories []string) ([]string, error) {
	model := NewSkillsTUIModel(items, categories)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run skills TUI: %w", err)
	}

	resModel, ok := finalModel.(*SkillsTUIModel)
	if !ok || !resModel.Confirmed {
		return nil, nil // Cancelado o salido
	}

	var selectedIDs []string
	for _, it := range resModel.Items {
		if it.Selected {
			selectedIDs = append(selectedIDs, it.ID)
		}
	}

	return selectedIDs, nil
}
