## Objectives

By the end of this session, you should:

- Understand the architecture of LevelDB, and where each of the components/ideas we've discussed so far fits
- Have a rough roadmap for the remaining project steps (aside from the details of implementing a Bloom filter)

In particular, after the session you might find it worthwhile to re-read the LSMT section from Chapter 3 of Designing Data-Intensive Applications, to see how much more you understand this time.

## Agenda

- Check-in on project progress / blockers
- Questions
	- When doing binary search on the index or the block where the key could be in, what's the goal?
		- How is the "index" used?
			- Anytime there's some read involving that SSTable
		- How often are you doing the linear time operation involving the "index"
			- Idea: keep an in-memory index forever
	- Keeping track of SSTable size
		- Bryan is keeping track of the approximate size:
			- When you do a `Put`, you can add the key/value sizes
				- Subtract length of old value, add length of new value
			- How do you handle `Delete`?
				- Subtract length of old value (maybe add 1 for tombstone length)?

- Discussion
	- Merging multiple iterators
	- Compaction of multiple SSTable files
	- Leveled Compaction
		- level 0 is a temporary level, no "structure"
		- for levels 1 and above:
			- each level has more data than the previous one (e.g. 10x more data)
			- within each level, each file has a *distinct* range of keys (no overlap within a level)

		- how do you move an sstable from level 0 to level 1?
			- 

- Open question
	- How do you enable compression for your sstables?
		- Can you avoid having to decompress the entire file just to read a block
		- Idea: in the index, store offsets where compressed data is?
			- Might require strong coupling with choice of compression algorithm / strong understanding of how exactly your compression algorithm works