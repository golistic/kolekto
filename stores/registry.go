// Copyright (c) 2022, Geert JM Vanderkelen

package stores

import (
	"github.com/golistic/kolekto/kolektor"
)

var registry = map[kolektor.StoreKind]func(dsn string) (kolektor.Storer, error){}

// New instantiated a certain kind of store using dsn for connecting.
func New(kind kolektor.StoreKind, dsn string) (kolektor.Storer, error) {
	store, err := registry[kind](dsn)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// Register registers a kind of store mapping it with it initialization
// function.
func Register(kind kolektor.StoreKind, fn func(dsn string) (kolektor.Storer, error)) {
	registry[kind] = fn
}

// Registered returns map of all registered stores. The key is the kind
// and the function with which a Store instance is created.
func Registered() map[kolektor.StoreKind]func(dsn string) (kolektor.Storer, error) {
	return registry
}
