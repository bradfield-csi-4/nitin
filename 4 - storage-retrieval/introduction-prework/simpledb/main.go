package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var dbFilename = "db.txt"

func main() {
	action := os.Args[1]
	if action == "set" {
		set(os.Args[2], os.Args[3])
	} else if action == "get" {
		get(os.Args[2])
	}
}

var dbFile *os.File

func init() {
	dbFile, _ = os.OpenFile(dbFilename, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
}

func get(key string) {
	var val string

	scanner := bufio.NewScanner(dbFile)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, ",")
		iKey := splitLine[0]
		iVal := splitLine[1]
		if key == iKey {
			val = iVal
		}
	}

	fmt.Println(val)
}

func set(key, val string) {
	_, err := dbFile.WriteString(fmt.Sprintf("%s,%s\n", key, val))
	if err != nil {
		fmt.Printf("Error setting key: %v\n", err)
	}
}
