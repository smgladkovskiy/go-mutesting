package models

import (
	"fmt"
	"sort"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
)

type MutatorLookup map[MutatorName]Mutator

// Register registers a mutator instance function with the given name.
func (ml *MutatorLookup) Register(name MutatorName, mutator Mutator) error {
	if mutator == nil {
		return errs.ErrNilMutatorFunction
	}

	if _, err := ml.GetByName(name); err == nil {
		return fmt.Errorf("%w %q", errs.ErrMutatorRegistered, name)
	}

	mlNew := *ml

	mlNew[name] = mutator

	*ml = mlNew

	return nil
}

// GetByName returns a new mutator instance given the registered name of the mutator.
// The error return argument is not nil, if the name does not exist in the registered mutator list.
func (ml MutatorLookup) GetByName(name MutatorName) (Mutator, error) {
	mutator, ok := ml[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", errs.ErrUnknownMutator, name)
	}

	return mutator, nil
}

// List returns a list of all registered mutator names.
func (ml MutatorLookup) List() []MutatorName {
	keyMutatorLookup := make([]MutatorName, 0, len(ml))

	for key := range ml {
		keyMutatorLookup = append(keyMutatorLookup, key)
	}

	sort.Slice(keyMutatorLookup, func(i, j int) bool {
		return keyMutatorLookup[i].String() > keyMutatorLookup[j].String()
	})

	return keyMutatorLookup
}
