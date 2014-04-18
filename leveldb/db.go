package leveldb

import (
	"encoding/json"
	"github.com/jmhodges/levigo"
)

const defaultFilterBits int = 10

type Config struct {
	Path string `json:"path"`

	Compression bool `json:"compression"`

	BlockSize       int `json:"block_size"`
	WriteBufferSize int `json:"write_buffer_size"`
	CacheSize       int `json:"cache_size"`
}

type DB struct {
	cfg *Config
	db  *levigo.DB

	opts *levigo.Options

	//for default read and write options
	readOpts     *levigo.ReadOptions
	writeOpts    *levigo.WriteOptions
	iteratorOpts *levigo.ReadOptions

	cache *levigo.Cache

	filter *levigo.FilterPolicy
}

func Open(configJson json.RawMessage) (*DB, error) {
	cfg := new(Config)
	err := json.Unmarshal(configJson, cfg)
	if err != nil {
		return nil, err
	}

	db := new(DB)
	db.cfg = cfg

	db.opts = db.initOptions(cfg)

	db.readOpts = levigo.NewReadOptions()
	db.writeOpts = levigo.NewWriteOptions()
	db.iteratorOpts = levigo.NewReadOptions()
	db.iteratorOpts.SetFillCache(false)

	db.db, err = levigo.Open(cfg.Path, db.opts)
	return db, err
}

func (db *DB) initOptions(cfg *Config) *levigo.Options {
	opts := levigo.NewOptions()

	opts.SetCreateIfMissing(true)

	if cfg.CacheSize > 0 {
		db.cache = levigo.NewLRUCache(cfg.CacheSize)
		opts.SetCache(db.cache)
	}

	//we must use bloomfilter
	db.filter = levigo.NewBloomFilter(defaultFilterBits)
	opts.SetFilterPolicy(db.filter)

	if !cfg.Compression {
		opts.SetCompression(levigo.NoCompression)
	}

	blockSize := cfg.BlockSize * 1024
	if blockSize > 0 {
		opts.SetBlockSize(blockSize)
	}

	writeBufferSize := cfg.WriteBufferSize * 1024 * 1024
	if writeBufferSize > 0 {
		opts.SetWriteBufferSize(writeBufferSize)
	}

	return opts
}

func (db *DB) Close() {
	db.opts.Close()

	if db.cache != nil {
		db.cache.Close()
	}

	if db.filter != nil {
		db.filter.Close()
	}

	db.readOpts.Close()
	db.writeOpts.Close()
	db.iteratorOpts.Close()

	db.db.Close()
	db.db = nil
}

func (db *DB) Destroy() {
	db.Close()
	opts := levigo.NewOptions()
	defer opts.Close()

	levigo.DestroyDatabase(db.cfg.Path, opts)
}

func (db *DB) Put(key, value []byte) error {
	return db.db.Put(db.writeOpts, key, value)
}

func (db *DB) Get(key []byte) ([]byte, error) {
	return db.db.Get(db.readOpts, key)
}

func (db *DB) Delete(key []byte) error {
	return db.db.Delete(db.writeOpts, key)
}

func (db *DB) NewWriteBatch() *WriteBatch {
	wb := new(WriteBatch)
	wb.wb = levigo.NewWriteBatch()
	wb.db = db
	return wb
}

//like c++ iterator, [begin, end)
//begin should less than end
//if begin is nil, we will seek to first
//if end is nil, we will next until read last
//limit <= 0, no limit
func (db *DB) Iterator(begin []byte, end []byte, limit int) *Iterator {
	return newIterator(db, db.iteratorOpts, begin, end, limit, forward)
}

//like c++ reverse_iterator, [rbegin, rend)
//rbegin should bigger than rend
//if rbegin is nil, we will seek to last
//if end is nil, we will next until read first
//limit <= 0, no limit
func (db *DB) ReverseIterator(rbegin []byte, rend []byte, limit int) *Iterator {
	return newIterator(db, db.iteratorOpts, rbegin, rend, limit, backward)
}

func (db *DB) NewSnapshot() *Snapshot {
	return newSnapshot(db)
}
