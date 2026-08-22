package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

// RenderInput renderiza el prompt de ingreso de nombre del proyecto
func RenderInput(ti textinput.Model) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Creating a new Koko project")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render("Project name")))
	b.WriteString(fmt.Sprintf("%s  %s\n", BarSymbol, ti.View()))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[Enter] Continuar • [Esc] Menú principal")))
	return b.String()
}
