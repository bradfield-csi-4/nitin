package main

type limitNode struct {
	input iterator
	limit int
	count int
}

func (n *limitNode) init() {
	n.input.init()
}

func newLimitNode(input iterator, limit int) *limitNode {
	return &limitNode{
		input: input,
		limit: limit,
	}
}

func (n *limitNode) next() *tuple {
	tuple := n.input.next()
	if tuple == nil || n.count >= n.limit {
		return nil
	}

	n.count++
	return tuple
}
