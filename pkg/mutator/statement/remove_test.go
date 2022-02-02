package statement_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/statement"
	"github.com/smgladkovskiy/go-mutesting/test"
)

func TestMutatorRemoveStatement(t *testing.T) {
	t.Parallel()

	test.Mutator(
		t,
		statement.MutatorRemoveStatement,
		"../../../testdata/statement/remove.go",
		17,
	)
}
