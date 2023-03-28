## Objectives

By the end of this session, you should:

- Understand the architecture of LevelDB at a high level (especially the steps involved in executing lookups)
- Be ready to further investigate a suitable data structure (skip list, balanced tree) and use it to implement the "memtable" component

## Agenda

- Follow-ups from last time
	- Better "covering scans", index-only scan / explanation
		- `\d movies`
		- `pg_filedump -f -i <index> | less`
		- `EXPLAIN SELECT title, genres FROM movies WHERE title > 'T' and title < 'Ta';`
	- CLUSTER command in Postgres vs. general "clustered index" concept
	- *Important* Discussion of the file API / abstraction
		- Programmer's view of working with files / disk
			- file is "pointer" to some location on disk
			- "index card" you can read / write from, and store
			- file has some "state" associated with it:
				- where you are in the file (an "offset")
					- can manipulate that offset via "seek"
			- in unix, files don't even have to be associated with the disk, they're just a "source of bits"

		- How does this differ from your view of "memory"
			- in main memory, any location is accessible
				- "random access" memory
			- for a disk, it depends: you may have to move the "physical head" from one location to another
			- main memory is constrained by pointer size

		- What's your "abstract view" of main memory?
			- large array of bytes
			- scratch pad
			- Nightcrawler (teleports)
		- What's your "abstract view" of disk / a file?
			- also large array of bytes
			- persistent

		- System calls
			- open
			- read
				- reads from wherever the "current offset" is (associated with the open file)
			- write
				- writes at wherever the "current offset" is
			- lseek
				- updates offset
		- More Golang API methods
			- ReadAt
			- ReadFrom
		- Performance considerations
			- random vs. sequential reads on a spinning disk
			- smallest read (analogous to "cache line")
			- fsync
		- Summary
			- File is just an array of bytes
			- You can read/write from a file similarly to memory
				- Caveat: it has an "offset" associated with it
				- (But you can easily move it with "seek")
			- There are some performance differences
			- It's persistent

- Discuss setup / progress so far, share ideas

- memtable options
	- map?
	- balanced binary tree?
	- skip list
		- Goal is NOT to go over every detail, but just give enough context for you to feel confident investigating further

- Goals for next time

- High-level overview of LevelDB architecture
	- files that are created
	- steps involved in Get / Put
	- handling deletion
	- possible discussion of LSM tree???

## Suggestion

```go
type Node struct {
	key []byte
	value []byte

	// Array of pointers to next nodes
	// next[0] is bottom level
	// next[maxLevel - 1] is top level ("express" lane)
	next [maxLevel]*Node

	// Max valid level of this node
	level int
}
```