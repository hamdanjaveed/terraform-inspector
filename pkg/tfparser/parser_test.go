package tfparser_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/hamdanjaveed/terraform-inspector/pkg/tfparser"
)

func TestParse(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "..", "files", "testdata", "oneofeach.tfplan"))
	assert.NoError(t, err)

	p, err := tfparser.Parse(f)
	assert.NoError(t, err)

	fmt.Printf("%+v\n", p)
}
