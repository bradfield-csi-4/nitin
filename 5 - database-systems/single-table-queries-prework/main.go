package main

func main() {

	selectStarFromMoviesWhereIdEquals2 := &selectionNode{
		input: &scanNode{
			table: "movies",
			idx:   0,
		},
		predicate: predicate{
			column: "id",
			value:  "2",
			op:     "EQUALS",
		},
	}

	selectIdFromMovies := &projectionNode{
		input:   selectStarFromMoviesWhereIdEquals2,
		columns: map[string]bool{"id": true},
	}

	tuples := execute(selectIdFromMovies)
	
	print(tuples)
}
