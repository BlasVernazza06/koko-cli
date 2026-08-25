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
	"github.com/BlasVernazza06/koko-cli/internal/validator"
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

// manualStepConfig represents the definition of a step in manual configuration
type manualStepConfig struct {
	title   string
	label   string
	options []views.SelectOption
}

type runnerStep struct {
	name string
	run  func(string, scaffold.ScaffoldConfig) error
}

// mainModel represents the global state and context of the TUI
type mainModel struct {
	state       sessionState
	prevState   sessionState
	versionInfo []views.VersionItem

	// Main Menu Options
	menuCursor  int
	menuOptions []views.SelectOption

	// Setup Mode Options
	modeCursor  int
	modeOptions []views.SelectOption

	// Recipe Options (Quick Setup)
	recipeCursor  int
	recipeOptions []views.SelectOption

	// Step-by-Step Manual Configuration
	manualSteps      []manualStepConfig
	manualStepIdx    int
	manualCursors    []int
	manualSelections []views.SelectOption

	// Initialization fields
	projectNameInput textinput.Model
	inputErr         string
	chosenName       string
	chosenMode       string
	chosenRecipe     string
	cancelledMsg     string
	scaffoldConfig   scaffold.ScaffoldConfig

	// Scaffolding Progress State
	currentStep int
	runnerSteps []runnerStep
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
			title: "Select Frontend Framework",
			label: "Frontend",
			options: []views.SelectOption{
				{Value: "nextjs", Label: "Next.js (App Router)", Hint: "React framework with SSR & Server Components"},
				{Value: "react", Label: "React + Vite", Hint: "Ultra-fast Single Page Application"},
				{Value: "vue", Label: "Vue 3 + Vite", Hint: "Progressive and reactive framework"},
				{Value: "svelte", Label: "SvelteKit", Hint: "Reactive compiler with no virtual DOM"},
				{Value: "none", Label: "None", Hint: "Backend / REST API only"},
			},
		},
		{
			title: "Select Backend Framework / Runtime",
			label: "Backend",
			options: []views.SelectOption{
				{Value: "express", Label: "Node.js / Express", Hint: "Lightweight REST API with TypeScript"},
				{Value: "fastapi", Label: "Python / FastAPI", Hint: "Async framework with Pydantic v2 validation"},
				{Value: "go_chi", Label: "Go / Chi Router", Hint: "High performance with strict types"},
				{Value: "nestjs", Label: "NestJS", Hint: "Enterprise modular architecture with TypeScript"},
				{Value: "none", Label: "None", Hint: "No dedicated backend (Server Actions or BaaS)"},
			},
		},
		{
			title: "Select Package Manager",
			label: "Package Manager",
			options: []views.SelectOption{
				{Value: "pnpm", Label: "PNPM", Hint: "Fast and disk space efficient (Recommended)"},
				{Value: "npm", Label: "NPM", Hint: "Standard Node package manager"},
				{Value: "bun", Label: "Bun", Hint: "All-in-one JavaScript runtime & package manager"},
			},
		},
		{
			title: "Select Database",
			label: "Database",
			options: []views.SelectOption{
				{Value: "postgres", Label: "PostgreSQL", Hint: "Standard relational database with Docker"},
				{Value: "mongodb", Label: "MongoDB", Hint: "NoSQL document database"},
				{Value: "mysql", Label: "MySQL / MariaDB", Hint: "Traditional SQL database"},
				{Value: "sqlite", Label: "SQLite", Hint: "Embedded lightweight database"},
				{Value: "none", Label: "None", Hint: "No database persistence"},
			},
		},
		{
			title: "Select ORM / Query Builder",
			label: "ORM / Tool",
			options: []views.SelectOption{
				{Value: "drizzle", Label: "Drizzle ORM", Hint: "Lightweight, type-safe with native SQL support"},
				{Value: "prisma", Label: "Prisma", Hint: "Next-gen ORM with auto type generation"},
				{Value: "sqlalchemy", Label: "SQLAlchemy / SQLModel", Hint: "Standard ORM for Python"},
				{Value: "gorm", Label: "GORM", Hint: "Feature-rich ORM for Go"},
				{Value: "none", Label: "None / Raw SQL", Hint: "Direct driver connection without ORM"},
			},
		},
		{
			title: "Select Authentication Provider",
			label: "Auth",
			options: []views.SelectOption{
				{Value: "better-auth", Label: "Better-Auth", Hint: "Comprehensive TypeScript-first auth"},
				{Value: "clerk", Label: "Clerk", Hint: "Complete user management & auth suite"},
				{Value: "none", Label: "None", Hint: "Skip authentication setup"},
			},
		},
		{
			title: "Select Addons / Tooling",
			label: "Addons",
			options: []views.SelectOption{
				{Value: "docker_cicd", Label: "Docker Compose + GitHub Actions", Hint: "Full containerization & CI/CD workflow"},
				{Value: "docker", Label: "Docker Compose", Hint: "Local containerized services"},
				{Value: "github_actions", Label: "GitHub Actions CI", Hint: "Automated linting and test workflows"},
				{Value: "none", Label: "None", Hint: "No extra tooling"},
			},
		},
		{
			title: "Initialize Git Repository?",
			label: "Git",
			options: []views.SelectOption{
				{Value: "yes", Label: "Yes", Hint: "Initialize a new Git repository (git init)"},
				{Value: "no", Label: "No", Hint: "Skip Git repository initialization"},
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
			{Value: "init", Label: "Initialize project", Hint: "Create a new application with Koko"},
			{Value: "version", Label: "System version", Hint: "View environment and CLI details"},
			{Value: "exit", Label: "Exit", Hint: "Quit the application"},
		},
		modeOptions: []views.SelectOption{
			{Value: "quick", Label: "Quick Setup", Hint: "Production-ready recipes ready to use"},
			{Value: "manual", Label: "Manual Configuration", Hint: "Choose stack step-by-step (Frontend, Backend, DB, ORM, Auth, etc.)"},
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
		spinner:          s,
	}
}

