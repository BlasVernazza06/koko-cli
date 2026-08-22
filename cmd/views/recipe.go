package views

import (
	"fmt"
	"strings"
)

// GetRecipeLabel obtiene el label amigable de una receta dada su clave
func GetRecipeLabel(recipeOptions []SelectOption, recipeVal string) string {
	for _, r := range recipeOptions {
		if r.Value == recipeVal {
			return r.Label
		}
	}
	return recipeVal
}

// GetPackageManager retorna el gestor de paquetes correspondiente a una receta
func GetPackageManager(recipeVal string) string {
	if recipeVal == "saas" || recipeVal == "pern" || recipeVal == "mern" {
		return "pnpm"
	}
	return "npm"
}

// RenderRecipe renderiza la selección de recetas de proyectos
func RenderRecipe(projectName string, options []SelectOption, cursor int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n%s  %s\n", TopSymbol, StyleHeader.Render("Creating a new Koko project")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render("Project name"), DotDivider, StyleValue.Render(projectName)))
	b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", CompletedDiamond, StylePromptTitle.Render("Setup mode"), DotDivider, StyleValue.Render("Setup Rápido")))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", ActiveDiamond, StylePromptTitle.Render("Choose recipe / stack")))
	b.WriteString(RenderOptions(options, cursor))
	b.WriteString(fmt.Sprintf("%s\n", BarSymbol))
	b.WriteString(fmt.Sprintf("%s  %s\n", BottomSymbol, StyleMuted.Render("[↑/↓] Navegar • [Enter] Crear proyecto • [Esc] Atrás")))
	return b.String()
}
