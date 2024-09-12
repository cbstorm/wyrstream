package utils

import (
	"math/rand"
	"sync"
	"time"
)

func RandomRange(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := r.Intn(max-min+1) + min
	return res
}

var counter *Counter
var counter_sync sync.Once

func GetCounter() *Counter {
	if counter == nil {
		counter_sync.Do(func() {
			counter = &Counter{val: map[string]int{}}
		})
	}
	return counter
}

type Counter struct {
	mu  sync.Mutex
	val map[string]int
}

func (c *Counter) Increase(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.val[key] == 256 {
		c.val[key] = 0
	}
	c.val[key]++
	return c.val[key]
}

func AssertDefaultInt(v int, d int) int {
	if v == 0 {
		return d
	}
	return v
}
