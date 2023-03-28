## Objectives

By the end of this session, you should understand:

- The problem that garbage collection is intended to solve, and two high-level approaches for solving that problem (reference counting, mark and sweep)
- How Go's concurrent collector avoids the problem of huge latencies with a naive "stop the world" approach
- How to identify when GC is the cause of performance problems in a Go program; ways to improve performance in that case

## Agenda

**Questions?**

- Go through some details of the GOGC doc
- Open question: is there a production GC that actually does "stop the world"?
- What are "roots", i.e. how do we start the mark process?
	- Enumerate everything that could possibly be accessible to the program
		- registers
		- entire contents of stack
		- global variables
	- float core_temperature; // 7000.27 K
		- maybe this value HAPPENS to look like a heap pointer
		- "conservative": interpret this as a heap pointer, even if it's not
			- mark that address as reachable
			- can't get a strong guarantee that the node is actually reachable
		- issue: malloc calls might not be in 1-1 correspondence with objects?
			- malloc(sizeof(struct my_data_structure));
				- We get back an address pointing to the start of the region of exactly the right size, that we can now use
				- The "heap blocks" are all going to be exactly the size that we asked for
				- (Or maybe they'll be bigger after we free and coalesce them)
			- How much virtual memory did we get from the OS as a result of this call?
			- When the malloc library function calls mmap or sbrk under the hood, what's the size of the thing we get back?
				- It's going to get back some number of pages
				- Some multiple of 4kb, or whatever the "page size" happens to be
		- packing data into low-order bits?
			- let's say next_ptr was supposed to be 0x10014800
			- let's say we know that all blocks are a multiple of 16
			- what happens if we see next_ptr = 0x10014807
				- Interpret the pointer as 0x10014800 (the value with the last 4 bits cleared)
				- 3 bits are on (we can use that for various boolean flags, or maybe store)

```c
struct heap_block_header {
	heap_block_header* next_ptr;
	heap_block_header* prev_ptr;
	int marked;
	int is_allocated;
	int block_size;
};
```

Why do you have this linked list of blocks?
- Many possible reasons:
	- Keep a linked list only of free blocks, so that you can quickly find one when user calls `malloc`
	- You have a mark and sweep collector and you want to be able to scan the entire heap
- Question:
	- When we do "sweep", do we go through all memory that was mmap'd, or all blocks that were malloc'd

						  v---- pointer to some field in the middle of that
allocated object: [                  ]

```c
struct node {
	struct node *left_child;
	struct node *right_child;
	void *data;
}
```

when you call malloc and get back a pointer
	that pointer points directly to the start of the USABLE region / payload

w                  v----- ptr that malloc returns
[heap_block_header][actual contents of block (e.g. 24 bytes)][heap_block_footer]
		- imagine a binary search tree of blocks, sorted by address

"Root nodes":
- stack variables, registers, global variables

**Testing without GC**

```
GOGC=off ./project

# In another terminal window:
hey -m POST -c 100 -n 1000000 "http://localhost:5000/search?term=topic&cnn=on&bbc=on&nyt=on"

# In a third terminal window:
top
```

- If we do `-n 10000` instead, how does performance compare with GC off vs. on?
- How long does it take for the program to crash?

**Examples**

- Reference counting in Python
	- Printing out reference count with `sys.getrefcount`
		- https://docs.python.org/dev/library/sys.html#sys.getrefcount
	- Adding a finalizer (`__del__`) to show when an object is destroyed
		- https://docs.python.org/3/reference/datamodel.html#object.__del__
	- References getting reassigned
	- Reference falling out of scope
	- Circular references
	- Running `gc.collect()`
	- How often does gc run?
		- "When the number of allocations minus the number of deallocations exceeds threshold0, collection starts."
		- https://docs.python.org/3/library/gc.html#gc.set_threshold

- Allocating a new object in go
	- `go tool compile -S newobject.go`
	- runtime.newobject
		- https://github.com/golang/go/blob/master/src/runtime/malloc.go#L1198
		- Calls into runtime.mallocgc
	- When does `free` happen?

- Finalizing objects in go
	- How many objects were collected?
	- What if we change `numTotal`?
	- What if we set the following environment variables:
		- `GODEBUG=gctrace ./finalize`
		- `GOGC=off GODEBUG=gctrace ./finalize`
		- `GOMEMLIMIT=10KiB GODEBUG=gctrace=1 ./finalize`
		- See https://pkg.go.dev/runtime#hdr-Environment_Variables

**Discussion Questions from Reading**

Now that you're thoroughly horrified of the dangers of manual memory management (from last time), let's talk about how garbage collection lets us avoid that!

**Profiling Exercise**

```
# Install
git clone https://github.com/ardanlabs/gotraining.git
cd gotraining/topics/go/profiling/project
go build
GODEBUG=gctrace=1 ./project

# Run
open http://localhost:5000/search
hey -m POST -c 100 -n 10000 "http://localhost:5000/search?term=topic&cnn=on&bbc=on&nyt=on"

# Profile
open http://localhost:5000/debug/pprof/
go tool pprof http://localhost:5000/debug/pprof/allocs
top -cum
list rssSearch
```