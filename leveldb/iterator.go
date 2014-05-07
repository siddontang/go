package leveldb

import (
	"bytes"
	"github.com/jmhodges/levigo"
)

const forward uint8 = 0
const backward uint8 = 1

type Iterator struct {
	it *levigo.Iterator

	start []byte
	stop  []byte
	limit int

	step int

	//0 for forward, 1 for backward
	direction uint8
}

func newIterator(db *DB, opts *levigo.ReadOptions, start []byte, stop []byte, limit int, direction uint8) *Iterator {
	it := new(Iterator)
	it.it = db.db.NewIterator(opts)

	it.start = start
	it.stop = stop
	it.limit = limit
	it.direction = direction
	it.step = 0

	if start == nil {
		if direction == forward {
			it.it.SeekToFirst()
		} else {
			it.it.SeekToLast()
		}
	} else {
		it.it.Seek(start)

		if it.Valid() && !bytes.Equal(it.Key(), start) {
			//for forward, key is the next bigger than start
			//for backward, key is the next bigger than start, so must go prev
			if direction == backward {
				it.it.Prev()
			}
		}
	}

	return it
}

func (it *Iterator) Valid() bool {
	if !it.it.Valid() {
		return false
	}

	if it.limit > 0 && it.step >= it.limit {
		return false
	}

	if it.direction == forward {
		if it.stop != nil && bytes.Compare(it.Key(), it.stop) > 0 {
			return false
		}
	} else {
		if it.stop != nil && bytes.Compare(it.Key(), it.stop) < 0 {
			return false
		}
	}

	return true
}

func (it *Iterator) GetError() error {
	return it.it.GetError()
}

func (it *Iterator) Next() {
	it.step++

	if it.direction == forward {
		it.it.Next()
	} else {
		it.it.Prev()
	}
}

func (it *Iterator) Key() []byte {
	return it.it.Key()
}

func (it *Iterator) Value() []byte {
	return it.it.Value()
}

func (it *Iterator) Close() {
	it.it.Close()
}

func (it *Iterator) IntValue() (int64, error) {
	return Int(it.Value(), nil)
}

func (it *Iterator) UintValue() (uint64, error) {
	return Uint(it.Value(), nil)
}

func (it *Iterator) FloatValue() (float64, error) {
	return Float(it.Value(), nil)
}

func (it *Iterator) StringValue() (string, error) {
	return String(it.Value(), nil)
}

func (it *Iterator) SliceValue() ([]byte, error) {
	return Slice(it.Value(), nil)
}
