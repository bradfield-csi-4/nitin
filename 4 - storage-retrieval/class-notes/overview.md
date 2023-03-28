## Objectives

By the end of this session, you should understand:

- The overall goal of the module / scope of the project

- How an abstract description of a data structure translates into actual bytes in physical storage (memory or disk)

- Project steps for next time

## Agenda

- Overview of module
	- This module is pretty different from the rest
		- There's not that much "content" but there is a LOT of design / implementation, and the goal is to end up with a nontrivial project by the end of it
		- The focus is definitely more on depth (actually implementing some ideas to get a really solid understanding) instead of breadth (getting a superficial exposure to lots of topics)

	- We will be encountering a few important data structures (skip list, B+ tree, log-structured merge tree, bloom filter) as well as more general ideas (sstables, indexes, write-ahead logs, read/write amplification), but all of these will be in the context of the project (instead of a random grab bag of data structures / ideas)

	- The goal is to make a clone of LevelDB
		- LevelDB is an "embedded key-value store"
		- A lot of the ideas used in LevelDB are mentioned in the BigTable paper, and are also used in Cassandra
		- RocksDB (a fork with a lot more features) is pretty widely used
		- There are even variants of MySQL (MyRocks) that use log-structured merge trees instead of B+ trees

- Questions about reading?
	- Quick overview of B+ trees
	- Hard to see concrete examples
	- "Writing to disk" still a fuzzy concept

- Approaches for `xkcd` exercise?
	- Confusion about two types of indexes
		- Primary storage id -> data
		- Secondary indexes (e.g. "posting list" type index)

- Discussion about abstract description of data structures (easy linked list example)

- Postgres case study
	- People who have taken SSBA will have seen this demo already
	- What is the data we have logically, and how does this get translated to what we actually have on disk?
		- `\d movies`
	- Quick look at heap file format
		- Fig. 1.4 of https://www.interdb.jp/pg/pgsql01.html
		- `SELECT pg_relation_filepath('movies');`
		- `xxd`
		- `pg_filedump -f -i -D int,varchar,text <file> | less`
	- Quick look at B+ tree index file format
		- Discussion of high-level process
		- Look for a movie "by hand"
			- What movie should we look for?
				- Terminator
			- `SELECT pg_relation_filepath('movies');`
			- `pg_filedump -f -i <file> | less`

- Invariant on B+ tree:
	- Every leaf node is at exactly the same depth

- More details on project
	- On-disk key/value storage
		- Open(path)
		- DB
			- Get(key)
			- Put(key, value)
			- Delete(key)
			- RangeScan(startKey, endKey) -> Iterator

		- Iterator
			- Valid() -> bool
			- Key() -> []byte
			- Value() -> []byte
			- Next()

	- Challenge: making it durable and efficient
		- Skip list for in-memory data structure
		- write-ahead log for durability
		- sstable format for on-disk lookups
			- index for quickly finding data
		- log-structured merge tree / compaction for efficient lookup
		- Bloom filter for quickly filtering out files to check

- `ldb_client` example

- Discussion of setup steps
	- Trivial implementation
	- Test cases?
		- Integrating `xkcd` exercise?
	- Benchmarks?
