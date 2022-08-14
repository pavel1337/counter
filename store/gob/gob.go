package gob

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type FileStore struct {
	path string
	lock sync.Mutex
}

// NewFileStore returns a new FileStore.
func NewFileStore(path string) (*FileStore, error) {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	file.Close()

	fs := &FileStore{
		path: path,
		lock: sync.Mutex{},
	}

	err = fs.migrate()
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// migrate checks the file store and migrates it if necessary.
func (s *FileStore) migrate() error {
	if fileExists(s.path) && !fileEmpty(s.path) {
		_, err := s.all()
		return err
	}
	return s.save([]time.Time{})
}

func (s *FileStore) Add(t time.Time) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ts, err := s.all()
	if err != nil {
		return err
	}
	ts = append(ts, t)
	return s.save(ts)
}

func (s *FileStore) Last(sec int) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	ts, err := s.all()
	if err != nil {
		return 0, err
	}
	if len(ts) == 0 {
		return 0, nil
	}

	var counter int

	var newTs []time.Time
	for _, t := range ts {
		if t.After(time.Now().Add(-time.Duration(sec) * time.Second)) {
			newTs = append(newTs, t)
			counter++
		}
	}

	return counter, s.save(newTs)
}

// all gets all timestamps from the file. Not safe for concurrent use.
func (s *FileStore) all() ([]time.Time, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	var ts []time.Time
	err = decoder.Decode(&ts)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file: %s", err)
	}
	return ts, nil
}

// save writes all timestamps to the file. Not safe for concurrent use.
func (s *FileStore) save(ts []time.Time) error {
	f, err := os.OpenFile(s.path, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)

	err = encoder.Encode(ts)
	if err != nil {
		return fmt.Errorf("failed to encode file: %s", err)
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
