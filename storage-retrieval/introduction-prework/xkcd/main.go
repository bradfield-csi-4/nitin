package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
The popular web comic _xkcd_ has a JSON interface. For example, a request to
[https://xkcd.com/571/info.0.json](https://xkcd.com/571/info.0.json) produces a detailed description of comic 571,
one of many favorites. Download each URL (once!) and build an offline index. Write a tool `xkcd` that, using this index,
prints the URL and transcript of each comic that matches a search term provided on the command line.
*/

var comicsToIndex = 571
var xkcdHost = "https://xkcd.com/"
var xkcdJsonPath = "/info.0.json"

var titleWordIndex inMemoryStringToIntsIndex
var idIndex inMemoryHeapIndex

func main() {
	//buildOnDiskIndexes()

	search("(sketch)")
}

func init() {
	gob.Register(xkcdResponse{})

	if titleWordIndex.Store == nil {
		titleWordIndex.Store = make(map[string][]int, 10000)
	}
	if idIndex.Store == nil {
		idIndex.Store = make(map[string]interface{}, 10000)
	}
}

type xkcdResponse struct {
	Num        int
	Title      string
	Transcript string
}

type inMemoryHeapIndex struct {
	Store map[string]interface{}
}

type inMemoryStringToIntsIndex struct {
	Store map[string][]int
}

func (idx *inMemoryHeapIndex) get(key string) interface{} {
	return idx.Store[key]
}

func (idx *inMemoryHeapIndex) set(key string, val interface{}) {
	idx.Store[key] = val
}

func (idx *inMemoryStringToIntsIndex) get(key string) []int {
	return idx.Store[key]
}

func (idx *inMemoryStringToIntsIndex) set(key string, val int) {
	idx.Store[key] = append(idx.Store[key], val)
}

func buildOnDiskIndexes() {
	responseBuffer := make([]byte, 10000)

	for i := 1; i <= comicsToIndex; i++ {
		resp, err := http.Get(xkcdHost + strconv.Itoa(i) + xkcdJsonPath)
		if err != nil {
			fmt.Printf("Error making GET request: %v\n", err)
			return
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Error response from GET request with Status for ID %v: %v\n", i, resp.Status)
			continue
		}

		for {
			nBytes, err := resp.Body.Read(responseBuffer)
			if nBytes > 0 {
				var respObj xkcdResponse
				err := json.Unmarshal(responseBuffer[:nBytes], &respObj)
				if err != nil {
					fmt.Printf("Error unmarshalling JSON for ID %v: %v\n", i, err)
					continue
				}

				id := respObj.Num

				for _, word := range strings.Split(respObj.Title, " ") {
					titleWordIndex.set(strings.ToLower(word), id)
				}

				idIndex.set(strconv.Itoa(id), respObj)
			}
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}
		}

		if i%20 == 0 {
			fmt.Printf("Finished %v comics\n", i)
		}
	}

	err := storeIndexToFile(titleWordIndex, "title_words.gob")
	if err != nil {
		panic(err)
	}
	err = storeIndexToFile(idIndex, "ids.gob")
	if err != nil {
		panic(err)
	}
}

func storeIndexToFile(index interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(index); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func search(query string) {
	var titleWordIndex inMemoryStringToIntsIndex
	var idIndex inMemoryHeapIndex

	err := loadIndexFromFile(&titleWordIndex, "title_words.gob")
	if err != nil {
		panic(err)
	}

	err = loadIndexFromFile(&idIndex, "ids.gob")
	if err != nil {
		panic(err)
	}

	matchingIds := titleWordIndex.Store[query]

	for _, id := range matchingIds {
		comicData := idIndex.Store[strconv.Itoa(id)].(xkcdResponse)
		fmt.Println(xkcdHost + strconv.Itoa(id))
		fmt.Println(comicData.Transcript)
		fmt.Println()
	}
}

func loadIndexFromFile(index interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(file)

	switch index := index.(type) {
	case *inMemoryStringToIntsIndex:
		err = decoder.Decode(index)
		if err != nil {
			return err
		}
	case *inMemoryHeapIndex:
		err = decoder.Decode(index)
		if err != nil {
			return err
		}
	}

	return nil
}
