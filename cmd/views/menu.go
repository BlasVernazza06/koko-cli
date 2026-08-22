package views

import (
	"fmt"
	"strings"
)

// RenderMenu renderiza la pantalla del menú principal
func RenderMenu(options []SelectOption, cursor int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Koko CLI")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render("¿Qué deseas hacer?")))
	b.WriteString(RenderOptions(options, cursor))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[↑/↓] Navegar • [Enter] Seleccionar • [Ctrl+C] Salir")))
	return b.String()
}
