package parser

import (
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hamdanjaveed/terraform-inspector/internal/tf"
)

var (
	noChanges = regexp.MustCompile(`No\s+changes\.\s+Your\s+infrastructure\s+matches\s+the\s+configuration\.`)

	// checkRegexp is used to checkRegexp two things:
	//   - Are we parsing a terraform plan?
	//   - Are the symbols for create/update/delete/replace what we expect them to be?
	checkRegexp = regexp.MustCompile(
		// Match text exactly.
		`Resource\s+actions\s+are\s+indicated\s+with\s+the\s+following\s+symbols:` +
			// Optionally match each of the following lines in order.
			// We could have any combination of them depending on which changes are present in the plan.
			`((\n  \+ create)?` +
			`(\n  ~ update in-place)?` +
			`(\n  - destroy)?` +
			`(\n-/\+ destroy and then create replacement)?)`)

	// ansiRegexp is used to strip the terraform plan of any ANSI color codes. These might be present
	// if the terraform plan was piped in from the terminal.
	ansiRegexp = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

	// changesRegexp captures all the resource changes from the terraform plan.
	changesRegexp = regexp.MustCompile(
		// multi line: ^ and $ match the beginning and end of each line.
		"(?m)" +
			// single line: Dot matches newlines.
			"(?s)" +
			// Match the # that denotes the beginning of a resource change. This # is then followed by the address.
			"^  #" +
			// Match everything until we see the final closing brace which we identify by looking for exactly 4 spaces.
			"(.*?^    })")

	// changeRegexp captures all the details about a single resource change.
	changeRegexp = regexp.MustCompile(
		// multi line: ^ and $ match the beginning and end of each line.
		`(?m)` +
			// single line: Dot matches newlines.
			`(?s)` +
			// From the beginning of the string, look for the #.
			`\A  # ` +
			// Capture the resource address type and name.
			`(?P<address>(\S+\.)?(?P<type>[\S]+)\.(?P<name>[\S]+)) ` +
			// Ignore everything up-till we get to the diff.
			`.*?` +
			// Capture the diff. Look for the resource change identifier.
			`(?P<diff>(?P<identifier>-/\+|[+~-])` +
			// Then capture all the characters remaining in the string.
			` resource.*)`)

	// summaryRegexp matches the summary of the terraform plan.
	summaryRegexp = regexp.MustCompile(`(?m)^Plan: ([\d]+) to add, ([\d]+) to change, ([\d]+) to destroy.$`)
)

func Parse(s string) (changes tf.ResourceChanges, summary string, err error) {
	s = ansiRegexp.ReplaceAllString(s, "")

	if err := validate(s); err != nil {
		return nil, "", err
	}

	so, sa, ok := strings.Cut(s, "Terraform will perform the following actions")
	if !ok {
		return nil, "", errors.Errorf("failed to find actions")
	}

	o, err := parseChanges(so, tf.OutsideChange)
	if err != nil {
		return nil, "", errors.Wrap(err, 0)
	}

	a, err := parseChanges(sa, tf.ActionChange)
	if err != nil {
		return nil, "", errors.Wrap(err, 0)
	}

	return append(o, a...), summaryRegexp.FindString(s), nil
}

func validate(s string) error {
	if !checkRegexp.MatchString(s) && !noChanges.MatchString(s) {
		return errors.Errorf("invalid terraform plan")
	}
	return nil
}

func parseChanges(s string, t tf.ChangeType) (tf.ResourceChanges, error) {
	rc := tf.ResourceChanges{}

	cs := changesRegexp.FindAllString(s, -1)

	for _, c := range cs {
		m := changeRegexp.FindStringSubmatch(c)
		if len(m) != 7 {
			return nil, errors.Errorf("expected 7 capturing groups, got %d", len(m))
		}
		g := make(map[string]string)
		for i, n := range changeRegexp.SubexpNames() {
			if i != 0 && n != "" {
				g[n] = m[i]
			}
		}

		as, err := tf.ActionsFromIdentifier(g["identifier"])
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		rc = append(rc, tf.ResourceChange{
			Type:         t,
			Address:      g["address"],
			ResourceType: g["type"],
			Name:         g["name"],
			Actions:      as,
			Diff:         g["diff"],
		})
	}

	return rc, nil
}
