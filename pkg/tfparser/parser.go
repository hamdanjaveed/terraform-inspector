package tfparser

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var tfPlanLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: `OutsideChanges`, Pattern: `Terraform\s+detected\s+the\s+following\s+changes\s+made\s+outside\s+of\s+Terraform\s+since\s+the\s+last\s+"terraform apply":`},
	{Name: `ws`, Pattern: `[ ]+`},
	{Name: `Any`, Pattern: `[\s\S]*`},
})

var tfParser = participle.MustBuild(
	&Plan{},
	participle.Lexer(tfPlanLexer),
)

type Plan struct {
	Before         string `@Any`
	OutsideChanges string `@OutsideChanges`
	After          string `@Any`
	// Resources []*Resource
}

type Resource struct {
}

func Parse(r io.Reader) (*Plan, error) {
	p := &Plan{}
	if err := tfParser.Parse("", r, p); err != nil {
		return nil, err
	}
	return p, nil
}
