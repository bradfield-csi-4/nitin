package main

import "fmt"

type movie struct {
	id, title string
}

type rating struct {
	movieId, rating string
}

func initializeTables() {
	var err error
	movies := getMovies()
	ratings := getRatings()

	ds, err = newDatastore()
	if err != nil {
		fmt.Println(err)
		return
	}

	ds.createTable("movies", []string{"id", "title"})

	for _, m := range movies {
		err = ds.appendRow("movies", []string{m.id, m.title})
		if err != nil {
			fmt.Println(err)
		}
	}

	ds.createIndex("movies", "title")

	ds.createTable("ratings", []string{"movie_id", "rating"})

	for _, r := range ratings {
		err = ds.appendRow("ratings", []string{r.movieId, r.rating})
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

func getRatings() []rating {
	return []rating{
		{"1", "4.1"},
		{"2", "4.4"},
		{"3", "4.9"},
		{"4", "4.3"},
		{"5", "4.2"},
		{"6", "4.4"},
		{"7", "4.6"},
		{"8", "3.5"},
		{"9", "3.9"},
		{"10", "4.5"},
		{"11", "4.2"},
		{"1", "4.8"},
		{"2", "4.6"},
		{"3", "4.1"},
		{"4", "4.0"},
		{"5", "4.0"},
		{"6", "3.2"},
		{"7", "3.5"},
		{"8", "4.7"},
		{"9", "4.0"},
		{"10", "4.4"},
		{"11", "4.6"},
		{"1", "4.2"},
		{"2", "4.1"},
		{"3", "4.0"},
		{"4", "4.5"},
		{"5", "4.7"},
		{"6", "4.7"},
		{"7", "3.6"},
		{"8", "3.9"},
		{"9", "4.0"},
		{"10", "4.8"},
		{"11", "4.2"},
	}
}
