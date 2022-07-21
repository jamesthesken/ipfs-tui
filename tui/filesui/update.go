package filesui

import (
	"context"
	"ipfs-tui/tui/constants"
	"ipfs-tui/tui/statusui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.filetree.SetSize(msg.Width/2, msg.Height)
		// m.list.SetWidth(msg.Width/2)
		m.list.SetSize(msg.Width/2, msg.Height/2)

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			cmd = statusui.SelectPage("Index")
		case key.Matches(msg, constants.Keymap.Tab):
			m.toggleBox()
			m.list.ResetSelected()
		case key.Matches(msg, constants.Keymap.Enter):
			// check if the filetree is active first
			if m.activeBox == 0 {
				ctx := context.Background()
				selectedFile := m.filetree.GetSelectedItem()
				listSize := len(m.list.Items())

				// add file to local IPFS node - see utils.go
				cid := addFile(ctx, selectedFile.FileName(), selectedFile.ShortName())
				m.list.InsertItem(listSize+1,
					item{
						fileName:    cid,
						description: selectedFile.Description(),
					})
			}
		default:
			// can't toggle focus on the list, so if the list isn't selected this makes sure it doesn't 'respond' to user input
			if m.activeBox == 0 {
				m.list.ResetSelected()
			}
			m.list, cmd = m.list.Update(msg)

		}
		cmds = append(cmds, cmd)
	}
	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// toggleBox toggles between the two boxes.
func (m *Model) toggleBox() {
	m.activeBox = (m.activeBox + 1) % 2
	if m.activeBox == 0 {
		m.deactivateAllBubbles()
		m.filetree.SetIsActive(true)
		m.filetree.SetBorderColor(lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"})
	} else {
		m.deactivateAllBubbles()
		selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
		m.setBorderColor(lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"})
	}
}

// deactivateALlBubbles sets all bubbles to inactive.
func (m *Model) deactivateAllBubbles() {
	m.filetree.SetIsActive(false)
	m.filetree.SetBorderColor(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"})
	m.setBorderColor(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"})
	// this is the best I could come up with for now, until a method is available to toggle focus on the list.
	// this simply sets the active item to white
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ffff"))
}
