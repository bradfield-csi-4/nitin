package main

import (
	"errors"
)

const treeOrder = 6

// B+ Tree Implementation
type treeNode struct {
	isLeaf bool
	parent *treeNode
	keys   []string // contains n keys

	// children exist only if isLeaf = false
	children []*treeNode // contains space for up to n + 1 children

	// bottom 3 fields exist (only if isLeaf = true)
	dataPointers []*dataPointer // contains n pointers (corresponding to n keys)
	prev         *treeNode
	next         *treeNode
}

type dataPointer struct {
	fileOffset int // currently just an offset to a single table file (i.e. no pages)
}

func (n *treeNode) insert(newKey string, newDataPointer *dataPointer) {
	node := n.findLeaf(newKey)

	node.insertAtLeaf(newKey, newDataPointer)

	if len(node.keys) >= treeOrder {
		// TODO: split leaf
	}
	return
}

func (n *treeNode) insertAtLeaf(newKey string, newDataPointer *dataPointer) {
	var targetIdx int

	// TODO: Use binary search instead
	for i, key := range n.keys {
		if newKey < key {
			targetIdx = i
			break
		}

		// Reached end of keys without finding a match, so we can just append to slices
		if i == len(n.keys) {
			n.keys = append(n.keys, newKey)
			n.dataPointers = append(n.dataPointers, newDataPointer)
			return
		}
	}

	// Create space for the new key/data at the target index
	n.keys = append(n.keys[:targetIdx+1], n.keys[targetIdx:]...)
	n.dataPointers = append(n.dataPointers[:targetIdx+1], n.dataPointers[targetIdx:]...)

	// Set the target index value with the new key/data
	n.keys[targetIdx] = newKey
	n.dataPointers[targetIdx] = newDataPointer
}

func (n *treeNode) find(target string) (*dataPointer, error) {
	var i int
	c := n.findLeaf(target)

	for i = 0; i < len(c.keys); i++ {
		if c.keys[i] == target {
			break
		}
	}

	// Reached end of keys without finding a match
	if i == len(c.keys) {
		return nil, errors.New("key not found")
	}

	return c.dataPointers[i], nil

}

// Find leaf node where the target key should reside
func (n *treeNode) findLeaf(target string) *treeNode {
	curr := n

	// Starting at the root, iteratively search each level until you reach a leaf node
	for !curr.isLeaf {
		// TODO: Use binary search instead
		for i, key := range curr.keys {
			// Iterate through sorted keys until we find the first key that the target is less than
			// And then we follow the pointer immediately to the left of that key (which will have the same index) in children
			if target < key {
				curr = n.children[i]
				break
			}

			// This indicates the target was greater than all keys, so we follow the last item in children
			if i == len(curr.keys) {
				curr = n.children[i+1]
			}
		}
	}

	return curr
}
