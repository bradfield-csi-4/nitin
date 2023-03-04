package main

type nestedLoopJoinOperator struct {
	input1  operator
	input2  operator
	column1 string
	column2 string
	table1  []*row
	table2  []*row
	output  []*row
	idx     int
}

func newNestedLoopJoinOperator(input1, input2 operator, column1, column2 string) *nestedLoopJoinOperator {
	return &nestedLoopJoinOperator{
		input1:  input1,
		input2:  input2,
		column1: column1,
		column2: column2,
	}
}

func (op *nestedLoopJoinOperator) init() {

	op.input1.init()
	op.input2.init()

	for {
		row := op.input1.next()
		if row == nil {
			break
		}
		op.table1 = append(op.table1, row)
	}

	for {
		row := op.input2.next()
		if row == nil {
			break
		}
		op.table2 = append(op.table2, row)
	}

	column1Idx := getColumnIndex(*op.table1[0].columns, op.column1)
	column2Idx := getColumnIndex(*op.table2[0].columns, op.column2)

	newOutputRowColumns := append(*op.table1[0].columns, *op.table2[0].columns...)

	for _, row1 := range op.table1 {
		for _, row2 := range op.table2 {
			if row1.values[column1Idx] == row2.values[column2Idx] {
				newRow := &row{
					columns: &newOutputRowColumns,
					values:  append(row1.values, row2.values...),
				}

				op.output = append(op.output, newRow)
			}
		}
	}

}

func (op *nestedLoopJoinOperator) next() *row {
	if op.idx >= len(op.output) {
		return nil
	}

	nextRow := op.output[op.idx]
	op.idx++
	return nextRow
}

func (op *nestedLoopJoinOperator) close() {
	return
}
