package infection

import (
	"go/ast"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

// ResultsCount returns the number of corresponding mutations for a given mutator.
// It traverses the AST of the given node and calls the method Check of the given
// mutator for every node and sums up the returned counts. After completion of the
// traversal the final counter is returned.
func ResultsCount(pkg *types.Package, info *types.Info, node ast.Node, m models.Mutator) int {
	w := &resultsCount{
		count:   0,
		mutator: m,
		pkg:     pkg,
		info:    info,
	}

	ast.Walk(w, node)

	return w.count
}

type resultsCount struct {
	count   int
	mutator models.Mutator
	pkg     *types.Package
	info    *types.Info
}

// Visit implements the Visit method of the ast.Visitor interface.
func (w *resultsCount) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}

	w.count += len(w.mutator(w.pkg, w.info, node))

	return w
}
