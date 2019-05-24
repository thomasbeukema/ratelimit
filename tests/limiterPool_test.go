package ratelimit_test

import (
	"testing"
	"time"

	"github.com/thomasbeukema/ratelimit"
)

func cleanWhen(id string, lmt *ratelimit.Limiter) bool {
	/*
		You are free to implement you own logic in here to decide which Limiter to keep
		and which to remove. One example could be to check the 'exp' key from a JWT token,
		and to remove the limiter if the token has expired.
	*/

	// For the sake of simplicity; this will only remove 'userThree'
	if len(id) > 8 {
		return true // Return true if Limiter should be removed
	}

	return false // Return false if Limiter shouldn't be removed
}

func PoolTest(t *testing.T) {
	pool := ratelimit.NewPool()

	userOneLimiter := ratelimit.NewLimiter(10, time.Minute) // Create new Limiter instance
	pool.Register("userOne", userOneLimiter)                // Register Limiter in pool with ID 'userOne'

	userTwoLimiter := ratelimit.NewLimiter(20, time.Minute) // Create new Limiter instance
	pool.Register("userTwo", userTwoLimiter)                // Register Limiter in pool with ID 'userTwo'

	userThreeLimiter := ratelimit.NewLimiter(30, time.Minute) // Create new Limiter instance
	pool.Register("userThree", userThreeLimiter)              // Register Limiter in pool with ID 'userThree'

	_, err := pool.Find("userOne") // Retrieve Limiter by ID
	if err != nil {
		t.Errorf("Error at pool.Find: %v", err)
	}
	err = pool.Remove("userTwo") // Remove limiter from the pool
	if err != nil {
		t.Errorf("Error at pool.Remove: %v", err)
	}

	err = pool.Clean(cleanWhen) // Function which loops over every Limiter, decides to remove based on 'cleanWhen'
	if err != nil {
		t.Errorf("Error at pool.Clean: %v", err)
	}
}
