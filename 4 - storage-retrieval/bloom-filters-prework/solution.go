package main

import (
	"encoding/binary"
	"hash/fnv"
	"unsafe"
)

const (
	BITS_PER_ELEM  = 64
	BYTES_PER_ELEM = 8
)

type basicBloomFilter struct {
	data    []uint64
	nHashes int

	// For convenience; equal to len(data) * BITS_PER_ELEM
	limit uint32
}

func newBasicBloomFilter(size, nHashes int) *basicBloomFilter {
	elems := size / BYTES_PER_ELEM

	return &basicBloomFilter{
		data:    make([]uint64, elems),
		nHashes: nHashes,

		limit: uint32(elems) * BITS_PER_ELEM,
	}
}

// Bloom filters rely on multiple hash functions.
//
// One way to do this is to literally choose different hash functions (FNV, murmur3, ...)
// but we'd run out pretty quickly.
//
// A better approach is to just pick a single hash function and tweak it. Here, we just
// add different prefixes to the item we're hashing.
func (b *basicBloomFilter) hash(i int, item string) uint32 {
	h := fnv.New32()

	// The basic way to convert `item` into a byte array is as follows:
	//
	//     bytes := []byte(item)
	//
	// However, this `unsafe.Pointer` hack lets us avoid copying data and
	// get a slight speed-up by directly referencing the underlying bytes
	// of `item`.
	buf := *(*[]byte)(unsafe.Pointer(&item))

	prefix := byte(i)
	h.Write([]byte{prefix})
	h.Write(buf)

	return h.Sum32() % b.limit
}

// Is the k-th bit in the b.data "bit-field" turned on?
func (b *basicBloomFilter) get(k uint32) bool {
	idx := k / BITS_PER_ELEM
	offset := k % BITS_PER_ELEM
	return b.data[idx]&(1<<offset) != 0
}

// Turn on the k-th bit in the b.data "bit-field"
func (b *basicBloomFilter) set(k uint32) {
	idx := k / BITS_PER_ELEM
	offset := k % BITS_PER_ELEM
	b.data[idx] |= (1 << offset)
}

func (b *basicBloomFilter) add(item string) {
	for i := 0; i < b.nHashes; i++ {
		k := b.hash(i, item)
		b.set(k)
	}
}

func (b *basicBloomFilter) maybeContains(item string) bool {
	for i := 0; i < b.nHashes; i++ {
		k := b.hash(i, item)
		if !b.get(k) {
			return false
		}
	}
	return true
}

func (b *basicBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
