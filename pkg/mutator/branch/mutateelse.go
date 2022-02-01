package branch

import (
	"go/ast"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/astutil"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

// MutatorElse implements a mutator for else branches.
func MutatorElse(pkg *types.Package, info *types.Info, node ast.Node) []models.Mutation {
	n, ok := node.(*ast.IfStmt)
	if !ok {
		return nil
	}
	// We ignore else ifs and nil blocks
	_, ok = n.Else.(*ast.IfStmt)
	if ok || n.Else == nil {
		return nil
	}

	old := n.Else

	return []models.Mutation{
		{
			Change: func() {
				n.Else = astutil.CreateNoopOfStatement(pkg, info, old)
			},
			Reset: func() {
				n.Else = old
			},
		},
	}
}
