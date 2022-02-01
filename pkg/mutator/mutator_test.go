package mutator_test

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator"
	"github.com/stretchr/testify/assert"
)

func mockMutator(_ *types.Package, _ *types.Info, _ ast.Node) []mutator.Mutation {
	// Do nothing
	return nil
}

func TestMockMutator(t *testing.T) {
	// t.Parallel()
	// Mock is not registered
	for _, name := range mutator.List() {
		if name == "mock" {
			assert.Fail(t, "mock should not be in the mutator list yet")
		}
	}

	m, err := mutator.New("mock")
	assert.Nil(t, m)
	assert.NotNil(t, err)

	// Register mock
	mutator.Register("mock", mockMutator)

	// Mock is registered
	found := false

	for _, name := range mutator.List() {
		if name == "mock" {
			found = true

			break
		}
	}

	assert.True(t, found)

	m, err = mutator.New("mock")
	assert.NotNil(t, m)
	assert.Nil(t, err)

	// Register mock a second time
	caught := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		mutator.Register("mock", mockMutator)
	}()

	assert.True(t, caught)

	// Register nil function
	caught = false

	func() {
		defer func() {
			if r := recover(); r != nil {
				caught = true
			}
		}()

		mutator.Register("mockachino", nil)
	}()
	assert.True(t, caught)
}
