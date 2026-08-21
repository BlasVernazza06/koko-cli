package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

// Clack-style color & character tokens
var (
	colorIndigo = lipgloss.Color("#5a4fc4")
	colorMuted  = lipgloss.Color("#7D7D7D")

	// Glyphs
	barSymbol        = lipgloss.NewStyle().Foreground(colorMuted).Render("│")
	topSymbol        = lipgloss.NewStyle().Foreground(colorMuted).Render("┌")
	bottomSymbol     = lipgloss.NewStyle().Foreground(colorMuted).Render("└")
	activeDiamond    = lipgloss.NewStyle().Foreground(colorIndigo).Render("◇")    // Vacío sin relleno en #5a4fc4 cuando NO está completado
	completedDiamond = lipgloss.NewStyle().Foreground(colorIndigo).Render("◆")   // Rellenado en #5a4fc4 cuando está completado
	dotDivider       = lipgloss.NewStyle().Foreground(colorMuted).Render("·")
	radioActive      = lipgloss.NewStyle().Foreground(colorIndigo).Render("●")   // Selección activa en #5a4fc4
	radioInactive    = lipgloss.NewStyle().Foreground(colorMuted).Render("○")
	checkSuccess     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF7F")).Render("✓")
	crossError       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render("✗")
	diamondPending   = lipgloss.NewStyle().Foreground(colorMuted).Render("◇")

	// Lipgloss styles
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	stylePromptTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF"))

	styleValue = lipgloss.NewStyle().
			Foreground(colorIndigo)

	styleActiveItem = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	styleInactiveItem = lipgloss.NewStyle().
				Foreground(colorMuted)

	styleHint = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	styleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleSuccess = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF7F"))

	styleError = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF4444"))
)

type SelectOption struct {
	Value string
	Label string
	Hint  string
}

type sessionState int

const (
	stateMenu sessionState = iota
	stateVersion
	stateInitInput
	stateInitMode
	stateInitRecipe
	stateInitManual
	stateInitRunning
	stateInitDone
	stateCancelled
)

type stepFinishedMsg struct {
	step int
	err  error
}

type mainModel struct {
	state       sessionState
	prevState   sessionState
	versionInfo []struct{ key, val string }

	// Menu options
	menuCursor  int
	menuOptions []SelectOption

	// Mode options
	modeCursor  int
	modeOptions []SelectOption

	// Recipe options
	recipeCursor  int
	recipeOptions []SelectOption

	// Project init fields
	projectNameInput textinput.Model
	chosenName       string
	chosenMode       string
	chosenRecipe     string
	cancelledMsg     string

	// Progress tracking
	currentStep int
	stepNames   []string
	stepStatus  []string // "pending", "running", "success", "error"
	spinner     spinner.Model
	startTime   time.Time
	elapsedTime time.Duration
	scaffoldErr error
}

func initialModel(initialState sessionState, initialProjectName string) mainModel {
	ti := textinput.New()
	ti.Placeholder = "my-app"
	ti.CharLimit = 64
	ti.Width = 30
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(colorIndigo).Bold(true)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(colorIndigo)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorIndigo)

	if initialProjectName != "" {
		ti.SetValue(initialProjectName)
	}
	if initialState == stateInitInput {
		ti.Focus()
	}

	return mainModel{
		state:     initialState,
		prevState: stateMenu,
		versionInfo: []struct{ key, val string }{
			{"Koko CLI", "v0.1.0"},
			{"OS / Arch", fmt.Sprintf("%s / %s", runtime.GOOS, runtime.GOARCH)},
			{"Go Runtime", runtime.Version()},
			{"Build Date", "2026-08-12"},
		},
		menuOptions: []SelectOption{
			{Value: "init", Label: "Inicializar proyecto", Hint: "Crear una nueva aplicación con Koko"},
			{Value: "version", Label: "Versión del sistema", Hint: "Ver detalles de entorno y CLI"},
			{Value: "exit", Label: "Salir", Hint: "Cerrar la herramienta"},
		},
		modeOptions: []SelectOption{
			{Value: "quick", Label: "Setup Rápido", Hint: "Recetas de producción listas para usar"},
			{Value: "manual", Label: "Configuración Manual", Hint: "Elegir stack paso a paso (Próximamente)"},
		},
		recipeOptions: []SelectOption{
			{Value: "saas", Label: "⚡ SaaS Starter", Hint: "Next.js + Drizzle + Better-Auth + Stripe"},
			{Value: "mern", Label: "💻 MERN Stack", Hint: "React + Express + MongoDB"},
			{Value: "pern", Label: "🚀 PERN Stack", Hint: "React + Express + PostgreSQL"},
			{Value: "fastapi_react", Label: "🐍 FastAPI + React", Hint: "FastAPI Backend + React SPA"},
		},
		projectNameInput: ti,
		stepNames: []string{
			"Estructurando directorios y copiando plantillas...",
			"Generando configuración de Docker y DB...",
			"Inicializando repositorio Git...",
			"Creando manifiesto koko.config.json...",
		},
		stepStatus: []string{"pending", "pending", "pending", "pending"},
		spinner:    s,
	}
}

