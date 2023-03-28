## Objectives

By the end of this session, you should understand the tools and mindset involved in writing concurrent Go programs.

## Agenda

Questions?

- Read/write locks

Discussion

- What's a "data race"? What are some strategies for avoiding them?
	- Say you have some multi-threaded process (goroutines / threads) working on shared data
	- A data race is when you have some conflicting access to that data, and you're not really able to make guarantees about that access
	- You fix it by synchronizing access to that data
	- Race condition: maybe it's a problem if correctness of program depends on how different threads / goroutines are interleaved?
- Let's leave "race condition" undefined for now, or think about it as general term for "concurrency problem"
- For now, let's focus on more narrow definition of "data race"

- What's meant by "don't communicate by sharing memory, share memory by communicating"?
	- To share data, use channels
	- If you're using mutable state, that's more error-prone than if we're reasoning about them sending messages to each other
	- When you share memory, you need mutual exclusion, but if you're sending messages, you don't have to worry about that

```go
// Example of "communicate by sharing memory"

var hasWork bool
var mu sync.Mutex

// Announce that there's work to do
go func() {
	mu.Lock()
	hasWork := true
	mu.Unlock()
}()

// If there's some work to do, then do it
go func() {
	var shouldDoWork bool
	mu.Lock()
	shouldDoWork = hasWork
	mu.Unlock()

	if shouldDoWork {
		doWork()
	}
}()

// Example of "share memory by communicating"

var workQueue chan bool

// Announce that there's work to do
go func() {
	workQueue <- true
}

go func() {
	<-workQueue
	doWork()
}
```

- When do `sync.Mutex` and `sync.RWMutex` do? When should you use one or the other?
	- Mutex gives "mutual exclusion": only one goroutine at a time can hold the lock
	- RWMutex allows either one writer or multiple readers into the lock at once, but never both (also never more than one writer)
	- What are some of the nuances around priority / "starvation"

```
10 readers currently hold the RWMutex
1 writer tries to call rwmutex.Lock()
	- it'll block because you can't have readers and writers at the same time
What should happen if another reader comes in and tries to call rwmutex.RLock()
	In Postgres, reader has to wait!

GenerateAnalyticsReport() <- long running read-only transaction

	RunMigration() <- long running write transaction
		This one has to wait for the read-only transaction to finish

	LoadPage() <- short read-only transaction
		While RunMigration is waiting, LoadPage can't run either

```

- What's a "recursive" or "re-entrant" mutex, and what's Go's reasoning for not supporting it?
	- For simplicity (reentrant locks sets you up for weidr bugs)

- What are some reasons you can have more goroutines than threads?
	- goroutine stacks are not fixed size
	- fun experiment to try: how many goroutines can you launch?
	- qualitative differences
		- not much to say here but recommend checking out https://www.youtube.com/watch?v=YHRO5WQGh0k

- This is an open-ended question: at a high level, how does concurrency in Go "feel different" from concurrency in C?
	- Like go approach! "go keyword" is very nice
	- Don't like channel syntax, but like the concept in general
		- chan -> value
		- chan <- value
	- Kubernetes vs. physically managing infra: can do same things but different perspective, certain ops are easier to reason about
	- Feels more "declarative" (describing how information is going to flow, rather than handling low-level details)
	- Seems like Go is written for world where concurrency is a fact of life (convenient abstractions was necessary)
		- vs. C, thinking about concurrency, but doesn't feel like it was built in a world where that was central
	- one weakness of goroutines: C is talking more closely to OS

---

Go over exercises

How do `sync/atomic` and `sync.Mutex` actually provide synchronization?

	GOOS=linux GOARCH=amd64 go tool compile -S counterservice.go

Tour of `sync` and `sync/atomic` packages