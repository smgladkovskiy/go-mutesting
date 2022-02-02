package models_test

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/stretchr/testify/assert"
)

// MutatorCase implements a mutator mock.
func mutatorMock(_ *types.Package, _ *types.Info, _ ast.Node) []models.Mutation {
	return make([]models.Mutation, 0)
}

func TestMutatorsMap_GetByName(t *testing.T) {
	type inStruct struct {
		ml   models.MutatorsMap
		name models.MutatorName
	}

	type expStruct struct {
		mutator models.Mutator
		err     error
	}

	type tc struct {
		name string
		in   inStruct
		exp  expStruct
	}

	tcs := []tc{
		{
			name: "existing mutator",
			in: inStruct{
				ml:   models.MutatorsMap{"mutatorName": mutatorMock},
				name: models.MutatorName("mutatorName"),
			},
			exp: expStruct{
				mutator: mutatorMock,
				err:     nil,
			},
		},
		{
			name: "not existing mutator",
			in: inStruct{
				ml:   models.MutatorsMap{"mutatorName": mutatorMock},
				name: models.MutatorName("unknownMutatorName"),
			},
			exp: expStruct{
				mutator: nil,
				err:     errs.ErrUnknownMutator,
			},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := tc.in.ml.GetByName(tc.in.name)
			if tc.exp.err != nil {
				if !assert.Error(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			if tc.exp.mutator == nil {
				assert.Nil(t, out)
			} else {
				assert.Equal(
					t,
					tc.exp.mutator(nil, nil, nil),
					out(nil, nil, nil),
				)
			}
		})
	}
}

func TestMutatorsMap_Register(t *testing.T) {
	type inStruct struct {
		ml      models.MutatorsMap
		name    models.MutatorName
		mutator models.Mutator
	}

	type expStruct struct {
		ml  models.MutatorsMap
		err error
	}

	type tc struct {
		name string
		in   inStruct
		exp  expStruct
	}

	tcs := []tc{
		{
			name: "register mutator",
			in: inStruct{
				ml:      models.MutatorsMap{},
				name:    "mutator1",
				mutator: mutatorMock,
			},
			exp: expStruct{
				ml:  models.MutatorsMap{"mutator1": mutatorMock},
				err: nil,
			},
		},
		{
			name: "register existing mutator",
			in: inStruct{
				ml:      models.MutatorsMap{"mutator1": mutatorMock},
				name:    "mutator1",
				mutator: mutatorMock,
			},
			exp: expStruct{
				ml:  models.MutatorsMap{"mutator1": mutatorMock},
				err: errs.ErrMutatorRegistered,
			},
		},
		{
			name: "nil mutator",
			in: inStruct{
				ml:      models.MutatorsMap{},
				name:    "mutator1",
				mutator: nil,
			},
			exp: expStruct{
				ml:  models.MutatorsMap{},
				err: errs.ErrNilMutatorFunction,
			},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.in.ml.Register(tc.in.name, tc.in.mutator)
			if tc.exp.err != nil {
				if !assert.Error(t, err) {
					t.FailNow()
				}
			} else {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, len(tc.exp.ml), len(tc.in.ml))

			for k1, m1 := range tc.exp.ml {
				m2 := tc.in.ml[k1]

				assert.Equal(t, m2(nil, nil, nil), m1(nil, nil, nil))
			}
		})
	}
}

func TestMutatorsMap_Names(t *testing.T) {
	type tc struct {
		name string
		in   models.MutatorsMap
		exp  []models.MutatorName
	}

	tcs := []tc{
		{
			name: "filled map",
			in:   models.MutatorsMap{"mutator1": mutatorMock, "mutator2": mutatorMock},
			exp:  []models.MutatorName{"mutator1", "mutator2"},
		},
		{
			name: "empty map",
			in:   models.MutatorsMap{},
			exp:  []models.MutatorName{},
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.exp, tc.in.Names())
		})
	}
}
