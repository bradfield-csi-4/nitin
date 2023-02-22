package main

type scanNode struct {
	tuples []tuple
	table  string
	idx    int
}

func (n *scanNode) init() {
	n.tuples = datastore[n.table]
}

func newScanNode(table string) *scanNode {
	return &scanNode{
		table: table,
		idx:   0,
	}
}

func (n *scanNode) next() *tuple {
	if n.idx >= len(n.tuples) {
		return nil
	}

	tuple := n.tuples[n.idx]
	n.idx++
	return &tuple
}
