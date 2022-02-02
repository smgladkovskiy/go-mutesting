package branch_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/branch"
	"github.com/smgladkovskiy/go-mutesting/test"
)

func TestMutatorCase(t *testing.T) {
	t.Parallel()

	test.Mutator(t, branch.MutatorCase, "../../../test/testdata/branch/mutatecase.go", 3)
}
