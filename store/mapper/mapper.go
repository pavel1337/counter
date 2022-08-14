package mapper

import (
	"encoding/gob"
	"io"
	"os"
	"sync"
	"time"
)

type store struct {
	lock sync.Mutex
	path string
}

// M is the map of times represented as Unix seconds as key and count as value.
type M map[int64]int

// New returns a new store.
func New(path string) (*store, error) {
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	f.Close()

	s := &store{
		lock: sync.Mutex{},
		path: path,
	}

	err = s.migrate()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// migrate checks the file store and migrates it if necessary.
func (s *store) migrate() error {
	if fileExists(s.path) && !fileEmpty(s.path) {
		_, err := s.load()
		return err
	}
	return s.save(M{})
}

// Add adds a time to the store.
func (s *store) Add(t time.Time) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	m, err := s.load()
	if err != nil {
		return err
	}
	m[t.Unix()]++

	return s.save(m)
}

// Last returns the number of hits that the store has seen within the
// last sec seconds.
func (s *store) Last(sec int) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	m, err := s.load()
	if err != nil {
		return 0, err
	}

	var counter int
	for i := time.Now().Unix(); i > time.Now().Unix()-int64(sec); i-- {
		counter += m[i]
	}

	return counter, nil
}

// load loads the store from the file. Not safe for concurrent use.
func (s *store) load() (M, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m M
	dec := gob.NewDecoder(f)
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// save saves the store to the file. Not safe for concurrent use.
func (s *store) save(m M) error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(m)
	if err != nil {
		return err
	}
	return nil
}

// fileExists returns true if the file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// fileEmpty returns true if the file is empty.
func fileEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Read(make([]byte, 1))
	return err == io.EOF
}
