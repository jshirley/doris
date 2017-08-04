package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jshirley/doris/api"
	"github.com/jshirley/doris/storage"
)

func main() {
	tmpfile, err := ioutil.TempFile("", "storage_test")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing storage and API handlers")
	store := storage.NewStorage(tmpfile.Name())
	api := api.New(store)

	log.Printf("Starting http server up\n")
	srv := &http.Server{
		Handler:      api.Router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
