package internal

import (
	"regexp"
	"strings"

	"github.com/go-errors/errors"
)

const check = `
  + create
  ~ update in-place
  - destroy
-/+ destroy and then create replacement
`

var changesRegexp = regexp.MustCompile(`(?m)(?s)^  #(.*?^    })`)
var changeRegexp = regexp.MustCompile(`(?m)(?s)\A  # ([\S]*?([\S]+)\.([\S]+)).*?((-/\+|[+~-]) resource.*)`)
var summaryRegexp = regexp.MustCompile(`(?m)^Plan: ([\d]+) to add, ([\d+]) to change, ([\d+]) to destroy.$`)

func Parse(s string) (
	outsideChanges ResourceChanges,
	actions ResourceChanges,
	summary string,
	err error,
) {
	if !strings.Contains(s, strings.TrimSpace(check)) {
		return nil, nil, "", errors.Errorf("failed to find check")
	}

	so, sa, ok := strings.Cut(s, "Terraform will perform the following actions")
	if !ok {
		return nil, nil, "", errors.Errorf("failed to find actions")
	}

	o, err := parseChanges(so)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, 0)
	}

	a, err := parseChanges(sa)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, 0)
	}

	return o, a, summaryRegexp.FindString(s), nil
}

func parseChanges(s string) (ResourceChanges, error) {
	rc := ResourceChanges{}

	cs := changesRegexp.FindAllString(s, -1)

	for _, c := range cs {
		m := changeRegexp.FindStringSubmatch(c)
		if len(m) != 6 {
			return nil, errors.Errorf("not enough capturing groups")
		}
		as, err := ActionsFromIdentifier(m[5])
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		rc = append(rc, ResourceChange{
			Address: m[1],
			Type:    m[2],
			Name:    m[3],
			Actions: as,
			Diff:    m[4],
		})
	}

	return rc, nil
}
