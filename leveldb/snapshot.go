package leveldb

import (
	"github.com/jmhodges/levigo"
	"github.com/siddontang/golib/hack"
	"strconv"
)

type Snapshot struct {
	db           *DB
	s            *levigo.Snapshot
	readOpts     *levigo.ReadOptions
	iteratorOpts *levigo.ReadOptions
}

func newSnapshot(db *DB) *Snapshot {
	s := new(Snapshot)
	s.db = db
	s.s = db.db.NewSnapshot()

	s.readOpts = levigo.NewReadOptions()
	s.readOpts.SetSnapshot(s.s)

	s.iteratorOpts = levigo.NewReadOptions()
	s.iteratorOpts.SetSnapshot(s.s)
	s.iteratorOpts.SetFillCache(false)

	return s
}

func (s *Snapshot) Close() {
	s.db.db.ReleaseSnapshot(s.s)

	s.iteratorOpts.Close()
	s.readOpts.Close()
}

func (s *Snapshot) Get(key []byte) ([]byte, error) {
	return s.db.db.Get(s.readOpts, key)
}

//same as db iterator and reverse iterator
func (s *Snapshot) Iterator(begin []byte, end []byte, limit int) *Iterator {
	return newIterator(s.db, s.iteratorOpts, begin, end, limit, forward)
}

func (s *Snapshot) ReverseIterator(rbegin []byte, rend []byte, limit int) *Iterator {
	return newIterator(s.db, s.iteratorOpts, rbegin, rend, limit, backward)
}

func (s *Snapshot) GetInt(key []byte) (int64, error) {
	v, err := s.Get(key)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseInt(hack.String(v), 10, 64)
}

func (s *Snapshot) GetUInt(key []byte) (uint64, error) {
	v, err := s.Get(key)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseUint(hack.String(v), 10, 64)
}

func (s *Snapshot) GetFloat(key []byte) (float64, error) {
	v, err := s.Get(key)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseFloat(hack.String(v), 64)
}

func (s *Snapshot) GetString(key []byte) (string, error) {
	v, err := s.Get(key)
	if err != nil {
		return "", err
	} else if v == nil {
		return "", nil
	}

	return hack.String(v), nil
}

func (s *Snapshot) GetSlice(key []byte) ([]byte, error) {
	v, err := s.Get(key)
	if err != nil {
		return nil, err
	} else if v == nil {
		return []byte{}, nil
	}

	return v, nil
}
