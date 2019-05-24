package ratelimit

import "time"

// Limiter represents a rate limiter
type Limiter struct {
	bucketSize  int
	bucketValue int
	interval    time.Duration
	resetTimer  time.Time
}

// NewLimiter returns limiter with max of @bucketSize tokens, reseting every @interval
func NewLimiter(bucketSize int, interval time.Duration) *Limiter {
	lmt := new(Limiter)
	lmt.bucketSize = bucketSize
	lmt.bucketValue = bucketSize
	lmt.interval = interval
	lmt.resetTimer = time.Now().Add(interval)

	return lmt
}

// Allowed checks if rate limit allows operation
func (limiter *Limiter) Allowed(cost int) bool {
	if time.Now().After(limiter.resetTimer) { // check if time.Now is later than resetTimer
		limiter.resetTimer = time.Now().Add(limiter.interval)
		limiter.bucketValue = limiter.bucketSize // reset the bucketValue to max
		return true
	}

	if limiter.bucketValue-cost >= 0 {
		limiter.bucketValue -= cost
		return true
	}

	return false
}
