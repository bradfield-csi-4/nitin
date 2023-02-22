package main

import "fmt"

type iterator interface {
	init()
	next() *tuple
}

type tuple struct {
	columns *[]string
	values  []string
}

func execute(iterator iterator) []tuple {
	var tuples []tuple

	iterator.init()

	for {
		tuple := iterator.next()
		if tuple == nil {
			break
		}
		tuples = append(tuples, *tuple)
	}

	return tuples
}

func print(tuples []tuple) {
	for _, v := range *(tuples[0].columns) {
		fmt.Printf("'%v' ", v)
	}

	for _, tuple := range tuples {
		for _, v := range tuple.values {
			fmt.Printf("'%v' ", v)
		}
	}
}
