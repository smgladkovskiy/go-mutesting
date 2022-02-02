package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	t.Parallel()

	assert.Equal(t, foo(), 16)
}
