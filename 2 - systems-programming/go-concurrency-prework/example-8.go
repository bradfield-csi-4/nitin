package main

import (
	"fmt"
	"sync"
	"time"
)

type dbService struct {
	lock       *sync.RWMutex
	connection string
}

func newDbService(connection string) *dbService {
	return &dbService{
		lock:       &sync.RWMutex{},
		connection: connection,
	}
}

func (d *dbService) logState() {
	fmt.Printf("connection %q is healthy\n", d.connection)
}

func (d *dbService) takeSnapshot() {
	d.lock.RLock()
	defer d.lock.RUnlock()

	fmt.Printf("Taking snapshot over connection %q\n", d.connection)

	// Simulate slow operation
	time.Sleep(time.Second)

	d.logState()
}

func (d *dbService) updateConnection(connection string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.connection = connection
}

func main() {
	d := newDbService("127.0.0.1:3001")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.takeSnapshot()
	}()

	// Simulate other DB accesses
	time.Sleep(200 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.updateConnection("127.0.0.1:8080")
	}()

	wg.Wait()
}

/*
Issue: The `takeSnapshot` call tries to acquire a read lock concurrently with the `updateConnection` trying to acquire
a write lock. Since there is a 200ms delay, `takeSnapshot` usually wins. Within `takeSnapshot`, the `logState` call
tries to acquire another read lock, and this also happens concurrently with `updateConnection`. The result is a race between
the two. If `logState` acquires the read lock first, then all is good. But since there is a 1 second delay, `updateConnection`
tries to acquire a write lock. With a RWMutex, a write lock is blocked until all existing locks are released. The next thing
that happens is that`logState` gets blocked trying to acquire a read lock. That means `takeSnapshot` will never release
its lock, and we have a deadlock.

Fix: We can remove the locking attempt from `logState` because the goroutine already holds a read lock. If needed, we can
create another `logState` function that does contain a lock for other use cases that require it.
*/
