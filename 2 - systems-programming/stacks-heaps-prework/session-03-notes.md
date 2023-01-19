## Objectives

By the end of this session, you should understand:

- "the stack" and "the heap" in the context of a program's execution and memory usage
- How function calls, parameters, and local variables interact with the stack
- How `malloc` and `free` interact with the heap
- Allocations in C vs. Go

## Agenda

**Questions?**

- Will you only use sbrk if you yourself are writing a memory allocator?
	- sbrk or mmap
		- sbrk: Can't imagine use case outside of memory allocation
		- mmap: Might want to efficiently access a file
	- What does it mean to "map a device into memory"?
		- Don't know! Open question
```c
char *buf = mmap(...) // map device into some region
buf[5] = b; // start writing data into device
```

**CS:APP 3.7**

- How does a function call work at the machine code level?
	- `call` instruction
		- Pushes the return location onto the stack
		- Jump to the start instruction of the function being called
	- Allocate new stack frame for that call?
- How do you pass parameters?
	- Registers, then if there are more, use stack
- How do you return a value?
	- Put the result in %rax register
- How and where are local variables handled?
	- You can use registers
	- Otherwise, use the stack
- How do you go back to what was previously executing, after a function call was done?
	- The "return address" gets pushed to the stack by the `call` instruction, popped and PC is set to that value by the `ret` instruction
- How does recursion work at the machine code level?
	- Step through visualization as needed
- Open question: what does `alloca` do when you pass in a dynamic value?
- Open question: can we write a sample function where "callee saved" or "caller saved" registers comes into play?

**CS:APP 9.9**

- Can you describe what's going on in Figure 9.33?
- What do mmap / munmap / sbrk do?
	- Get more memory from the operating system
- What do malloc / free do?
	- mmap, etc. are system calls
	- malloc and free are standard library functions
		- They are an intermediary between us and the OS
		- One malloc call doesn't necessary result in one mmap call
			- malloc can split / reuse memory from previous mmap calls
	- What about calloc and realloc?
		- calloc is like malloc but initializes
		- realloc lets you grow or shrink some allocated region
			- Open question: Does realloc ACTUALLY respect requests to shrink a region?
- Why do we need dynamic memory allocation in the first place?
	- If the amount of memory you need cannot be predicted beforehand
	- If you want a value to stick around for a long time (outside the context of a particular function)?
	- Local variables are usually passed around "by value", but maybe you want to avoid copying a large buffer when passing it around
	- Heap also lets you handle larger sizes (usually)
- What are some requirements / design goals of a storage allocator?
	- Want it to be fast (don't spend too much time looking through the heap)
	- Also want high memory utilization (reuse memory that's previously been allocated)
		- Ideally we would want some sort of "coalescing" to merge adjacent blocks that have been freed
		- We want to avoid leaving lots of small gaps that aren't useful but still take up space
- How would you design an allocator that's really fast but really poor utilization?
	- Never coalesce, always grow the heap forward
- How would you design an allocator that's really good utilization but really slow?
	- Do a linear scan of the entire heap every time you try to allocate, to try to find the PERFECT sized block to return
	- Maybe use some sort of "moving" / "compaction" strategy
		- Might need double pointers to make that work?!

**CS:APP 9.11**

- Run through `memory_errors.c` on Digital Ocean instance, using `valgrind`

**Understanding Allocations**

- Returning a pointer to a local variable
	- `local_pointer.c`
		- Note different behavior on OSX compared to Linux ("undefined behavior")
	- `local_pointer.go`
- Escape analysis
	- `escape.go`

**Exploration**

- Go through results in `prework/solution`