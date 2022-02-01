package models

import (
	"strings"

	"github.com/smgladkovskiy/go-mutesting/pkg/mutator"
	log "github.com/spacetab-io/logs-go/v2"
)

type MutatorItem struct {
	Name    string
	Mutator mutator.Mutator
}

func GetMutators(disableMutators []string) []MutatorItem {
	var mutators []MutatorItem

MUTATOR:
	for _, name := range mutator.List() {
		if len(disableMutators) > 0 {
			for _, d := range disableMutators {
				pattern := strings.HasSuffix(d, "*")

				if (pattern && strings.HasPrefix(name, d[:len(d)-2])) || (!pattern && name == d) {
					continue MUTATOR
				}
			}
		}

		log.Debug().Msgf("mutator %s enabled", name)

		m, _ := mutator.New(name)
		mutators = append(mutators, MutatorItem{
			Name:    name,
			Mutator: m,
		})
	}

	return mutators
}