func (m mainModel) Init() tea.Cmd {
	return textinput.Blink
}

func setupRunnerSteps(cfg scaffold.ScaffoldConfig) []runnerStep {
	steps := []runnerStep{
		{
			name: "Scaffolding directories and copying templates...",
			run: func(p string, c scaffold.ScaffoldConfig) error {
				return scaffold.CopyTemplates(p, c)
			},
		},
		{
			name: "Generating Docker and Database configuration...",
			run: func(p string, c scaffold.ScaffoldConfig) error {
				return scaffold.GenerateDockerAndDB(p, c)
			},
		},
	}

	if cfg.InitGit {
		steps = append(steps, runnerStep{
			name: "Initializing Git repository...",
			run: func(p string, c scaffold.ScaffoldConfig) error {
				return scaffold.InitGit(p)
			},
		})
	}

	steps = append(steps, runnerStep{
		name: "Creating koko.config.json manifest...",
		run: func(p string, c scaffold.ScaffoldConfig) error {
			return kokoConfig.GenerateConfig(p, c)
		},
	})

	return steps
}

func runStepCmd(stepIdx int, projectName string, cfg scaffold.ScaffoldConfig, fn func(string, scaffold.ScaffoldConfig) error) tea.Cmd {
	return func() tea.Msg {
		var err error
		if fn != nil {
			err = fn(projectName, cfg)
		}
		time.Sleep(450 * time.Millisecond) // micro-animation
		return stepFinishedMsg{step: stepIdx, err: err}
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancelledMsg = "Operation cancelled."
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
					m.cancelledMsg = "Goodbye!"
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
				m.inputErr = ""
				m.state = stateMenu
			case "enter":
				name := strings.TrimSpace(m.projectNameInput.Value())
				if name == "" {
					name = "my-app"
					m.projectNameInput.SetValue("my-app")
				}
				if err := validator.Validate(name); err != nil {
					m.inputErr = err.Error()
					return m, nil
				}
				m.inputErr = ""
				m.chosenName = name
				m.state = stateInitMode
			default:
				m.inputErr = ""
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
				m.chosenMode = m.modeOptions[m.modeCursor].Value
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
				// Save current selection
				selectedOpt := currentStep.options[m.manualCursors[m.manualStepIdx]]
				m.manualSelections[m.manualStepIdx] = selectedOpt

				// Advance to next step or complete configuration
				if m.manualStepIdx < len(m.manualSteps)-1 {
					m.manualStepIdx++
				} else {
					// All steps completed -> initialize project
					initGit := m.manualSelections[7].Value == "yes"
					m.scaffoldConfig = scaffold.ScaffoldConfig{
						ProjectName:    m.chosenName,
						Recipe:         "",
						InitGit:        initGit,
						Frontend:       m.manualSelections[0].Value,
						Backend:        m.manualSelections[1].Value,
						PackageManager: m.manualSelections[2].Value,
						Database:       m.manualSelections[3].Value,
						ORM:            m.manualSelections[4].Value,
						Auth:           m.manualSelections[5].Value,
						Addons:         m.manualSelections[6].Value,
					}

					m.chosenRecipe = ""
					m.runnerSteps = setupRunnerSteps(m.scaffoldConfig)
					m.stepNames = make([]string, len(m.runnerSteps))
					m.stepStatus = make([]string, len(m.runnerSteps))
					for i, s := range m.runnerSteps {
						m.stepNames[i] = s.name
						m.stepStatus[i] = "pending"
					}

					m.state = stateInitRunning
					m.currentStep = 0
					m.stepStatus[0] = "running"
					m.startTime = time.Now()
					return m, tea.Batch(m.spinner.Tick, runStepCmd(0, m.chosenName, m.scaffoldConfig, m.runnerSteps[0].run))
				}

			case "esc":
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
				m.scaffoldConfig = scaffold.ScaffoldConfig{
					ProjectName: m.chosenName,
					Recipe:      m.chosenRecipe,
					InitGit:     true,
				}

				m.runnerSteps = setupRunnerSteps(m.scaffoldConfig)
				m.stepNames = make([]string, len(m.runnerSteps))
				m.stepStatus = make([]string, len(m.runnerSteps))
				for i, s := range m.runnerSteps {
					m.stepNames[i] = s.name
					m.stepStatus[i] = "pending"
				}

				m.state = stateInitRunning
				m.currentStep = 0
				m.stepStatus[0] = "running"
				m.startTime = time.Now()
				return m, tea.Batch(m.spinner.Tick, runStepCmd(0, m.chosenName, m.scaffoldConfig, m.runnerSteps[0].run))
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
			m.stepStatus[msg.step] = "error"
			m.scaffoldErr = msg.err
			m.state = stateInitDone
			return m, tea.Quit
		}

		m.stepStatus[msg.step] = "success"
		if msg.step+1 < len(m.runnerSteps) {
			m.currentStep = msg.step + 1
			m.stepStatus[m.currentStep] = "running"
			return m, runStepCmd(m.currentStep, m.chosenName, m.scaffoldConfig, m.runnerSteps[m.currentStep].run)
		}

		m.elapsedTime = time.Since(m.startTime)
		m.state = stateInitDone
		return m, tea.Quit
	}

	return m, nil
}

