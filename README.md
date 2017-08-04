# Doris

tl;dr: A tool so you can go to http://go/some-short-name and it redirects to http://some-longer-and-harder.com/to-remember/webpage.

[![Coverage Status](https://coveralls.io/repos/github/jshirley/doris/badge.svg?branch=master)](https://coveralls.io/github/jshirley/doris?branch=master)
[![Build Status](https://travis-ci.org/jshirley/doris.svg?branch=master)](https://travis-ci.org/jshirley/doris)

## In Prose

Doris, from Dôron and Zôros, meaning "abundance" and "pure and unmixed",
representing the bountiful fertility of the ocean. The web is our ocean, the
fish the morsels of knowledge that may elude us with their cunning. With the
gracious help of Doris, we may improve our harvest and thus our knowledge.

# Working with Doris

## Storage

Doris uses a local storage file, which makes it a stateful service. I may make
an option to store somewhere else. Who knows! It uses
[Bolt](https://github.com/boltdb/bolt) to store things, since it is really just
a key/value store. I could use [Datastore](https://cloud.google.com/datastore/)
for GCP installs, which would let me add some search functionality.

## Setup & Testing

I'm not checking in dependencies, and using
[`dep`](https://github.com/golang/dep). Hopefully this is the right choice!

```
go get -u github.com/golang/dep/cmd/dep
dep ensure
# If you want to test everything but the vendors:
go test -race -v $(go list ./... | grep -v '/vendor/')
```

There's CI at [Travis](https://travis-ci.org/jshirley/doris), and
[coverage reports](https://coveralls.io/github/jshirley/doris).
