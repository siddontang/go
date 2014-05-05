package leveldb

import (
	"github.com/jmhodges/levigo"
)

type WriteBatch struct {
	db *DB
	wb *levigo.WriteBatch
}

func (wb *WriteBatch) Put(key, value []byte) {
	wb.wb.Put(key, value)
}

func (wb *WriteBatch) Delete(key []byte) {
	wb.wb.Delete(key)
}

func (wb *WriteBatch) Commit() error {
	return wb.db.db.Write(wb.db.writeOpts, wb.wb)
}

func (wb *WriteBatch) Rollback() {
	wb.wb.Clear()
}

func (wb *WriteBatch) Close() {
	if wb.wb == nil {
		return
	}

	wb.wb.Close()
	wb.wb = nil
}
