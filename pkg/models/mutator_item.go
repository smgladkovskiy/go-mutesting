package models

import (
	"strings"

	log "github.com/spacetab-io/logs-go/v2"
)

type MutatorItem struct {
	Name    MutatorName
	Mutator Mutator
}

func GetMutators(ml MutatorsInterface, disabledMutators []MutatorName) []MutatorItem {
	mutators := make([]MutatorItem, 0)

mutatorsLoop:
	for _, name := range ml.Names() {
		if len(disabledMutators) > 0 {
			for _, d := range disabledMutators {
				pattern := strings.HasSuffix(d.String(), "*")

				if (pattern && strings.HasPrefix(name.String(), d[:len(d)-2].String())) || (!pattern && name == d) {
					continue mutatorsLoop
				}
			}
		}

		log.Debug().Msgf("mutator %s enabled", name)

		m, _ := ml.GetByName(name)
		mutators = append(mutators, MutatorItem{Name: name, Mutator: m})
	}

	return mutators
}
