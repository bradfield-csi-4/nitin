package main

import (
	"io"
	"log"
)

type seqScanOperator struct {
	rows  []row
	table string
	idx   int
}

func newSeqScanOperator(table string) *seqScanOperator {
	return &seqScanOperator{
		table: table,
		idx:   0,
	}
}

func (op *seqScanOperator) init() {
	var err error
	catalogEntry := ds.systemCatalog[op.table]

	nextOffset := 0
	var bytesRead int

	var rows []row

	for {
		var row row

		row.values, bytesRead, err = ds.readRow(op.table, nextOffset)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading row while initializing scan operator", err)
		}

		nextOffset += bytesRead

		row.columns = &catalogEntry.columnNames
		rows = append(rows, row)
	}

	op.rows = rows
}

func (op *seqScanOperator) next() *row {
	if op.idx >= len(op.rows) {
		return nil
	}

	row := op.rows[op.idx]
	op.idx++
	return &row
}

func (op *seqScanOperator) close() {
	return
}
