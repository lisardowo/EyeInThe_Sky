package TUI

import (
	connection "EyeInThe_Sky/createConnection"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (

	purple = lipgloss.Color("#5F00FF")
	red = lipgloss.Color("#FF0000")
	darkGray = lipgloss.Color("#242424")

	leftBoxStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).
	BorderForeground(purple).
	Padding(1)

	rightBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(1)
	
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).
	Background(purple).
	Bold(true).
	Padding(0,1)

)

type HomeState struct {

	Operator     string
	VLAN         int
	Uptime       int 
	
	
}

const sigilASCII = `
 		      .      .     .
     		   \    /   //
      	  \		\  /   //
⠀⠀⠀⠀⠀   \⠀⠀⠀⠀⠀⢀/ ⢀⣀//⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀	⠀⠀⠀⢀\⣴⣶⠟⠃/\⠀⠀//⠙⢶⣤⡀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⢀⢤⣺⣿⣿\⠇⠀/ ⠀\//⣠⣤⣬⣿⣿⣷⣄⣂⠀⠀⠀⠀⠀
---(⠀⠀⢖⢼⡺⣿⣿⣿⣿⣿\/ ⠀⠀//⠈⠹⠿⠿⣿⣿⣿⣿⣚⠸⣣⠆⠀⠀)---
	⠀⠀⠀⠀⠑⠨⢟⣿⣿⣿/\⠠⡀//⠀\⠀⠀⢠⣿⣿⠿⠣⠙⠁⠁⠀⠀⠀
⠀⠀⠀⠀	⠀⠀⠀⠀⠉/ ⠻\//⡲⢖⣌\⠴⠟⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀/ ⠀⠀//⠀⠀⠀⠀⠈\⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		   /   // \     \
		  /   //   \     \
		 /   /'     '     '
=========================  EYE IN THE SKY
`

func renderHome(m HomeState, TerminalHeight int, TerminalWidth int, TrustLevel connection.TrustLevel, Uptime int) string {
	
	operator := m.Operator
	if operator == "" {
		operator = "operator"
	}

	vlan := m.VLAN
	if vlan == 0 {
		vlan = 10
	}
	
	leftContent := fmt.Sprintf(
		"%s\n\n%s\n\nSpecs:\n- Master: %s\n- Mode: %s\n- Uptime: %ds", //TODO implement a dynamic transform 4 minutes/hours(?), get uptime from servers 
		titleStyle.Render("WELCOME OPERATOR"),
		sigilASCII,
		operator,
		TrustLevel.TrustToString(),
		Uptime,
	)

rightContent := fmt.Sprintf(
	"Flags\n------------\n"+
		"Press Ctrl to set flags then enter and login\n"+
        "[u]      - Force USB Transport\n" +
        "[s]      - Force SSH Fallback\n" +
        "[m]      - Toggle Analysis Mode\n" +
	"Utility Info\n------------\n" +
		"[enter]  - Start Handshake\n" +
		"[Ctrl + q]      - Abort / Exit Session",
	
	)


	availableWidth := TerminalWidth
	if availableWidth <= 0 {
		availableWidth = 120
	}
// TODO calculations below can be moved to a utility getSize function
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

	// Home view does not consider the height of the screen to construct the render while dash does
	
	borderColor := TrustBorderColor(TrustLevel)
	leftStyle := leftBoxStyle.BorderForeground(borderColor).Width(leftWidth)
	rightStyle := rightBoxStyle.BorderForeground(borderColor).Width(rightWidth)
	
	left := leftStyle.Render(leftContent)
	right := rightStyle.Render(rightContent)

	if stacked {
		return lipgloss.JoinVertical(lipgloss.Left, left, lipgloss.NewStyle().Height(1).Render(""), right)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, left, lipgloss.NewStyle().Width(gap).Render(""), right)
}


