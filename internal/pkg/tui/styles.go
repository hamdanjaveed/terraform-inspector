package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	App lipgloss.Style

	Header            lipgloss.Style
	Body              lipgloss.Style
	CreateChangeBody  lipgloss.Style
	DeleteChangeBody  lipgloss.Style
	ReplaceChangeBody lipgloss.Style
	UpdateChangeBody  lipgloss.Style

	ListPaginator lipgloss.Style
}

func NewStyles() *Styles {
	s := new(Styles)

	s.App = lipgloss.NewStyle().Margin(1, 2)

	s.Header = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(1, 0).
		BorderStyle(lipgloss.RoundedBorder())
	s.Body = lipgloss.NewStyle()
	s.CreateChangeBody = s.Body.Copy().
		Foreground(lipgloss.Color("2"))
	s.DeleteChangeBody = s.Body.Copy().
		Foreground(lipgloss.Color("1"))
	s.ReplaceChangeBody = s.DeleteChangeBody.Copy()
	s.UpdateChangeBody = s.Body.Copy().
		Foreground(lipgloss.Color("3"))

	s.ListPaginator = lipgloss.NewStyle().
		Margin(0)

	return s
}
