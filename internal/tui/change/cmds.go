package change

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BackMsg struct{}

func back() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}
