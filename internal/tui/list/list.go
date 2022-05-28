package list

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hamdanjaveed/terraform-inspector/internal/pkg/tui"
	"github.com/hamdanjaveed/terraform-inspector/internal/tf"
)

type item struct {
	tf.ResourceChange
	expanded bool
}

func (i item) Title() string {
	return i.Address
}

func (i item) FilterValue() string {
	return i.Title()
}

type itemDelegate struct {
	style *tui.Styles
}

func (d itemDelegate) Height() int {
	// s := renderItem(d)
	return 3
}

func (d itemDelegate) Spacing() int {
	return 3
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	io.WriteString(w, renderItem(i, index, m.Index(), d.style))
}

func renderItem(item item, index, mIndex int, style *tui.Styles) string {
	var s string
	if index == mIndex {
		s = style.CreateChangeBody.Render(item.Title())
	} else {
		s = style.Body.Render(item.Title())
	}

	if item.expanded {
		s += item.Diff
	}

	return s
}

type Bubble struct {
	changes tf.ResourceChanges
	summary string

	width        int
	widthMargin  int
	height       int
	heightMargin int
	style        *tui.Styles

	list list.Model
}

func NewBubble(
	rc tf.ResourceChanges,
	summary string,
	styles *tui.Styles,
	width, wm, height, hm int,
) *Bubble {
	items := make([]list.Item, len(rc))
	for i, c := range rc {
		items[i] = item{c, false}
	}
	l := list.New(items, itemDelegate{style: styles}, width-wm, height-hm)
	l.SetShowFilter(false)
	l.SetShowHelp(true)
	l.SetShowPagination(true)
	l.SetShowStatusBar(false)
	l.SetShowTitle(true)
	l.SetFilteringEnabled(true)
	l.Title = summary

	b := &Bubble{
		changes: rc,
		summary: summary,

		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		style:        styles,

		list: l,
	}
	b.SetSize(width, height)
	return b
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
		switch msg.String() {
		case "enter", " ":
			cmd := selectChange(b.changes[b.list.Index()])
			cmds = append(cmds, cmd)
		case "x":
			i, ok := b.list.SelectedItem().(item)
			if !ok {
				println("asdf")
				break
			}
			i.expanded = !i.expanded
			b.list.SetItem(b.list.Index(), i)
		}
	}

	m, cmd := b.list.Update(msg)
	b.list = m
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

func (b Bubble) View() string {
	return b.list.View()
}

func (b *Bubble) SetSize(w, h int) {
	b.width = w
	b.height = h
	b.list.SetSize(w-b.widthMargin, h-b.heightMargin)
	b.list.Styles.PaginationStyle = b.style.ListPaginator.Copy().Width(w - b.widthMargin)
}
