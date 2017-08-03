package storage

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var store *Storage

func setup() {
	tmpfile, err := ioutil.TempFile("", "storage_test")
	if err != nil {
		log.Fatal(err)
	}
	store = NewStorage(tmpfile.Name())
}

func teardown() {
	defer os.Remove(store.File) // clean up
}

func TestGetLink(t *testing.T) {
	setup()
	defer teardown()

	url, err := store.GetLink("foo")
	if url != nil {
		t.Errorf("Expected an error, got: %+v", err)
	}

	err = store.StoreLink("foo", "bogus url")
	if err == nil || err.Error() != "invalid_absolute_url" {
		t.Errorf("Expected to reject bogus URL: %+v", err)
	}

	testUrl := "http://google.com"
	err = store.StoreLink("foo", testUrl)
	if err != nil {
		t.Errorf("Expected to store legitimate URL: %+v", err)
	}

	url, err = store.GetLink("foo")
	if err != nil {
		t.Errorf("Expected a URL, got: %+v", err)
	}

	if url.String() != testUrl {
		t.Errorf("Expected %v, got %v", testUrl, url)
	}

	testUrl2 := "http://google.com/updated"
	err = store.StoreLink("foo", testUrl2)
	if err != nil {
		t.Errorf("Expected to store legitimate URL: %+v", err)
	}
	url, _ = store.GetLink("foo")
	if url.String() != testUrl2 {
		t.Errorf("Expected store to update to %v, got %v", testUrl2, url)
	}
}
