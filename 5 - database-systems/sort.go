package main

import "sort"

type sortOperator struct {
	input     operator
	tuples    []row
	column    string
	direction string
	idx       int
}

func (op *sortOperator) init() {
	op.input.init()

	for {
		tuple := op.input.next()
		if tuple == nil {
			break
		}

		op.tuples = append(op.tuples, *tuple)
	}

	colIdx := getColumnIndex(*op.tuples[0].columns, op.column)

	sort.Slice(op.tuples, func(i, j int) bool {
		if op.direction == "ASC" {
			return op.tuples[i].values[colIdx] < op.tuples[j].values[colIdx]
		} else {
			return op.tuples[i].values[colIdx] > op.tuples[j].values[colIdx]
		}
	})
}

func newSortOperator(input operator, column, direction string) *sortOperator {
	return &sortOperator{
		input:     input,
		column:    column,
		direction: direction,
	}
}

func (op *sortOperator) next() *row {
	if op.idx >= len(op.tuples) {
		return nil
	}

	tuple := op.tuples[op.idx]
	op.idx++
	return &tuple
}

func (op *sortOperator) close() {
	return
}
