package main

import (
	"fmt"
	"sync"
)

const (
	numGoroutines = 100
	numIncrements = 100
)

type counter struct {
	count int
}

func safeIncrement(lock *sync.Mutex, c *counter) {
	lock.Lock()
	defer lock.Unlock()

	c.count += 1
}

func main() {
	var globalLock sync.Mutex
	c := &counter{
		count: 0,
	}

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < numIncrements; j++ {
				safeIncrement(&globalLock, c)
			}
		}()
	}

	wg.Wait()
	fmt.Println(c.count)
}

/*
Issue: In the call to `safeIncrement`, the globalLock variable is passed. In Go, all functions are pass-by-value, so a new
lock is actually created for every call which defeats the purpose of using the lock.

Fix: By instead passing a pointer to globalLock, the same lock is shared between all goroutines which synchronizes the calls
and prevents the data race. The result is 10000 increments.
*/
