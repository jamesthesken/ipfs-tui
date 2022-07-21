package filesui

import tea "github.com/charmbracelet/bubbletea"

// asynchronously fetch the number of connected peers

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.filetree.Init(), getIpfsStats)
}
