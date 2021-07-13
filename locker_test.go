package locker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDelay struct{}

func (d *testDelay) Next() time.Duration {
	return 0
}

func TestNewLocker(t *testing.T) {
	ctx := context.Background()
	var acquired, released bool
	acquire := func(ctx context.Context, key interface{}) error {
		acquired = true
		return nil
	}
	release := func(ctx context.Context, key interface{}) error {
		released = true
		return nil
	}
	// Create new locker
	locker := New(acquire, release, WithDelay(&testDelay{}))
	assert.IsType(t, &testDelay{}, locker.delay)
	assert.False(t, acquired)
	assert.False(t, released)

	// Acquiring lock
	locker.acquire(ctx, "")
	assert.True(t, acquired)
	assert.False(t, released)

	// Releasing lock
	locker.release(ctx, "")
	assert.True(t, released)
}

func TestAcquire(t *testing.T) {
	release := func(ctx context.Context, key interface{}) error {
		return nil
	}
	// Regular test case
	t.Run("regular", func(t *testing.T) {
		ctx := context.Background()
		var acquired bool
		acquire := func(ctx context.Context, key interface{}) error {
			acquired = true
			return nil
		}
		var optCalled bool
		lockOpt := func(*Lock) {
			optCalled = true
		}
		locker := New(acquire, release, WithDelay(&testDelay{}))
		lock, err := locker.Acquire(ctx, "key", lockOpt)
		require.NoError(t, err)
		assert.True(t, acquired)
		assert.Equal(t, "key", lock.key)
		assert.True(t, optCalled)
	})

	t.Run("already_locked", func(t *testing.T) {
		ctx := context.Background()
		var acquired, alreadylocked bool
		i := 0
		// This acquire function will return ErrAlreadyLocked the first time it's called and set acquired to true on next calls
		acquire := func(ctx context.Context, key interface{}) error {
			if i == 0 {
				i = i + 1
				alreadylocked = true
				return ErrAlreadyLocked
			}
			acquired = true
			return nil
		}
		locker := New(acquire, release, WithDelay(&testDelay{}))
		_, err := locker.Acquire(ctx, "")
		require.NoError(t, err)
		assert.True(t, acquired)
		assert.True(t, alreadylocked)
	})

	t.Run("acquire_error", func(t *testing.T) {
		ctx := context.Background()
		acquire := func(ctx context.Context, key interface{}) error {
			return fmt.Errorf("error")
		}
		locker := New(acquire, release)
		_, err := locker.Acquire(ctx, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot acquire lock for key")
	})
}

func TestRelease(t *testing.T) {
	ctx := context.Background()
	var acquired, released bool
	acquire := func(ctx context.Context, key interface{}) error {
		acquired = true
		return nil
	}
	release := func(ctx context.Context, key interface{}) error {
		released = true
		return nil
	}
	locker := New(acquire, release)
	lock, err := locker.Acquire(ctx, "")
	require.NoError(t, err)
	require.True(t, acquired)
	err = lock.Release(ctx)
	require.NoError(t, err)
	assert.True(t, released)
}

func TestAttempts(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 5
	attempts := 0
	acquire := func(ctx context.Context, key interface{}) error {
		attempts++
		return ErrAlreadyLocked
	}
	release := func(ctx context.Context, key interface{}) error {
		return nil
	}
	locker := New(acquire, release, WithAttempts(maxAttempts))
	_, err := locker.Acquire(ctx, "")
	assert.Equal(t, ErrMaxAttemptReached, err)
	assert.Equal(t, maxAttempts, attempts)
}
