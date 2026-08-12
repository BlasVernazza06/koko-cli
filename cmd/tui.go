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

func getRecipeName(recipe string) string {
	switch recipe {
	case "saas":
		return "SaaS Starter (Next.js + Drizzle + Better-Auth + Stripe)"
	case "mern":
		return "MERN Stack (React + Express + MongoDB)"
	case "pern":
		return "PERN Stack (React + Express + PostgreSQL)"
	case "fastapi_react":
		return "FastAPI + React SPA"
	default:
		return recipe
	}
}

func getPackageManager(recipe string) string {
	if recipe == "saas" || recipe == "pern" || recipe == "mern" {
		return "pnpm"
	}
	return "npm"
}

type sessionState int

const (
	stateMenu sessionState = iota
	stateVersion
	stateInitInput
	stateInitRecipe
	stateInitRunning
	stateInitDone
)

// Styling tokens
var (
	accentColor = lipgloss.Color("#8A2BE2") // Purple/violet
	cyanColor   = lipgloss.Color("#00FFFF")
	greenColor  = lipgloss.Color("#00FF7F")
	grayColor   = lipgloss.Color("#7D7D7D")
	redColor    = lipgloss.Color("#FF3333")

	titleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			MarginLeft(2).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			MarginLeft(2).
			MarginTop(1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(cyanColor).
				Bold(true).
				PaddingLeft(2)

	successStyle = lipgloss.NewStyle().
			Foreground(greenColor).
			Bold(true)

	grayMutedStyle = lipgloss.NewStyle().
			Foreground(grayColor)
)

type stepFinishedMsg struct {
	step int
	err  error
}

type mainModel struct {
	state       sessionState
	cursor      int
	choices     []string
	versionInfo string

	// Project init fields
	projectNameInput textinput.Model
	recipeCursor     int
	recipeChoices    []string
	recipeLabels     []string

	// Progress tracking
	currentStep  int
	stepNames    []string
	stepStatus   []string // "pending", "running", "success", "error"
	spinner      spinner.Model
	startTime    time.Time
	elapsedTime  time.Duration
	scaffoldErr  error
}

func initialModel() mainModel {
	ti := textinput.New()
	ti.Placeholder = "my-project"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(cyanColor)

	return mainModel{
		state:   stateMenu,
		choices: []string{"🚀  Inicializar un nuevo proyecto", "ℹ️   Versión e información del sistema", "❌  Salir"},
		versionInfo: fmt.Sprintf(
			"Koko-CLI v0.1.0\nOS/Arch:    %s/%s\nBuild Date: 2026-08-12",
			runtime.GOOS, runtime.GOARCH,
		),
		projectNameInput: ti,
		recipeChoices:    []string{"saas", "mern", "pern", "fastapi_react"},
		recipeLabels: []string{
			"⚡ SaaS Starter (Next.js + Drizzle + Better-Auth + Stripe)",
			"💻 MERN Stack (React + Express + MongoDB)",
			"🚀 PERN Stack (React + Express + PostgreSQL)",
			"🐍 FastAPI + React SPA",
		},
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

		time.Sleep(500 * time.Millisecond) // micro-animation/premium delay
		return stepFinishedMsg{step: step, err: err}
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.state {
		case stateMenu:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				switch m.cursor {
				case 0:
					m.state = stateInitInput
					m.projectNameInput.Focus()
					return m, nil
				case 1:
					m.state = stateVersion
				case 2:
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
					m.projectNameInput.SetValue("my-project")
				}
				m.state = stateInitRecipe
			default:
				m.projectNameInput, cmd = m.projectNameInput.Update(msg)
				return m, cmd
			}

		case stateInitRecipe:
			switch msg.String() {
			case "esc":
				m.state = stateInitInput
			case "up", "k":
				if m.recipeCursor > 0 {
					m.recipeCursor--
				}
			case "down", "j":
				if m.recipeCursor < len(m.recipeChoices)-1 {
					m.recipeCursor++
				}
			case "enter":
				m.state = stateInitRunning
				m.currentStep = 1
				m.stepStatus[0] = "running"
				m.startTime = time.Now()
				return m, tea.Batch(m.spinner.Tick, runStepCmd(1, m.projectNameInput.Value(), m.recipeChoices[m.recipeCursor]))
			}

		case stateInitDone:
			switch msg.String() {
			case "esc", "q", "enter":
				// Reset model to menu
				m.state = stateMenu
				m.projectNameInput.SetValue("")
				m.stepStatus = []string{"pending", "pending", "pending", "pending"}
				m.scaffoldErr = nil
			}
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
			return m, runStepCmd(m.currentStep, m.projectNameInput.Value(), m.recipeChoices[m.recipeCursor])
		}

		m.elapsedTime = time.Since(m.startTime)
		m.state = stateInitDone
		return m, tea.Quit
	}

	return m, nil
}

