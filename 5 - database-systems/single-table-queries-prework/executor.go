package main

import "fmt"

func execute(operator operator) []row {
	var rows []row

	operator.init()

	for {
		row := operator.next()
		if row == nil {
			break
		}
		rows = append(rows, *row)
	}

	return rows
}

func print(rows []row) {
	for _, v := range *(rows[0].columns) {
		fmt.Printf("'%v' ", v)
	}

	for _, row := range rows {
		for _, v := range row.values {
			fmt.Printf("'%v' ", v)
		}
	}
}
