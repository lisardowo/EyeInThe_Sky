package TUI

import (
	//connection "EyeInThe_Sky/createConnection"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	homeScreen int = iota
	dashScreen
)

// Main Welcome model with a counter to decide if its either un home or dash


type Model struct {
	CurrentMode int
	Home HomeState
	Dash DashState
	Width,Height	int
	LastKey      string	
	LastAction   string // Last change
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
		m.LastKey = msg.String()

		if msg.Type == tea.KeyCtrlQ {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyEnter:
			m.CurrentMode = dashScreen
			m.LastAction = "start handshake"
			return m, nil
		case tea.KeyEsc:
			m.CurrentMode = homeScreen
			m.LastAction = "return home"
			return m, nil
		
		}

		if(m.CurrentMode == dashScreen){
			switch msg.String() {
		
		
				case "l":
					currentValueBefore := fmt.Sprintf("currentValue After first 3 L : %d", m.Dash.FocusedPanel)
					m.Dash.LogsBuffer = append(m.Dash.LogsBuffer, currentValueBefore)
					if(m.Dash.FocusedPanel < 2){
						m.Dash.FocusedPanel += 1
					} else
					{
						currentValueAfter := fmt.Sprintf("currentValue Before first 3 L : %d", m.Dash.FocusedPanel)
						m.Dash.LogsBuffer = append(m.Dash.LogsBuffer, currentValueAfter)
						m.Dash.FocusedPanel = 0
					}
				
				case "h":
					currentValueBefore := fmt.Sprintf("currentValue After first 3 h : %d", m.Dash.FocusedPanel)
					m.Dash.LogsBuffer = append(m.Dash.LogsBuffer, currentValueBefore)
					if(m.Dash.FocusedPanel == 0){
						m.Dash.FocusedPanel = 2
					} else
					{
						currentValueAfter := fmt.Sprintf("currentValue Before first 3 h : %d", m.Dash.FocusedPanel)
						m.Dash.LogsBuffer = append(m.Dash.LogsBuffer, currentValueAfter)
						m.Dash.FocusedPanel -= 1 
					}
				default: 
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	m.Home.Uptime = time.Since(m.Home.BootAt)

	if m.CurrentMode == dashScreen {
		return renderDash(m.Dash, m.Height, m.Width)
	}

	return renderHome(m.Home, m.Height, m.Width)
}
