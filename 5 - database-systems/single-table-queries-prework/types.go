package main

type operator interface {
	init()
	next() *row
	close()
}

type row struct {
	columns *[]string
	values  []string
}
