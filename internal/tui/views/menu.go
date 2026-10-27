package views

import (
	"fmt"
	"strings"
)

// RenderMenu renders the main menu screen
func RenderMenu(options []SelectOption, cursor int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Koko CLI")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render("What would you like to do?")))
	b.WriteString(RenderOptions(options, cursor))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[↑/↓] Navigate • [Enter] Select • [Ctrl+C] Exit")))
	return b.String()
}
