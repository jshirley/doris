package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jshirley/doris/storage"
)

func setup() *storage.Storage {
	tmpfile, err := ioutil.TempFile("", "storage_test")
	if err != nil {
		log.Fatal(err)
	}
	return storage.NewStorage(tmpfile.Name())
}

func teardown(store *storage.Storage) {
	store.Close()
	defer os.Remove(store.File) // clean up
}

var apiTests = []struct {
	path   string
	method string
	body   io.Reader
	out    string
}{
	{"/api/links", "GET", nil, linksResult(true, "ok", []storage.Link{}, 0, 1)},
	{"/api/links", "POST", makeLink("test", "https://www.google.com"), linksResult(true, "created a link", []storage.Link{{Name: "test", Url: "https://www.google.com"}}, 1, 1)},
}

func linksResult(ok bool, message string, links []storage.Link, count int, page int) string {
	res, err := json.Marshal(LinkResult{Ok: ok, Message: message, Links: links, Count: count, Page: page})
	if err != nil {
		log.Fatalf("Unable to marshal link response: %+v", err)
		return ""
	}
	return string(res)
}

func makeLink(name string, link string) io.Reader {
	l := storage.Link{Name: name, Url: link}
	s, _ := json.Marshal(l)
	return bytes.NewBuffer(s)
}

func TestLinkRequests(t *testing.T) {
	store := setup()
	defer teardown(store)

	apiObj := New(mux.NewRouter(), store)

	for _, tt := range apiTests {
		res := httptest.NewRecorder()

		req, _ := http.NewRequest(tt.method, tt.path, tt.body)
		req.Header.Set("Content-Type", "application/json")

		apiObj.Router.ServeHTTP(res, req)
		if res.Body.String() != tt.out {
			t.Errorf("Expected result from %s %s to be %s and got %s", tt.method, tt.path, tt.out, res.Body.String())
		}
	}
}
