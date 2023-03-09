package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initializeTables()
	code := m.Run()
	os.Exit(code)
}

func TestSelectStarFromMovies(t *testing.T) {
	scanNode := newSeqScanOperator("movies")
	results := execute(scanNode)

	// Test
	expectedCols := []string{"id", "title"}
	for i, label := range *results[0].columns {
		if expectedCols[i] != label {
			t.Errorf("SelectStarFromMovies: expected %s, got %s", expectedCols[i], label)
		}
	}

	// Just checking the first and final values, but results should have all records
	var values = [][]string{
		{"1", "Saving Private Ryan"},
		{"11", "Back to the Future"},
	}
	for i, val := range results[0].values {
		if values[0][i] != val {
			t.Errorf("SelectStarFromMovies: expected %s, got %s", values[0][i], val)
		}
	}
	for i, val := range results[len(results)-1].values {
		if values[1][i] != val {
			t.Errorf("SelectStarFromMovies: expected %s, got %s", values[1][i], val)
		}
	}
}

func TestSelectTitleWhereID5FromMovies(t *testing.T) {
	scanNode := newSeqScanOperator("movies")
	selectionNode := newSelectionNode(scanNode, "id", "5", "EQUALS")
	projectionNode := newProjectionOperator(selectionNode, []string{"title"})
	results := execute(projectionNode)

	if len(results) != 1 || len(results[0].values) != 1 || results[0].values[0] != "Enemy of the State" || (*results[0].columns)[0] != "title" {
		t.Error("SelectTitleWhereID5FromMovies: Expected single record, single column with title Enemy of the State")
	}
}

func TestSelectStarWhereTitleSpotlightUsingIndexFromMovies(t *testing.T) {
	indexScanOperator := newIndexOperator("movies", "title", "Spotlight")
	results := execute(indexScanOperator)

	if len(results) != 1 || len(results[0].values) != 2 || results[0].values[1] != "Spotlight" || (*results[0].columns)[1] != "title" {
		t.Error("SelectStarWhereTitleSpotlightUsingIndexFromMovies: Expected single record, single column with title Spotlight")
	}
}

func TestSelectStarLimit8Movies(t *testing.T) {
	limit := 8
	scanNode := newSeqScanOperator("movies")
	limitNode := newLimitOperator(scanNode, limit)
	results := execute(limitNode)

	if len(results) != 8 {
		t.Error(fmt.Sprintf("SelectStarLimit8Movies: expected %v records", limit))
	}
}

func TestSelectFirst3MoviesSortedByTitle(t *testing.T) {
	limit := 3
	scanNode := newSeqScanOperator("movies")
	sortNode := newSortOperator(scanNode, "title", "ASC")
	limitNode := newLimitOperator(sortNode, limit)
	results := execute(limitNode)

	if len(results) != 3 {
		t.Error(fmt.Sprintf("SelectFirst3MoviesSortedByTitle: expected %v records", limit))
	}
	var expectedTuples = [][]string{
		{"6", "3 Idiots"},
		{"11", "Back to the Future"},
		{"5", "Enemy of the State"},
	}

	for i, tuple := range results {
		for j, val := range tuple.values {
			if expectedTuples[i][j] != val {
				t.Errorf("SelectFirst3MoviesSortedByTitle expected %s, got %s", expectedTuples[i][j], val)
			}
		}
	}
}

func TestSelectStarNestedLoopsJoin(t *testing.T) {
	moviesScanOperator := newSeqScanOperator("movies")
	ratingsScanOperator := newSeqScanOperator("ratings")
	joinTable1 := joinTable{
		input:  moviesScanOperator,
		column: "id",
		name:   "movies",
	}

	joinTable2 := joinTable{
		input:  ratingsScanOperator,
		column: "movie_id",
		name:   "ratings",
	}
	joinOperator := newNestedLoopsJoinOperator(joinTable1, joinTable2)
	results := execute(joinOperator)

	// Test
	expectedCols := []string{"id", "title", "movie_id", "rating"}
	for i, label := range *results[0].columns {
		if expectedCols[i] != label {
			t.Errorf("SelectStarNestedLoopsJoin: expected %s, got %s", expectedCols[i], label)
		}
	}

	// Just checking the first and final values, but results should have all records
	var values = [][]string{
		{"1", "Saving Private Ryan", "1", "4.1"},
		{"11", "Back to the Future", "11", "4.2"},
	}

	for i, val := range results[0].values {
		if values[0][i] != val {
			t.Errorf("SelectStarNestedLoopsJoin: expected %s, got %s", values[0][i], val)
		}
	}

	for i, val := range results[len(results)-1].values {
		if values[1][i] != val {
			t.Errorf("SelectStarNestedLoopsJoin: expected %s, got %s", values[1][i], val)
		}
	}
}

func TestSelectStarHashJoin(t *testing.T) {
	moviesScanOperator := newSeqScanOperator("movies")
	ratingsScanOperator := newSeqScanOperator("ratings")
	joinTable1 := joinTable{
		input:  moviesScanOperator,
		column: "id",
		name:   "movies",
	}

	joinTable2 := joinTable{
		input:  ratingsScanOperator,
		column: "movie_id",
		name:   "ratings",
	}
	joinOperator := newHashJoinOperator(joinTable1, joinTable2)
	results := execute(joinOperator)

	// Test
	expectedCols := []string{"id", "title", "movie_id", "rating"}
	for i, label := range *results[0].columns {
		if expectedCols[i] != label {
			t.Errorf("SelectStarHashJoin: expected %s, got %s", expectedCols[i], label)
		}
	}

	// Just checking the first and final values, but results should have all records
	var values = [][]string{
		{"1", "Saving Private Ryan", "1", "4.1"},
		{"11", "Back to the Future", "11", "4.2"},
	}

	for i, val := range results[0].values {
		if values[0][i] != val {
			t.Errorf("SelectStarHashJoin: expected %s, got %s", values[0][i], val)
		}
	}

	for i, val := range results[len(results)-1].values {
		if values[1][i] != val {
			t.Errorf("SelectStarHashJoin: expected %s, got %s", values[1][i], val)
		}
	}
}
