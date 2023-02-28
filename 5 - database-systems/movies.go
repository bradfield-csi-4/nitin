package main

import "fmt"

type movie struct {
	id, title string
}

func createAndPopulateMoviesTable() {
	var err error
	movies := getMovies()

	ds, err = newDatastore()
	if err != nil {
		fmt.Println(err)
		return
	}

	ds.createTable("movies", []string{"id", "title"})

	for _, movie := range movies {
		err = ds.appendRow("movies", []string{movie.id, movie.title})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getMovies() []movie {
	return []movie{
		{"1", "Saving Private Ryan"},
		{"2", "Shawshank Redemption"},
		{"3", "Nick & Norah's Infinite Playlist"},
		{"4", "Good Will Hunting"},
		{"5", "Enemy of the State"},
		{"6", "3 Idiots"},
		{"7", "Spotlight"},
		{"8", "Long Weekend"},
		{"9", "The Big Sick"},
		{"10", "Truman Show"},
		{"11", "Back to the Future"},
	}
}
