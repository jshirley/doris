package storage

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/boltdb/bolt"
)

type Storage struct {
	File    string
	Bucket  string
	db      *bolt.DB
	linkMap map[string]Link
}

type Link struct {
	Version     string `json:"version",omitempty`
	Name        string `json:"name"`
	Description string `json:"description",omitempty`
	Url         string `json:"url"`
}

func NewStorage(file string) *Storage {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// Setup the default bucket.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Links"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{file, "Links", db, make(map[string]Link)}
}

func (store *Storage) Close() {
	defer store.db.Close()
}

func (store *Storage) ListAll() ([]Link, error) {
	links := []Link{}
	var link Link

	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(store.Bucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := json.Unmarshal(v, &link); err != nil {
				return err
			}
			links = append(links, link)
		}

		return nil
	})

	if err != nil {
		return []Link{}, err
	}
	return links, err
}

func (store *Storage) GetLink(name string) (*url.URL, error) {
	var link Link

	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(store.Bucket))
		v := b.Get([]byte(name))
		if err := json.Unmarshal(v, &link); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return url.Parse(link.Url)
}

func (store *Storage) StoreLink(name string, incomingUrl string) error {
	parsedUrl, err := url.Parse(incomingUrl)
	if err != nil {
		return err
	}

	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return errors.New("invalid_absolute_url")
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		link := Link{Name: name, Url: incomingUrl}
		bkt := tx.Bucket([]byte(store.Bucket))

		if buf, err := json.Marshal(link); err != nil {
			return err
		} else if err := bkt.Put([]byte(name), buf); err != nil {
			return err
		}
		return nil
	})
}
