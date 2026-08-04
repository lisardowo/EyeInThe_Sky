package TUI

import (
	connection "EyeInThe_Sky/createConnection"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	homeScreen int = iota
	dashScreen
)

// Main Welcome model with a counter to decide if its either un home or dash

type Model struct {
	AnalysisMode connection.TrustLevel
	Operator     string
	VLAN         int
	BootAt       time.Time
	Uptime       time.Duration
	Width, Height int
	lastKey      string	
	lastAction   string // Last change
	currentMode  int
}

var modelWelcome = Model{Operator: "mock", VLAN: 10,
	AnalysisMode: connection.TrustLevel(2), Uptime: 67, Width: 120, Height: 120} // TODO debug model

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		m.lastKey = msg.String()

		if msg.Type == tea.KeyCtrlQ {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyEnter:
			m.currentMode = dashScreen
			m.lastAction = "start handshake"
			return m, nil
		case tea.KeyEsc:
			m.currentMode = homeScreen
			m.lastAction = "return home"
			return m, nil
		
		}

		switch msg.String() {
    		case "l":
			case "h":
			default: 
		}
	}

	return m, nil
}

func (m Model) View() string {
	m.Uptime = time.Since(m.BootAt)

	if m.currentMode == dashScreen {
		return renderDash(dashStateFromModel(m))
	}

	return renderHome(m)
}

func dashStateFromModel(m Model) DashState {
	width := m.Width
	if width <= 0 {
		width = 120
	}

	height := m.Height
	if height <= 0 {
		height = 40
	}

	// Construct the dash state for the first change from home to dash
	return DashState{
		TerminalWidth:  width,
		TerminalHeight: height,
		IsSecure:       m.AnalysisMode == connection.Secure,
		FocusedPanel:  PanelTelemetry, //This would be better keeping in memory last known panel..?
		CPUUsage:       92.4,
		RAMUsage:       21.2,
		LogsBuffer:     []string{"boot sequence ready", "press esc to return home"},
		Width:          width,
	}
}