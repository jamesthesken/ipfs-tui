package constants

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// ErrStyle provides styling for error messages
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

// AlertStyle provides styling for alert messages
var AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render

// TODO: Might be better to implement pageID's as an integer instead?
type NavMsg struct {
	PageTitle string
}

type keymap struct {
	Navigate   key.Binding
	Enter      key.Binding
	Back       key.Binding
	ToggleFile key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Navigate: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "nav"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	ToggleFile: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle info"),
	),
}
