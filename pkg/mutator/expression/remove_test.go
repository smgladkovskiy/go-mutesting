package expression_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/expression"
	"github.com/smgladkovskiy/go-mutesting/test"
)

func TestMutatorRemoveTerm(t *testing.T) {
	t.Parallel()
	test.Mutator(
		t,
		expression.MutatorRemoveTerm,
		"../../../test/testdata/expression/remove.go",
		6,
	)
}
