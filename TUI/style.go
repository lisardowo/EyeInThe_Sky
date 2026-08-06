package TUI

import (
	connection "EyeInThe_Sky/createConnection"

	"github.com/charmbracelet/lipgloss"
)


func TrustBorderColor(level connection.TrustLevel) lipgloss.Color {
	if level == connection.Secure {
		return colorSecure
	}
	return colorUnsecure
}
