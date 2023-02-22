package main

import "strconv"

var datastore map[string][]tuple

func init() {
	datastore = make(map[string][]tuple)
	datastore["movies"] = loadMoviesTable()
}

func loadMoviesTable() []tuple {
	titles := []string{
		"Saving Private Ryan",
		"Shawshank Redemption",
		"Nick & Norah's Infinite Playlist",
		"Good Will Hunting",
		"Enemy of the State",
		"3 Idiots",
		"Spotlight",
		"Long Weekend",
		"The Big Sick",
		"Truman Show",
		"Back to the Future",
	}

	var movies []tuple

	cols := []string{"id", "title"}

	for i, title := range titles {
		movieTuple := tuple{
			columns: &cols,
			values:  []string{strconv.Itoa(i + 1), title},
		}
		movies = append(movies, movieTuple)
	}

	return movies
}
