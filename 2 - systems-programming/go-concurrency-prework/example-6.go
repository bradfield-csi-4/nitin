package main

import (
	"fmt"
	"sync"
)

type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

func (c *coordinator) logState() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) setLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	c.leader = leader
	c.lock.Unlock()

	if shouldLog {
		c.logState()
	}
}

func main() {
	c := newCoordinator("us-east")
	c.logState()
	c.setLeader("us-west", true)
}

/*
Issue: In the setLeader function, a write lock is obtained, and the corresponding unlock call is deferred to the end
of the function. However, before the function ends, the `logState` function acquires a read lock. When a thread acquires
multiple locks in Go, the result is a deadlock.

Fix: Instead of deferring the unlock to the end of the function, we can call it immediately after the write operation.
*/
