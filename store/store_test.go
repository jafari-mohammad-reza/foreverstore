package store

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: HashPathTransformFunc,
	}
	pathName := "testDir"
	store := NewStore(opts)
	writtenData := []byte("some data")
	storeData := bytes.NewReader(writtenData)
	if err := store.writeStream(pathName, storeData); err != nil {
		t.Error(err)
	}
	_, err := store.readStream(pathName)
	if err != nil {
		t.Error(err)
	}
	file, err := store.read(pathName)
	if err != nil {
		t.Error(err)
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, fileData, writtenData)
	exists := store.exists(pathName)
	assert.True(t, exists)
	err = store.delete(pathName)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		os.RemoveAll(store.opts.Root)
	}()
}

func TestPathTransformFunc(t *testing.T) {
	key := "testDir"
	pathname := HashPathTransformFunc(key)
	assert.Equal(t, "b3732/ef50a/d2121/88a6a/50826/ec8e3/c6db7/0c41b", pathname.Pathname)
	assert.Equal(t, "b3732ef50ad212188a6a50826ec8e3c6db70c41b", pathname.Filename)
}
