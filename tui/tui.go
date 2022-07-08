package tui

import (
	"fmt"
	"ipfs-tui/tui/filesui"
	"ipfs-tui/tui/statusui"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

type sessionState int

const (
	statusView sessionState = iota
	filesView
)

type MainModel struct {
	state      sessionState
	status     tea.Model
	file       tea.Model
	windowSize tea.WindowSizeMsg
}

func StartTea() {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	m := New()
	p = tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// New initialize the main model for your program
func New() MainModel {
	return MainModel{
		state:  statusView,
		status: statusui.New(),
		file:   filesui.New(p),
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.file.Init()
}

// Update handle IO and commands
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg // pass this along to the entry view so it uses the full window size when it's initialized
	case statusui.SelectMsg:
		m.state = filesView
	}

	switch m.state {
	case statusView:
		newStatus, newCmd := m.status.Update(msg)
		statusModel, ok := newStatus.(statusui.Model)
		if !ok {
			panic("could not perform assertion on projectui model")
		}
		m.status = statusModel
		cmd = newCmd

	case filesView:
		// m.file = filesui.New(p)
		newFile, newCmd := m.file.Update(msg)
		fileModel, ok := newFile.(filesui.Model)
		if !ok {
			panic("could not perform assertion on projectui model")
		}
		m.file = fileModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m MainModel) View() string {
	switch m.state {
	case filesView:
		return m.file.View()
	default:
		return m.status.View()
	}
}
