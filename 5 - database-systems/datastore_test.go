package main

import (
	"fmt"
	"testing"
)

func TestWriteReadFromMovies(t *testing.T) {
	movies := getMovies()

	initializeTables()

	nextOffset := 0

	for i := 0; i < len(movies); i++ {
		row, bytesRead, err := ds.readRow("movies", nextOffset)
		if err != nil {
			fmt.Println(err)
			return
		}

		nextOffset += bytesRead

		m := getMovies()[i]

		if m.id != row[0] && m.title != row[1] {
			t.Errorf("TestWriteReadFromMovies - expected %s, got %s", m, row)
		}
	}
}
