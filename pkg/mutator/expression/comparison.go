package expression

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

var comparisonMutations = map[token.Token]token.Token{
	token.LSS: token.LEQ,
	token.LEQ: token.LSS,
	token.GTR: token.GEQ,
	token.GEQ: token.GTR,
}

// MutatorComparison implements a mutator to change comparisons.
func MutatorComparison(_ *types.Package, _ *types.Info, node ast.Node) []models.Mutation {
	n, ok := node.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	o := n.Op

	r, ok := comparisonMutations[n.Op]
	if !ok {
		return nil
	}

	return []models.Mutation{
		{
			Change: func() {
				n.Op = r
			},
			Reset: func() {
				n.Op = o
			},
		},
	}
}
