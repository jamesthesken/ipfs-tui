package filesui

import (
	"context"
	"fmt"
	"io"
	"ipfs-tui/tui/constants"
	"ipfs-tui/tui/statusui"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/filetree"

	shell "github.com/ipfs/go-ipfs-api"
)

type sessionState int
type item string
type itemDelegate struct{}

const (
	uiState sessionState = iota
	showPeerState
	listHeight = 14
)

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

type Model struct {
	filetree  filetree.Bubble
	list      list.Model
	activeBox int
	p         *tea.Program
}

func New(p *tea.Program) *Model {
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

	m := Model{list: l, filetree: filetreeModel}
	// m.p = p

	return &m
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

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

var sh *shell.Shell
var ncalls int

var _ = time.ANSIC

type peerInfo []shell.SwarmConnInfo

func generateSwarmMsg(m Model) func() tea.Msg {
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
func (m Model) getSwarmPeers() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return generateSwarmMsg(m)()
	})

}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.filetree.Init(),
		generateSwarmMsg(m))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case tea.WindowSizeMsg:
		m.filetree.SetSize(msg.Width/2, msg.Height)
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			cmd = statusui.SelectPage("Index")

		default:
			m.list, cmd = m.list.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
	// return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.filetree.View(), m.list.View())
	// formatted := lipgloss.JoinHorizontal(lipgloss.Top, m.filetree.View(), m.list.View())
	// return constants.DocStyle.Render(formatted)
}

// deactivateALlBubbles sets all bubbles to inactive.
func (m *Model) deactivateAllBubbles() {
	m.filetree.SetIsActive(false)
	m.list.SetShowFilter(false)
}

// toggleBox toggles between the two panes.
func (m *Model) toggleBox() {
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
