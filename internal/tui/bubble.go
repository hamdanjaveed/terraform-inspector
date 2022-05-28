package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hamdanjaveed/terraform-inspector/internal/pkg/tui"
	"github.com/hamdanjaveed/terraform-inspector/internal/tf"
	"github.com/hamdanjaveed/terraform-inspector/internal/tui/change"
	"github.com/hamdanjaveed/terraform-inspector/internal/tui/list"
)

type state int

const (
	listState state = iota
	changeState
)

type Bubble struct {
	Changes tf.ResourceChanges
	Summary string

	state        state
	height       int
	heightMargin int
	width        int
	widthMargin  int
	style        *tui.Styles
	boxes        []tea.Model
}

func NewBubble(
	rc tf.ResourceChanges,
	summary string,
	width, wm, height, hm int,
) *Bubble {
	b := &Bubble{
		Changes: rc,
		Summary: summary,

		state:        listState,
		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		style:        tui.NewStyles(),
		boxes:        make([]tea.Model, 2),
	}

	heightMargin := hm + lipgloss.Height(b.headerView())
	b.boxes[listState] = list.NewBubble(rc, summary, b.style, b.width, wm, b.height, heightMargin)
	b.boxes[changeState] = change.NewBubble(b.style, b.width, wm, b.height, heightMargin)

	return b
}

func (b *Bubble) Init() tea.Cmd {
	return b.setupCmd
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
		for i, bx := range b.boxes {
			m, cmd := bx.Update(msg)
			b.boxes[i] = m
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case list.SelectMsg:
		b.state = changeState
	case change.BackMsg:
		b.state = listState
	}

	m, cmd := b.boxes[b.state].Update(msg)
	b.boxes[b.state] = m
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

func (b *Bubble) SetSize(w, h int) {
	b.width = w
	b.height = h
}

func (b *Bubble) headerView() string {
	return ""
}

func (b *Bubble) View() string {
	return b.headerView() + b.boxes[b.state].View()
}

func (b *Bubble) setupCmd() tea.Msg {
	cmds := make([]tea.Cmd, 0)
	for _, bs := range b.boxes {
		if bs != nil {
			initCmd := bs.Init()
			// TODO: error handling
			// if initCmd != nil {
			// 	msg := initCmd()
			// 	switch msg := msg.(type) {
			// 	case CustomErrMsg:
			// 		return msg
			// 	}
			// }
			cmds = append(cmds, initCmd)
		}
	}
	return tea.Batch(cmds...)
}
