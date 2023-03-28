## Objectives

By the end of this session, you should:

- Understand the three join algorithms (and variations) available in Postgres
	- Be able to sketch out pseudocode for each
- Be comfortable reading query plans for complicated queries involving multiple joins

## Agenda

- Code samples for DB

- Q's
	- What's the point of a join / what does it do?
		- Two places where data is being held (e.g. two tables)
		- You want to bring related data together
	- Have you ever implemented one "by hand" (either intentionally or unintentionally)?
		- Aggregating baseball stats from multiple CSV files / Excel files
		- Have seen this in "bad code"
		- Have seen joining in memory for a performance reason
	- Have you actually used one?
		- Pretty often

- Review of the three join algorithms / variations
	- Nested loop join
		- Basic approach:
			- For each member of one table:
				- Go through each member of other table
					- Check predicate
		- Optimizations / variations
			- Chunk nested loops join
				- Taking multiple chunks at a time
			- Index nested loops join
				- Rather than doing a sequential scan, use index to replace inner loop

	- Sort-merge join
		- Sort both tables
		- Zip the records together when you iterate
		- Only needs "one pass" through each table (in total)

	- Hash join
		- Add the contents of one table to a hashmap
		- Iterate over the other table, check against the hashmap

	- Of particular interest:
		- Handling duplicates in merge join (state machine)
		- Hybrid hash join with skew

## Pseudocode sketches

```python
# regular nested loop
for book_tuple in books:
	for author_tuple in authors:
		if book_tuple.author_id == author_tuple.id:
			yield combine_tuples(book_tuple, author_tuple)

# index nested loop
# want scenario (A) for "nice locality"
for book_tuple in books:
	author_tuple = author_pk_index.lookup(book_tuple.author_id)
	yield combine_tuples(book_tuple, author_tuple)

# sort-merge join
sorted_books = external_sort(books, books.author_id)
sorted_authors = external_sort(authors, authors.id)
bp = 0
ap = 0
while bp < len(sorted_books) and ap < len(sorted_authors):
	next_book = ...
	next_author = ...
	# if book.id < author.id: go to next book
	# if author.id < book.id: go to next author
	# if same: output a combined tuple, increment both

# index scan of authors
# want scenario (B) for "nice locality""
for pointer in author_pk_index.index_scan():
	author_tuple = authors.get(pointer)

# sort-merge join using a B+ tree index?
# replace one of the "external_sort" with just an index scan
sorted_books = external_sort(books, books.author_id)
sorted_authors = index_scan(author_pk_index)
bp = 0
ap = 0
while bp < len(sorted_books) and ap < len(sorted_authors):
	next_book = ...
	next_author = ...
	# if book.id < author.id: go to next book
	# if author.id < book.id: go to next author
	# if same: output a combined tuple, increment both

# hash join
hashmap = Hashmap()
for author_tuple in authors:
	hashmap.put(author_tuple.id, author_tuple)
	# if it's not a primary key:
	# hashmap[book_tuple.author_id].append(book_tuple)

for book_tuple in books:
	if hashmap.contains(book_tuple.author_id):
		matching_author_tuple = hashmap.get(book_tuple.author_id)
			yield combine_tuples(book_tuple, matching_author_tuple)

# hash join, if neither table fits into memory
n_chunks = estimate_chunks()
for author_tuple in authors:

	# Special case for prolific authors
	skew_hashmap = Hashmap()
	if author_tuple.id in very_prolific_authors:
		skew_hashmap.put(author_tuple.id, author_tuple)
		continue

	hashmap = Hashmap()
	idx = hash_function(author_tuple.id) % n_chunks
	if idx == n_chunks - 1:
		# add it to hashmap right away
	else:
		authors_chunks[idx].write(author_tuple)

for book_tuple in books:

	# Special case for prolific authors
	if book_tuple.author_id in very_prolific_authors:
		# process it right away
		continue

	idx = hash_function(book_tuple.author_id) % n_chunks
	if idx == n_chunks - 1:
		# process it right away
	else:
		books_chunks[idx].write(book_tuple)

for idx in range(n_chunks - 1):
	tmp_authors = authors_chunks[idx]
	tmp_books = books_chunks[idx]
	# Now we hash join tmp_authors and tmp_books
	hashmap = Hashmap()
	for author_tuple in tmp_authors:
		hashmap.put(author_tuple.id, author_tuple)

	# V1, inefficient but works
	"""
	for book_tuple in books:
		if hashmap.contains(book_tuple.author_id):
			matching_author_tuple = hashmap.get(book_tuple.author_id)
				yield combine_tuples(book_tuple, matching_author_tuple)
	"""
	# V2, more efficient
	for book_tuple in tmp_books:
		if hashmap.contains(book_tuple.author_id):
			matching_author_tuple = hashmap.get(book_tuple.author_id)
				yield combine_tuples(book_tuple, matching_author_tuple)
```

Scenario:
- books was recently clustered using a pkey on books.id
- authors, we know nothing about the order of rows in the heap file
- authors has a pkey on authors.id
- idea:
	- process books so that books with the same author id are processed together (either a sorting or a hashing step)

Exercise to try at home:
- Create these tables, see which join gets chosen

## Explore queries

```sql
-- Cartesian product
select * from movies as m, movies as n;

-- Ratings for Toy Story
select * from ratings as r, movies as m
where r.movieid = m.id and m.title = 'Toy Story (1995)';

-- Ratings for all movies, not just Toy Story
select m.title, r.rating
from movies as m, ratings as r
where m.id = r.movieid;

-- Top rated movies
select r.movieid, m.title, avg(r.rating) as avg_rating
from movies as m, ratings as r
where m.id = r.movieid
group by r.movieid, m.title
order by avg_rating desc
limit 10;

select r.movieid, avg(r.rating) as avg_rating
from movies as m, ratings as r
where m.id = r.movieid
group by r.movieid
order by avg_rating desc
limit 10;

-- Remove the order
select r.movieid, avg(r.rating) as avg_rating
from movies as m, ratings as r
where m.id = r.movieid
group by r.movieid
limit 10;

-- Other things to try
DROP INDEX idx_movies_title;
CREATE INDEX idx_movies_title ON movies (title);
CREATE INDEX idx_movies_title ON movies using hash (title);

SET enable_hashjoin=off;
SET enable_mergejoin=off;
SET enable_nestloop=off;

-- turn off parallel queries (default is 2)
SET max_parallel_workers_per_gather = 0;

-- update catalog stats
analyze;
```
