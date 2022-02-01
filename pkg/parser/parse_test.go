package parser_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseAndTypeCheckFileTypeCheckWholePackage(t *testing.T) {
	// t.Parallel()
	_, _, _, _, err := parser.ParseAndTypeCheckFile("astutil/create.go") // nolint: dogsled
	assert.Nil(t, err)
}
