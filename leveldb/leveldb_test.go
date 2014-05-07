package leveldb

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
)

var testConfigJson = []byte(`
    {
        "path" : "./testdb",
        "compression":true,
        "block_size" : 32768,
        "write_buffer_size" : 2097152,
        "cache_size" : 20971520
    }
    `)

var testOnce sync.Once
var testDB *DB

func getTestDB() *DB {
	f := func() {
		var err error
		testDB, err = Open(testConfigJson)
		if err != nil {
			println(err.Error())
			panic(err)
		}
	}

	testOnce.Do(f)
	return testDB
}

func TestSimple(t *testing.T) {
	db := getTestDB()

	key := []byte("key")
	value := []byte("hello world")
	if err := db.Put(key, value); err != nil {
		t.Fatal(err)
	}

	if v, err := db.Get(key); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(v, value) {
		t.Fatal("not equal")
	}

	if err := db.Delete(key); err != nil {
		t.Fatal(err)
	}

	if v, err := db.Get(key); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Fatal("must nil")
	}
}

func TestBatch(t *testing.T) {
	db := getTestDB()

	key1 := []byte("key1")
	key2 := []byte("key2")

	value := []byte("hello world")

	db.Put(key1, value)
	db.Put(key2, value)

	wb := db.NewWriteBatch()
	defer wb.Close()

	wb.Delete(key2)
	wb.Put(key1, []byte("hello world2"))

	if err := wb.Commit(); err != nil {
		t.Fatal(err)
	}

	if v, err := db.Get(key2); err != nil {
		t.Fatal(err)
	} else if v != nil {
		t.Fatal("must nil")
	}

	if v, err := db.Get(key1); err != nil {
		t.Fatal(err)
	} else if string(v) != "hello world2" {
		t.Fatal(string(v))
	}

	wb.Delete(key1)

	wb.Rollback()

	if v, err := db.Get(key1); err != nil {
		t.Fatal(err)
	} else if string(v) != "hello world2" {
		t.Fatal(string(v))
	}

	db.Delete(key1)
}

func TestIterator(t *testing.T) {
	db := getTestDB()
	for it := db.Iterator(nil, nil, 0); it.Valid(); it.Next() {
		db.Delete(it.Key())
	}

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		db.Put(key, value)
	}

	step := 0
	var it *Iterator
	for it = db.Iterator(nil, nil, 0); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()

		if string(key) != fmt.Sprintf("key_%d", step) {
			t.Fatal(string(key), step)
		}

		if string(value) != fmt.Sprintf("value_%d", step) {
			t.Fatal(string(value), step)
		}

		step++
	}

	it.Close()

	step = 2
	for it = db.Iterator([]byte("key_2"), nil, 3); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()

		if string(key) != fmt.Sprintf("key_%d", step) {
			t.Fatal(string(key), step)
		}

		if string(value) != fmt.Sprintf("value_%d", step) {
			t.Fatal(string(value), step)
		}

		step++
	}
	it.Close()

	if step != 5 {
		t.Fatal("invalid step", step)
	}

	step = 2
	for it = db.Iterator([]byte("key_2"), []byte("key_5"), 0); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()

		if string(key) != fmt.Sprintf("key_%d", step) {
			t.Fatal(string(key), step)
		}

		if string(value) != fmt.Sprintf("value_%d", step) {
			t.Fatal(string(value), step)
		}

		step++
	}
	it.Close()

	if step != 6 {
		t.Fatal("invalid step", step)
	}

	step = 2
	for it = db.Iterator([]byte("key_5"), []byte("key_2"), 0); it.Valid(); it.Next() {
		step++
	}
	it.Close()

	if step != 2 {
		t.Fatal("must 0")
	}

	step = 9
	for it = db.ReverseIterator(nil, nil, 0); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()

		if string(key) != fmt.Sprintf("key_%d", step) {
			t.Fatal(string(key), step)
		}

		if string(value) != fmt.Sprintf("value_%d", step) {
			t.Fatal(string(value), step)
		}

		step--
	}
	it.Close()

	step = 5
	for it = db.ReverseIterator([]byte("key_5"), nil, 3); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()

		if string(key) != fmt.Sprintf("key_%d", step) {
			t.Fatal(string(key), step)
		}

		if string(value) != fmt.Sprintf("value_%d", step) {
			t.Fatal(string(value), step)
		}

		step--
	}
	it.Close()

	if step != 2 {
		t.Fatal("invalid step", step)
	}

	step = 5
	for it = db.ReverseIterator([]byte("key_5"), []byte("key_2"), 0); it.Valid(); it.Next() {
		key := it.Key()
		value := it.Value()

		if string(key) != fmt.Sprintf("key_%d", step) {
			t.Fatal(string(key), step)
		}

		if string(value) != fmt.Sprintf("value_%d", step) {
			t.Fatal(string(value), step)
		}

		step--
	}
	it.Close()

	if step != 1 {
		t.Fatal("invalid step", step)
	}

	step = 5
	for it = db.ReverseIterator([]byte("key_2"), []byte("key_5"), 0); it.Valid(); it.Next() {
		step--
	}
	it.Close()

	if step != 5 {
		t.Fatal("must 5")
	}
}

