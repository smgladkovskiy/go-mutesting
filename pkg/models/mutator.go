package models

import (
	"go/ast"
	"go/types"
)

type MutatorName string

func (n MutatorName) String() string {
	return string(n)
}

// Mutator defines a mutator for mutation testing by returning a list of possible
// mutations for the given node.
type Mutator func(pkg *types.Package, info *types.Info, node ast.Node) []Mutation
