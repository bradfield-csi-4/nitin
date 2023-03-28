## Objectives

By the end of this session, you should understand:

* Useful tools for interacting with DB implementation details
* Code structure and control flow of query execution ("tree of iterators")
* Algorithms for external sorting and hashing

## Agenda

- Discuss overall goals of module (didn't really talk about it last time)
	- Not as clear-cut as previous module: databases is a very large topic
	- Understand your relational database as a piece of software, rather than as an abstraction
		- Algorithms / data structures (both in-memory and on disk)
		- Components / division of responsibilities
	- Add a lot of ideas to your toolbox for both understanding existing systems and designing your own systems

- Tools we can use
	- postgres / psql
		- SHOW data_directory;
		- SELECT pg_relation_filepath('...');
		- EXPLAIN <QUERY>;
		- EXPLAIN ANALYZE <QUERY>;
	- xxd / pg_filedump for inspecting files
		- https://wiki.postgresql.org/wiki/Pg_filedump
	- xxd / pg_waldump for inspecting write-ahead log
	- Wireshark for inspecting protocol
		- But we won't really be doing this for the rest of the module
	- lldb (especially for stepping through query execution)
		- But you still need to figure out where to set breakpoints (by browsing repo)

- Better dataset
	- Source: https://bigmachine.io/products/a-curious-moon/

- Go over prework
	- 

- Discussion of "sorting and hashing" in the context of queries

## Q. How do we sort a table that doesn't fit in memory?

High-level idea:
- Sort chunks that DO fit into memory
- Merge all the chunks in a "streaming" way
	- Merging multiple chunks does not require loading the entire chunk into memory

Open question 1:
- After we have multiple sorted chunks, should we do another pass to turn it into a SINGLE sorted chunk?

	- Assuming "YES" to the question above:
		- How do we do this in the first place?
			- We can employ "merging"

		- When merging small sorted chunks into larger sorted chunks, how many should we merge at a time, and in what order?
			- Merge one in at a time
				- O(n^2) time
			- Merge them all at once
				- Much more efficient
				- Load one page from each chunk at a time
				- Use a min heap decide which is the smallest among all chunks
			- Merge them in pairs

		- Should the larger sorted chunk be a single file, or split up into multiple files (e.g. similar to the LevelDB strategy)
			- John question

	- Actually, the answer is "NO" to the question above!
		- Greg: "it would be redundant work"

- "Hashing"
	- 

```
id,name,favorite_letter
1,alice,A
2,bob,Z
3,carol,D
4,dave,A
5,eve,D
6,fred,G
7,george,Z

SORT ON favorite_letter
1,alice,A
4,dave,A
3,carol,D
5,eve,D
6,fred,G
2,bob,Z
7,george,Z

HASH ON favorite_letter
7,george,Z
2,bob,Z
1,alice,A
4,dave,A
3,carol,D
5,eve,D
6,fred,G
```

```go
m := make(map[ColValue][]Row)
for _, row := range rows {
	c := getHashKey(row)
	if _, ok := m[c]; !ok {
		m[c] = make([]Row)
	}
	m[c] = append(m[c], row)
}

var output []Row
for colValue, rowList := range m {
	output = append(output, rowList...)
}
```

What if the table doesn't fit in memory, how can we do this?


- Query plans involving:
	- in-memory quicksort
	- top-n heapsort
	- external sort?
	- no sorting because we index-scanned?

- Step through query executor code in Postgres
	- simple query (SELECT)
	- limit
	- sort
	- index scan

- (If time allows) Live code some stretch goals?

- TODO: Figure out where the temp files are created