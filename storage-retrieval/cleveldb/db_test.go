package cleveldb

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

var db DB
var keys [][]byte

func init() {
	db = getEmptyDB()

	rand.Seed(time.Now().Unix())
	benchmarkSeedSize := 5_000_000

	for i := 0; i < benchmarkSeedSize; i++ {
		key := randStr(20)
		keys = append(keys, key)
		_ = db.Put(key, randStr(300))
	}
}

func getEmptyDB() DB {
	db, err := NewNaiveDB()
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
		return nil
	}
	return db
}

func Test_PutGetReturnCorrectValue(t *testing.T) {
	db = getEmptyDB()

	_ = db.Put([]byte("firstName"), []byte("nitin"))

	var tests = []struct {
		key   string
		value string
	}{
		{"firstName", "nitin"},
		{"lastName", ""},
	}

	for _, test := range tests {
		actualValue, _ := db.Get([]byte(test.key))
		expectedValue := []byte(test.value)

		if bytes.Compare(expectedValue, actualValue) != 0 {
			t.Errorf(`storage.Get("%s") returns unexpected value: "%s"`, test.key, actualValue)
		}
	}
}

func Test_DeleteRemovesValue(t *testing.T) {
	db := getEmptyDB()

	key := []byte("name")
	val := []byte("nitin")

	_ = db.Put(key, val)
	actual, _ := db.Get(key)
	if bytes.Compare(val, actual) != 0 {
		t.Errorf(`storage.Get("%s") returns unexpected value: "%s"`, key, actual)
	}

	_ = db.Delete(key)
	actual, _ = db.Get(key)
	if bytes.Compare([]byte(""), actual) != 0 {
		t.Errorf(`storage.Get("%s") returns unexpected value: "%s"`, key, actual)
	}
}

func Test_RangeScanAndNextReturnCorrectOrderedValues(t *testing.T) {
	db := getEmptyDB()

	keys := [][]byte{[]byte("b"), []byte("c"), []byte("a"), []byte("f"), []byte("d")}
	vals := [][]byte{[]byte("nitin"), []byte("neha"), []byte("cassie"), []byte("karli"), []byte("david")}

	_ = db.Put(keys[0], vals[0])
	_ = db.Put(keys[1], vals[1])
	_ = db.Put(keys[2], vals[2])
	_ = db.Put(keys[3], vals[3])
	_ = db.Put(keys[4], vals[4])

	iter, _ := db.RangeScan([]byte("b"), []byte("d"))

	expectedKeys := [][]byte{[]byte("b"), []byte("c"), []byte("d")}
	expectedVals := [][]byte{[]byte("nitin"), []byte("neha"), []byte("david")}

	for i := 0; i < len(expectedKeys); i++ {
		expectedVal := string(expectedVals[i])
		expectedKey := string(expectedKeys[i])
		actualVal := string(iter.Value())
		actualKey := string(iter.Key())

		if expectedKey != actualKey || expectedVal != actualVal {
			t.Errorf(`storage.RangeScan returns unexpected key/value: "%s: %s" at index: %v`, actualKey, actualVal, i)
		}
		iter.Next()
	}
}

func randStr(length int) []byte {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	byteSlice := make([]byte, length)
	for i := range byteSlice {
		byteSlice[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return byteSlice
}

//
//func Benchmark_Put(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = db.Put(randStr(20), randStr(300))
//	}
//
//	fmt.Printf("Ran %d times\n", b.N)
//}
//
//func Benchmark_Get(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_, _ = db.Get(keys[i])
//	}
//
//	fmt.Printf("Ran %d times\n", b.N)
//}
//
//func Benchmark_Delete(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = db.Delete(keys[i])
//	}
//
//	fmt.Printf("Ran %d times\n", b.N)
//}

func Benchmark_RangeScan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = db.RangeScan([]byte("l"), []byte("p"))
	}

	fmt.Printf("Ran %d times\n", b.N)
}
