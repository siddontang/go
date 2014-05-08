package leveldb

import (
	"bytes"
	"github.com/jmhodges/levigo"
)

const (
	IteratorForward  uint8 = 0
	IteratorBackward uint8 = 1
)

const (
	RangeClose uint8 = 0x00
	RangeLOpen uint8 = 0x01
	RangeROpen uint8 = 0x10
	RangeOpen  uint8 = 0x11
)

//min must less or equal than max
//range type:
//close: [min, max]
//open: (min, max)
//lopen: (min, max]
//ropen: [min, max)
type Range struct {
	Min []byte
	Max []byte

	Type uint8
}

func NewRange(min []byte, max []byte, tp uint8) *Range {
	return &Range{min, max, tp}
}

type Iterator struct {
	it *levigo.Iterator

	r *Range

	offset int
	limit  int

	step int

	//0 for IteratorForward, 1 for IteratorBackward
	direction uint8
}

func newIterator(db *DB, opts *levigo.ReadOptions, r *Range, offset int, limit int, direction uint8) *Iterator {
	it := new(Iterator)
	it.it = db.db.NewIterator(opts)

	it.r = r
	it.offset = offset
	it.limit = limit
	it.direction = direction

	it.step = 0

	if offset < 0 {
		return it
	}

	if direction == IteratorForward {
		if r.Min == nil {
			it.it.SeekToFirst()
		} else {
			it.it.Seek(r.Min)

			if r.Type&RangeLOpen > 0 {
				if it.it.Valid() && bytes.Equal(it.it.Key(), r.Min) {
					it.it.Next()
				}
			}
		}
	} else {
		if r.Max == nil {
			it.it.SeekToLast()
		} else {
			it.it.Seek(r.Max)

			if !it.it.Valid() {
				it.it.SeekToLast()
			} else {
				if !bytes.Equal(it.it.Key(), r.Max) {
					it.it.Prev()
				}
			}

			if r.Type&RangeROpen > 0 {
				if it.Valid() && bytes.Equal(it.Key(), r.Max) {
					it.it.Prev()
				}
			}
		}
	}

	for i := 0; i < offset; i++ {
		if it.it.Valid() {
			if it.direction == IteratorForward {
				it.it.Next()
			} else {
				it.it.Prev()
			}
		}
	}

	return it
}

func (it *Iterator) Valid() bool {
	if it.offset < 0 {
		return false
	} else if !it.it.Valid() {
		return false
	} else if it.limit >= 0 && it.step >= it.limit {
		return false
	}

	if it.direction == IteratorForward {
		if it.r.Max != nil {
			r := bytes.Compare(it.Key(), it.r.Max)
			if it.r.Type&RangeROpen > 0 {
				return !(r >= 0)
			} else {
				return !(r > 0)
			}
		}
	} else {
		if it.r.Min != nil {
			r := bytes.Compare(it.Key(), it.r.Min)
			if it.r.Type&RangeLOpen > 0 {
				return !(r <= 0)
			} else {
				return !(r < 0)
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

	if it.direction == IteratorForward {
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
