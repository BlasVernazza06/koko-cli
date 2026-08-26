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
	"github.com/BlasVernazza06/koko-cli/internal/compatibility"
	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/BlasVernazza06/koko-cli/internal/validator"
	"github.com/BlasVernazza06/koko-cli/internal/vfs"
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

// findNextEnabled returns the index of the next enabled option in the given direction (delta: +1 for down, -1 for up).
func findNextEnabled(options []views.SelectOption, currentIdx int, delta int) int {
	if len(options) == 0 {
		return 0
	}
	next := currentIdx + delta
	for next >= 0 && next < len(options) {
		if !options[next].Disabled {
			return next
		}
		next += delta
	}
	return currentIdx
}

// findFirstEnabled returns the index of the first enabled option.
func findFirstEnabled(options []views.SelectOption) int {
	for i, opt := range options {
		if !opt.Disabled {
			return i
		}
	}
	return 0
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
			title:   "Select Frontend Framework",
			label:   "Frontend",
			options: compatibility.BaseOptions(compatibility.StepFrontend),
		},
		{
			title:   "Select Backend Framework / Runtime",
			label:   "Backend",
			options: compatibility.BaseOptions(compatibility.StepBackend),
		},
		{
			title:   "Select Package Manager",
			label:   "Package Manager",
			options: compatibility.BaseOptions(compatibility.StepPackageManager),
		},
		{
			title:   "Select Database",
			label:   "Database",
			options: compatibility.BaseOptions(compatibility.StepDatabase),
		},
		{
			title:   "Select ORM / Query Builder",
			label:   "ORM / Tool",
			options: compatibility.BaseOptions(compatibility.StepORM),
		},
		{
			title:   "Select Addons / Tooling",
			label:   "Addons",
			options: compatibility.BaseOptions(compatibility.StepAddons),
		},
		{
			title:   "Initialize Git Repository?",
			label:   "Git",
			options: compatibility.BaseOptions(compatibility.StepGit),
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
	var cachedVFS *vfs.VFS

	steps := []runnerStep{
		{
			name: "Generating templates and configuration in memory...",
			run: func(p string, c scaffold.ScaffoldConfig) error {
				var err error
				cachedVFS, err = scaffold.GenerateVFS(c)
				return err
			},
		},
		{
			name: "Writing files safely to disk...",
			run: func(p string, c scaffold.ScaffoldConfig) error {
				if cachedVFS == nil {
					var err error
					cachedVFS, err = scaffold.GenerateVFS(c)
					if err != nil {
						return err
					}
				}
				return scaffold.WriteTree(cachedVFS, p)
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
					firstOpts := compatibility.GetStepOptions(0, m.manualSelections)
					m.manualCursors[0] = findFirstEnabled(firstOpts)
				}
			}

		case stateInitManual:
			currentOptions := compatibility.GetStepOptions(m.manualStepIdx, m.manualSelections)
			cursor := m.manualCursors[m.manualStepIdx]
			if cursor >= len(currentOptions) || currentOptions[cursor].Disabled {
				cursor = findFirstEnabled(currentOptions)
				m.manualCursors[m.manualStepIdx] = cursor
			}

			switch msg.String() {
			case "up", "k":
				m.manualCursors[m.manualStepIdx] = findNextEnabled(currentOptions, cursor, -1)
			case "down", "j":
				m.manualCursors[m.manualStepIdx] = findNextEnabled(currentOptions, cursor, 1)
			case "enter":
				if cursor >= 0 && cursor < len(currentOptions) && !currentOptions[cursor].Disabled {
					selectedOpt := currentOptions[cursor]
					m.manualSelections[m.manualStepIdx] = selectedOpt

					// Advance to next step or complete configuration
					if m.manualStepIdx < len(m.manualSteps)-1 {
						m.manualStepIdx++
						nextOpts := compatibility.GetStepOptions(m.manualStepIdx, m.manualSelections)
						m.manualCursors[m.manualStepIdx] = findFirstEnabled(nextOpts)
					} else {
						// All steps completed -> initialize project
						initGit := m.manualSelections[compatibility.StepGit].Value == "yes"
						m.scaffoldConfig = scaffold.ScaffoldConfig{
							ProjectName:    m.chosenName,
							Recipe:         "",
							InitGit:        initGit,
							Frontend:       m.manualSelections[compatibility.StepFrontend].Value,
							Backend:        m.manualSelections[compatibility.StepBackend].Value,
							PackageManager: m.manualSelections[compatibility.StepPackageManager].Value,
							Database:       m.manualSelections[compatibility.StepDatabase].Value,
							ORM:            m.manualSelections[compatibility.StepORM].Value,
							Addons:         m.manualSelections[compatibility.StepAddons].Value,
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
				}

			case "esc":
				if m.manualStepIdx > 0 {
					m.manualStepIdx--
					prevOpts := compatibility.GetStepOptions(m.manualStepIdx, m.manualSelections)
					if m.manualCursors[m.manualStepIdx] >= len(prevOpts) || prevOpts[m.manualCursors[m.manualStepIdx]].Disabled {
						m.manualCursors[m.manualStepIdx] = findFirstEnabled(prevOpts)
					}
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

func (m mainModel) getSummaryItems() []views.SummaryItem {
	var summary []views.SummaryItem
	summary = append(summary, views.SummaryItem{
		Label: "Project name",
		Value: m.chosenName,
	})

	if m.chosenRecipe != "" {
		recipeLabel := views.GetRecipeLabel(m.recipeOptions, m.chosenRecipe)
		pkgManager := m.scaffoldConfig.PackageManager
		if pkgManager == "" {
			pkgManager = views.GetPackageManager(m.chosenRecipe)
		}
		summary = append(summary, views.SummaryItem{
			Label: "Recipe / Stack",
			Value: recipeLabel,
		})
		summary = append(summary, views.SummaryItem{
			Label: "Package Manager",
			Value: pkgManager,
		})
	} else {
		for i, step := range m.manualSteps {
			if i < len(m.manualSelections) && m.manualSelections[i].Label != "" {
				summary = append(summary, views.SummaryItem{
					Label: step.label,
					Value: m.manualSelections[i].Label,
				})
			}
		}
	}
	return summary
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
		currentOptions := compatibility.GetStepOptions(m.manualStepIdx, m.manualSelections)
		cursor := m.manualCursors[m.manualStepIdx]
		if cursor >= len(currentOptions) || currentOptions[cursor].Disabled {
			cursor = findFirstEnabled(currentOptions)
		}
		return views.RenderManual(m.chosenName, history, currentStep.title, currentOptions, cursor)

	case stateInitRecipe:
		return views.RenderRecipe(m.chosenName, m.recipeOptions, m.recipeCursor)
	case stateInitRunning:
		return views.RenderRunning(m.getSummaryItems(), m.stepNames, m.stepStatus, m.spinner)
	case stateInitDone:
		pkgManager := m.scaffoldConfig.PackageManager
		if pkgManager == "" {
			pkgManager = views.GetPackageManager(m.chosenRecipe)
		}
		return views.RenderDone(m.getSummaryItems(), m.chosenName, pkgManager, m.stepNames, m.stepStatus, m.elapsedTime, m.scaffoldErr)
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
