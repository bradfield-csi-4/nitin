## Objectives

By the end of this session, you should understand:

- How a skip list works
- The overall architecture of LevelDB
- Upcoming project steps

## Agenda

- How did implementation go so far?
- Blockers?
	- Go ecosystem issues:
		- Wrestling with go packages
	- Rusty on Go (+1)
	- Go pointers vs. C pointers, "does it live in heap or stack, do you actually need a pointer here"
- Strategies for benchmarking?
 
```go
start := time.Now()
// Do 10,000 gets
elapsed := time.Since(start)
// Print how long 10,000 gets took

// make(type, len, cap)
next := make([]*node, 0, MAX_LEVEL)



type Node struct {
	// fields of Node
}

func (n *Node) Get(key []byte) {

}
```

- Recap of high level idea?
	- Linked list with additional nodes that point far into the list
	- Starting at the top is kind of like doing a binary search (skipping much more)

- TODO: What does LevelDB's "snapshot" feature do with the in-memory skip list?

- Skip list implementation notes
	- `next []*node` vs `next [MAX_LEVEL]*node`
	- Role of randomness / chance of imbalance
		- Compared to binary search tree / hashmap?
	- Probabilistic runtime analysis out of scope for today (we can talk about it if there's interest)

- Design brainstorming for adding some more features (before looking at actual implementation)
	- Goals?
		- Persisting it
		- If it crashes, we don't want to lose data
		- If the dataset gets too big to completely fit in memory, it should still work
	- Stretch goals
		- Reverse iteration
		- Snapshots
- Ideas
	- Each node is a file
	- Create a binary format
	- What if we only need to restore the in-memory structure
		- Append-only event log
	- Also periodically write the in-memory state to disk as a "snapshot"
		- Could depend on size of keys and values of in-memory structure
		- Every time it reaches some threshold (10 mb or whatever), you do something
			- Flush the data in the skiplist out to disk
	- Let's say it's been running for a while and you have 1000 files on disk that were the result of periodic "flushing"

- One binary search on a 1 gb file is faster than 1000 separate binary searches on 1000 1 mb files
	- Approximately how many lookups do you need on a 1 gb file?
		- 1 billion items in it
		- What's log2(1 billion)?
			- 2^10 is about 1000
			- 2^30 is about 10^9
		- What's log2(1000000)
			- 2^20
			- 20 lookups PER FILE
			- 20,000 lookups

- Go over LevelDB architecture
	- Revisiting `ldbclient` example / on-disk files
	- Path that a Get / Put takes
	- Handling deletes
	- Different components that are involved

- Goals for next time