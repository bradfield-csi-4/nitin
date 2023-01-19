## Objectives

By the end of this session, you should understand:

- Building blocks you can use to write concurrent programs in C
- Common pitfalls of concurrent programming

## Agenda

- Hello
	- How was break / previous module?
		- Everything was new (coming from web development background)
		- Fun to think about
			- Hierarchy
			- Not micro-optimizing
		- Interrupts at every layer!
		- Binary representation of data
		- Thinking about portability

	- Overview of module
		- Used to be called "Advanced Programming" (worst name ever)
		- Could've also been called "Runtimes and Languages"
		- Two parts:
			- Systems programming environment
				- Kind of an extension of Systems, but more emphasis on language / library / OS rather than hardware side
			- Tour of languages / compilers

- Polls about production experience with concurrency
	- What kinds of concurrent programming have you done (e.g. multiple threads, multiple goroutines, asynchronous tasks)?
		- multiple goroutines for RPC calls
		- process checking for repair data, reach out to different nodes
	- What synchronization primitives have you used (e.g. channel, mutex)?
		- mutexes, channels
		- promises (syntactic sugar for passing callback functions)
	- What concurrency errors have you run into?
		- deadlock on go prep work
		- sidekiq - job queue

- Interleaving demo
	- Take some time to predict what will get printed
	- Until 5:30PM pacific
	- Ideas?
		- Should never both be 0:
			- Reasoning: if first statement that runs is x = 1, later r2 = 1; if first statement that runs is y = 1, later r1 = 1
		- Assuming they're atomic, there are 6 paths it can take
			- In 4 of them, they're both 1's
			- In the other 2, they're 1/0 or 0/1?
		- All of them are possible?
		- Compiler might reorder reads and writes within a thread?
			- Do you know a way to prevent that?
			- volatile?
	- Open question: is "x = 1" atomic? (We can look at disassembly)

r1 = 0, r2 = 0: 4
r1 = 0, r2 = 1: 978566
r1 = 1, r2 = 0: 21430
r1 = 1, r2 = 1: 0

- Review questions
	- Can you list some situations where concurrency would be useful?
		- All of our systems (from lowest level up to highest), we're wrecked if we have to wait
		- Managing access to any shared resources
		- Breaking problems down into independent subproblems
		- Delegating independent work that can be done at the same time

	- If someone on your team watched a Rob Pike talk and got confused by the phrase "concurrency is not parallelism", what would you tell them?
		- Concurrency as a design situation
		- Parallelism as a state of running
		- Concurrency as "interleaving"

	- For each of the three approaches described in CS:APP (processes, threads, I/O multiplexing):
		- How does it work at a high level? What kind of "shared state" is there (if any)?
			- processes
				- share file table, but otherwise independent virtual address space
			- threads
			- I/O multiplexing
		- What are the pros and cons of this approach?
			- processes
				- pros: impossible to overwrite shared state
				- cons: difficult if you do need to share state; slower, more costly (context for a process is fairly heavy)
				- cons: IPC (e.g. message passing): if there's data dependency relationship, you've got the illusion of indepedent state, but they're not really independent
			- threads
				- pros: lighter weight than processes (since you're sharing virtual memory, don't need a new address space)
				- pros: can create more threads than processes
				- cons: have to worry about shared data
			- I/O multiplexing
				- you don't get the parallelism
				- state management is tricky
		- What are some systems calls and/or library functions associated with this approach?
			- processes
				- fork
				- wait_pid
				- execve
			- threads
			- I/O multiplexing
				- select
					- takes file descriptor set
				- "watch ALL of these file descriptors"
				- when at least one of them has data, return
				- still just one thread
				- read(conn, timeout)

```
fd_set = {conn1, conn2, stdin}
# This will block until at least one of the items in the set has data available to read
# It will also return the set of items that has data available
ready = select(fd_set) 
# Maybe now ready is {conn2}, or {stdin}, or {conn1, stdin}, etc.
```

Open question: how is this implemented?
	- Polling?
	- Signal?

In Node.js:
	- There's a thread doing I/O multiplexing
	- When you encounter something you want to defer (e.g. blocking system call), it spawns another thread
		- When that other thread is done, it sends some signal that indicates to main event loop that some data is ready

- Examples
	- Deadlock
		- How do we fix the example?

	- Race condition
		- What goes wrong in the badcnt.c example?
			- Why does it work for small counts?
			- Look at disassembly
				- `cc --target=x86_64-apple-darwin-macho *.c`
				- `objdump --disassemble-symbols=_thread a.out`
			- What if we compile it with O2?
				- https://c9x.me/x86/html/file_module_x86_id_140.html
			- What happens if we use `atomic_long`?
				- Make sure to `#include <stdatomic.h>`
				- What happens to the disassembly?
				- What was the point of the volatile keyword?
					- Can we fix the original example with 

	- Matrix multiplication example
		- Poll: What kinds of speedup did you get previously?
		- How can we use thread-level parallelism to speed things up?
			- What kind of speedup do you expect?
				- `sysctl -n hw.ncpu`
				- About This Mac -> System Report -> Hardware
		- `clock` measures elapsed CPU time across all threads, whereas `gettimeofday` measures elapsed wall clock time. What do you think would happen if we use `clock` to measure the parallel strategy?
		
	- pstree / vmmap
	- Bonus: What concurrency model does Postgres use to handle multiple connections?

## Resources

- https://csapp.cs.cmu.edu/3e/code.html
- https://en.wikipedia.org/wiki/CPU_time#POSIX_functions_clock()_and_getrusage()
