package views

import (
	"fmt"
	"strings"
)

type VersionItem struct {
	Key string
	Val string
}

// RenderVersion renderiza la pantalla de información del sistema
func RenderVersion(items []VersionItem) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Koko CLI · Información del Sistema")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	for _, item := range items {
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render(item.Key), DotDivider, StyleValue.Render(item.Val)))
	}
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[Esc / Enter] Volver al menú")))
	return b.String()
}
