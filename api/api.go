package api

import (
	"encoding/json"
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

	s.HandleFunc("/links", api.ListLinksHandler).Methods("GET").Name("ListLinks")
	s.HandleFunc("/links", api.CreateLinkHandler).Methods("POST").Name("CreateLink")

	return api
}

type LinkResult struct {
	Ok      bool           `json:"ok"`
	Message string         `json:"message"`
	Links   []storage.Link `json:"links", omitempty`
	Count   int            `json:"count", omitempty`
	Page    int            `json:"page", omitempty`
}

func (api *API) ListLinksHandler(w http.ResponseWriter, r *http.Request) {
	links, err := api.Store.ListAll()
	if err != nil {
		respondWithError(w, http.StatusOK, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, LinkResult{Ok: true, Message: "ok", Links: links, Count: len(links), Page: 1})
	}
}

func (api *API) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var data storage.Link
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = api.Store.StoreLink(data.Name, data.Url)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	link, err := api.Store.GetLink(data.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, LinkResult{Ok: true, Message: "created a link", Links: []storage.Link{*link}, Count: 1, Page: 1})
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, LinkResult{Ok: false, Message: message})
}
