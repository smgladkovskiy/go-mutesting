package models

import (
	"fmt"
	"sort"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
)

type MutatorsList []MutatorItem

// Register registers a mutator instance function with the given name.
func (ml *MutatorsList) Register(name MutatorName, mutator Mutator) error {
	if mutator == nil {
		return errs.ErrNilMutatorFunction
	}

	if _, err := ml.GetByName(name); err == nil {
		return fmt.Errorf("%w %q", errs.ErrMutatorRegistered, name)
	}

	*ml = append(*ml, MutatorItem{Name: name, Mutator: mutator})

	return nil
}

// GetByName returns a new mutator instance given the registered name of the mutator.
// The error return argument is not nil, if the name does not exist in the registered mutator list.
func (ml MutatorsList) GetByName(name MutatorName) (Mutator, error) {
	for _, m := range ml {
		if m.Name == name {
			return m.Mutator, nil
		}
	}

	return nil, fmt.Errorf("%w: %q", errs.ErrUnknownMutator, name)
}

// Names returns a list of all registered mutator names.
func (ml MutatorsList) Names() []MutatorName {
	mutatorNames := make([]MutatorName, 0, len(ml))

	for _, m := range ml {
		mutatorNames = append(mutatorNames, m.Name)
	}

	sort.Slice(mutatorNames, func(i, j int) bool {
		return mutatorNames[i].String() < mutatorNames[j].String()
	})

	return mutatorNames
}
