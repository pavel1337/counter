package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Store interface {
	Add(t time.Time) error
	Last(sec int) (int, error)
}

type app struct {
	store Store
	port  int
	lock  sync.Mutex
}

// New returns a new app.
func New(store Store, port int) *app {
	return &app{
		store: store,
		port:  port,
		lock:  sync.Mutex{},
	}
}

// Server returns a new http.Server.
func (a *app) Server() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.allHandler)

	return &http.Server{
		Addr:    ":" + strconv.Itoa(a.port),
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}
}

// allHandler is a handler that returns a counter of the total number of requests
// that it has received during the previous 60 seconds (moving window).
func (a *app) allHandler(w http.ResponseWriter, r *http.Request) {
	// a.lock.Lock()
	// defer a.lock.Unlock()

	count, err := a.store.Last(60)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = a.store.Add(time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%d requests in the last 60 seconds", count)

	w.Write([]byte(strconv.Itoa(count)))
}
