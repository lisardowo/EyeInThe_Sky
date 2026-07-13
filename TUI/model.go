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
    lastKey             string
    lastAction          string // TODO debug variable
            
}


var modelWelcome = Model{Operator: "mock", VLAN: 10, 
AnalysisMode: connection.TrustLevel(2) , Uptime: 67, Width:  120, Height: 120 }
// Model is the exported TUI model used by main.


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

        //this reads a signal to force quit the program in case anything happens
        
        if msg.Type == tea.KeyCtrlQ {
            return m, tea.Quit
        } 
        
        switch msg.String(){
            case "h"
            case "l"
            case "enter"
        }
    
    }

    return m, nil
}


func (m Model) View() string {
    return renderHome(m) //TODO debug home view
    // TODO harcorded debug dash model
    /* debugArray := []string{"mock", "mock", "mock", "mock"} //TODO delete
    return renderDash(DashState{
    TerminalWidth:  150,
	TerminalHeight: 13,
	IsSecure:       true,
	FocusedPanel:   "telemetry",   // "telemetry", "commands", "logs" | Maybe use and enum, when h (-1), l (+1) regenerate each time that keys are pressed 
	CPUUsage:       89.5,
	RAMUsage:       95.2,
	LogsBuffer:     debugArray,
    Width:         120,} ) */

    

}

func (m Model) uptime() string {
    if m.BootAt.IsZero() {
        return "0s"
    }

    return time.Since(m.BootAt).Truncate(time.Second).String()
}