func TestIterator_2(t *testing.T) {
	db := getTestDB()
	for it := db.Iterator(nil, nil, 0); it.Valid(); it.Next() {
		db.Delete(it.Key())
	}

	db.Put([]byte("key_1"), []byte("value_1"))
	db.Put([]byte("key_7"), []byte("value_9"))
	db.Put([]byte("key_9"), []byte("value_9"))

	it := db.Iterator([]byte("key_0"), []byte("key_8"), 0)
	if !it.Valid() {
		t.Fatal("must valid")
	}

	if string(it.Key()) != "key_1" {
		t.Fatal(string(it.Key()))
	}

	it.Close()

	it = db.ReverseIterator([]byte("key_8"), []byte("key_0"), 0)
	if !it.Valid() {
		t.Fatal("must valid")
	}

	if string(it.Key()) != "key_7" {
		t.Fatal(string(it.Key()))
	}

	it.Close()

	for it = db.Iterator(nil, nil, 0); it.Valid(); it.Next() {
		db.Delete(it.Key())
	}

	it.Close()

	it = db.Iterator([]byte("key_0"), []byte("key_8"), 0)
	if it.Valid() {
		t.Fatal("must not valid")
	}

	it.Close()

	it = db.ReverseIterator([]byte("key_8"), []byte("key_0"), 0)
	if it.Valid() {
		t.Fatal("must not valid")
	}

	it.Close()
}

func TestSnapshot(t *testing.T) {
	db := getTestDB()

	key := []byte("key")
	value := []byte("hello world")

	db.Put(key, value)

	s := db.NewSnapshot()
	defer s.Close()

	db.Put(key, []byte("hello world2"))

	if v, err := s.Get(key); err != nil {
		t.Fatal(err)
	} else if string(v) != string(value) {
		t.Fatal(string(v))
	}

	found := false
	var it *Iterator
	for it = s.Iterator(nil, nil, 0); it.Valid(); it.Next() {
		if string(it.Key()) == string(key) {
			found = true
			break
		}
	}

	it.Close()

	if !found {
		t.Fatal("must found")
	}

	found = false
	for it = s.ReverseIterator(nil, nil, 0); it.Valid(); it.Next() {
		if string(it.Key()) == string(key) {
			found = true
			break
		}
	}

	it.Close()

	if !found {
		t.Fatal("must found")
	}

}

func TestDestroy(t *testing.T) {
	db := getTestDB()

	db.Destroy()

	if _, err := os.Stat(db.cfg.Path); !os.IsNotExist(err) {
		t.Fatal("must not exist")
	}
}
