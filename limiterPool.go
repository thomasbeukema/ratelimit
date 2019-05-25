package ratelimit

import (
	"fmt"
	"sync"
)

// LimiterPool manages multiple pools
type LimiterPool struct {
	limiters map[string]*Limiter
	mut      *sync.Mutex
}

// NewPool returns new instance of a LimiterPool
func NewPool() *LimiterPool {
	pool := new(LimiterPool)                     // instantiate LimiterPool
	pool.limiters = make(map[string]*Limiter, 0) // instantiate map of limiters
	pool.mut = new(sync.Mutex)
	return pool
}

// Register allows you to register a new limiter to be managed by the pool
func (pool *LimiterPool) Register(id string, lmt *Limiter) {
	pool.mut.Lock()         // lock the pool for concurrency
	defer pool.mut.Unlock() // when the function ends, mutex unlocks

	pool.limiters[id] = lmt // assign Limiter to id in pool
}

// Find returns Limiter stored in pool, throws error if Limiter doesn't exist
func (pool *LimiterPool) Find(id string) (*Limiter, error) {
	pool.mut.Lock()         // lock the pool for concurrency
	defer pool.mut.Unlock() // when the function ends, mutex unlocks

	if lmt, exists := pool.limiters[id]; exists {
		return lmt, nil
	}
	return nil, fmt.Errorf("unknown id: %v", id)
}

// Remove removes an Limiter from the pool
func (pool *LimiterPool) Remove(id string) error {
	pool.mut.Lock()         // lock the pool for concurrency
	defer pool.mut.Unlock() // when the function ends, mutex unlocks

	if _, exists := pool.limiters[id]; exists {
		delete(pool.limiters, id)
	}

	return fmt.Errorf("unknown id: %v", id)
}

// cleanerFunction for pool.Clean()
type cleanerFunction func(string, *Limiter) bool

// Clean deletes every Limiter from the pool to which provided func evaluates to false
func (pool *LimiterPool) Clean(needsCleaning cleanerFunction) error {
	pool.mut.Lock()         // lock the pool for concurrency
	defer pool.mut.Unlock() // when the function ends, mutex unlocks

	for id, lmt := range pool.limiters { // loop over pool items
		if needsCleaning(id, lmt) { // check with user provided function if item has to be deleted
			err := pool.Remove(id) // remove from pool
			if err != nil {
				return err
			}
		}
	}
	return nil
}
