package main

import (
	"fmt"
	"testing"
)

func TestWriteReadFromMovies(t *testing.T) {
	initializeMoviesTable()

	nextOffset := 0

	for i := 0; i < len(getMovies()); i++ {
		row, bytesRead, err := ds.readRow("movies", nextOffset)
		if err != nil {
			fmt.Println(err)
			return
		}

		nextOffset += bytesRead

		m := movie{row[0], row[1]}

		if getMovies()[i] != m {
			t.Errorf("TestWriteReadFromMovies - expected %s, got %s", getMovies()[0], m)
		}
	}
}