// View delegates rendering to cmd/views
func (m mainModel) View() string {
	switch m.state {
	case stateCancelled:
		return views.RenderCancelled(m.cancelledMsg)
	case stateMenu:
		return views.RenderMenu(m.menuOptions, m.menuCursor)
	case stateVersion:
		return views.RenderVersion(m.versionInfo)
	case stateInitInput:
		return views.RenderInput(m.projectNameInput, m.inputErr)
	case stateInitMode:
		return views.RenderMode(m.chosenName, m.modeOptions, m.modeCursor)
	case stateInitManual:
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
		pkgManager := m.scaffoldConfig.PackageManager
		if pkgManager == "" {
			pkgManager = views.GetPackageManager(m.chosenRecipe)
		}
		recipeLabel := "Manual Configuration"
		if m.chosenRecipe != "" {
			recipeLabel = views.GetRecipeLabel(m.recipeOptions, m.chosenRecipe)
		}
		return views.RenderRunning(m.chosenName, recipeLabel, pkgManager, m.stepNames, m.stepStatus, m.spinner)
	case stateInitDone:
		pkgManager := m.scaffoldConfig.PackageManager
		if pkgManager == "" {
			pkgManager = views.GetPackageManager(m.chosenRecipe)
		}
		recipeLabel := "Manual Configuration"
		if m.chosenRecipe != "" {
			recipeLabel = views.GetRecipeLabel(m.recipeOptions, m.chosenRecipe)
		}
		return views.RenderDone(m.chosenName, recipeLabel, pkgManager, m.stepNames, m.stepStatus, m.elapsedTime, m.scaffoldErr)
	}
	return ""
}

func RunTUI() {
	p := tea.NewProgram(initialModel(stateMenu, ""))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func RunTUIInit(projectName string) {
	state := stateInitInput
	inputErr := ""
	if projectName != "" {
		if err := validator.Validate(projectName); err != nil {
			inputErr = err.Error()
			state = stateInitInput
		} else {
			state = stateInitMode
		}
	}
	model := initialModel(state, projectName)
	model.inputErr = inputErr
	if projectName != "" {
		if inputErr == "" {
			model.chosenName = projectName
		} else {
			model.projectNameInput.SetValue(projectName)
		}
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
