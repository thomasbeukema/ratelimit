# RateLimit

The ratelimit package allows for easy rate limiting in your Go application. It is meant for applications which need to put a rate limit on every user/ip/uuid/etc. It is designed with concurrency in mind, and managing a lot of Limiters has been made easy thanks to the LimiterPool.

| Donations: | Address |
| ---------- | ------- |
| BTC | 1Bi1W26FE9gSS3SoY7fPHd1C9fMk8E11z2 |
| ETH | 0xc406ac84b93bd3fec4801f6dd77aff243cdc574a |
| XLM | GA6BIUPHPS476B227MUXDAN5TA32QG6J73XTIZS2CS2D2AWJABCIMZOI |
| ETN | etnk5waqRk725J6ybevSQY9BU2ggJuj1Whskypw9e4pZ85Hmcki2dJBfbF31aJZuBn8ZEou6cFuFCW4G2iYnUcze55V27ycAS3 |

## Installation

```bash
go get github.com/thomasbeukema/ratelimiter
```


## Usage
### Limiter
The Limiter type handles all RateLimit logic. Before every operation you want to be rate limited, add an Allowed function. The function takes an argument of cost, which can be any integer.

```go
package main

import "fmt"
import "time" // Needed for Duration
import "github.com/thomasbeukema/ratelimiter" // Import the package

func main() {
	// Create a new Limiter with a capacity of 60 for an interval of 1 minute
	lmt := ratelimiter.NewLimiter(60, time.Minute)
	
	// Loop 69 times
	for i:=0; i<=69; i++ {
		// lmt.Allowed subs the provided cost from the bucket
		// This operation has a cost of 1
		if lmt.Allowed(1) {
			fmt.Println(i)
		} else {
			fmt.Println("Not Allowed")
		}
	}
}
```
### LimiterPool
LimiterPool allows you to easily manage any amount of Limiters. It is designed to be used when you want to apply a rate limit for every user, for every IP address, for every ... It is completely safe to use with concurrency, since it implements a mutex.
```go
package main

import "time"
import "github.com/thomasbeukema/ratelimiter"

// Function to decide which Limiters to remove from pool
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

func main() {
	pool := ratelimiter.NewPool() // Create new pool instance
	
	/*
		We give every Limiter in the pool an id of userNumber. In a production environment 
		you'd use a username, JWT token, IP address, etc as ID.
	*/
	
	userOneLimiter := ratelimiter.NewLimiter(10, time.Minute) // Create new Limiter instance
	pool.Register("userOne", userOneLimiter) // Register Limiter in pool with ID 'userOne'
	
	userTwoLimiter := ratelimiter.NewLimiter(20, time.Minute) // Create new Limiter instance
	pool.Register("userTwo", userTwoLimiter) // Register Limiter in pool with ID 'userTwo'
	
	userThreeLimiter := ratelimiter.NewLimiter(30, time.Minute) // Create new Limiter instance
	pool.Register("userThree", userThreeLimiter) // Register Limiter in pool with ID 'userThree'

	limiterOne := pool.Find("userOne") // Retrieve Limiter by ID
	err := pool.Remove("userTwo") // Remove limiter from the pool
	if err != nil {
		panic(err)
	}

	err = pool.Clean(cleanWhen) // Function which loops over every Limiter, decides to remove based on 'cleanWhen'
	if err != nil {
		panic(err)
	}
}
```
## License
This project is licensed under the MIT License. This license can be found in the LICENSE file.