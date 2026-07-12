package TUI

import (
	connection "EyeInThe_Sky/createConnection"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    AnalysisMode connection.TrustLevel
    Operator      string
    VLAN          int
    BootAt        time.Time
    Uptime			time.Duration
    Width, Height         int
            
}


var modelWelcome = Model{Operator: "mock", VLAN: 10, 
AnalysisMode: connection.TrustLevel(2) , Uptime: 67, Width:  120, Height: 120 }
// Model is the exported TUI model used by main.


func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.Width = msg.Width
        m.Height = msg.Height
        return m, nil
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlQ {
            return m, tea.Quit
        }
    }

    return m, nil
}

func (m Model) View() string {
    return renderHome(m)
}

func (m Model) uptime() string {
    if m.BootAt.IsZero() {
        return "0s"
    }

    return time.Since(m.BootAt).Truncate(time.Second).String()
}


