package ratelimit_test

import (
	"testing"
	"time"

	"github.com/thomasbeukema/ratelimit"
)

func TestLimiter(t *testing.T) {
	lmt := ratelimit.NewLimiter(60, time.Minute)

	if !lmt.Allowed(10) {
		t.Error("Not Allowed, should be Allowed")
	}

	if lmt.Allowed(51) {
		t.Error("Shouldn't be allowed, except it is")
	}

	time.Sleep(time.Minute)

	if !lmt.Allowed(10) {
		t.Error("Not Allowed, should be Allowed")
	}
}

func TestLooping(t *testing.T) {
	// Create a new Limiter with a capacity of 60 for an interval of 1 minute
	lmt := ratelimit.NewLimiter(60, time.Minute)

	allowed := 0
	notAllowed := 0

	// Loop 69 times
	for i := 0; i <= 69; i++ {
		// lmt.Allowed subs the provided cost from the bucket
		// This operation has a cost of 1
		if lmt.Allowed(1) {
			allowed++
		} else {
			notAllowed++
		}
	}

	if allowed != 60 {
		t.Errorf("Not enough allowed: %v", allowed)
	}

	if notAllowed != 10 {
		t.Errorf("Not enough not allowed: %v", notAllowed)
	}
}