func (m mainModel) View() string {
	logo := "\033[36m" + `
			 __ _    __ _ __      __
			/  _ | |/ _' |\ \ /\ / /
			| (__| | (_| | \ V  V / 
			\___ |_|\__,_|  \_/\_/  v0.1.0` + "\033[0m\n"

	switch m.state {
	case stateMenu:
		s := logo + titleStyle.Render("Selecciona una opción:") + "\n\n"
		for i, choice := range m.choices {
			if m.cursor == i {
				s += selectedMenuItemStyle.Render("> " + choice) + "\n"
			} else {
				s += menuItemStyle.Render("  " + choice) + "\n"
			}
		}
		s += "\n" + grayMutedStyle.Render("  [↑/↓] Navegar • [Enter] Seleccionar") + "\n"
		return s

	case stateVersion:
		infoBox := boxStyle.Render(m.versionInfo)
		return logo + titleStyle.Render("Información de Sistema:") + "\n" + infoBox + "\n\n" + grayMutedStyle.Render("  [Esc/q/Enter] Volver al menú") + "\n"

	case stateInitInput:
		inputBox := boxStyle.Render(
			fmt.Sprintf("Nombre del proyecto:\n\n%s", m.projectNameInput.View()),
		)
		return logo + titleStyle.Render("Inicializar nuevo proyecto") + "\n" + inputBox + "\n\n" + grayMutedStyle.Render("  [Enter] Confirmar • [Esc] Cancelar") + "\n"

	case stateInitRecipe:
		s := logo + titleStyle.Render("Selecciona una receta de producción:") + "\n\n"
		for i, label := range m.recipeLabels {
			if m.recipeCursor == i {
				s += selectedMenuItemStyle.Render("> " + label) + "\n"
			} else {
				s += menuItemStyle.Render("  " + label) + "\n"
			}
		}
		s += "\n" + grayMutedStyle.Render("  [↑/↓] Navegar • [Enter] Empezar • [Esc] Atrás") + "\n"
		return s

	case stateInitRunning:
		pkgManager := getPackageManager(m.recipeChoices[m.recipeCursor])
		recipeName := getRecipeName(m.recipeChoices[m.recipeCursor])

		s := logo
		s += fmt.Sprintf("\n\033[90m┌\033[0m  \033[1mCreando un nuevo proyecto Koko\033[0m\n")
		s += fmt.Sprintf("\033[90m│\033[0m\n")
		s += fmt.Sprintf("\033[90m│\033[0m  \033[36m◇\033[0m  \033[1mNombre del proyecto\033[0m   \033[90m⋅\033[0m  \033[32m%s\033[0m\n", m.projectNameInput.Value())
		s += fmt.Sprintf("\033[90m│\033[0m  \033[36m◇\033[0m  \033[1mReceta de producción\033[0m  \033[90m⋅\033[0m  \033[35m%s\033[0m\n", recipeName)
		s += fmt.Sprintf("\033[90m│\033[0m  \033[36m◇\033[0m  \033[1mGestor de paquetes\033[0m    \033[90m⋅\033[0m  \033[33m%s\033[0m\n", pkgManager)
		s += fmt.Sprintf("\033[90m│\033[0m\n")

		for i, stepName := range m.stepNames {
			switch m.stepStatus[i] {
			case "pending":
				s += fmt.Sprintf("\033[90m│\033[0m  \033[90m◇  %s\033[0m\n", stepName)
			case "running":
				s += fmt.Sprintf("\033[90m│\033[0m  \033[36m%s\033[0m  %s\n", m.spinner.View(), stepName)
			case "success":
				s += fmt.Sprintf("\033[90m│\033[0m  \033[32m✓\033[0m  %s\n", stepName)
			case "error":
				s += fmt.Sprintf("\033[90m│\033[0m  \033[31m✗\033[0m  %s\n", stepName)
			}
		}
		return s

	case stateInitDone:
		pkgManager := getPackageManager(m.recipeChoices[m.recipeCursor])
		recipeName := getRecipeName(m.recipeChoices[m.recipeCursor])

		s := logo
		s += fmt.Sprintf("\n\033[90m┌\033[0m  \033[1mCreando un nuevo proyecto Koko\033[0m\n")
		s += fmt.Sprintf("\033[90m│\033[0m\n")
		s += fmt.Sprintf("\033[90m│\033[0m  \033[36m◇\033[0m  \033[1mNombre del proyecto\033[0m   \033[90m⋅\033[0m  \033[32m%s\033[0m\n", m.projectNameInput.Value())
		s += fmt.Sprintf("\033[90m│\033[0m  \033[36m◇\033[0m  \033[1mReceta de producción\033[0m  \033[90m⋅\033[0m  \033[35m%s\033[0m\n", recipeName)
		s += fmt.Sprintf("\033[90m│\033[0m  \033[36m◇\033[0m  \033[1mGestor de paquetes\033[0m    \033[90m⋅\033[0m  \033[33m%s\033[0m\n", pkgManager)
		s += fmt.Sprintf("\033[90m│\033[0m\n")

		for i, stepName := range m.stepNames {
			if m.stepStatus[i] == "success" {
				s += fmt.Sprintf("\033[90m│\033[0m  \033[32m✓\033[0m  %s\n", stepName)
			} else if m.stepStatus[i] == "error" {
				s += fmt.Sprintf("\033[90m│\033[0m  \033[31m✗\033[0m  %s\n", stepName)
			} else {
				s += fmt.Sprintf("\033[90m│\033[0m  \033[90m◇  %s\033[0m\n", stepName)
			}
		}

		s += fmt.Sprintf("\033[90m│\033[0m\n")

		if m.scaffoldErr != nil {
			s += fmt.Sprintf("\033[90m└\033[0m  \033[31m\033[1mError en la creación: %v\033[0m\n\n", m.scaffoldErr)
			s += grayMutedStyle.Render("  [Enter/Esc] Volver al menú") + "\n"
			return s
		}

		s += fmt.Sprintf("\033[90m└\033[0m  \033[32m\033[1m¡Proyecto creado con éxito en %.2f segundos!\033[0m\n\n", m.elapsedTime.Seconds())

		s += fmt.Sprintf("\033[1mSiguientes pasos:\033[0m\n")
		s += fmt.Sprintf("  1. cd %s\n", m.projectNameInput.Value())
		if pkgManager == "pnpm" {
			s += fmt.Sprintf("  2. pnpm install\n")
			s += fmt.Sprintf("  3. pnpm dev\n\n")
		} else {
			s += fmt.Sprintf("  2. npm install\n")
			s += fmt.Sprintf("  3. npm run dev\n\n")
		}

		s += grayMutedStyle.Render("  [Enter/Esc] Volver al menú") + "\n"
		return s
	}

	return ""
}

func RunTUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al arrancar la TUI: %v", err)
		os.Exit(1)
	}
}
