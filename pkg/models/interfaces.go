package models

type MutatorsInterface interface {
	Register(name MutatorName, mutator Mutator) error
	GetByName(name MutatorName) (Mutator, error)
	Names() []MutatorName
}
