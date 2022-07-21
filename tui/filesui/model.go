package filesui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/filetree"
)

type sessionState int

type peer string

type item struct {
	fileName    string
	description string
	Type        uint8
	Size        uint64
	Hash        string
}

type itemDelegate struct{}

const (
	uiState sessionState = iota
	navFiles
	listHeight = 14
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ffff"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type Model struct {
	filetree       filetree.Bubble
	list           list.Model
	state          sessionState
	activeBox      int
	connectedPeers int
	borderColor    lipgloss.AdaptiveColor
}

func New() tea.Model {

	filetreeModel := filetree.New(
		true,
		false,
		"",
		"",
		lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"},
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "63", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
	)

	items := getIpfsFiles()
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "IPFS Files"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := Model{
		list:     l,
		filetree: filetreeModel,
	}

	m.filetree.SetSize(5, 10)
	m.list.SetWidth(10)

	m.setBorderColor(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"})

	return m
}

func (i item) FilterValue() string { return "" }

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.fileName)

	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("  " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}
