package expression_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/expression"
	"github.com/smgladkovskiy/go-mutesting/test"
)

func TestMutatorComparison(t *testing.T) {
	t.Parallel()
	test.Mutator(
		t,
		expression.MutatorComparison,
		"../../../test/testdata/expression/comparison.go",
		4,
	)
}
