package storage

import (
	"io/ioutil"
	"log"
	"net/url"
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
	store.Close()
	defer os.Remove(store.File) // clean up
}

func TestStorage(t *testing.T) {
	setup()
	defer teardown()

	link, err := store.GetLink("foo")
	if link != nil {
		t.Errorf("Expected an error, got: %+v", err)
	}

	err = store.StoreLink("foo", "not an absolute url")
	if err == nil || err.Error() != "invalid_absolute_url" {
		t.Errorf("Expected to reject bogus URL: %+v", err)
	}

	err = store.StoreLink("foo", "bogus%ZZurl")
	if serr, ok := err.(*url.EscapeError); ok {
		t.Errorf("Expected to reject bogus URL: %+v", serr)
	}

	testUrl := "http://google.com"
	err = store.StoreLink("foo", testUrl)
	if err != nil {
		t.Errorf("Expected to store legitimate URL: %+v", err)
	}

	link, err = store.GetLink("foo")
	if err != nil {
		t.Errorf("Expected a URL, got: %+v", err)
	}

	if link.String() != testUrl {
		t.Errorf("Expected %v, got %v", testUrl, link)
	}

	testUrl2 := "http://google.com/updated"
	err = store.StoreLink("foo", testUrl2)
	if err != nil {
		t.Errorf("Expected to store legitimate URL: %+v", err)
	}
	link, _ = store.GetLink("foo")
	if link.String() != testUrl2 {
		t.Errorf("Expected store to update to %v, got %v", testUrl2, link)
	}
}
