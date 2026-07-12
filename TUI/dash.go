package TUI

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

type DashState struct {
	TerminalWidth  int
	TerminalHeight int
	IsSecure       bool
	FocusedPanel   string   // "telemetry", "commands", "logs"
	CPUUsage       float64
	RAMUsage       float64
	LogsBuffer     []string
}

func renderDash(state DashState) string {
	
	topHalfHeight := (state.TerminalHeight / 2) - 2
	bottomHalfHeight := (state.TerminalHeight / 2) - 2
	
	leftPanelWidth := (state.TerminalWidth / 3) - 2
	rightPanelWidth := ((state.TerminalWidth * 2) / 3) - 2

	
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
		Width(leftPanelWidth).
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
		renderProgressBar(state.CPUUsage, leftPanelWidth-10), state.CPUUsage,
		renderProgressBar(state.RAMUsage, leftPanelWidth-10), state.RAMUsage,
	)

	// PANEL B
	commandsBorder := lipgloss.NormalBorder()
	if state.FocusedPanel == "commands" {
		commandsBorder = lipgloss.DoubleBorder()
	}

	commandsStyle := lipgloss.NewStyle().
		Border(commandsBorder).
		BorderForeground(ifThenColor(state.FocusedPanel == "commands", activeBorderColor, colorMuted)).
		Width(rightPanelWidth).
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

	return lipgloss.JoinVertical(lipgloss.Left, topHalf, bottomHalf)
}

// HELP

func ifThenColor(cond bool, t, f lipgloss.Color) lipgloss.Color {
	if cond {
		return t
	}
	return f
}

func renderProgressBar(percent float64, width int) string {
	bars := int((percent / 100.0) * float64(width))
	if bars < 0 { bars = 0 }
	if bars > width { bars = width }
	return strings.Repeat("█", bars) + strings.Repeat(" ", width-bars)
}