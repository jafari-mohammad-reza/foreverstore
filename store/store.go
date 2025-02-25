package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"strings"
)

type PathTransformFunc func(string) PathKey

func HashPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashString) / blockSize
	paths := make([]string, sliceLen)
	for i := range sliceLen {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}
	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashString,
	}
}

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}
type PathKey struct {
	Pathname string
	Filename string
}

func (p *PathKey) FullPath(rootDir string) string {
	return path.Join(rootDir, p.Pathname, p.Filename)
}

func (p *PathKey) FirstPath() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

type Store struct {
	opts StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if len(opts.Root) == 0 {
		opts.Root = "forever-store"
	}
	return &Store{
		opts: opts,
	}
}

func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.opts.PathTransformFunc(key)
	return os.Open(pathKey.FullPath(s.opts.Root))
}

func (s *Store) Delete(key string) error {
	pathKey := s.opts.PathTransformFunc(key)
	return os.RemoveAll(pathKey.FullPath(s.opts.Root))
}

func (s *Store) Exists(key string) bool {
	pathKey := s.opts.PathTransformFunc(key)
	_, err := os.Stat(pathKey.FullPath(s.opts.Root))
	if os.ErrNotExist == err {
		return false
	}
	return true
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	defer func() {
		err := f.Close()
		if err != nil {
			slog.Error("Store read", "Error", err.Error())
		}
	}()
	if err != nil {
		return nil, err
	}
	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, f)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathKey := s.opts.PathTransformFunc(key)
	if err := os.MkdirAll(path.Join(s.opts.Root, pathKey.Pathname), os.ModePerm); err != nil {
		return err
	}
	pathAndFileName := pathKey.FullPath(s.opts.Root)
	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", n, pathAndFileName)
	return nil
}
