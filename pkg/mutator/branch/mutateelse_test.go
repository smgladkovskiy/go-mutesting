package branch_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/branch"
	"github.com/smgladkovskiy/go-mutesting/test"
)

func TestMutatorElse(t *testing.T) {
	t.Parallel()
	test.Mutator(
		t,
		branch.MutatorElse,
		"../../../test/testdata/branch/mutateelse.go",
		1,
	)
}
