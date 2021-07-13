package locker

import "time"

// Delay interface defines the delay algorithm used in retry mechanism
type Delay interface {

	// Next returns the duration for the next attempt
	// Implemtation should also increase an attempt counter
	Next() time.Duration
}

var defaultDelay = &FixedDelay{
	Duration: 5 * time.Millisecond,
}

// FixedDelay structure defines a delay algorithm always returning a fixed duration
type FixedDelay struct {
	Duration time.Duration
}

// Next returns the duration for the next attempt
func (b *FixedDelay) Next() time.Duration {
	return b.Duration
}
