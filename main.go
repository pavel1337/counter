package main

import (
	"counter/server"
	"counter/store/gob"
	"counter/store/mapper"
	"flag"
	"log"
)

func main() {
	imp := flag.String("imp", "", "implementation of the store (gob (array of timestamps) or mapper (map of timestamps in unix seconds format))")
	path := flag.String("path", "", "path to the store")

	flag.Parse()

	var s server.Store
	var err error

	switch *imp {
	case "gob":
		s, err = gob.NewFileStore(*path)
		if err != nil {
			log.Fatal(err)
		}

	case "mapper":
		s, err = mapper.New(*path)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("unknown implementation")
	}

	app := server.New(s, 8080)
	server := app.Server()
	log.Println("Listening on port 8080")
	log.Fatal(server.ListenAndServe())

}
