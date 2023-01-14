package cleveldb

import (
	"log"
)

type ClevelDB struct {
}

func NewClevelDB() (*ClevelDB, error) {
	return &ClevelDB{}, nil
}

func (clevelDb *ClevelDB) Get(key []byte) ([]byte, error) {
	log.Fatal("not yet implemented")
	return nil, nil
}

func (clevelDb *ClevelDB) Put(key, value []byte) error {
	log.Fatal("not yet implemented")
	return nil

}

func (clevelDb *ClevelDB) Delete(key []byte) error {
	log.Fatal("not yet implemented")
	return nil
}

func (clevelDb *ClevelDB) RangeScan(start, limit []byte) (Iterator, error) {
	log.Fatal("not yet implemented")
	return nil, nil
}

type ClevelIterator struct {
}

func (i *ClevelIterator) Next() bool {
	log.Fatal("not yet implemented")
	return false
}

func (i *ClevelIterator) Error() error {
	log.Fatal("not yet implemented")
	return nil
}

func (i *ClevelIterator) Key() []byte {
	log.Fatal("not yet implemented")
	return nil
}

func (i *ClevelIterator) Value() []byte {
	log.Fatal("not yet implemented")
	return nil
}
