### Task

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to return the correct numbers after restarting it, by persisting data to a file.

#### Overview

The app consists of the server pkg which defines (and uses) **Store** interface. And there are two implementations of this interface in a *store* folder. One is *gob* -- which is an array of timestamps written to disk on every request. Another one called *mapper* -- is a map of timestamps in Unixseconds format as a key and a number of hits (counter) as a value.


#### Some tests and benchmarks

```bash
# clone the repo
git clone

# to run all tests
go test -v ./...

# to run all benchmarks
go test -v ./... -bench=. -run=xxx -benchmem

```

### How to install

```
# build
