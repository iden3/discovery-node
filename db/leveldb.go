package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Db struct {
	ldb *leveldb.DB
	prefix []byte
}

func New(path string) (*Db, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err!=nil {
		return nil, err
	}
	return &Db{ldb, []byte{}}, nil
}

func (db *Db) WithPrefix(prefix []byte) *Db {
	return &Db{db.ldb, append(db.prefix, prefix...)}
}

func (db *Db) Put(key, value []byte) error {
	err := db.ldb.Put(append(db.prefix, key[:]...), value, nil)
	return err
}

func (db *Db) Get(key []byte) ([]byte, error) {
	v, err := db.ldb.Get(append(db.prefix, key[:]...), nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}

