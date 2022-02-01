package branch

import (
	"go/ast"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/astutil"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

// MutatorIf implements a mutator for if and else if branches.
func MutatorIf(pkg *types.Package, info *types.Info, node ast.Node) []models.Mutation {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil
	}

	old := n.Body.List

	return []models.Mutation{
		{
			Change: func() {
				n.Body.List = []ast.Stmt{
					astutil.CreateNoopOfStatement(pkg, info, n.Body),
				}
			},
			Reset: func() {
				n.Body.List = old
			},
		},
	}
}
