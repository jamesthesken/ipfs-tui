package filesui

import (
	"context"
	"fmt"
	"io"
	"ipfs-tui/tui/constants"
	"ipfs-tui/tui/statusui"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/filetree"

	shell "github.com/ipfs/go-ipfs-api"
	files "github.com/ipfs/go-ipfs-files"
)

type sessionState int

type item struct {
	fileName    string
	description string
}

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
	filetree filetree.Bubble
	list     list.Model
}

func New(p *tea.Program) *Model {
	items := []list.Item{}
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "IPFS Files"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	filetreeModel := filetree.New(
		true,
		true,
		"./",
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

	str := fmt.Sprintf("%s", i.fileName)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

/*
	addFile() returns the CID of a file added to IPFS
*/
// TODO: Clean up error handling!
func addFile(ctx context.Context, path string, fileName string) string {
	// Opens a 'shell' to the local IPFS node
	sh := shell.NewShell("localhost:5001")
	file, err := os.Open(path)
	if err != nil {
		fmt.Print(err)
	}
	fileReader := files.NewReaderFile(file)

	cid, err := sh.Add(fileReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	// When referencing files added to ipfs, you can get them at /ipfs/<cid>
	ipfsPath := fmt.Sprintf("/ipfs/%s", cid)

	// Human-readable name in the / directory of the MFS
	// TODO: allow the user to select which directory this file will go to
	newPath := fmt.Sprintf("/%s", fileName)

	// TODO: Look into utilizing FilesWrite to write directly to IPFS' MFS
	// FilesCp copies the IPFS file we just added to the MFS.
	err = sh.FilesCp(ctx, ipfsPath, newPath)
	if err != nil {
		err = fmt.Errorf("error: %s", err)
		fmt.Print(err)
	}

	return cid
}

func (m Model) Init() tea.Cmd {
	return m.filetree.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.filetree.SetSize(msg.Width/2, msg.Height)
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			cmd = statusui.SelectPage("Index")
		case key.Matches(msg, constants.Keymap.Enter):
			ctx := context.Background()
			selectedFile := m.filetree.GetSelectedItem()
			listSize := len(m.list.Items())
			// ctx = context.WithValue(ctx, "path", selectedFile)
			cid := addFile(ctx, selectedFile.FileName(), selectedFile.ShortName())
			m.list.InsertItem(listSize+1,
				item{
					fileName:    cid,
					description: selectedFile.Description(),
				})
		case key.Matches(msg, constants.Keymap.ToggleFile):

		default:
			m.list, cmd = m.list.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	formatted := lipgloss.JoinHorizontal(lipgloss.Top, m.filetree.View(), m.list.View())
	return constants.DocStyle.Render(formatted)
}
