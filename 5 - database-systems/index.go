package main

import (
	"log"
	"sort"
)

type indexOperator struct {
	table       string
	predicate   predicate
	index       *index
	columnNames []string
	matches     []*indexEntry
	matchIdx    int
}

func (op *indexOperator) init() {
	op.index = ds.indexes[op.table]
	op.columnNames = ds.systemCatalog[op.table].columnNames

	values := op.index.values

	// binary search to find index of first match
	idx := sort.Search(len(values), func(i int) bool {
		return values[i].value >= op.predicate.value
	})

	// Scan from first match until there's mismatch
	for _, match := range values[idx:] {
		if match.value != op.predicate.value {
			break
		}
		op.matches = append(op.matches, match)
	}
}

func newIndexOperator(table string, column, value string) *indexOperator {
	return &indexOperator{
		table: table,
		predicate: predicate{
			column: column,
			value:  value,
		},
	}
}

func (op *indexOperator) next() *row {
	if op.matchIdx == len(op.matches) {
		return nil
	}

	targetValue := op.matches[op.matchIdx]

	// Use offset to read from heap file
	rowValues, _, err := ds.readRow(op.table, targetValue.offset)
	if err != nil {
		log.Fatal("Error reading row while using index scan operator", err)
	}

	op.matchIdx++

	return &row{
		columns: &op.columnNames,
		values:  rowValues,
	}
}

func (op *indexOperator) close() {
	return
}
