package internal

import "github.com/charmbracelet/lipgloss"

type styles struct {
	app lipgloss.Style
}

func newStyles() *styles {
	return &styles{
		app: lipgloss.NewStyle().Margin(1, 2),
	}
}
