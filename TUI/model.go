package TUI

import (
	//connection "EyeInThe_Sky/createConnection"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	homeScreen int = iota
	dashScreen
)

// Main Welcome model with a counter to decide if its either un home or dash


type Model struct {
	currentMode int
	home HomeState
	dash DashState
	Width,Height	int
	lastKey      string	
	lastAction   string // Last change
}


//var modelWelcome = Model.HomeState{Operator: "mock", VLAN: 10,
	//AnalysisMode: connection.TrustLevel(2), Uptime: 67, Width: 120, Height: 120} // TODO debug model

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
				if (m.dash.FocusedPanel <= 3){
					m.dash.FocusedPanel = 0
				} else {
					m.dash.FocusedPanel += 1
				}
			case "h":
				m.dash.FocusedPanel -= 1
			default: 
		}
	}

	return m, nil
}

func (m Model) View() string {
	m.home.Uptime = time.Since(m.home.BootAt)

	if m.currentMode == dashScreen {
		return renderDash(m.dash)
	}

	return renderHome(m.home)
}
