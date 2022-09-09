package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const url = "https://xkcd.com/"
const subPath = "/info.0.json"
const maxID = 10

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

type Index map[string][]int

type ComicData map[int]ComicInfoResult

type ComicInfoResult struct {
	Num        int
	Title      string
	Transcript string
}

func SearchXKCD(searchTerm string) error {
	index := make(Index)
	data := make(ComicData)

	indexFile, _ := os.OpenFile("index.json", os.O_RDONLY, os.ModePerm)
	defer indexFile.Close()
	err := json.NewDecoder(indexFile).Decode(&index)
	if err != nil {
		return fmt.Errorf("json decoding failed: %v", err)
	}

	dataFile, _ := os.OpenFile("data.json", os.O_RDONLY, os.ModePerm)
	defer dataFile.Close()
	err = json.NewDecoder(dataFile).Decode(&data)
	if err != nil {
		return fmt.Errorf("json decoding failed: %v", err)
	}

	ids := index[searchTerm]

	for _, id := range ids {
		result := data[id]
		fmt.Printf("\n%v -- %v\n", result.Num, result.Title)
		fmt.Printf("%v\n", result.Transcript)
		fmt.Printf("%v\n", getFullURL(id))
	}

	return nil
}

func StoreDataToFile() error {
	results, err := retrieveResultsFromAPI()
	if err != nil {
		return err
	}

	data := buildDataMapFromResults(results)

	err = writeDataToFile(data)
	if err != nil {
		return err
	}

	index := buildIndexFromResults(results)

	err = writeIndexToFile(index)
	if err != nil {
		return err
	}

	return nil
}

func writeIndexToFile(index *Index) error {
	file, _ := os.OpenFile("index.json", os.O_WRONLY, os.ModePerm)
	defer file.Close()
	err := json.NewEncoder(file).Encode(*index)
	if err != nil {
		return fmt.Errorf("json encoding failed: %v", err)
	}
	return nil
}

func writeDataToFile(data *ComicData) error {
	file, _ := os.OpenFile("data.json", os.O_WRONLY, os.ModePerm)
	defer file.Close()
	err := json.NewEncoder(file).Encode(data)
	if err != nil {
		return fmt.Errorf("json encoding failed: %v", err)
	}
	return nil
}

func retrieveResultsFromAPI() ([]ComicInfoResult, error) {
	var results []ComicInfoResult

	for id := 1; id <= maxID; id++ {
		result, err := getResultById(id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		results = append(results, *result)
	}
	return results, nil
}

func buildIndexFromResults(results []ComicInfoResult) *Index {
	index := make(Index)

	for _, result := range results {
		for _, word := range strings.Split(result.Title, " ") {
			word := clearString(word)
			index[word] = append(index[word], result.Num)
		}
	}
	return &index
}

func buildDataMapFromResults(results []ComicInfoResult) *ComicData {
	data := make(ComicData)

	for _, result := range results {
		data[result.Num] = result
	}
	return &data
}

func getResultById(id int) (*ComicInfoResult, error) {
	resp, err := http.Get(getFullURL(id))
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("http request failed for id: %d err: %v", id, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s id: %d", resp.Status, id)
	}

	var result ComicInfoResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("json decoding failed for id: %d err: %v", id, err)
	}

	return &result, nil
}

func getFullURL(id int) string {
	return url + strconv.Itoa(id) + subPath
}

func clearString(str string) string {
	return strings.ToLower(nonAlphanumericRegex.ReplaceAllString(str, ""))
}
