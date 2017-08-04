package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jshirley/doris/storage"
)

type API struct {
	Router *mux.Router
	Store  *storage.Storage
}

func New(store *storage.Storage) *API {
	router := mux.NewRouter()
	s := router.PathPrefix("/api").Subrouter()

	api := &API{router, store}

	s.HandleFunc("/links", api.LinksHandler).Methods("GET", "POST").Name("Links")

	return api
}

func (api *API) LinksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Fetching from API storage all the links")
	links, err := api.Store.ListAll()
	if err != nil {
		respondWithError(w, http.StatusOK, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, links)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
