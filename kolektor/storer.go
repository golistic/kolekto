// Copyright (c) 2022, Geert JM Vanderkelen

package kolektor

import "context"

// Storer defines methods which must be implemented by data
// stores types.
type Storer interface {
	Name() string
	GetObject(obj Modeler, field string, value any) error
	StoreObject(obj Modeler) (*Meta, error)
	RemoveCollection(model Modeler) error
	InitCollection(model Modeler) error
	Connection(ctx context.Context) (any, error)
}

type Indexer interface {
	Indexes(kind StoreKind) []Index
}

type Index struct {
	Name       string
	Unique     bool
	Expression string
}
