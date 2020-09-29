package locker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedBackoff(t *testing.T) {
	d := 10 * time.Millisecond
	b := FixedBackoff{
		Duration: d,
	}
	assert.Equal(t, d, b.Next())
}
