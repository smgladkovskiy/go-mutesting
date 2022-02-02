package models_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetMutators(t *testing.T) {
	type inStruct struct {
		ml              models.MutatorsInterface
		disableMutators []models.MutatorName
	}

	type tc struct {
		name string
		in   inStruct
		exp  []models.MutatorItem
	}

	tcs := []tc{
		{
			name: "filled mutators, 0 disabled",
			in: inStruct{
				ml:              &models.MutatorsList{{"mutator1", mutatorMock}, {"mutator2", mutatorMock}},
				disableMutators: nil,
			},
			exp: []models.MutatorItem{{"mutator1", mutatorMock}, {"mutator2", mutatorMock}},
		},
		{
			name: "filled mutators, 1 disabled",
			in: inStruct{
				ml:              &models.MutatorsList{{"mutator1", mutatorMock}, {"mutator2", mutatorMock}},
				disableMutators: []models.MutatorName{"mutator1"},
			},
			exp: []models.MutatorItem{{"mutator2", mutatorMock}},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := models.GetMutators(tc.in.ml, tc.in.disableMutators)

			assert.Equal(t, len(tc.exp), len(out))

			for _, mi := range tc.exp {
				m, err := tc.in.ml.GetByName(mi.Name)
				if !assert.NoError(t, err) {
					t.FailNow()
				}

				assert.Equal(t, mi.Mutator(nil, nil, nil), m(nil, nil, nil))
			}
		})
	}
}
