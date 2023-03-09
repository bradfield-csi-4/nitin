package main

type nestedLoopsJoinOperator struct {
	joinTable1 joinTable
	joinTable2 joinTable
	output     []*row
	idx        int
}

type joinTable struct {
	input  operator
	column string
	name   string
	data   []*row
}

func newNestedLoopsJoinOperator(joinTable1, joinTable2 joinTable) *nestedLoopsJoinOperator {
	return &nestedLoopsJoinOperator{
		joinTable1: joinTable1,
		joinTable2: joinTable2,
	}
}

func (op *nestedLoopsJoinOperator) init() {
	table1 := op.joinTable1
	table2 := op.joinTable2

	table1.input.init()
	table2.input.init()

	table1Columns := ds.systemCatalog[table1.name].columnNames
	table2Columns := ds.systemCatalog[table2.name].columnNames

	column1Idx := getColumnIndex(table1Columns, table1.column)
	column2Idx := getColumnIndex(table2Columns, table2.column)

	for {
		row := table1.input.next()
		if row == nil {
			break
		}
		table1.data = append(table1.data, row)
	}

	for {
		row := table2.input.next()
		if row == nil {
			break
		}
		table2.data = append(table2.data, row)
	}

	newOutputRowColumns := append(table1Columns, table2Columns...)

	for _, row1 := range table1.data {

		for _, row2 := range table2.data {
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

func (op *nestedLoopsJoinOperator) next() *row {
	if op.idx >= len(op.output) {
		return nil
	}

	nextRow := op.output[op.idx]
	op.idx++
	return nextRow
}

func (op *nestedLoopsJoinOperator) close() {
	return
}
