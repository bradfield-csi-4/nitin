package main

type selectionNode struct {
	input     iterator
	predicate predicate
}

type predicate struct {
	column string
	value  string
	op     string
}

func (n *selectionNode) init() {}

func newSelectionNode(input iterator, column, value, op string) *selectionNode {
	return &selectionNode{
		input: input,
		predicate: predicate{
			column: column,
			value:  value,
			op:     op,
		},
	}
}

func (n *selectionNode) next() *tuple {
	for {
		input := n.input
		input.init()
		tuple := input.next()
		if tuple == nil {
			return nil
		}

		idx := getColumnIndex(*tuple.columns, n.predicate.column)

		if tuple.values[idx] == n.predicate.value {
			return tuple
		}
	}
}

func getColumnIndex(columns []string, targetColumn string) int {
	for i, col := range columns {
		if col == targetColumn {
			return i
		}
	}

	return -1
}
