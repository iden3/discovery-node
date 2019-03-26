package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Db is the database with the specified prefix
type Db struct {
	ldb    *leveldb.DB
	prefix []byte
}

// New creates a new Db database
func New(path string) (*Db, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Db{ldb, []byte{}}, nil
}

// WithPrefix returns a subdatabase with the specified prefix
func (db *Db) WithPrefix(prefix []byte) *Db {
	return &Db{db.ldb, append(db.prefix, prefix...)}
}

// Put adds the key value to the database
func (db *Db) Put(key, value []byte) error {
	err := db.ldb.Put(append(db.prefix, key[:]...), value, nil)
	return err
}

// Get retreives a value from the database for a given key
func (db *Db) Get(key []byte) ([]byte, error) {
	v, err := db.ldb.Get(append(db.prefix, key[:]...), nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Iterate iterates over the database
func (db *Db) Iterate(f func([]byte, []byte)) error {
	snapshot, err := db.ldb.GetSnapshot()
	if err != nil {
		return err
	}
	iter := snapshot.NewIterator(nil, nil)
	for iter.Next() {
		f(iter.Key(), iter.Value())
	}
	iter.Release()
	err = iter.Error()
	return err
}
