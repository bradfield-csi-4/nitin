package main

type selectionOperator struct {
	input     operator
	predicate predicate
}

type predicate struct {
	column string
	value  string
	op     string
}

func (op *selectionOperator) init() {
	op.input.init()
}

func newSelectionNode(input operator, column, value, op string) *selectionOperator {
	return &selectionOperator{
		input: input,
		predicate: predicate{
			column: column,
			value:  value,
			op:     op,
		},
	}
}

func (op *selectionOperator) next() *row {

	for {
		tuple := op.input.next()
		if tuple == nil {
			return nil
		}

		idx := getColumnIndex(*tuple.columns, op.predicate.column)

		if tuple.values[idx] == op.predicate.value {
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

func (op *selectionOperator) close() {
	return
}
