package retry

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestFailingRetry(t *testing.T) {
	n := 0
	// Always returning true to try again, should
	// eventually reach the max retries.
	X(4, 1*time.Millisecond, func() bool {
		n++
		return true
	})
	assert.Equal(t, 4, n)
}

func TestSuccessfulRetry(t *testing.T) {
	n := 0
	// Returing false from the function should
	// terminate the retries.
	X(4, 1*time.Millisecond, func() bool {
		n++
		if n == 2 {
			return false
		}
		return true
	})
	assert.Equal(t, 2, n)
}

func TestBackoff(t *testing.T) {
	const max = 8 * time.Second

	// A value of i less than 1 should be set to 1.
	// Large values of i should never return a
	// duration larger than max.
	for i := -10; i < 1000; i++ {
		assert.T(t, max >= Backoff(i, max))
	}
}

func TestTailBackoff(t *testing.T) {
	const max = 8 * time.Second

	// Test that beyond the third try,
	// the max duration is returned.
	third := Backoff(3, max)
	for i := 4; i < 1000; i++ {
		assert.Equal(t, third, Backoff(i, max))
	}
}
