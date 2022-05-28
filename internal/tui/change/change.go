package change

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hamdanjaveed/terraform-inspector/internal/pkg/tui"
	"github.com/hamdanjaveed/terraform-inspector/internal/tui/list"
)

type keyMap struct {
	Back key.Binding
	Help key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Back: key.NewBinding(key.WithKeys("esc", "enter"), key.WithHelp("esc/enter", "go back")),
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

type Bubble struct {
	content string

	width        int
	widthMargin  int
	height       int
	heightMargin int
	style        *tui.Styles

	vp   viewport.Model
	keys keyMap
	help help.Model
}

func NewBubble(
	styles *tui.Styles,
	width, wm, height, hm int,
) *Bubble {
	return &Bubble{
		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		style:        styles,

		vp:   viewport.Model{},
		keys: keys,
		help: help.New(),
	}
}

func (b Bubble) Init() tea.Cmd {
	return nil
}

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keys.Back):
			cmds = append(cmds, back())
		case key.Matches(msg, b.keys.Help):
			b.help.ShowAll = !b.help.ShowAll
		}
	case list.SelectMsg:
		b.content = msg.Change.Diff
		b.vp.SetContent(msg.Change.Diff)
	}

	m, cmd := b.vp.Update(msg)
	b.vp = m
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	vpView := b.vp.View()
	helpView := b.help.View(b.keys)
	return lipgloss.JoinVertical(lipgloss.Left, vpView, helpView)
}

func (b *Bubble) SetSize(w, h int) {
	b.width = w
	b.height = h
	b.vp.Width = w - b.widthMargin
	b.vp.Height = h - b.heightMargin - 1
	b.help.Width = w
}
