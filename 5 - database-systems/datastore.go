package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

var ds *datastore

func init() {
	var err error

	ds, err = newDatastore()
	if err != nil {
		fmt.Println(err)
	}
}

type tableMetadata struct {
	name            string
	columnNames     []string
	file            *os.File
	freeSpaceOffset int
}

type index struct {
	key    string
	values []*indexEntry
}

type indexEntry struct {
	value  string
	offset int
}

type datastore struct {
	systemCatalog map[string]*tableMetadata
	indexes       map[string]*index
}

func newDatastore() (*datastore, error) {
	return &datastore{
		systemCatalog: make(map[string]*tableMetadata),
		indexes:       make(map[string]*index),
	}, nil
}

func (ds *datastore) createTable(name string, columns []string) {
	file, err := os.OpenFile(fmt.Sprintf("storage/%v.db", name), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating file for new table", err)
	}

	table := &tableMetadata{
		name:            name,
		file:            file,
		columnNames:     columns,
		freeSpaceOffset: 0,
	}

	table.appendHeader()

	ds.systemCatalog[name] = table
}

func (ds *datastore) createIndex(table string, key string) {
	//file, err := os.OpenFile(fmt.Sprintf("storage/%v.idx", name), os.O_RDWR|os.O_CREATE, os.ModePerm)
	//if err != nil {
	//	log.Fatal("Error creating file for new index", err)
	//}

	colIdx := getIndex(ds.systemCatalog[table].columnNames, key)

	nextOffset := 0
	var bytesRead int
	var err error

	var indexValues []*indexEntry
	var row []string

	for {
		row, bytesRead, err = ds.readRow(table, nextOffset)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading row while creating index", err)
		}

		nextOffset += bytesRead

		indexValues = append(indexValues, &indexEntry{
			value:  row[colIdx],
			offset: nextOffset - bytesRead,
		})
	}

	sort.Slice(indexValues, func(i, j int) bool {
		return indexValues[i].value < indexValues[j].value
	})

	ds.indexes[table] = &index{
		key:    key,
		values: indexValues,
	}
}

func getIndex(strings []string, targetString string) int {
	for i, str := range strings {
		if str == targetString {
			return i
		}
	}
	return -1
}

func (table *tableMetadata) appendHeader() {
	// TODO: Implement this to enable persistence (i.e. the ability to load a table from disk)
}

func (ds *datastore) appendRow(table string, row []string) error {
	tableMetadata := ds.systemCatalog[table]
	file := tableMetadata.file
	columns := tableMetadata.columnNames

	// Seek to the offset where free space starts
	_, err := file.Seek(int64(tableMetadata.freeSpaceOffset), io.SeekStart)
	if err != nil {
		log.Fatal("Error seeking to start of file", err)
	}

	if len(columns) != len(row) {
		return errors.New("invalid number of columns")
	}

	var rowBytes []byte

	for _, field := range row {
		fieldLen := uint32(len(field))

		rowBytes = binary.BigEndian.AppendUint32(rowBytes, fieldLen)
		rowBytes = append(rowBytes, field...)
	}

	_, err = file.Write(rowBytes)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	// Get/set the new offset where "free space" starts
	offset, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Fatal("Error seeking to start of file", err)
	}

	tableMetadata.freeSpaceOffset = int(offset)

	return nil
}

// readRow : reads single row from the specified table
func (ds *datastore) readRow(table string, offset int) ([]string, int, error) {
	tableMetadata := ds.systemCatalog[table]
	file := tableMetadata.file
	columns := tableMetadata.columnNames

	_, err := file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		log.Fatal("Error seeking to start of file", err)
	}

	var row []string

	var bytesRead int

	for i := 0; i < len(columns); i++ {
		fieldLenBytes := make([]byte, 4)

		n1, err := file.Read(fieldLenBytes)
		if err != nil {
			return nil, 0, err
		}

		fieldLen := binary.BigEndian.Uint32(fieldLenBytes)

		fieldBytes := make([]byte, fieldLen)

		n2, err := file.Read(fieldBytes)
		if err != nil {
			return nil, 0, err
		}

		bytesRead += n1 + n2

		row = append(row, string(fieldBytes))
	}

	return row, bytesRead, nil
}
