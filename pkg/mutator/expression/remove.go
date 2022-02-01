package expression

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator"
)

func InitRemove() {
	mutator.Register("expression/remove", MutatorRemoveTerm)
}

// MutatorRemoveTerm implements a mutator to remove expression terms.
func MutatorRemoveTerm(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	if n.Op != token.LAND && n.Op != token.LOR {
		return nil
	}

	var r *ast.Ident

	// nolint: exhaustive
	switch n.Op {
	case token.LAND:
		r = ast.NewIdent("true")
	case token.LOR:
		r = ast.NewIdent("false")
	}

	x := n.X
	y := n.Y

	return []mutator.Mutation{
		{
			Change: func() {
				n.X = r
			},
			Reset: func() {
				n.X = x
			},
		},
		{
			Change: func() {
				n.Y = r
			},
			Reset: func() {
				n.Y = y
			},
		},
	}
}