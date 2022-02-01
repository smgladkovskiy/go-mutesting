package sub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaz(t *testing.T) {
	// t.Parallel()
	assert.Equal(t, baz(), 2)
}
