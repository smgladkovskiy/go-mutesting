package models_test

import (
	"testing"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMutationStats_Total(t *testing.T) {
	type tc struct {
		name string
		in   models.MutationStats
		exp  int
	}

	tcs := []tc{
		{
			name: "11 escaped + 9 killed",
			in: models.MutationStats{
				MutantsEscaped: 11,
				MutantsKilled:  9,
			},
			exp: 20,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.exp, tc.in.Total())
		})
	}
}

func TestMutationStats_Score(t *testing.T) {
	type tc struct {
		name string
		in   models.MutationStats
		exp  float64
	}

	tcs := []tc{
		{
			name: "Total 20, escaped: 9, killed: 11",
			in: models.MutationStats{
				MutantsEscaped: 11,
				MutantsKilled:  9,
				MutantsSkipped: 8,
			},
			exp: 0.450000,
		},
		{
			name: "Total 21, escaped 10, killed 11",
			in: models.MutationStats{
				MutantsEscaped: 11,
				MutantsKilled:  10,
				MutantsSkipped: 8,
			},
			exp: 0.47619047619047616,
		},
		{
			name: "Total 0",
			in:   models.MutationStats{},
			exp:  0,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.exp, tc.in.Score())
		})
	}
}
