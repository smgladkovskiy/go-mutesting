package infection

import (
	"go/ast"
	"strings"

	log "github.com/spacetab-io/logs-go/v2"
)

// Results traverses the AST of the given node and prints every node to STDOUT.
func Results(node ast.Node) {
	w := &result{
		level: 0,
	}

	ast.Walk(w, node)
}

type result struct {
	level int
}

// Visit implements the Visit method of the ast.Visitor interface.
func (w *result) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		w.level++

		log.Info().Msgf("%s(%p)%#v\n", strings.Repeat("\t", w.level), node, node)
	} else {
		w.level--
	}

	return w
}
