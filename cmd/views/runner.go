package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
)

// RenderRunning renders the project creation progress with spinners
func RenderRunning(summary []SummaryItem, stepNames, stepStatus []string, s spinner.Model) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Creating a new Koko project")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	for _, item := range summary {
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render(item.Label), DotDivider, StyleValue.Render(item.Value)))
	}
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))

	for i, stepName := range stepNames {
		switch stepStatus[i] {
		case "pending":
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, DiamondPending, StyleMuted.Render(stepName)))
		case "running":
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, s.View(), StyleActiveItem.Render(stepName)))
		case "success":
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, CheckSuccess, StylePromptTitle.Render(stepName)))
		case "error":
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, CrossError, StyleError.Render(stepName)))
		}
	}
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("Scaffolding in progress...")))
	return b.String()
}
