package leveldb

import (
	"bytes"
	"github.com/jmhodges/levigo"
)

const forward uint8 = 0
const backward uint8 = 1

//min must less or equal than max
//MinEx if true, range is left open interval (min, ...
//MaxEx if true, range is right open interval ..., max)
//Default range is close interval
type Range struct {
	Min []byte
	Max []byte

	MinEx bool
	MaxEx bool
}

func NewRange(min []byte, max []byte) *Range {
	return &Range{min, max, false, false}
}

func NewOpenRange(min []byte, max []byte) *Range {
	return &Range{min, max, true, true}
}

func NewLOpenRange(min []byte, max []byte) *Range {
	return &Range{min, max, true, false}
}

func NewROpenRange(min []byte, max []byte) *Range {
	return &Range{min, max, true, true}
}

type Iterator struct {
	it *levigo.Iterator

	r *Range

	limit int

	step int

	//0 for forward, 1 for backward
	direction uint8
}

func newIterator(db *DB, opts *levigo.ReadOptions, r *Range, limit int, direction uint8) *Iterator {
	it := new(Iterator)
	it.it = db.db.NewIterator(opts)

	it.r = r
	it.limit = limit
	it.direction = direction
	it.step = 0

	if direction == forward {
		if r.Min == nil {
			it.it.SeekToFirst()
		} else {
			it.it.Seek(r.Min)

			if r.MinEx {
				if it.Valid() && bytes.Equal(it.Key(), r.Min) {
					it.it.Next()
				}
			}
		}
	} else {
		if r.Max == nil {
			it.it.SeekToLast()
		} else {
			it.it.Seek(r.Max)
			if it.Valid() && !bytes.Equal(it.Key(), r.Max) {
				//key must bigger than max, so we must go prev
				it.it.Prev()
			}

			if r.MaxEx {
				if it.Valid() && bytes.Equal(it.Key(), r.Max) {
					it.it.Prev()
				}
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
		if it.r.Max != nil {
			r := bytes.Compare(it.Key(), it.r.Max)
			if !it.r.MaxEx {
				return !(r > 0)
			} else {
				return !(r >= 0)
			}
		}
	} else {
		if it.r.Min != nil {
			r := bytes.Compare(it.Key(), it.r.Min)
			if !it.r.MinEx {
				return !(r < 0)
			} else {
				return !(r <= 0)
			}
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
