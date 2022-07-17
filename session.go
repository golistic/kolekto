// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"context"

	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/kolekto/stores"
)

// Session wraps around a store object to manage JSON collections.
type Session struct {
	store kolektor.Storer
}

// NewSession instantiates a new Session using a certain kind of
// data store. The dsn or data source name (DSN) is used to connect.
// The format of the DSN depends on the kind of store used.
func NewSession(kind kolektor.StoreKind, dsn string) (*Session, error) {
	ses := &Session{}
	var err error

	ses.store, err = stores.New(kind, dsn)
	if err != nil {
		return nil, err
	}

	return ses, nil
}

// newSession takes data source name as dsn, and the function used to
// instantiate the store.
// This is mostly used for testing for looping of all registered stores.
func newSession(dsn string, fn func(dsn string) (kolektor.Storer, error)) (*Session, error) {
	ses := &Session{}
	var err error

	ses.store, err = fn(dsn)
	if err != nil {
		return nil, err
	}
	return ses, nil
}

// Collection returns an instance that can be used to store and retrieve
// objects which are based on the provided model.
// If the collection is not yet available in the data store, it is created.
func (ses *Session) Collection(model kolektor.Modeler) (*Collection, error) {
	if err := ses.store.InitCollection(model); err != nil {
		return nil, err
	}

	return newCollection(ses)
}

// RemoveCollection will destroy the collection baed on the provided model.
// Without warning, without remorse. If you had no backups, and you did this
// by mistake, you can consider yourself screwed.
func (ses *Session) RemoveCollection(model kolektor.Modeler) error {
	return ses.store.RemoveCollection(model)
}

// Connection returns a connection to the store in use by this session.
// The caller is responsible for type asserting the result to the appropriated
// type for this store.
func (ses *Session) Connection(ctx context.Context) (any, error) {
	return ses.store.Connection(ctx)
}
