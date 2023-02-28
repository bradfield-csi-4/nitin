package main

func main() {

	selectStarFromMoviesWhereIdEquals2 := &selectionOperator{
		input: &seqScanOperator{
			table: "movies",
			idx:   0,
		},
		predicate: predicate{
			column: "id",
			value:  "2",
			op:     "EQUALS",
		},
	}

	selectIdFromMovies := &projectionOperator{
		input:   selectStarFromMoviesWhereIdEquals2,
		columns: map[string]bool{"id": true},
	}

	rows := execute(selectIdFromMovies)

	print(rows)
}
