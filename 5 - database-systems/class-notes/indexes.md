## Objectives

By the end of this session, you should understand:

- Different ways you can scan a table with (or without) an index
	- Sequential, index, index only, bitmap
- Additional features / use case
	- Clustered indexes, composite indexes, partial indexes, bulk loading
- How table statistics are involved in deciding which kind of scan to use

## Agenda

- Questions about prework?
	- TODO: Sample prework solution? (Will post in Slack afterwards)
	- Q on B+ tree:
		1. One B+ tree is associated with one table
		2. Pages are fixed size, so it's easy to jump to a particular page in the table (go to an offset at `PAGE_NUMBER * PAGE_SIZE` bytes)
		3. "Page" could refer to two things:
			- A 8kb region of the heap file (which contains rows)
			- A 8kb region of the B+ tree file (which contains an encoding of a single wide node of the B+ tree)

- Cluster
	1. Leaf of B+ tree contains (page #, item #) "UNCLUSTERED INDEX"
	2. Leaf of B+ tree contains full row (id, title, age, salary, ...) "CLUSTERED INDEX"
		NOT SUPPORTED in Postgres
		SQLite? TBD
		MySQL InnoDB
			- Source: https://dev.mysql.com/doc/refman/5.7/en/innodb-index-types.html
	3. Sort the heap file so that the order of the records in the heap file matches the order of the records in the B+ tree
		- Entries of the B+ tree look like
			- (1, 1), (1, 2), (1, 3), ... if you go traverse the leafs in order
		"CLUSTER" keyword in Postgres

- We've already stepped through the on-disk B+ tree format in depth during the first session of Data Structures
	- Recommendation: try this exercise yourself!

- Exploration of different ways we can use an index (or not)
	- Index scans
		- `index_scan.sql`
		- What query do you expect?

	- Bitmap index scan
		- `bitmap_heap_scan.md`
		- How does it work?
			1. Traverse the index (ignoring the heap file) to collect:
				- All the tuple pointers
				- For each page, all of the items on that page that match our query
			2. Efficiently visit all of the tuples in the heap file, in sequential order
		- Where does the "bitmap" come in?
		- Why is this a potential performance improvement compared to a standard index scan?

	- Clustered indexes
		- `cluster.sql`
		- Review: what does it mean for an index to be "clustered"?
		- How many indexes can you cluster a table on?
		- How do you ensure the index stays clustered as you update the table?
		- See https://www.postgresql.org/docs/current/sql-cluster.html

	- Composite indexes
		- `composite.sql`
		- What's a "composite index", and why is it useful?
		- How might you implement this?
		- Open question: how are delimiters handled in a composite index?

	- Partial indexes
		- `partial.sql`
		- What's a "partial index", why might it be useful?
			- See https://www.heap.io/blog/speeding-up-postgresql-queries-with-partial-indexes for more details

	- Bulk loading
		- `bulk_load.sql`
		- What happens if you insert a lot of data on a table that has an index?
		- What's "bulk loading" and why is it better than the alternative?
		- What's the "fill factor" / when do you start a new page in the bulk loading process?
			- See comment at top of `nbtsort.c`
		- Open question: can Postgres detect "bulk load" cases

```sql
SET max_parallel_workers_per_gather = 0;

EXPLAIN SELECT m.title, r.rating
FROM movies AS m, ratings AS r
WHERE m.id = r.movieid
ORDER BY m.id;
```

- How did the planner make these cost estimates?
	- Constants are from https://www.postgresql.org/docs/current/runtime-config-query.html#RUNTIME-CONFIG-QUERY-CONSTANTS

```sql
SELECT relname, relkind, reltuples, relpages FROM pg_class WHERE relname = 'movies';
SHOW cpu_tuple_cost;
SHOW cpu_operator_cost;
SHOW seq_page_cost;

EXPLAIN SELECT * FROM movies;
-- cpu_tuple_cost * N_tuples + seq_page_cost * N_pages
SELECT 0.01 * 27281 + 1 * 267;

EXPLAIN SELECT * FROM movies WHERE genres LIKE '%Comedy%';
-- (cpu_tuple_cost + cpu_operator_cost) * N_tuples + seq_page_cost * N_pages
SELECT (0.01 + 0.0025) * 27281 + 1 * 267;

EXPLAIN SELECT * FROM movies as m1, movies as m2;
-- seq_scan_cost    = cpu_tuple_cost * N_tuples + seq_page_cost * N_pages
-- materialize_cost = seq_scan_cost + 2 * cpu_operator_cost * N_tuples
-- rescan_cost      = cpu_operator_cost * N_tuples
-- 
-- cpu_operator_cost * N_inner * N_outer +
-- 	seq_scan_cost +
--	materialize_cost +
--	rescan_cost * (N_tuples - 1)
SELECT 0.01 * 27281 + 1 * 267;      -- seq_scan_cost
SELECT 539.81 + 2 * 0.0025 * 27281; -- materialize_cost
SELECT 0.0025 * 27281;              -- rescan_cost
SELECT 0.01 * 27281 * 27281 + 539.81 + 676.215 + 68.2025 * (27281 - 1);
```

- What statistics are available to the planner?
	- See https://www.postgresql.org/docs/current/planner-stats.html

```sql
-- pg_class values updated by VACUUM, ANALYZE, CREATE INDEX, etc.
SELECT relname, relkind, reltuples, relpages FROM pg_class WHERE relname = 'movies';

-- pg_statistic updated by ANALYZE and VACUUM ANALYZE
-- pg_stats is a human-readable view on top of pg_statistic
\d pg_stats
-- n_distinct is negative if it's a fraction, positive if it's a count
-- correlation is between logical ordering and physical ordering
-- inherited is for foreign keys (unused for us)
-- the elem columns are for collection types (unused for us)
\x
SELECT * FROM pg_stats WHERE tablename = 'movies' AND attname = 'id';
SELECT * FROM pg_stats WHERE tablename = 'movies' AND attname = 'title';
\x
```

- Why would CLUSTER followed by ANALYZE make a difference in the query plan?

```sql
SELECT attname, correlation FROM pg_stats WHERE tablename = 'movies';
```

- How would you estimate selectivity?
	- of `<`?
	- of `=`?

```sql
EXPLAIN SELECT * FROM movies WHERE title < '36th Chamber';
EXPLAIN SELECT * FROM movies WHERE title < 'Absentia';
-- Where did the estimated number of rows come from?
SHOW default_statistics_target;
SELECT histogram_bounds FROM pg_stats WHERE tablename = 'movies' AND attname = 'title';

-- Where did the estimated number of rows come from?
EXPLAIN SELECT * FROM movies WHERE genres = 'Comedy';
SELECT most_common_vals, most_common_freqs FROM pg_stats WHERE tablename = 'movies' AND attname = 'genres';
-- What if the value you're looking for isn't one of the most common?
```

- What if you have multiple conditions?

```sql
EXPLAIN SELECT * FROM movies WHERE title < 'Absentia' AND genres = 'Comedy';
```

- What if you have correlated conditions?
	- See https://www.postgresql.org/docs/current/multivariate-statistics-examples.html

```sql
\d multivar
SELECT * FROM multivar LIMIT 5;

EXPLAIN ANALYZE SELECT * FROM multivar WHERE a = 1;
EXPLAIN ANALYZE SELECT * FROM multivar WHERE b = 1;
-- estimate matches actual very well

EXPLAIN ANALYZE SELECT * FROM multivar WHERE a = 1 AND b = 1;
-- fail

CREATE STATISTICS stts (dependencies) ON a, b FROM multivar;
ANALYZE multivar;
EXPLAIN ANALYZE SELECT * FROM multivar WHERE a = 1 AND b = 1;
-- works now
```