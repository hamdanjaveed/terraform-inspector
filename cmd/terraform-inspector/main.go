package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-errors/errors"
	"github.com/hamdanjaveed/terraform-inspector/internal"
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

	os, as, s, err := internal.Parse(string(b))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	p := tea.NewProgram(
		internal.Bubble{
			OutsideChanges: os,
			Actions:        as,
			Summary:        s,
			ShowingDetail:  nil,
			Cursor:         -1,
			// Selected:       make(map[int]struct{}),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if err := p.Start(); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
