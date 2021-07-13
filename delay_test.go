package locker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedDelay(t *testing.T) {
	d := 10 * time.Millisecond
	b := FixedDelay{
		Duration: d,
	}
	assert.Equal(t, d, b.Next())
}