func (m mainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m mainModel) getRecipeLabel(recipeVal string) string {
	for _, r := range m.recipeOptions {
		if r.Value == recipeVal {
			return r.Label
		}
	}
	return recipeVal
}

func (m mainModel) getPackageManager(recipeVal string) string {
	if recipeVal == "saas" || recipeVal == "pern" || recipeVal == "mern" {
		return "pnpm"
	}
	return "npm"
}

func runStepCmd(step int, projectName string, recipe string) tea.Cmd {
	return func() tea.Msg {
		cfg := scaffold.ScaffoldConfig{
			ProjectName: projectName,
			Recipe:      recipe,
		}

		var err error
		switch step {
		case 1:
			err = scaffold.CopyTemplates(projectName, cfg)
		case 2:
			err = scaffold.GenerateDockerAndDB(projectName, cfg)
		case 3:
			err = scaffold.InitGit(projectName)
		case 4:
			err = kokoConfig.GenerateConfig(projectName, cfg)
		}

		time.Sleep(450 * time.Millisecond) // micro-animation delay
		return stepFinishedMsg{step: step, err: err}
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancelledMsg = "Operación cancelada."
			m.state = stateCancelled
			return m, tea.Quit
		}

		switch m.state {
		case stateMenu:
			switch msg.String() {
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case "down", "j":
				if m.menuCursor < len(m.menuOptions)-1 {
					m.menuCursor++
				}
			case "enter":
				selected := m.menuOptions[m.menuCursor].Value
				switch selected {
				case "init":
					m.state = stateInitInput
					m.projectNameInput.Focus()
					return m, textinput.Blink
				case "version":
					m.state = stateVersion
				case "exit":
					m.cancelledMsg = "Hasta pronto."
					m.state = stateCancelled
					return m, tea.Quit
				}
			}

		case stateVersion:
			switch msg.String() {
			case "esc", "q", "enter":
				m.state = stateMenu
			}

		case stateInitInput:
			switch msg.String() {
			case "esc":
				m.state = stateMenu
			case "enter":
				name := strings.TrimSpace(m.projectNameInput.Value())
				if name == "" {
					name = "my-app"
					m.projectNameInput.SetValue("my-app")
				}
				m.chosenName = name
				m.state = stateInitMode
			default:
				m.projectNameInput, cmd = m.projectNameInput.Update(msg)
				return m, cmd
			}

		case stateInitMode:
			switch msg.String() {
			case "esc":
				m.state = stateInitInput
				m.projectNameInput.Focus()
			case "up", "k":
				if m.modeCursor > 0 {
					m.modeCursor--
				}
			case "down", "j":
				if m.modeCursor < len(m.modeOptions)-1 {
					m.modeCursor++
				}
			case "enter":
				m.chosenMode = m.modeOptions[m.modeCursor].Label
				if m.modeOptions[m.modeCursor].Value == "quick" {
					m.state = stateInitRecipe
				} else {
					m.state = stateInitManual
				}
			}

		case stateInitManual:
			switch msg.String() {
			case "esc", "enter", "q":
				m.state = stateInitMode
			}

		case stateInitRecipe:
			switch msg.String() {
			case "esc":
				m.state = stateInitMode
			case "up", "k":
				if m.recipeCursor > 0 {
					m.recipeCursor--
				}
			case "down", "j":
				if m.recipeCursor < len(m.recipeOptions)-1 {
					m.recipeCursor++
				}
			case "enter":
				m.chosenRecipe = m.recipeOptions[m.recipeCursor].Value
				m.state = stateInitRunning
				m.currentStep = 1
				m.stepStatus[0] = "running"
				m.startTime = time.Now()
				return m, tea.Batch(m.spinner.Tick, runStepCmd(1, m.chosenName, m.chosenRecipe))
			}

		case stateInitDone:
			switch msg.String() {
			case "esc", "q", "enter":
				return m, tea.Quit
			}

		case stateCancelled:
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		return m, spinnerCmd

	case stepFinishedMsg:
		if msg.err != nil {
			m.stepStatus[msg.step-1] = "error"
			m.scaffoldErr = msg.err
			m.state = stateInitDone
			return m, tea.Quit
		}

		m.stepStatus[msg.step-1] = "success"
		if msg.step < 4 {
			m.currentStep = msg.step + 1
			m.stepStatus[msg.step] = "running"
			return m, runStepCmd(m.currentStep, m.chosenName, m.chosenRecipe)
		}

		m.elapsedTime = time.Since(m.startTime)
		m.state = stateInitDone
		return m, tea.Quit
	}

	return m, nil
}

