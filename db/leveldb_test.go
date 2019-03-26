package db

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

func TestNew(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	_, err = New(dir)
	assert.Nil(t, err)
}
func TestPutAndGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := New(dir)
	assert.Nil(t, err)

	k := []byte("test key")
	v := []byte("test value")

	err = db.Put(k, v)
	assert.Nil(t, err)

	vg, err := db.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, v, vg)

}

func TestPrefix(t *testing.T) {
	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)
	db, err := New(dir)
	assert.Nil(t, err)

	dbPrefix := db.WithPrefix([]byte("prefix1"))

	err = dbPrefix.Put([]byte("k1"), []byte("v1"))
	assert.Nil(t, err)

	_, err = db.Get([]byte("k1"))
	assert.Equal(t, err, errors.ErrNotFound)

	v1, err := dbPrefix.Get([]byte("k1"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("v1"), v1)
}
