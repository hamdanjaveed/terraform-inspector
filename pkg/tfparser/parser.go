package tfparser

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var tfPlanLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: `OutsideChanges`, Pattern: `Terraform\s+detected\s+the\s+following\s+changes\s+made\s+outside\s+of\s+Terraform\s+since\s+the\s+last\s+"terraform apply":`},
	{Name: `ws`, Pattern: `[ ]`},
	{Name: `Resource`, Pattern: `resource`},
	{Name: `Any`, Pattern: `[\s\S]*`},
	{Name: `Change`, Pattern: `~|-\/\+|-|\+|<=`},
})

var tfParser = participle.MustBuild(
	&Plan{},
	participle.Lexer(tfPlanLexer),
)

type Plan struct {
	// OutsideChanges          string      `parser:"(@OutsideChanges"`
	// OutsideChangedResources []*Resource `parser:"@@+)?"`
	Resources []*Resource `parser:"@@*"`
}

type Resource struct {
	Change string `parser:"@Change"`
	Name   string `parser:"@Resource"`
	Diff   string `parser:""`
}

func Parse(r io.Reader) (*Plan, error) {
	p := &Plan{}
	if err := tfParser.Parse("", r, p); err != nil {
		return nil, err
	}
	return p, nil
}
