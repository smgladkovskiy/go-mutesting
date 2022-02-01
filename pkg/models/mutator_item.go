package models

import (
	"strings"

	log "github.com/spacetab-io/logs-go/v2"
)

type MutatorItem struct {
	Name    MutatorName
	Mutator Mutator
}

func GetMutators(ml MutatorLookup, disableMutators []MutatorName) []MutatorItem {
	var mutators []MutatorItem

MUTATOR:
	for _, name := range ml.List() {
		if len(disableMutators) > 0 {
			for _, d := range disableMutators {
				pattern := strings.HasSuffix(d.String(), "*")

				if (pattern && strings.HasPrefix(name.String(), d[:len(d)-2].String())) || (!pattern && name == d) {
					continue MUTATOR
				}
			}
		}

		log.Debug().Msgf("mutator %s enabled", name)

		m, _ := ml.GetByName(name)
		mutators = append(mutators, MutatorItem{
			Name:    name,
			Mutator: m,
		})
	}

	return mutators
}
