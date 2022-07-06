package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/filetree"

	shell "github.com/ipfs/go-ipfs-api"
)

/// --- TUI
const listHeight = 14

var (
	titleStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type sessionState int

const (
	uiState sessionState = iota
	showPeerState
)

type model struct {
	filetree  filetree.Bubble
	list      list.Model
	choice    string
	quitting  bool
	activeBox int
	state     sessionState
}

var sh *shell.Shell
var ncalls int

var _ = time.ANSIC

type peerInfo []shell.SwarmConnInfo

func generateSwarmMsg(m model) func() tea.Msg {
	return func() tea.Msg {
		sh = shell.NewShell("localhost:5001")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		swarm, err := sh.SwarmPeers(ctx)
		if err != nil {
			fmt.Errorf("Error, %s", err)
		}

		return peerInfo(swarm.Peers)
	}
}

// doTick()
func (m model) getSwarmPeers() tea.Cmd {

	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return generateSwarmMsg(m)()
	})

}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.filetree.Init(),
		generateSwarmMsg(m))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case peerInfo:
		peers := []list.Item{}
		for i := 0; i < len(msg); i++ {
			peers = append(peers, item(msg[i].Peer))
		}
		m.list.SetItems(peers)
		m.list.Update(msg)
		cmds = append(cmds, m.getSwarmPeers())
		// return m, m.getSwarmPeers()

	case tea.WindowSizeMsg:
		m.filetree.SetSize(msg.Width/2, msg.Height)
		m.list.SetWidth(msg.Width)
		// return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			m.toggleBox()
		case "ctrl+c":
			m.quitting = true
			cmds = append(cmds, tea.Quit)
			// return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			cmds = append(cmds, tea.Quit)
			// return m, tea.Quit
		default:
			m.list, cmd = m.list.Update(msg)
			// cmds = append(cmds, cmd)
		}
	}

	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
	// return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.filetree.View(), m.list.View())
}

// deactivateALlBubbles sets all bubbles to inactive.
func (m *model) deactivateAllBubbles() {
	m.filetree.SetIsActive(false)
	m.list.SetShowFilter(false)
}

// toggleBox toggles between the two panes.
func (m *model) toggleBox() {
	m.activeBox = (m.activeBox + 1) % 2
	if m.activeBox == 0 {
		m.deactivateAllBubbles()
		m.filetree.SetIsActive(true)
		print(m.list.FilterValue())
	} else {
		m.deactivateAllBubbles()
		// switch m.state {
		// case showPeerState:
		// 	m.deactivateAllBubbles()
		// 	m.list.SetShowFilter(true)
		// }
	}

}

func main() {
	items := []list.Item{}
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Connected Peers"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	filetreeModel := filetree.New(
		true,
		true,
		"",
		"",
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "63", Dark: "63"},
		lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
	)

	m := model{list: l, filetree: filetreeModel}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
