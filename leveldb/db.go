package leveldb

import (
	"encoding/json"
	"github.com/jmhodges/levigo"
	"github.com/siddontang/golib/hack"
	"strconv"
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

func (db *DB) GetInt(key []byte) (int64, error) {
	v, err := db.Get(key)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseInt(hack.String(v), 10, 64)
}

func (db *DB) GetUInt(key []byte) (uint64, error) {
	v, err := db.Get(key)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseUint(hack.String(v), 10, 64)
}

func (db *DB) GetFloat(key []byte) (float64, error) {
	v, err := db.Get(key)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return strconv.ParseFloat(hack.String(v), 64)
}

func (db *DB) GetString(key []byte) (string, error) {
	v, err := db.Get(key)
	if err != nil {
		return "", err
	} else if v == nil {
		return "", nil
	}

	return hack.String(v), nil
}

func (db *DB) GetSlice(key []byte) ([]byte, error) {
	v, err := db.Get(key)
	if err != nil {
		return nil, err
	} else if v == nil {
		return []byte{}, nil
	}

	return v, nil
}
