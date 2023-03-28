## Objectives

By the end of this session, you should:

- Understand the steps involved in processing a query, and use that as a perspective for understanding the overall architecture of a database system

- Be familiar with tools you can use for interacting with Postgres internals

## Agenda

- What are you most interested in learning about databases?
	- Performance tuning! +1001
		- DB access seems like the core of a lot of problems we're solving
	- Performance analysis, understanding characteristics of different queries and indices
		- Performance cost of retrieving too many columns?
		- Default settings for vacuuming / analyzing, benefits of changing that?
			- Heuristics for helping planner: how much do they matter
	- Understanding why relational databases were non-starter at larger companies, better understand why that's the case
		- Uber and Twitter for example
	- Given a use case, is there a good way to choose a database?
	- Being able to understand get insight into performance / resource usage

- Polls / anecdotes:
	- Production experience with relational databases?
		- 
	- Used an ORM?
		- Ran into ORM-related issues?
	- Write SQL queries directly in production code?
		- Ran into issues related to this?
	- Make ad-hoc queries against staging or production DB?
		- 
	- Performed basic DB administration (add/modify tables or indexes, perform backup/restore, update configuration)?
	- Debugged slow queries?
	- Ran into crazy bugs?

- Discussion of steps mentioned in readings
	- Network connection / protocol
		- What architectures have you seen?
	- Query parsing
		- What else do we need to do besides parsing?
	- Rewrite system
		- Example of view (and rewritten)
	- Planner / optimizer
		- EXPLAIN
		- What's selectivity?
	- Executor
		- Iterators

BP tuples vs. M tuples

```python
class FileScanIterator:
	def __init__(self, filename):
		...

	def next():
		# return next row from file

class SortIterator:
	def __init__(self, child):
		records = []
		while True:
			row = child.next()
			if row is None:
				break
			records.append(row)
		records.sort()
		records.reverse()
		self.records = records

	def next():
		# We reversed it, so pop gives us the smallest item
		return self.records.pop()

# Construct all of the iterators, put them into the "tree" structure
# What are the rows we're going to return to the client?

rowsToReturn = []
while True:
	row = rootIterator.next()
	if row is None:
		break
	rowsToReturn.append(row)
return rowsToReturn
```

- Tools we can use
	- postgres / psql
		- SHOW data_directory;
		- SELECT pg_relation_filepath('...');
		- EXPLAIN <QUERY>;
		- EXPLAIN ANALYZE <QUERY>;
	- xxd / pg_filedump for inspecting files
		- https://wiki.postgresql.org/wiki/Pg_filedump
	- xxd / pg_waldump for inspecting write-ahead log
	- lldb (especially for stepping through query execution)
		- But you still need to figure out where to set breakpoints (by browsing repo)
	- Wireshark for inspecting protocol
		- But we won't really be doing this for the rest of the module
