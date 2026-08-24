package views

import (
	"fmt"
	"strings"
)

// RenderManual renders the interactive step-by-step view for manual configuration.
// It dynamically renders the accumulated summary history and current step options.
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

	// 1. Project name summary
	b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render("Project name"), DotDivider, StyleValue.Render(projectName)))
	b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render("Setup mode"), DotDivider, StyleValue.Render("Manual Configuration")))

	// 2. History of previously chosen options
	for _, item := range history {
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render(item.Label), DotDivider, StyleValue.Render(item.Value)))
	}

	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))

	// 3. Current interactive step
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render(currentStepTitle)))
	b.WriteString(RenderOptions(options, cursor))

	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[↑/↓] Navigate • [Enter] Next • [Esc] Go back")))

	return b.String()
}
