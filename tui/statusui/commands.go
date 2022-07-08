package statusui

import (
	"ipfs-tui/tui/constants"

	tea "github.com/charmbracelet/bubbletea"
)

// SelectPage() returns the selected page back to the MainModel
func SelectPage(pageTitle string) tea.Cmd {
	return func() tea.Msg {
		return constants.NavMsg{PageTitle: pageTitle}
	}
}
