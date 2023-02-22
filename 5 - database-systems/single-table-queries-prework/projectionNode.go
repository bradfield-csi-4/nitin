package main

type projectionNode struct {
	input   iterator
	columns map[string]bool
}

func (n *projectionNode) init() {
	n.input.init()
}

func newProjectionNode(input iterator, columns []string) *projectionNode {
	columnsMap := make(map[string]bool)

	for _, col := range columns {
		columnsMap[col] = true
	}

	return &projectionNode{
		input:   input,
		columns: columnsMap,
	}
}

func (n *projectionNode) next() *tuple {
	inputTuple := n.input.next()
	if inputTuple == nil {
		return nil
	}

	var columns []string
	outputTuple := &tuple{columns: &columns}

	for i, col := range *inputTuple.columns {
		if n.columns[col] {
			*outputTuple.columns = append(*outputTuple.columns, col)
			outputTuple.values = append(outputTuple.values, inputTuple.values[i])
		}
	}

	return outputTuple
}

func remove(items []string, s int) []string {
	return append(items[:s], items[s+1:]...)
}
