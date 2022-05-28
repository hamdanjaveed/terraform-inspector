package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hamdanjaveed/terraform-inspector/internal/tf"
)

type SelectMsg struct {
	Change tf.ResourceChange
}

func selectChange(c tf.ResourceChange) tea.Cmd {
	return func() tea.Msg {
		return SelectMsg{Change: c}
	}
}
