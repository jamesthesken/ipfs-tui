package filesui

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	borderStyle := lipgloss.NewStyle().
		PaddingRight(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.borderColor)
	formatted := lipgloss.JoinHorizontal(lipgloss.Top, m.filetree.View(), borderStyle.Render(m.list.View()))
	return formatted
}

// setBorderColor modifies a component's border color
func (m *Model) setBorderColor(borderColor lipgloss.AdaptiveColor) {
	m.borderColor = borderColor
}
