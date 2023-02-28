package main

type projectionOperator struct {
	input   operator
	columns map[string]bool
}

func (op *projectionOperator) init() {
	op.input.init()
}

func newProjectionOperator(input operator, columns []string) *projectionOperator {
	columnsMap := make(map[string]bool)

	for _, col := range columns {
		columnsMap[col] = true
	}

	return &projectionOperator{
		input:   input,
		columns: columnsMap,
	}
}

func (op *projectionOperator) next() *row {
	inputTuple := op.input.next()
	if inputTuple == nil {
		return nil
	}

	var columns []string
	outputTuple := &row{columns: &columns}

	for i, col := range *inputTuple.columns {
		if op.columns[col] {
			*outputTuple.columns = append(*outputTuple.columns, col)
			outputTuple.values = append(outputTuple.values, inputTuple.values[i])
		}
	}

	return outputTuple
}

func remove(items []string, s int) []string {
	return append(items[:s], items[s+1:]...)
}

func (op *projectionOperator) close() {
	return
}