func renderOptions(options []SelectOption, cursor int) string {
	var b strings.Builder
	for i, opt := range options {
		if i == cursor {
			b.WriteString(fmt.Sprintf("%s  %s %s", barSymbol, radioActive, styleActiveItem.Render(opt.Label)))
			if opt.Hint != "" {
				b.WriteString(fmt.Sprintf("  %s", styleHint.Render(opt.Hint)))
			}
		} else {
			b.WriteString(fmt.Sprintf("%s  %s %s", barSymbol, radioInactive, styleInactiveItem.Render(opt.Label)))
			if opt.Hint != "" {
				b.WriteString(fmt.Sprintf("  %s", styleHint.Render(opt.Hint)))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m mainModel) View() string {
	switch m.state {
	case stateCancelled:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Koko CLI")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n\n", bottomSymbol, styleMuted.Render(m.cancelledMsg)))
		return b.String()

	case stateMenu:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Koko CLI")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", activeDiamond, stylePromptTitle.Render("¿Qué deseas hacer?")))
		b.WriteString(renderOptions(m.menuOptions, m.menuCursor))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("[↑/↓] Navegar • [Enter] Seleccionar • [Ctrl+C] Salir")))
		return b.String()

	case stateVersion:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Koko CLI · Información del Sistema")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		for _, item := range m.versionInfo {
			b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render(item.key), dotDivider, styleValue.Render(item.val)))
		}
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("[Esc / Enter] Volver al menú")))
		return b.String()

	case stateInitInput:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Creating a new Koko project")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", activeDiamond, stylePromptTitle.Render("Project name")))
		b.WriteString(fmt.Sprintf("%s  %s\n", barSymbol, m.projectNameInput.View()))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("[Enter] Continuar • [Esc] Menú principal")))
		return b.String()

	case stateInitMode:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Creating a new Koko project")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Project name"), dotDivider, styleValue.Render(m.chosenName)))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", activeDiamond, stylePromptTitle.Render("Choose setup mode")))
		b.WriteString(renderOptions(m.modeOptions, m.modeCursor))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("[↑/↓] Navegar • [Enter] Continuar • [Esc] Cambiar nombre")))
		return b.String()

	case stateInitManual:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Creating a new Koko project")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Project name"), dotDivider, styleValue.Render(m.chosenName)))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", activeDiamond, stylePromptTitle.Render("Configuración Manual")))
		b.WriteString(fmt.Sprintf("%s  %s\n", barSymbol, styleHint.Render("El modo interactivo modular estará disponible en la próxima versión.")))
		b.WriteString(fmt.Sprintf("%s  %s\n", barSymbol, styleHint.Render("Por favor selecciona 'Setup Rápido' con nuestras recetas listas para producción.")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("[Esc / Enter] Volver a selección de modo")))
		return b.String()

	case stateInitRecipe:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Creating a new Koko project")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Project name"), dotDivider, styleValue.Render(m.chosenName)))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Setup mode"), dotDivider, styleValue.Render("Setup Rápido")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", activeDiamond, stylePromptTitle.Render("Choose recipe / stack")))
		b.WriteString(renderOptions(m.recipeOptions, m.recipeCursor))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("[↑/↓] Navegar • [Enter] Crear proyecto • [Esc] Atrás")))
		return b.String()

	case stateInitRunning:
		pkgManager := m.getPackageManager(m.chosenRecipe)
		recipeLabel := m.getRecipeLabel(m.chosenRecipe)

		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Creating a new Koko project")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Project name"), dotDivider, styleValue.Render(m.chosenName)))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Recipe"), dotDivider, styleValue.Render(recipeLabel)))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Package Manager"), dotDivider, styleValue.Render(pkgManager)))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))

		for i, stepName := range m.stepNames {
			switch m.stepStatus[i] {
			case "pending":
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, diamondPending, styleMuted.Render(stepName)))
			case "running":
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, m.spinner.View(), styleActiveItem.Render(stepName)))
			case "success":
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, checkSuccess, stylePromptTitle.Render(stepName)))
			case "error":
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, crossError, styleError.Render(stepName)))
			}
		}
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s\n", bottomSymbol, styleMuted.Render("Scaffolding in progress...")))
		return b.String()

	case stateInitDone:
		pkgManager := m.getPackageManager(m.chosenRecipe)
		recipeLabel := m.getRecipeLabel(m.chosenRecipe)

		var b strings.Builder
		b.WriteString(fmt.Sprintf("\n%s  %s\n", topSymbol, styleHeader.Render("Creating a new Koko project")))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Project name"), dotDivider, styleValue.Render(m.chosenName)))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Recipe"), dotDivider, styleValue.Render(recipeLabel)))
		b.WriteString(fmt.Sprintf("%s  %s  %s  %s\n", completedDiamond, stylePromptTitle.Render("Package Manager"), dotDivider, styleValue.Render(pkgManager)))
		b.WriteString(fmt.Sprintf("%s\n", barSymbol))

		for i, stepName := range m.stepNames {
			if m.stepStatus[i] == "success" {
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, checkSuccess, stylePromptTitle.Render(stepName)))
			} else if m.stepStatus[i] == "error" {
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, crossError, styleError.Render(stepName)))
			} else {
				b.WriteString(fmt.Sprintf("%s  %s  %s\n", barSymbol, diamondPending, styleMuted.Render(stepName)))
			}
		}

		b.WriteString(fmt.Sprintf("%s\n", barSymbol))

		if m.scaffoldErr != nil {
			b.WriteString(fmt.Sprintf("%s  %s\n\n", bottomSymbol, styleError.Render(fmt.Sprintf("Error en la creación: %v", m.scaffoldErr))))
			return b.String()
		}

		b.WriteString(fmt.Sprintf("%s  %s\n\n", bottomSymbol, styleSuccess.Render(fmt.Sprintf("¡Proyecto creado con éxito en %.2fs!", m.elapsedTime.Seconds()))))

		b.WriteString(fmt.Sprintf("  %s\n", styleHeader.Render("Próximos pasos:")))
		b.WriteString(fmt.Sprintf("  1. %s\n", styleValue.Render(fmt.Sprintf("cd %s", m.chosenName))))
		if pkgManager == "pnpm" {
			b.WriteString(fmt.Sprintf("  2. %s\n", styleValue.Render("pnpm install")))
			b.WriteString(fmt.Sprintf("  3. %s\n\n", styleValue.Render("pnpm dev")))
		} else {
			b.WriteString(fmt.Sprintf("  2. %s\n", styleValue.Render("npm install")))
			b.WriteString(fmt.Sprintf("  3. %s\n\n", styleValue.Render("npm run dev")))
		}

		return b.String()
	}

	return ""
}

func RunTUI() {
	p := tea.NewProgram(initialModel(stateMenu, ""))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al arrancar la TUI: %v\n", err)
		os.Exit(1)
	}
}

func RunTUIInit(projectName string) {
	state := stateInitInput
	if projectName != "" {
		state = stateInitMode
	}
	model := initialModel(state, projectName)
	if projectName != "" {
		model.chosenName = projectName
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al arrancar la TUI: %v\n", err)
		os.Exit(1)
	}
}
