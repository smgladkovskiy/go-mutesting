package mutator

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
)

// Mutator defines a mutator for mutation testing by returning a list of possible mutations for the given node.
type Mutator func(pkg *types.Package, info *types.Info, node ast.Node) []Mutation

// todo remove globals.
var mutatorLookup = make(map[string]Mutator)

// New returns a new mutator instance given the registered name of the mutator.
// The error return argument is not nil, if the name does not exist in the registered mutator list.
func New(name string) (Mutator, error) {
	mutator, ok := mutatorLookup[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", errs.ErrUnknownMutator, name)
	}

	return mutator, nil
}

// List returns a list of all registered mutator names.
func List() []string {
	keyMutatorLookup := make([]string, 0, len(mutatorLookup))

	for key := range mutatorLookup {
		keyMutatorLookup = append(keyMutatorLookup, key)
	}

	sort.Strings(keyMutatorLookup)

	return keyMutatorLookup
}

// Register registers a mutator instance function with the given name.
func Register(name string, mutator Mutator) {
	if mutator == nil {
		panic("mutator function is nil")
	}

	if _, ok := mutatorLookup[name]; ok {
		panic(fmt.Sprintf("mutator %q already registered", name))
	}

	mutatorLookup[name] = mutator
}
