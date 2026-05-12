package TUI

import (
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
=========================EYE IN THE SKY
`

func (m Model) welcomeView() string {
	operator := m.Operator
	if operator == "" {
		operator = "operator"
	}

	vlan := m.VLAN
	if vlan == 0 {
		vlan = 10
	}

	leftContent := fmt.Sprintf(
		"%s\n\n%s\n\nSpecs:\n- Master: %s\n- Mode: %s\n- Uptime: %s",
		titleStyle.Render("WELCOME OPERATOR"),
		sigilASCII,
		operator,
		m.modeLabel(),
		m.uptime(),
	)

	rightContent := fmt.Sprintf(
		"UTILITY INFO\n------------\n"+
			"Press [ENTER] to Init Handshake\n"+
			"Press [Ctrl+Q]     to Exit Analysis\n\n"+
			"VLAN: %d (Isolated)\n"+
			"Terminal: %dx%d",
		vlan,
		m.Width,
		m.Height,
	)

	availableWidth := m.Width
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

	leftStyle := leftBoxStyle.Width(leftWidth)
	rightStyle := rightBoxStyle.Width(rightWidth) 
	
	left := leftStyle.Render(leftContent)
	right := rightStyle.Render(rightContent)

	if stacked {
		return lipgloss.JoinVertical(lipgloss.Left, left, lipgloss.NewStyle().Height(1).Render(""), right)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, left, lipgloss.NewStyle().Width(gap).Render(""), right)
}