package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"time"
)

const wordsPath = "/usr/share/dict/words"

func main() {
	words, err := loadWords(wordsPath)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	var b bloomFilter = newBloomFilter(1600000, 3)

	// Add every other word (even indices)
	for i := 0; i < len(words); i += 2 {
		b.add(words[i])
	}

	// Make sure there are no false negatives
	for i := 0; i < len(words); i += 2 {
		word := words[i]
		if !b.maybeContains(word) {
			log.Fatalf("false negative for word %q\n", word)
		}
	}

	falsePositives := 0
	numChecked := 0

	// None of the words at odd indices were added, so whenever
	// maybeContains returns true, it's a false positive
	for i := 1; i < len(words); i += 2 {
		if b.maybeContains(words[i]) {
			falsePositives++
		}
		numChecked++
	}

	falsePositiveRate := float64(falsePositives) / float64(numChecked)

	fmt.Printf("Elapsed time: %s\n", time.Since(start))
	fmt.Printf("Memory usage: %d bytes\n", b.memoryUsage())
	fmt.Printf("False positive rate: %0.2f%%\n", 100*falsePositiveRate)
}

func loadWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type BloomFilter struct {
	data      []byte
	size      int // in bytes
	numHashes int
}

func newBloomFilter(size int, numHashes int) *BloomFilter {
	return &BloomFilter{
		data:      make([]byte, size),
		size:      size,
		numHashes: numHashes,
	}
}

func (b *BloomFilter) add(item string) {
	var bitIdx, targetByteIdx, byteBitIdx int

	for i := 0; i < b.numHashes; i++ {
		// The index of the bit to flip on
		bitIdx = b.hash(append([]byte(item), byte(i))) % (b.size * 8)
		targetByteIdx = bitIdx / 8
		byteBitIdx = bitIdx % 8

		// Set the bit within the target byte (and update the byte slice)
		b.data[targetByteIdx] = setBit(b.data[targetByteIdx], uint(7-byteBitIdx))
	}
}

func (b *BloomFilter) hash(item []byte) int {
	h := fnv.New32a()

	_, err := h.Write(item)
	if err != nil {
		log.Fatal(err)
	}

	return int(h.Sum32())
}

func (b *BloomFilter) maybeContains(item string) bool {
	var bitIdx, targetByteIdx, byteBitIdx int

	for i := 0; i < b.numHashes; i++ {
		// The index of the bit that should be on
		bitIdx = b.hash(append([]byte(item), byte(i))) % (b.size * 8)

		// The index of the byte that contains the bit (that should be on)
		targetByteIdx = bitIdx / 8

		// The index of the bit within the byte (that contains the bit that should be on)
		byteBitIdx = bitIdx % 8

		if !hasBit(b.data[targetByteIdx], uint(7-byteBitIdx)) {
			return false
		}
	}

	return true
}

func (b *BloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}

func setBit(n byte, pos uint) byte {
	n |= 1 << pos
	return n
}

func hasBit(n byte, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

type trivialBloomFilter struct {
	data []uint64
}

func newTrivialBloomFilter(size int) *trivialBloomFilter {
	return &trivialBloomFilter{
		data: make([]uint64, size),
	}
}

func (b *trivialBloomFilter) add(item string) {
	// Do nothing
}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	return true
}

func (b *trivialBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
