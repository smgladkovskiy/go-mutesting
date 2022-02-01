package statement

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/astutil"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

// MutatorRemoveStatement implements a mutator to remove statements.
func MutatorRemoveStatement(pkg *types.Package, info *types.Info, node ast.Node) []models.Mutation {
	var l []ast.Stmt

	switch n := node.(type) {
	case *ast.BlockStmt:
		l = n.List
	case *ast.CaseClause:
		l = n.Body
	}

	var mutations []models.Mutation

	for i, ni := range l {
		if checkRemoveStatement(ni) {
			li := i
			old := l[li]

			mutations = append(mutations, models.Mutation{
				Change: func() {
					l[li] = astutil.CreateNoopOfStatement(pkg, info, old)
				},
				Reset: func() {
					l[li] = old
				},
			})
		}
	}

	return mutations
}

func checkRemoveStatement(node ast.Stmt) bool {
	switch n := node.(type) {
	case *ast.AssignStmt:
		if n.Tok != token.DEFINE {
			return true
		}
	case *ast.ExprStmt, *ast.IncDecStmt:
		return true
	}

	return false
}
