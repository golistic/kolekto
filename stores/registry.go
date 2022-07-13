// Copyright (c) 2022, Geert JM Vanderkelen

package stores

import (
	"github.com/golistic/kolekto/kolektor"
)

var registry = map[StoreKind]func(dsn string) (kolektor.Storer, error){}

// New instantiated a certain kind of store using dsn for connecting.
func New(kind StoreKind, dsn string) (kolektor.Storer, error) {
	store, err := registry[kind](dsn)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Register registers a kind of store mapping it with it initialization
// function.
func Register(kind StoreKind, fn func(dsn string) (kolektor.Storer, error)) {
	registry[kind] = fn
}
