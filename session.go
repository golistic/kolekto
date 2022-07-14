// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
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
