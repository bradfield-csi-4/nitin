package main

type limitOperator struct {
	input operator
	limit int
	count int
}

func (op *limitOperator) init() {
	op.input.init()
}

func newLimitOperator(input operator, limit int) *limitOperator {
	return &limitOperator{
		input: input,
		limit: limit,
	}
}

func (op *limitOperator) next() *row {
	tuple := op.input.next()
	if tuple == nil || op.count >= op.limit {
		return nil
	}

	op.count++
	return tuple
}

func (op *limitOperator) close() {
	return
}
