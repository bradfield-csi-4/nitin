## Objectives

By the end of this session, you should:

- Understand the role of a Bloom Filter and be ready to add one to your implementation
- Clarify any remaining blockers in your design

## Agenda

- What's a good benchmark / milestone to aim for?
	- Having single SSTable
	- Being able to handle memtable + one or more SSTables
	- (Optional) Having the ability to compact multiple SSTables into one
	- (Major stretch goal) Compaction with multiple levels, tracking ranges for each SSTable
	- (Major stretch goal) Background flushing (with two memtables), background compaction

- Bloom Filter discussion
	- Interface / use case?
		- Checking existence in a set
		- Why is it better than just using a Python set?
			- Much smaller (only store ~3 bits per key)
		- How is it worse?
			- Probabilistic
				- If an element is in the set, you will always get "yes it's there"
				- If an element is not in the set, you could still get  "yes it's there"
			- If your dataset is too big for your configuration, you could get to the point where everything is a "false positive" (always says yes)
	- How are we going to use this in our system?
		- Disk reads are expensive
		- We can be 100% confident that something is *not* present
		- We can make sure they're small enough to keep in memory
		- Goal: speed up Get operation
	- Implementation details
		- What sort of performance benefit does it have? (And when?)
			- Would especially have benefit if you're trying to call Get on data that's not there
			- Concretely, when would this happen?
				- If you have lots of deletes, this could happen?
				- What if you have a lot of SSTables in Level 0
			- Situations / use cases where Bloom filter wouldn't be a win?
				- If you rarely look up things that aren't there (e.g. reads usually follow writes to same key within short period of time)
				- High write, low reads
					- e.g. logging events, telemetry
				- If you're constantly just range scanning
		- How does it change the on-disk format
			- Serialize the bloom filter somewhere in the file as well (e.g. after the index)
	- How does the bloom filter actually work
		- Inserting
			- Take a few hash functions (e.g. 3)
				- Prefix your key with 1, 2, 3
				- Could also use one and split it into three pieces
			- Hash the key with all of the hash functions (and mod the size of your bloom filter)
			- Turn on all those "bits"
		- Searching
			- Run key through same 3 hash functions
			- Make sure the "bits" are all on
		- What's the data type (in memory)?
			- If there's not a "bitset" type, then an array / slice
			- Byte / uint8 / or could use larger sizes (e.g. uint32, uint64)
		- Can it be resizable?
			- No!
		- Fun question: is there a way you can modify the bloom filter to support delete?
	- Tuning
		- Parameters:
			- n (amount of memory used / total # of bits)
			- number of hash functions that we use
			- number of keys we're inserting
				- max length of a key?
			- chance of getting a false positive

```go
// The Get for a single SSTable
Get(key []byte) ([]byte, error) {
	if !filter.MaybeContains(key) {
		return nil, KeyNotPresentError
	}

	// check in index where to read
	// seek to the location and read bytes from disk
	// parse individual items until you've reached (or moved past) key
}
```

- TODO: Could there be a benefit to adding a "meta" bloom filter (e.g. per level)?
	- If you're in level 1, 2, ... then you only need to check a single file
	- If you're in level 0 you might have to check every file

- Scenarios:
	- One bloom filter per file in Level 0, plus one for all of Level 0	
		- If "meta bloom" lookup says "no it's not there", we avoid 10 "per-file bloom" look ups
		- If "per-file" bloom look up says "no it's not there", we avoid a disk read

	- Just one bloom filter to all of Level 0

---

- Check-in on progress
	- RangeScan issue with deleted values when merging:
		- Underlying iterator skips deleted values (because the interface isn't supposed to show values that are deleted)
		- goleveldb handles this by having a RawIterator that contains additional metadata
	- Get(key)
		- How do you signal the difference between "deleted" and "not present"?
		- Key is present and value is empty?
	- Overarching themes
		- Was difficult to maintain modularity (as opposed to depending on internals)

- Challenges / observations?
- Goals for next time
- Stretch goals (for the future)?