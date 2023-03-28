## Objectives

By the end of this session, you should understand:

* Durability consequences of having a buffer cache
* Role of the write-ahead log
* Features enabled by the WAL (point-in-time recovery, streaming replication)

## Agenda

- Clarification on terminology:
	- "serializability"
		- It is in the context of *transactions*
		- Executing transactions concurrently is equivalent to *some* "serial" ordering (executing the transactions one by one, in order)

	- "isolation level"
		- "serializable" is one of the possible isolation levels
		- what parts of other transactions can a given transaction see

**Follow-up on 2PL vs MVCC / isolation levels**

- What is "two phase locking" and how do we use it to get different isolation levels?
	- "two phase" means once you release a (exclusive) lock, you can't acquire any more
	- any time you update a row, you must first acquire an exclusive lock on that row
	- isolation level depends on how you handle **shared** locks
		- READ_UNCOMMITTED: no shared locks at all

- How do you get different isolation levels using MVCC
	- Recall that each tuple on disk has a tmin and a tmax
		- tmin: transaction id that inserted it
		- tmax: transaction id that deleted it
	- Postgres stores status of every transaction
		- IN_PROGRESS, ABORTED, COMMITTED
	- `txid_current_snapshot()` behavior depends on the isolation level of the transaction
		- Why do we need to know which transactions are currently in progress?

	- Visibility rules
	- `txid_current()` stays the same for the duration of a transaction
		- `READ_COMMITTED`: re-read for every query in the transaction
		- `REPEATABLE READ`: read once at the start of the transaction

		- `READ_UNCOMMITTED` is identical to `READ_COMMITTED`
		- `SERIALIZABLE` has some additional logic, out of scope
	- Difference between REPEATABLE READ and SERIALIZABLE:
		- "Serialization anomalies" could happen in REPEATABLE READ, but not in SERIALIZABLE

**Exploration of WAL**

```sql
SELECT pg_current_wal_lsn();
-- Or use pg_controldata
SELECT pg_walfile_name('...');

BEGIN;
INSERT INTO foo (id, name, age) VALUES (10, 'Neo', 37);
SELECT pg_relation_filepath('foo');
-- pg_filedump -f -i -D int,varchar,smallint <file>
-- What happens if the DB crashes?
COMMIT;
-- (pg_filedump again)
-- What happens if the DB crashes now?
CHECKPOINT;
-- (pg_filedump again)
```

**UNDO vs. REDO logging**

- What's the point of UNDO logging and REDO logging in the context of the lectures?
	- i.e. why was it so complicated in the Berkeley lectures?
	- What is the benefit of the version in the lecture?
		- No need for vacuuming
		- Space usage is less efficient in Postgres

- Why doesn't Postgres need UNDO logging?
	- It already has old versions in the heap file (MVCC)
- Why does Postgres need REDO logging?
	- Data for committed transactions might not be in the heap file yet

- What's STEAL / FORCE?
	- Buffer pool page management, whether a page is "dirty" or "clean"
		- FORCE: "when you commit a transaction, do you have to flush all involved buffer pool pages to heap file?"
			- if there's NO FORCE you need some sort of REDO logging
		- STEAL: "if a buffer pool page has data from an uncommitted transaction, are you allowed to evict it / flush it to heap file"
			- if there's STEAL you need some sort of UNDO logging (OR multiple versions)

- How do you rollback an aborted transaction in Postgres (MVCC)?

**WAL buffer**

- Step through:
	- https://www.slideshare.net/suzuki_hironobu/fig-902 (WAL updates)
	- https://www.slideshare.net/suzuki_hironobu/fig-903 (recovery)
- Look at config at https://www.postgresql.org/docs/current/runtime-config-wal.html
	- fsync
	- synchronous_commit
	- commit_delay / commit_siblings
	- wal_writer_delay
	- wal_sync_method
	- full_page_writes
	- checkpoint_timeout
- See also https://www.postgresql.org/docs/current/runtime-config-resource.html
	- shared_buffers
	- work_mem
- Fun fact: pg_xlog renamed to pg_wal because people were deleting it
	https://github.com/postgres/postgres/commit/f82ec32ac30ae7e3ec7c84067192535b2ff8ec0e

## Useful Commands

```sql
-- Transaction IDs
SELECT txid_current_snapshot();
SELECT txid_status(...);

-- Transactions levels
BEGIN TRANSACTION ISOLATION LEVEL ...;
-- either `ROLLBACK` or `COMMIT`

-- CRUD
CREATE TABLE foo (id integer, age smallint, name varchar(50));
INSERT INTO foo VALUES (1, 30, 'Neo');
UPDATE foo SET name='Dade Murphy' WHERE id=1;
DELETE FROM foo WHERE id=1;
CHECKPOINT; -- flush changes to disk

-- What locks are held?
SELECT locktype, relation::regclass, mode, transactionid AS tid,
virtualtransaction AS vtid, pid, granted
FROM pg_catalog.pg_locks l LEFT JOIN pg_catalog.pg_database db
ON db.oid = l.database WHERE (db.datname = 'bradfield' OR db.datname IS NULL) 
AND NOT pid = pg_backend_pid();

-- Indexes
CREATE INDEX idx_foo_name ON foo (name);

-- Buffer cache
SELECT n.nspname, c.relname, count(*) AS buffers
             FROM pg_buffercache b JOIN pg_class c
             ON b.relfilenode = pg_relation_filenode(c.oid) AND
                b.reldatabase IN (0, (SELECT oid FROM pg_database
                                      WHERE datname = current_database()))
             JOIN pg_namespace n ON n.oid = c.relnamespace
             WHERE n.nspname = 'public'
             GROUP BY n.nspname, c.relname
             ORDER BY 3 DESC;
-- Source: https://www.postgresql.org/docs/current/pgbuffercache.html
```

## Resources

- https://www.postgresql.org/docs/current/backup.html
- https://www.postgresql.org/docs/current/wal.html
- https://www.postgresql.org/docs/current/sql-checkpoint.html
- Database Hardware Selection Guidelines: https://www.youtube.com/watch?v=uXTsyiN_xJE
- Problems with writing to disk
	- https://www.postgresql.org/docs/current/wal-reliability.html
		- https://www.postgresql.org/docs/current/pgtestfsync.html
		- diskchecker: https://brad.livejournal.com/2116715.html

## Feedback

- Looked at Postgres for whole module, would be interesting to compare to other DBs (MySQL, SQLite, etc.)
	- Digging through systems that use 2PL would be helpful, e.g. https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html
- Would be interesting to drop some hints about NoSQL throughout the module (how the design differs, tradeoffs)
- How do AWS and Google Cloud implement products like Aurora (larger, more scalable databases)
- Move some demos over to prework!
	- Similar to optimizer exploration; suggest things to try
- Maybe some way to loop in interpreter type stuff for query parsing / ASTs?
- Upload .md notes to Bradfield courseware website