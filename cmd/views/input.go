package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

// RenderInput renders the project name input prompt with optional validation error
func RenderInput(ti textinput.Model, errMsg string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Creating a new Koko project")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render("Project name")))
	b.WriteString(fmt.Sprintf("%s  %s\n", BarSymbol, ti.View()))
	if errMsg != "" {
		b.WriteString(fmt.Sprintf("%s  %s %s\n", BarSymbol, CrossError, StyleError.Render(errMsg)))
	}
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[Enter] Continue • [Esc] Main menu")))
	return b.String()
}
