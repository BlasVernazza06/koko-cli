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

	"github.com/BlasVernazza06/koko-cli/cmd/views"
	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

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

// manualStepConfig representa la definición de un paso en la configuración manual
type manualStepConfig struct {
	title   string
	label   string
	options []views.SelectOption
}

// mainModel representa el estado global y contexto de la TUI
type mainModel struct {
	state       sessionState
	prevState   sessionState
	versionInfo []views.VersionItem

	// Opciones del Menú Principal
	menuCursor  int
	menuOptions []views.SelectOption

	// Opciones de Modo
	modeCursor  int
	modeOptions []views.SelectOption

	// Opciones de Receta (Setup Rápido)
	recipeCursor  int
	recipeOptions []views.SelectOption

	// Configuración Manual Paso a Paso
	manualSteps      []manualStepConfig
	manualStepIdx    int
	manualCursors    []int
	manualSelections []views.SelectOption

	// Campos de inicialización
	projectNameInput textinput.Model
	chosenName       string
	chosenMode       string
	chosenRecipe     string
	cancelledMsg     string

	// Estado y progreso del Scaffolding
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
	ti.TextStyle = views.StyleValue.Bold(true)
	ti.PlaceholderStyle = views.StyleHint
	ti.Cursor.Style = views.StyleValue

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = views.StyleValue

	if initialProjectName != "" {
		ti.SetValue(initialProjectName)
	}
	if initialState == stateInitInput {
		ti.Focus()
	}

	manualSteps := []manualStepConfig{
		{
			title: "Selecciona el Frontend Framework",
			label: "Frontend",
			options: []views.SelectOption{
				{Value: "nextjs", Label: "Next.js (App Router)", Hint: "React framework con SSR y Server Components"},
				{Value: "react", Label: "React + Vite", Hint: "Single Page Application ultrarrápida"},
				{Value: "vue", Label: "Vue 3 + Vite", Hint: "Framework reactivo y progresivo"},
				{Value: "svelte", Label: "SvelteKit", Hint: "Compilador reactivo sin virtual DOM"},
				{Value: "none", Label: "Sin Frontend", Hint: "Solo Backend / API REST"},
			},
		},
		{
			title: "Selecciona el Backend Framework / Runtime",
			label: "Backend",
			options: []views.SelectOption{
				{Value: "express", Label: "Node.js / Express", Hint: "API REST ligera con TypeScript"},
				{Value: "fastapi", Label: "Python / FastAPI", Hint: "Asíncrono con validación Pydantic v2"},
				{Value: "go_chi", Label: "Go / Chi Router", Hint: "Alto rendimiento y tipado estricto"},
				{Value: "nestjs", Label: "NestJS", Hint: "Arquitectura empresarial modular con TypeScript"},
				{Value: "none", Label: "Sin Backend dedicado", Hint: "Usar Server Actions o BaaS"},
			},
		},
		{
			title: "Selecciona la Base de Datos",
			label: "Database",
			options: []views.SelectOption{
				{Value: "postgres", Label: "PostgreSQL", Hint: "Base de datos relacional estándar con Docker"},
				{Value: "mongodb", Label: "MongoDB", Hint: "Base de datos de documentos NoSQL"},
				{Value: "mysql", Label: "MySQL / MariaDB", Hint: "Base de datos SQL tradicional"},
				{Value: "sqlite", Label: "SQLite", Hint: "Base de datos embebida en archivo local"},
				{Value: "none", Label: "Ninguna", Hint: "Sin persistencia de datos"},
			},
		},
		{
			title: "Selecciona el ORM / Query Builder",
			label: "ORM / Tool",
			options: []views.SelectOption{
				{Value: "drizzle", Label: "Drizzle ORM", Hint: "Type-safe y liviano con soporte SQL nativo"},
				{Value: "prisma", Label: "Prisma", Hint: "ORM popular con auto-generación de tipos"},
				{Value: "sqlalchemy", Label: "SQLAlchemy / SQLModel", Hint: "ORM estándar para Python"},
				{Value: "gorm", Label: "GORM", Hint: "ORM completo para Go"},
				{Value: "none", Label: "Ninguno / Raw SQL", Hint: "Conexión directa sin ORM"},
			},
		},
	}

	return mainModel{
		state:     initialState,
		prevState: stateMenu,
		versionInfo: []views.VersionItem{
			{Key: "Koko CLI", Val: "v0.1.0"},
			{Key: "OS / Arch", Val: fmt.Sprintf("%s / %s", runtime.GOOS, runtime.GOARCH)},
			{Key: "Go Runtime", Val: runtime.Version()},
			{Key: "Build Date", Val: "2026-08-12"},
		},
		menuOptions: []views.SelectOption{
			{Value: "init", Label: "Inicializar proyecto", Hint: "Crear una nueva aplicación con Koko"},
			{Value: "version", Label: "Versión del sistema", Hint: "Ver detalles de entorno y CLI"},
			{Value: "exit", Label: "Salir", Hint: "Cerrar la herramienta"},
		},
		modeOptions: []views.SelectOption{
			{Value: "quick", Label: "Setup Rápido", Hint: "Recetas de producción listas para usar"},
			{Value: "manual", Label: "Configuración Manual", Hint: "Elegir stack paso a paso (Frontend, Backend, DB, ORM)"},
		},
		recipeOptions: []views.SelectOption{
			{Value: "saas", Label: "⚡ SaaS Starter", Hint: "Next.js + Drizzle + Better-Auth + Stripe"},
			{Value: "mern", Label: "💻 MERN Stack", Hint: "React + Express + MongoDB"},
			{Value: "pern", Label: "🚀 PERN Stack", Hint: "React + Express + PostgreSQL"},
			{Value: "fastapi_react", Label: "🐍 FastAPI + React", Hint: "FastAPI Backend + React SPA"},
		},
		manualSteps:      manualSteps,
		manualStepIdx:    0,
		manualCursors:    make([]int, len(manualSteps)),
		manualSelections: make([]views.SelectOption, len(manualSteps)),
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

		time.Sleep(450 * time.Millisecond) // micro-animación
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
					m.manualStepIdx = 0
				}
			}

		case stateInitManual:
			currentStep := m.manualSteps[m.manualStepIdx]
			cursor := m.manualCursors[m.manualStepIdx]

			switch msg.String() {
			case "up", "k":
				if cursor > 0 {
					m.manualCursors[m.manualStepIdx]--
				}
			case "down", "j":
				if cursor < len(currentStep.options)-1 {
					m.manualCursors[m.manualStepIdx]++
				}
			case "enter":
				// Guardar la selección actual
				selectedOpt := currentStep.options[m.manualCursors[m.manualStepIdx]]
				m.manualSelections[m.manualStepIdx] = selectedOpt

				// Avanzar al siguiente paso o finalizar configuración
				if m.manualStepIdx < len(m.manualSteps)-1 {
					m.manualStepIdx++
				} else {
					// Todos los pasos completados -> Iniciar creación
					m.chosenRecipe = "saas" // o configuración manual generada
					m.state = stateInitRunning
					m.currentStep = 1
					m.stepStatus[0] = "running"
					m.startTime = time.Now()
					return m, tea.Batch(m.spinner.Tick, runStepCmd(1, m.chosenName, m.chosenRecipe))
				}

			case "esc":
				// Retroceder paso sin perder la selección previa
				if m.manualStepIdx > 0 {
					m.manualStepIdx--
				} else {
					m.state = stateInitMode
				}
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

// View actúa como enrutador central y delega a las funciones de cmd/views
func (m mainModel) View() string {
	switch m.state {
	case stateCancelled:
		return views.RenderCancelled(m.cancelledMsg)
	case stateMenu:
		return views.RenderMenu(m.menuOptions, m.menuCursor)
	case stateVersion:
		return views.RenderVersion(m.versionInfo)
	case stateInitInput:
		return views.RenderInput(m.projectNameInput)
	case stateInitMode:
		return views.RenderMode(m.chosenName, m.modeOptions, m.modeCursor)
	case stateInitManual:
		// Construir historial dinámico de los pasos anteriores ya respondidos
		var history []views.SummaryItem
		for i := 0; i < m.manualStepIdx; i++ {
			history = append(history, views.SummaryItem{
				Label: m.manualSteps[i].label,
				Value: m.manualSelections[i].Label,
			})
		}
		currentStep := m.manualSteps[m.manualStepIdx]
		cursor := m.manualCursors[m.manualStepIdx]
		return views.RenderManual(m.chosenName, history, currentStep.title, currentStep.options, cursor)

	case stateInitRecipe:
		return views.RenderRecipe(m.chosenName, m.recipeOptions, m.recipeCursor)
	case stateInitRunning:
		pkgManager := views.GetPackageManager(m.chosenRecipe)
		recipeLabel := views.GetRecipeLabel(m.recipeOptions, m.chosenRecipe)
		return views.RenderRunning(m.chosenName, recipeLabel, pkgManager, m.stepNames, m.stepStatus, m.spinner)
	case stateInitDone:
		pkgManager := views.GetPackageManager(m.chosenRecipe)
		recipeLabel := views.GetRecipeLabel(m.recipeOptions, m.chosenRecipe)
		return views.RenderDone(m.chosenName, recipeLabel, pkgManager, m.stepNames, m.stepStatus, m.elapsedTime, m.scaffoldErr)
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
