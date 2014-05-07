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

func OpenWithConfig(cfg *Config) (*DB, error) {
	db := new(DB)
	db.cfg = cfg

	db.opts = db.initOptions(cfg)

	db.readOpts = levigo.NewReadOptions()
	db.writeOpts = levigo.NewWriteOptions()
	db.iteratorOpts = levigo.NewReadOptions()
	db.iteratorOpts.SetFillCache(false)

	var err error
	db.db, err = levigo.Open(cfg.Path, db.opts)
	return db, err
}

func Open(configJson json.RawMessage) (*DB, error) {
	cfg := new(Config)
	err := json.Unmarshal(configJson, cfg)
	if err != nil {
		return nil, err
	}

	return OpenWithConfig(cfg)
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

	if cfg.BlockSize > 0 {
		opts.SetBlockSize(cfg.BlockSize)
	}

	if cfg.WriteBufferSize > 0 {
		opts.SetWriteBufferSize(cfg.WriteBufferSize)
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

func (db *DB) Clear() {
	it := db.Iterator(nil, nil, 0, 0, -1)
	for ; it.Valid(); it.Next() {
		db.Delete(it.Key())
	}
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

//limit < 0, unlimit
//offset must >= 0, if < 0, will get nothing
func (db *DB) Iterator(min []byte, max []byte, rangeType uint8, offset int, limit int) *Iterator {
	return newIterator(db, db.iteratorOpts, NewRange(min, max, rangeType), offset, limit, IteratorForward)
}

//limit < 0, unlimit
//offset must >= 0, if < 0, will get nothing
func (db *DB) RevIterator(min []byte, max []byte, rangeType uint8, offset int, limit int) *Iterator {
	return newIterator(db, db.iteratorOpts, NewRange(min, max, rangeType), offset, limit, IteratorBackward)
}

func (db *DB) NewSnapshot() *Snapshot {
	return newSnapshot(db)
}

func (db *DB) GetInt(key []byte) (int64, error) {
	return Int(db.Get(key))
}

func (db *DB) GetUInt(key []byte) (uint64, error) {
	return Uint(db.Get(key))
}

func (db *DB) GetFloat(key []byte) (float64, error) {
	return Float(db.Get(key))
}

func (db *DB) GetString(key []byte) (string, error) {
	return String(db.Get(key))
}

func (db *DB) GetSlice(key []byte) ([]byte, error) {
	return Slice(db.Get(key))
}
