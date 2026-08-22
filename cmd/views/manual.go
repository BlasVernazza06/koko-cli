package views

import (
	"fmt"
	"strings"
)

// RenderManual renderiza la vista interactiva paso a paso de la configuración manual.
// Renderiza dinámicamente el historial acumulado y las opciones del paso actual.
func RenderManual(
	projectName string,
	history []SummaryItem,
	currentStepTitle string,
	options []SelectOption,
	cursor int,
) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Creating a new Koko project")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))

	// 1. Resumen de Nombre del Proyecto
	b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render("Project name"), DotDivider, StyleValue.Render(projectName)))
	b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render("Setup mode"), DotDivider, StyleValue.Render("Configuración Manual")))

	// 2. Historial de tecnologías ya seleccionadas (iterativo sin if anidados)
	for _, item := range history {
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render(item.Label), DotDivider, StyleValue.Render(item.Value)))
	}

	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))

	// 3. Paso actual interactivo
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render(currentStepTitle)))
	b.WriteString(RenderOptions(options, cursor))

	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[↑/↓] Navegar • [Enter] Siguiente • [Esc] Volver atrás")))

	return b.String()
}
