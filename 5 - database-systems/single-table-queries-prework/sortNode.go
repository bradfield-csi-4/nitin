package main

import "sort"

type sortNode struct {
	input     iterator
	tuples    []tuple
	column    string
	direction string
	idx       int
}

func (n *sortNode) init() {

	input := n.input
	input.init()

	for {
		tuple := input.next()
		if tuple == nil {
			break
		}

		n.tuples = append(n.tuples, *tuple)
	}

	colIdx := getColumnIndex(*n.tuples[0].columns, n.column)

	sort.Slice(n.tuples, func(i, j int) bool {
		if n.direction == "ASC" {
			return n.tuples[i].values[colIdx] < n.tuples[j].values[colIdx]
		} else {
			return n.tuples[i].values[colIdx] > n.tuples[j].values[colIdx]
		}
	})
}

func newSortNode(input iterator, column, direction string) *sortNode {
	return &sortNode{
		input:     input,
		column:    column,
		direction: direction,
	}
}

func (n *sortNode) next() *tuple {
	if n.idx >= len(n.tuples) {
		return nil
	}

	tuple := n.tuples[n.idx]
	n.idx++
	return &tuple
}
