package internal

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/samber/lo"
)

type ResourceChanges []ResourceChange

func (rc ResourceChanges) ByAction() (
	create ResourceChanges,
	update ResourceChanges,
	delete ResourceChanges,
	replace ResourceChanges,
) {
	create = lo.Filter(rc, func(rc ResourceChange, _ int) bool { return rc.Actions.IsCreate() })
	update = lo.Filter(rc, func(rc ResourceChange, _ int) bool { return rc.Actions.IsUpdate() })
	delete = lo.Filter(rc, func(rc ResourceChange, _ int) bool { return rc.Actions.IsDelete() })
	replace = lo.Filter(rc, func(rc ResourceChange, _ int) bool { return rc.Actions.IsReplace() })
	return
}

type ResourceChange struct {
	Address string
	Type    string
	Name    string
	Actions Actions
	Diff    string
}

type Action string

const (
	CreateAction Action = "create"
	DeleteAction Action = "delete"
	NoopAction   Action = "no-op"
	ReadAction   Action = "read"
	UpdateAction Action = "update"
)

func ActionsFromIdentifier(i string) (Actions, error) {
	switch i {
	case "+":
		return Actions{CreateAction}, nil
	case "~":
		return Actions{UpdateAction}, nil
	case "-":
		return Actions{DeleteAction}, nil
	case "-/+":
		return Actions{DeleteAction, CreateAction}, nil
	default:
		return nil, errors.Errorf("unrecognized action %s", i)
	}
}

type Actions []Action

func (as Actions) String() string {
	if as.IsCreate() {
		return "+"
	} else if as.IsUpdate() {
		return "~"
	} else if as.IsDelete() {
		return "-"
	} else if as.IsReplace() {
		return "-/+"
	} else {
		s := ""
		for _, a := range as {
			s += fmt.Sprintf("%v,", a)
		}
		return fmt.Sprintf("UNRECOGNIZED %s", s)
	}
}

func (as Actions) IsCreate() bool {
	if len(as) != 1 {
		return false
	}
	return as[0] == CreateAction
}

func (as Actions) IsDelete() bool {
	if len(as) != 1 {
		return false
	}
	return as[0] == DeleteAction
}

func (as Actions) IsDestructive() bool {
	return lo.Contains(as, DeleteAction)
}

func (as Actions) IsReplace() bool {
	if len(as) != 2 {
		return false
	}
	return (as[0] == DeleteAction && as[1] == CreateAction) || (as[0] == CreateAction && as[1] == DeleteAction)
}

func (as Actions) IsUpdate() bool {
	if len(as) != 1 {
		return false
	}
	return as[0] == UpdateAction
}
