package statusui

import tea "github.com/charmbracelet/bubbletea"

func selectPage() tea.Cmd {
	return func() tea.Msg {
		return SelectMsg{pageID: 4}
	}
}
