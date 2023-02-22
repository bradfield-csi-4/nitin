package main

type limitNode struct {
	input iterator
	limit int
	count int
}

func (n *limitNode) init() {}

func newLimitNode(input iterator, limit int) *limitNode {
	return &limitNode{
		input: input,
		limit: limit,
	}
}

func (n *limitNode) next() *tuple {
	for {
		input := n.input
		input.init()
		tuple := input.next()
		if tuple == nil || n.count >= n.limit {
			return nil
		}

		n.count++
		return tuple
	}
}
