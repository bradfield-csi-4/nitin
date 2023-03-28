## Objectives

By the end of this session, you should have a high-level understanding of how the query planner turns a parsed / analyzed query into an "executable" query plan.

## Agenda

Discuss a few loose ends relating to statistics / cost estimates

- Estimating the cost of joins

```sql
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

- Combining multiple WHERE clauses / correlated columns

```sql
EXPLAIN SELECT * FROM movies WHERE title < 'Absentia';
EXPLAIN SELECT * FROM movies WHERE genres = 'Comedy';
EXPLAIN SELECT * FROM movies WHERE title < 'Absentia' AND genres = 'Comedy';
```

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

- Go over how query planning would work for a single table
	- Briefly go over diagrams / structs from interdb.jp

- Ideas
	- Obvious search space is indexes
		- When would an index be usable?
	- Always can just try a sequential scan
	- Doesn't seem like you don't have that many options!
	- Would stats help you if you don't have an index?
	- What are ALL the options we have?
		- Do we need to do some sort / aggregation step afterwards
	- What are ALL the options we have for *accessing the data* in the first place?
		- Read the heap file
		- Do an index scan (access heap file via index)
			- Directly traverse index
			- Bitmap index scan
		- Do an index-only scan (don't even touch the heap file)

- Optimizing the order of joins
	- What if you just have two tables you're joining?
		- What decisions do you need to make?
			- Decide on type of join
				- Nested loop
					- Materialize or not?
				- Merge join
				- Hash join
			- Which relation is the outer vs. inner
			- When do you do selection and projection?

	- What if we have more than two tables?
		- Which table to join with which table first?
			- 

1. We're talking about a hypothetical naive strategy
2. We're not talking about joining 2 tables, we're talking about joining N tables
	tables A, B, C, D, E
	A and B, then join C into that, then join D, into that, then join E into that

- Dynamic programming algorithm that's actually used
	- 


	- Heuristic for >= 12 tables

## Code paths for query optimization

```c
// postgres.c:4415
exec_simple_query(query_string);

// postgres.c:1104
plantree_list = pg_plan_queries(querytree_list, query_string,
								CURSOR_OPT_PARALLEL_OK, NULL);
// postgres.c:912
stmt = pg_plan_query(query, query_string, cursorOptions,
					 boundParams);

// postgres.c:821
plan = planner(querytree, query_string, cursorOptions, boundParams);

// planner.c:273
result = standard_planner(parse, query_string, cursorOptions, boundParams);

// planner.c:403
root = subquery_planner(glob, parse, NULL,
						false, tuple_fraction);

// planner.c:1024
grouping_planner(root, false, tuple_fraction);

// planner.c:2057
current_rel = query_planner(root, standard_qp_callback, &qp_extra);

// planmain.c:269
final_rel = make_one_rel(root, joinlist);

/* Single table queries */
// allpaths.c:222
set_base_rel_pathlists(root);

// allpaths.c:352
set_rel_pathlist(root, rel, rti, root->simple_rte_array[rti]);

// allpaths.c:500
set_plain_rel_pathlist(root, rel, rte);

/* Joining multiple tables */
// allpaths.c:227
rel = make_rel_from_joinlist(root, joinlist);

// allpaths.c:2950
return standard_join_search(root, levels_needed, initial_rels);

// allpaths.c:3019
join_search_one_level(root, lev);

// joinrels.c:123
make_rels_by_clause_joins(root,
						  old_rel,
						  other_rels_list,
						  other_rels);

// joinrels.c:312
(void) make_join_rel(root, old_rel, other_rel);

// joinrels.c:760
populate_joinrel_with_paths(root, rel1, rel2, joinrel, sjinfo,
							restrictlist);

// joinrels.c:807
add_paths_to_joinrel(root, joinrel, rel1, rel2,
					 JOIN_INNER, sjinfo,
					 restrictlist);

// joinpath.c
```