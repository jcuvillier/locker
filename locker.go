package locker

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// ErrAlreadyLocked is error when a lock is already acquired
var ErrAlreadyLocked = errors.New("resource already locked")

var ErrMaxAttemptReached = errors.New("maximum number of attempt reached")

// AcquireFunc is the function used to acquire a new lock
// When the resource is already locked, this function should return a ErrAlreadyLocked error
type AcquireFunc func(ctx context.Context, key interface{}) error

// ReleaseFunc is the function used to release a lock
type ReleaseFunc func(ctx context.Context, key interface{}) error

// New creates a new Locker with given AcquireFunc and ReleaseFunc
func New(acquire AcquireFunc, release ReleaseFunc, options ...Option) *Locker {
	locker := Locker{
		acquire:    acquire,
		release:    release,
		delay:      defaultDelay,
		maxAttempt: 6,
	}
	for _, opt := range options {
		opt(&locker)
	}
	return &locker
}

// Locker structure defines how the lock is acquired and released
type Locker struct {
	acquire    AcquireFunc
	release    ReleaseFunc
	delay      Delay
	maxAttempt int
}

// Option defines options to be applied when acquiring lock
type Option func(*Locker)

// WithBackoff allows to specify its own Backoff implementation
func WithDelay(d Delay) Option {
	return func(l *Locker) {
		l.delay = d
	}
}

// WithAttempts defines the max number of attempts to acquire the lock
func WithAttempts(attempts int) Option {
	return func(l *Locker) {
		l.maxAttempt = attempts
	}
}

// Acquire call the acquire function to create a new lock
// Options are then applied to the returned lock
//
// The lock is acquired with a backoff retry defined by its backoff algorithm
// Backoff algorithm can be specified by using the WithBackoff() option when instantiating locker
// Default backoff algorithm is an exponential backoff
func (l *Locker) Acquire(ctx context.Context, key interface{}, options ...LockOption) (*Lock, error) {
	// Acquire lock
	acquired := false
	for i := 0; i < l.maxAttempt; i++ {
		err := l.acquire(ctx, key)
		if err != nil {
			if errors.Is(err, ErrAlreadyLocked) {
				dur := l.delay.Next()
				time.Sleep(dur)
				continue
			}
			return nil, errors.Wrapf(err, "cannot acquire lock for key %v", key)
		}
		acquired = true
		break
	}

	if !acquired {
		return nil, ErrMaxAttemptReached
	}

	// Create lock and apply options
	lock := Lock{
		key:     key,
		release: l.release,
	}
	for _, opt := range options {
		opt(&lock)
	}

	return &lock, nil
}

// Lock structure defines the lock itself and how to release it
type Lock struct {
	key     interface{}
	release ReleaseFunc
}

// Release releases the lock
func (l *Lock) Release(ctx context.Context) error {
	return l.release(ctx, l.key)
}

// LockOption defines options to be applied on lock when acquired
type LockOption func(*Lock)
