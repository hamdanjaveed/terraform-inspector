package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-errors/errors"
	"github.com/hamdanjaveed/terraform-inspector/internal/parser"
	"github.com/hamdanjaveed/terraform-inspector/internal/tui"
	"golang.org/x/term"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.(*errors.Error).ErrorStack())
	}
}

func run() error {
	r := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	cs, s, err := parser.Parse(string(b))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	p := tea.NewProgram(
		tui.NewBubble(cs, s, w, 0, h, 0),
		tea.WithAltScreen(),
		// tea.WithMouseCellMotion(),
	)
	if err := p.Start(); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
