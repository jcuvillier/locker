package locker

import "time"

// Backoff interface defines the backoff algorithm used in retry mechanism
type Backoff interface {

	// Next returns the duration for the next attempt
	// Implemtation should also increase an attempt counter
	Next() time.Duration
}

var defaultBackoff = &FixedBackoff{
	Duration: 5 * time.Millisecond,
}

// FixedBackoff structure defines a backoff algorithm allways returning a fixed duration
type FixedBackoff struct {
	Duration time.Duration
}

// Next returns the duration for the next attempt
func (b *FixedBackoff) Next() time.Duration {
	return b.Duration
}
