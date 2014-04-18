package leveldb

import (
	"github.com/jmhodges/levigo"
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
