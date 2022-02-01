package infection

import (
	"go/ast"
	"go/types"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

// Launch mutates the given node with the given mutator returning a channel
// to control the mutation steps.
// It traverses the AST of the given node and calls the method Check of the given
// mutator to verify that a node can be mutated by the mutator. If a node can be
// mutated the method Mutate of the given mutator is executed with the node and
// the control channel. After completion of the traversal the control channel is closed.
func Launch(pkg *types.Package, info *types.Info, node ast.Node, m models.Mutator) chan bool {
	w := &infection{
		changed: make(chan bool),
		mutator: m,
		pkg:     pkg,
		info:    info,
	}

	go func() {
		ast.Walk(w, node)

		close(w.changed)
	}()

	return w.changed
}

type infection struct {
	changed chan bool
	mutator models.Mutator
	pkg     *types.Package
	info    *types.Info
}

// Visit implements the Visit method of the ast.Visitor interface.
func (w *infection) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return w
	}

	for _, m := range w.mutator(w.pkg, w.info, node) {
		m.Change()
		w.changed <- true
		<-w.changed

		m.Reset()
		w.changed <- true
		<-w.changed
	}

	return w
}
