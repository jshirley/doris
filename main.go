package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jshirley/doris/api"
	"github.com/jshirley/doris/storage"
)

type Doris struct {
	store  *storage.Storage
	api    *api.API
	router *mux.Router
}

func NewDoris(file string) *Doris {
	router := mux.NewRouter()
	store := storage.NewStorage(file)

	// Registers subroutes under /api/
	// I am not sure I like this pattern of mutating router, but I wasn't sure
	// how to return the api and then attach api.Router to the top level router.
	api := api.New(router, store)

	doris := &Doris{store: store, api: api, router: router}

	router.HandleFunc("/{token}", doris.TokenHandler).Methods("GET")
	router.HandleFunc("/{token}/{args}", doris.TokenHandler).Methods("GET")

	return doris
}

func (doris *Doris) TokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	args := vars["args"]
	if args != "" {
		w.Write([]byte(fmt.Sprintf("%s with args: %s", token, args)))
	} else {
		w.Write([]byte(token))
	}
}

func (doris *Doris) Start() {
	srv := &http.Server{
		Handler:      doris.router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func main() {
	var file string
	flag.StringVar(&file, "file", "", "Path to link storage, just a file!")

	flag.Parse()

	if file == "" {
		log.Fatal("Missing -file argument. Tell me where to store your links")
		return
	}

	doris := NewDoris(file)
	doris.Start()
}
