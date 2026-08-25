package views

import (
	"fmt"
	"strings"
	"time"
)

// RenderCancelled renders the output when the user cancels the operation
func RenderCancelled(msg string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Koko CLI")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n\n", BottomSymbol, StyleMuted.Render(msg)))
	return b.String()
}

// RenderDone renders the final summary (success/error) and next steps
func RenderDone(summary []SummaryItem, projectName, pkgManager string, stepNames, stepStatus []string, elapsed time.Duration, scaffoldErr error) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Creating a new Koko project")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	for _, item := range summary {
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render(item.Label), DotDivider, StyleValue.Render(item.Value)))
	}
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))

	for i, stepName := range stepNames {
		if stepStatus[i] == "success" {
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, CheckSuccess, StylePromptTitle.Render(stepName)))
		} else if stepStatus[i] == "error" {
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, CrossError, StyleError.Render(stepName)))
		} else {
			b.WriteString(fmt.Sprintf("%s  %s  %s\n", BarSymbol, DiamondPending, StyleMuted.Render(stepName)))
		}
	}

	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))

	if scaffoldErr != nil {
		b.WriteString(fmt.Sprintf("%s  %s\n\n", BottomSymbol, StyleError.Render(fmt.Sprintf("Creation error: %v", scaffoldErr))))
		return b.String()
	}

	b.WriteString(fmt.Sprintf("%s  %s\n\n", BottomSymbol, StyleSuccess.Render(fmt.Sprintf("Project created successfully in %.2fs!", elapsed.Seconds()))))

	b.WriteString(fmt.Sprintf("  %s\n", StyleHeader.Render("Next steps:")))
	b.WriteString(fmt.Sprintf("  1. %s\n", StyleValue.Render(fmt.Sprintf("cd %s", projectName))))
	if pkgManager == "pnpm" {
		b.WriteString(fmt.Sprintf("  2. %s\n", StyleValue.Render("pnpm install")))
		b.WriteString(fmt.Sprintf("  3. %s\n\n", StyleValue.Render("pnpm dev")))
	} else if pkgManager == "bun" {
		b.WriteString(fmt.Sprintf("  2. %s\n", StyleValue.Render("bun install")))
		b.WriteString(fmt.Sprintf("  3. %s\n\n", StyleValue.Render("bun dev")))
	} else {
		b.WriteString(fmt.Sprintf("  2. %s\n", StyleValue.Render("npm install")))
		b.WriteString(fmt.Sprintf("  3. %s\n\n", StyleValue.Render("npm run dev")))
	}

	return b.String()
}
