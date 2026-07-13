package TUI

import (
	connection "EyeInThe_Sky/createConnection"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
    
    homeScreen int = iota // TODO Typedef the int as a screen.. ?
    dashScreen       
)


type mainModel struct{ 
    currentMode     int // int representing home(0) or dash(1)
    home            Model
    dash            DashState

}

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
    debugArray := []string{"mock", "mock", "mock", "mock"} 
    mainModel := mainModel{currentMode: 0, 
        home: Model{AnalysisMode: connection.Secure, Operator: "mock", VLAN: 1, BootAt: time.Now(), Uptime: 67, Width: 120, Height: 240, lastKey: "mock", lastAction: "mock",},
        dash: DashState{TerminalWidth: 120, TerminalHeight:  240, IsSecure:  true, FocusedPanel:  "telemetry", CPUUsage: 92.4, RAMUsage: 21.2, LogsBuffer: debugArray, Width: 240},
    }
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
            case "h":
            
            case "l":
            
            case "enter":
                 // workin
                 //println(mainModel)
                //renderDashTEST(mainModel.dashTest)
                //return mainModel.dash, nil  
                _ = mainModel
        }
    
    }

    return m, nil
}


func renderDashTEST(state DashState) string {
	
	topHalfHeight := (state.TerminalHeight / 2) - 2
	bottomHalfHeight := (state.TerminalHeight / 2) - 2

	panelAwidth := (state.TerminalWidth / 2) - 2
	//panelBwidth := (state.TerminalWidth / 2) - 2
	
	var activeBorderColor lipgloss.Color
	if state.IsSecure {
		activeBorderColor = colorSecure
	} else {
		activeBorderColor = colorUnsecure
	}

	// PANEL A
	telemetryBorder := lipgloss.NormalBorder()
	if state.FocusedPanel == "telemetry" {
		telemetryBorder = lipgloss.DoubleBorder()
	}
	
	telemetryStyle := lipgloss.NewStyle().
		Border(telemetryBorder).
		BorderForeground(ifThenColor(state.FocusedPanel == "telemetry", activeBorderColor, colorMuted)).
		Width(panelAwidth).
		Height(topHalfHeight).
		Padding(1)

	statusStr := "SECURE"
	if !state.IsSecure {
		statusStr = " UNSECURE (SANDBOX ACTIVE)"
	}

	panelAContent := fmt.Sprintf(
		"%s\n\nSTATUS: %s\nLINK: PHYSICAL_USB\n\n%s\nCPU: [%s] %.0f%%\nRAM: [%s] %.0f%%",
		headerStyle.Render("SYSTEM METRICS"),
		statusStr,
		headerStyle.Render("RESOURCES"),
		renderProgressBar(state.CPUUsage, panelAwidth-10), state.CPUUsage,
		renderProgressBar(state.RAMUsage, panelAwidth-10), state.RAMUsage,
	)

	// PANEL B
	commandsBorder := lipgloss.NormalBorder()
	if state.FocusedPanel == "commands" {
		commandsBorder = lipgloss.DoubleBorder()
	}

	commandsStyle := lipgloss.NewStyle().
		Border(commandsBorder).
		BorderForeground(ifThenColor(state.FocusedPanel == "commands", activeBorderColor, colorMuted)).
		Width(panelAwidth).
		Height(topHalfHeight).
		Padding(1)

	//
	var commandMenu string
	if state.IsSecure {
		commandMenu = "> [1] Deploy Update\n> [2] Network Diagnostics\n> [3] Open SSH Terminal Session"
	} else {
		commandMenu = lipgloss.NewStyle().Foreground(colorUnsecure).Render(
			"> [1] ISOLATE VLAN CRITICAL\n> [2] KILL NETWORK INTERFACE\n> [3] DUMP MEMORY BUFFER TO DISK",
		)
	}

	panelBContent := fmt.Sprintf(
		"%s\n\n%s\n\n\n%s\n[h/l] Switch Focus  |  [L] Lock Session",
		headerStyle.Render("AVAILABLE ACTIONS"),
		commandMenu,
		lipgloss.NewStyle().Foreground(colorTextMute).Render("Navigation Help:"),
	)

	// PANEL C
	logsBorder := lipgloss.NormalBorder()
	if state.FocusedPanel == "logs" {
		logsBorder = lipgloss.DoubleBorder()
	}

	logsStyle := lipgloss.NewStyle().
		Border(logsBorder).
		BorderForeground(ifThenColor(state.FocusedPanel == "logs", activeBorderColor, colorMuted)).
		Width(state.TerminalWidth - 2).
		Height(bottomHalfHeight).
		Padding(0, 1)

	logsContent := fmt.Sprintf(
		"%s\n%s",
		headerStyle.Render("REAL-TIME EVENT STREAM (PROCESS-AS-YOU-GO)"),
		strings.Join(state.LogsBuffer, "\n"),
	)

	topHalf := lipgloss.JoinHorizontal(lipgloss.Top, telemetryStyle.Render(panelAContent), commandsStyle.Render(panelBContent))
	bottomHalf := logsStyle.Render(logsContent)

	availableWidth := state.Width
	if availableWidth <= 0 {
		availableWidth = 120
	}
	minPaneWidth := 28
	gap := 2
	stacked := availableWidth < 90

	leftWidth := availableWidth - 4
	rightWidth := leftWidth - 4
	if !stacked {
		leftWidth = int(float64(availableWidth) * 0.58)
		rightWidth = availableWidth - leftWidth - gap - 6
		if leftWidth < 34 {
			leftWidth = 34
		}
		if rightWidth < 24 {
			rightWidth = 24 - 6
		}
		if leftWidth+rightWidth+gap > availableWidth {
			rightWidth = availableWidth - leftWidth - gap - 6
		}
		if rightWidth < minPaneWidth {
			rightWidth = minPaneWidth
		}
	} else {
		if leftWidth < minPaneWidth {
			leftWidth = minPaneWidth
		}
		rightWidth = leftWidth - 6
	}

	if stacked {
		return lipgloss.JoinVertical(lipgloss.Left, topHalf, bottomHalf)
	}

	return lipgloss.JoinVertical(lipgloss.Left, topHalf, bottomHalf)

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


