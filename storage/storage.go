package storage

import (
	//"encoding/json"
	"errors"
	"net/url"
)

type Storage struct {
	File    string
	linkMap map[string]Link
}

type Link struct {
	Version     string
	Name        string
	Description string
	URL         string
}

func NewStorage(file string) *Storage {
	return &Storage{file, make(map[string]Link)}
}

func (store *Storage) GetLink(name string) (*url.URL, error) {
	link, exists := store.linkMap[name]
	if exists == false {
		return nil, errors.New("invalid_link_specified")
	}

	return url.Parse(link.URL)
}

func (store *Storage) StoreLink(name string, incomingUrl string) error {
	parsedUrl, err := url.Parse(incomingUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return errors.New("invalid_absolute_url")
	}

	link, exists := store.linkMap[name]
	if exists {
		link.URL = parsedUrl.String()
	} else {
		store.linkMap[name] = Link{Name: name, URL: parsedUrl.String()}
	}
	return nil
}
