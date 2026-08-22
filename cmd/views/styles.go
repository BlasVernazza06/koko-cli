package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Tokens de color y caracteres
var (
	ColorViolet      = lipgloss.Color("#8A2BE2") // Violeta para el paso activo
	ColorGreenCheck  = lipgloss.Color("#00FF7F") // Verde tipo check para el paso completado
	ColorLightViolet = lipgloss.Color("#A78BFA") // Violeta más claro para el value / texto ingresado
	ColorMuted       = lipgloss.Color("#7D7D7D")

	// Glifos
	BarSymbol        = lipgloss.NewStyle().Foreground(ColorMuted).Render("│")
	TopSymbol        = lipgloss.NewStyle().Foreground(ColorMuted).Render("┌")
	BottomSymbol     = lipgloss.NewStyle().Foreground(ColorMuted).Render("└")
	ActiveDiamond    = lipgloss.NewStyle().Foreground(ColorViolet).Render("◇")     // Rombo sin relleno y violeta cuando se está completando
	CompletedDiamond = lipgloss.NewStyle().Foreground(ColorGreenCheck).Render("◆") // Rombo relleno color verde tipo check cuando está completado
	DotDivider       = lipgloss.NewStyle().Foreground(ColorMuted).Render("·")
	RadioActive      = lipgloss.NewStyle().Foreground(ColorViolet).Render("●")
	RadioInactive    = lipgloss.NewStyle().Foreground(ColorMuted).Render("○")
	CheckSuccess     = lipgloss.NewStyle().Foreground(ColorGreenCheck).Render("✓")
	CrossError       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render("✗")
	DiamondPending   = lipgloss.NewStyle().Foreground(ColorMuted).Render("◇")

	// Estilos Lipgloss
	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	StylePromptTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF"))

	StyleValue = lipgloss.NewStyle().
			Foreground(ColorLightViolet)

	StyleActiveItem = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	StyleInactiveItem = lipgloss.NewStyle().
				Foreground(ColorMuted)

	StyleHint = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleSuccess = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGreenCheck)

	StyleError = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF4444"))
)

// RenderOptions dibuja una lista de opciones con el cursor y estilo de selección
func RenderOptions(options []SelectOption, cursor int) string {
	var b strings.Builder
	for i, opt := range options {
		if i == cursor {
			b.WriteString(fmt.Sprintf("%s  %s %s", BarSymbol, RadioActive, StyleActiveItem.Render(opt.Label)))
			if opt.Hint != "" {
				b.WriteString(fmt.Sprintf("  %s", StyleHint.Render(opt.Hint)))
			}
		} else {
			b.WriteString(fmt.Sprintf("%s  %s %s", BarSymbol, RadioInactive, StyleInactiveItem.Render(opt.Label)))
			if opt.Hint != "" {
				b.WriteString(fmt.Sprintf("  %s", StyleHint.Render(opt.Hint)))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
