package branch_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/branch"
	"github.com/smgladkovskiy/go-mutesting/test"
)

func TestMutatorIf(t *testing.T) {
	t.Parallel()
	test.Mutator(
		t,
		branch.MutatorIf,
		"../../../test/testdata/branch/mutateif.go",
		2,
	)
}
