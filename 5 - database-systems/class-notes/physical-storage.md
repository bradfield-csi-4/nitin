## Objectives

By the end of this session, you should understand how Postgres actually stores data on disk.

(We'll later build on this understanding when discussing indexes, performance, and transactions / MVCC.)

## Agenda

- What is the problem we're trying to solve (with our "file format design" / "storage solution")?
	- Making sure we have consistency / durability, data isn't lost, transactions happen all at once
	- Want to be able to write and lookup data with reasonable performance

What I said earlier:
- It would be safe if you didn't update the WAL until immediately before you try to "COMMIT" a transaction
	- Because it's acceptable to lose any "uncommitted" data

```sql
BEGIN TRANSACTION;
INSERT INTO ...
UPDATE ...
UPDATE ...
DELETE ...
INSERT ...
```

If you have data for an uncommitted transaction?
1. Is the data visible to someone else using postgres (e.g. some SELECT statement)?
	- Definitely not!
2. Is the data inside some file on disk somewhere?
	- Maybe in the WAL
	- Maybe!

Downside of having table data as CSV?
- Can have delimiting issues
- Almost always have to read the entire file
- Updates are difficult (unless you append only)

Downside of having table data as SSTables?
- Performing compaction takes up resources
- Individual files are immutable -> more tricky to handle updates (can't handle it within the context of 1 file)
- One key could occur in multiple SSTables

What's Postgres's approach to this?
- "slotted pages"
	- Fixed size "page" that contains fixed-size pointers to variable length data
	- File is going to be a bunch of these pages
- How do you add / update / delete rows?
	- Whenever you add a row to a page, you have to rewrite the whole page to disk
		- But you can defer rewriting the whole page (because the updates are already in the write-ahead log)

- Performance of writing to WAL vs. writing to slotted page
	- Where / what kind of write are we doing?
		- WAL: sequential
		- slotted page: random
	- How much data are we writing?
		- WAL: Just whatever update we're making (e.g. 1 row)
		- slotted page: rewrite entire page (8kb)

- At the leaf of a B+ tree, how do we refer to items?
	- (page, id of the item pointer)
- High-level discussion of slotted page structure
	- Tradeoffs vs. SSTables?
- Inspection of docs / live pages / struct definitions

- TODO if time allows: try to figure out how buffer manager handles shared memory (e.g. shmget, mmap)
- TODO: Post slide deck about commit process
- TODO: Why exactly did our "Bob" update not need more than 1 page (even after thousands of updates)?
	- See https://www.postgresql.org/docs/current/storage-hot.html
	
## Inspecting pages

```sql
SELECT pg_relation_filepath('foo');
xxd <file> | less
pg_filedump -f -i -D int,varchar,varchar <file> | less

# Checksums
\x
select * from pg_settings where name ~ 'checksum';
\x

CREATE TABLE foo (id int, name varchar(255), age smallint);
SELECT pg_relation_filepath('foo');

BEGIN;
SELECT txid_current();
INSERT INTO foo (id, name, age) VALUES (...);
ABORT;

BEGIN;
SELECT txid_current();
INSERT INTO foo (id, name, age) VALUES (...);
COMMIT;

CHECKPOINT;

CREATE INDEX idx_foo_name ON foo (name);
SELECT pg_relation_filepath('idx_foo_name');

INSERT INTO foo (id, name, age) VALUES(...);

CHECKPOINT;

UPDATE foo SET name = '...' WHERE id = ...;

CHECKPOINT;
```

## Polls / Discussion

- In your design, how do you load data from disk into memory when scanning a table?
	- All at once
	- Streaming, one row at a time
	- Streaming, one "page" at a time
	- Random access

- In your design, where is your schema stored?
	- Specified in the DB code
	- Inline with each table
	- Separately in some config file or table

- Based on your format, what data would you need to START writing a table to disk?
	- Schema
	- Number of rows
	- Actual row data for every row
	- None of the above

- Which of these features did you exclude for your design?
	- Split data into pages
	- Compression
	- Error detection / correction
	- Immutability

- Which of these features did you exclude from your design?
	- Split data into pages
	- Compression
	- Error detection / correction
	- Immutability

- If you're designing a database, how would you want to handle rows larger than the page size?
	- Don't split files into pages in the first place
	- Require row data to be smaller than page size
	- Spill row data onto subsequent pages if needed
	- Store column data as indirect pointers

## Inspecting buffer usage

* Experiment: restart Postgres (to wipe cache), run a query and see what gets populated

```sql
SHOW shared_buffers;

-- From https://www.postgresql.org/docs/current/pgbuffercache.html
SELECT n.nspname, c.relname, count(*) AS buffers
             FROM pg_buffercache b JOIN pg_class c
             ON b.relfilenode = pg_relation_filenode(c.oid) AND
                b.reldatabase IN (0, (SELECT oid FROM pg_database
                                      WHERE datname = current_database()))
             JOIN pg_namespace n ON n.oid = c.relnamespace
             WHERE n.nspname = 'public'
             GROUP BY n.nspname, c.relname
             ORDER BY 3 DESC;
```

## FAQ

- How does Postgres handle consistency in the face of power loss / crashes?
	- https://www.postgresql.org/docs/current/wal-reliability.html
		- Many levels of caching
			- OS buffer cache
			- disk drive controller cache
			- disk drive cache
	- From command line: `pg_test_fsync`
	- https://www.postgresql.org/docs/current/runtime-config-wal.html
		- Relevant flags:
			- `fsync`
			- `wal_sync_method`
	- https://github.com/postgres/postgres/blob/master/src/backend/storage/file/fd.c
		- It makes system calls like `open` and `fcntl`
	- https://pgpedia.info/d/direct-i-o.html
		- PostgreSQL does NOT write directly to the storage device that directly bypasses the OS

- How does the Postgres file format handle rows that are larger than the page size?
	- https://wiki.postgresql.org/wiki/TOAST
	- https://www.postgresql.org/docs/current/storage-toast.html
	- Compression + replacing "wide" columns with pointers to external data

## Resources

How Postgres handles very large fields:
	https://wiki.postgresql.org/wiki/TOAST

Considerations around ensuring durability when writing to disk:
	https://www.postgresql.org/docs/current/wal-reliability.html

E-book with details about file structure
	https://www.interdb.jp/pg/pgsql01.html

https://malisper.me/the-file-layout-of-postgres-tables/

Heap Only Tuples
	https://www.interdb.jp/pg/pgsql07.html
	https://github.com/postgres/postgres/blob/master/src/backend/access/heap/README.HOT
