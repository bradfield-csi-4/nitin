## Objectives

By the end of this session, you should understand:

- Concurrency control in general
- 2PL vs. MVCC
- SQL isolation levels, phenomena that can occur at various levels
- How Postgres implements concurrency control
	- Snapshot isolation
	- Locks in Postgres

## Quick discussion

- What's the problem we're trying to solve (with "concurrency control")?
	- If you have more than one **transaction** trying to access the same resource, one of them writing, you've got potential for problems
	- Why don't we just use existing techniques, such as `sync.RWMutex`?
		- This is solving a different problem, on a different timescale

## Demos

- Isolation levels
	- READ UNCOMMITTED
	- READ COMMITTED <- Default
	- REPEATABLE READ
	- SERIALIZABLE

- What is 2PL?
	- Locking happens in two phases: acquire, release
	- Once you start releasing **EXCLUSIVE LOCKS**, you can't acquire any more in that transaction
		- You're allowed to acquire and release **shared locks**

	- But WHAT locks do you need to acquire as you execute your transaction?
		- Whenever you UPDATE a row, you always acquire an EXCLUSIVE lock on that row
		- Whenever you READ a row, you acquire a SHARED lock on that row???
			- *The rules for when you acquire/release shared locks depends on the isolation level*

	- Policies for different isolation levels:
		- READ UNCOMMITTED:
			- No shared locks at all
			- You can get a phenomena called "dirty reads"
		- READ COMMITTED:
			- When you're reading a row, acquire the shared lock then release it after you're done reading
			- You can get a phenomena called "nonrepeatable reads"
		- REPEATABLE READ:
			- When you're reading a row, acquire the shared lock then *hold it for the rest of the transaction*

	- Dirty reads can't happen in Postgres
	- Nonrepeatable reads can happen in `READ_COMMITTED` but not `REPEATABLE READ`
	- Phantom reads can't happen in Postgres `REPEATABLE READ`
		- Why could it happen in 2PL but not MVCC?
	- Serialization anomaly
		- "dots" example from https://wiki.postgresql.org/wiki/SSI

- Locks
	- Updating a row doesn't interfere with reading that row (UNLIKE 2PL!)
	- Two updates in `READ_COMMITTED` block each other
	- Two updates in `REPEATABLE READ` not only block each other, but could cause an abort
		- TODO: Why do we have these semantics?
	- Deadlock between two transactions
	- Manually locking a table blocks reads
		- `LOCK TABLE foo IN ACCESS EXCLUSIVE MODE;`
		- What operations automatically acquire this lock?
	- See what locks are held
		- https://big-elephants.com/2013-09/exploring-query-locks-in-postgres/

- MVCC
	- `txid_current()` stays the same for the duration of a transaction
	- `txid_current_snapshot()` behavior depends on the isolation level of the transaction
		- `READ_COMMITTED`: re-read for every query in the transaction
		- `REPEATABLE READ`: read once at the start of the transaction

## Topics to discuss

- Semantics of reader/writer locks
- Locks at different "levels" (OS vs. DB)
- 2PL vs. MVCC
	- WARNING: 2PL is NOT used in Postgres!
	- Still useful to study for understanding isolation levels
	- Postgres still acquires locks, but there is one key difference: **readers and writers don't block each other**
	- DDL commands / explicit LOCK command can still block readers
- Ridiculous naming:
	- ROW SHARE / ROW EXCLUSIVE are modes of table-level locks
- Isolation levels
- Transaction snapshots

## Useful Commands

```sql
-- Transaction IDs
SELECT txid_current();
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
```

```shell
# Inspect index
pg_filedump -f -i ... | less 

# Inspect heap file
pg_filedump -f -i -D int,varchar,text ... | less
```

## Resources

- https://www.interdb.jp/pg/pgsql05.html
- Isolation levels
	- https://www.postgresql.org/docs/current/transaction-iso.html
	- https://docs.actian.com/zen/v14/index.html#page/adonet%2Fisolation.htm%23
	- https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels
- Locking
	- https://www.postgresql.org/docs/current/explicit-locking.html
	- https://www.postgresql.org/docs/current/sql-lock.html
	- Lock monitoring
		- https://wiki.postgresql.org/wiki/Lock_Monitoring
		- https://big-elephants.com/2013-09/exploring-query-locks-in-postgres/
	- Avoiding migration downtime
		- https://gocardless.com/blog/zero-downtime-postgres-migrations-the-hard-parts/#fn-1
		- https://medium.com/paypal-tech/postgresql-at-scale-database-schema-changes-without-downtime-20d3749ed680
- Transactional DDL: https://wiki.postgresql.org/wiki/Transactional_DDL_in_PostgreSQL:_A_Competitive_Analysis
