package TUI

import (
	connection "EyeInThe_Sky/createConnection"
	"EyeInThe_Sky/helpers"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type FocusPanel int

var (
	// TODO TEST FILE MOVE TO GENERAL FILE TO REUSE
	colorSecure   = lipgloss.Color("#5F00FF") // 
	colorUnsecure = lipgloss.Color("#FF0000") // 
	colorMuted    = lipgloss.Color("#ffffff") //#242424
	colorTextMute = lipgloss.Color("#ffffff")     //240 

	
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Underline(true)
)

const (
    PanelTelemetry FocusPanel = iota //0
    PanelCommands //1
    PanelLogs //2
)

type DashState struct {

	FocusedPanel    FocusPanel
	CPUUsage        float64
	RAMUsage        float64
	LogsBuffer      RingBuffer[helpers.LogEntry]
	ActiveFilter    *helpers.LogCategory
	CommandMode		bool
	CommandInput	string
	
}

func renderDash(state DashState, TerminalHeight int, TerminalWidth int, TrustLevel connection.TrustLevel) string { //TODO make a function that unifies all the different entries into one alone
	
	topHalfHeight := (TerminalHeight / 2) - 2
	bottomHalfHeight := (TerminalHeight / 2) - 2

	panelAwidth := (TerminalWidth / 2) - 2
	//panelBwidth := (state.TerminalWidth / 2) - 2
	activeBorderColor := TrustBorderColor(TrustLevel)

	// PANEL A
	telemetryBorder := lipgloss.NormalBorder()
	if state.FocusedPanel == PanelTelemetry {
		telemetryBorder = lipgloss.DoubleBorder()
	}
	
	telemetryStyle := lipgloss.NewStyle().
		Border(telemetryBorder).
		BorderForeground(onFocus(state.FocusedPanel == PanelTelemetry, activeBorderColor, colorMuted)).
		Width(panelAwidth).
		Height(topHalfHeight).
		Padding(1)

	statusStr := "SECURE"
	if TrustLevel == connection.Unsecure {
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
	if state.FocusedPanel == PanelCommands {
		commandsBorder = lipgloss.DoubleBorder()
	}

	commandsStyle := lipgloss.NewStyle().
		Border(commandsBorder).
		BorderForeground(onFocus(state.FocusedPanel == PanelCommands, activeBorderColor, colorMuted)).
		Width(panelAwidth).
		Height(topHalfHeight).
		Padding(1)

	
	var commandMenu string
	if TrustLevel == connection.Secure{
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
	if state.FocusedPanel == PanelLogs {
		logsBorder = lipgloss.DoubleBorder()
	}

	logsStyle := lipgloss.NewStyle().
		Border(logsBorder).
		BorderForeground(onFocus(state.FocusedPanel == PanelLogs, activeBorderColor, colorMuted)).
		Width(TerminalWidth - 2).
		Height(bottomHalfHeight).
		Padding(0,1)

	filtersLabel := "ALL"
	if state.ActiveFilter != nil {
		filtersLabel = state.ActiveFilter.IntToStringCategory()
	}

	lines := helpers.FormatFiltered(state.LogsBuffer.GetEntries(), state.ActiveFilter)
	logsContent := fmt.Sprintf("%s [filter: %s]\n\n%s", headerStyle.Render("REAL-TIME EVENT STREAM (PROCESS-AS-YOU-GO)"),
	filtersLabel,
	strings.Join(lines, "\n"),
)

	topHalf := lipgloss.JoinHorizontal(lipgloss.Top, telemetryStyle.Render(panelAContent), commandsStyle.Render(panelBContent))
	bottomHalf := logsStyle.Render(logsContent)
	if state.CommandMode {
		commandBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#1a1a1a")).
		Width(TerminalWidth).
		Padding(1,0)
	
		commandBar := commandBarStyle.Render(state.CommandInput)
	
		return lipgloss.JoinVertical(lipgloss.Left, topHalf, bottomHalf, commandBar)
	
	}

	return lipgloss.JoinVertical(lipgloss.Left, topHalf, bottomHalf)

}


	



// HELP

func onFocus(isFocused bool, focused, notFocused lipgloss.Color) lipgloss.Color { // color if focussed
	if isFocused {
		return focused
	}
	return notFocused
}

func renderProgressBar(percent float64, width int) string {
	bars := int((percent / 100.0) * float64(width))
	if bars < 0 { bars = 0 }
	if bars > width { bars = width }
	return strings.Repeat("█", bars) + strings.Repeat(" ", width-bars)
}

