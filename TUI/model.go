package TUI

import (
	//connection "EyeInThe_Sky/createConnection"
	connection "EyeInThe_Sky/createConnection"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	HomeScreen int = iota
	DashScreen
)

type tickMsg time.Time

type Model struct {
	WhichScreen int
	TrustLevel connection.TrustLevel
	Uptime			int // Made this a pointer so it modifies outside of the msg ?
	Width,Height	int
	LastKey      string	
	LastAction   string // Last change
	Home HomeState
	Dash DashState
}




func (m Model) Init() tea.Cmd {
	 return tea.Batch(tea.WindowSize(), tickCmd())
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	
	case tickMsg:
		
		m.Home.Uptime++
		return m, tickCmd()

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
			m.WhichScreen = DashScreen
			m.LastAction = "start handshake"
			return m, nil
		case tea.KeyEsc:
			m.WhichScreen = HomeScreen
			m.LastAction = "return home"
			return m, nil
		
		}
		//TODO move this to a wrapping fuction to keep update clean
		if(m.WhichScreen == DashScreen){
			switch msg.String() {
		
		
				case "l":
					currentValueBefore := fmt.Sprintf("currentValue After first 3 L : %d", m.Dash.FocusedPanel)
					m.Dash.LogsBuffer.Add(currentValueBefore) //TODO debug = append(m.Dash.LogsBuffer, currentValueBefore)
					if(m.Dash.FocusedPanel < 2){
						m.Dash.FocusedPanel += 1
					} else
					{
						currentValueAfter := fmt.Sprintf("currentValue Before first 3 L : %d", m.Dash.FocusedPanel)
						m.Dash.LogsBuffer.Add(currentValueAfter)
						m.Dash.FocusedPanel = 0
					}
				
				case "h":
					currentValueBefore := fmt.Sprintf("currentValue After first 3 h : %d", m.Dash.FocusedPanel)
					m.Dash.LogsBuffer.Add(currentValueBefore)
					if(m.Dash.FocusedPanel == 0){
						m.Dash.FocusedPanel = 2
					} else
					{
						currentValueAfter := fmt.Sprintf("currentValue Before first 3 h : %d", m.Dash.FocusedPanel)
						m.Dash.LogsBuffer.Add(currentValueAfter)
						m.Dash.FocusedPanel -= 1 
					}
				default: // TODO implement an "warning" notification letter not recognized, this goes over the layout in the corner of commands
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	//m.Home.Uptime = time.Since(m.Home.BootAt)
	
	if m.WhichScreen == DashScreen {
		return renderDash(m.Dash, m.Height, m.Width, m.TrustLevel)
	}

	return renderHome(m.Home, m.Height, m.Width, m.TrustLevel, m.Home.Uptime)  
}

func tickCmd() tea.Cmd {
	return tea.Tick(1 * time.Second, func(t time.Time) tea.Msg{
		return tickMsg(t)
	})
}