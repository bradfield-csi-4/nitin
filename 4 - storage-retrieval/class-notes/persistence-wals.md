## Objectives

By the end of this session, you should:

- How to use a write-ahead log to ensure durability of a data store
- Be able to work with binary data

In addition, you should be ready to design and implement an SSTable format

## Agenda

- Somewhat open-ended discussion today
- Look at code path for working with write-ahead log in LevelDB
	- What should you do when you first load/open the database?
		- Read the write-ahead log
		- Apply all the writes to the database

- Up next:
	- Where / how are deletions handled?
		- Guesses?
			- Write a nil value
			- "Tombstones"?
				- Indicate that entry has been deleted
				- TODO: Talk about why this is important for SSTables
		- What actually happens?
			- Different "keyType"
			- Delete has a nil value

	- What's the format / how are we parsing it?
		- How do you know when a record stops / when the next one starts?
			- If keys/values are fixed length, then you just know how much to increment
			- If not, you can use header bytes (first 4 bytes tell you key size, last 4 bytes tell you value size)
			- What would go wrong if you use a delimiter to separate keys/values?
				- Then you can't have arbitrary byte sequences (cause those might contain delimiters)

	- What are "sessions" in LevelDB?

- Discuss next steps for project / SSTables
	- What could go wrong with a long-running process, given the current implementation?
		- Ops could get slower
		- Run out of memory!
	- First off, what is an "SSTable"?
		- Sorted Strings table
		- 

## Notes

Code path when replaying journal:

	db, err := leveldb.OpenFile("path/to/db", nil)
		Open
		openDB
		recoverJournal
		decodeBatchToMem

Code path when writing:

	err = db.Put([]byte("key"), []byte("value"), nil)
		putRec
		writeLocked
		writeJournal

