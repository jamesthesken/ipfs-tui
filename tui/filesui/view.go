package filesui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	borderStyle := lipgloss.NewStyle().
		PaddingRight(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.borderColor)

	statsTitle := m.list.Styles.Title.Copy()

	statsBox := borderStyle.Copy().PaddingLeft(2).BorderForeground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}).Width(m.list.Width() / 2).
		Height(m.list.Height() / 2)

	statsItem := lipgloss.NewStyle().PaddingTop(1).Render

	totalPeers := fmt.Sprintf("Connected Peers: %d", m.connectedPeers)

	stats := lipgloss.JoinHorizontal(lipgloss.Top, statsBox.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			statsTitle.Render("Stats"),
			statsItem(totalPeers)),
	))

	formatted := lipgloss.JoinHorizontal(lipgloss.Top, m.filetree.View(), lipgloss.JoinVertical(lipgloss.Right, borderStyle.Render(m.list.View()), stats))
	return formatted
}

// setBorderColor modifies a component's border color
func (m *Model) setBorderColor(borderColor lipgloss.AdaptiveColor) {
	m.borderColor = borderColor
}
