package main

type hashJoinOperator struct {
	joinTable1 joinTable
	joinTable2 joinTable
	output     []*row
	idx        int
}

func newHashJoinOperator(joinTable1, joinTable2 joinTable) *hashJoinOperator {
	return &hashJoinOperator{
		joinTable1: joinTable1,
		joinTable2: joinTable2,
	}
}

func (op *hashJoinOperator) init() {
	table1 := op.joinTable1
	table2 := op.joinTable2

	table1.input.init()
	table2.input.init()

	table1Columns := ds.systemCatalog[table1.name].columnNames
	table2Columns := ds.systemCatalog[table2.name].columnNames

	column1Idx := getColumnIndex(table1Columns, table1.column)
	column2Idx := getColumnIndex(table2Columns, table2.column)

	hashMap := make(map[string][]*row)

	// Load first table into in-memory hash table
	for {
		row1 := table1.input.next()
		if row1 == nil {
			break
		}

		keyVal := row1.values[column1Idx]
		hashMap[keyVal] = append(hashMap[keyVal], row1)
	}

	newOutputRowColumns := append(table1Columns, table2Columns...)

	for {
		row2 := table2.input.next()
		if row2 == nil {
			break
		}

		keyVal := row2.values[column1Idx]
		row1s := hashMap[keyVal]

		for _, row1 := range row1s {
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

func (op *hashJoinOperator) next() *row {
	if op.idx >= len(op.output) {
		return nil
	}

	nextRow := op.output[op.idx]
	op.idx++
	return nextRow
}

func (op *hashJoinOperator) close() {
	return
}
