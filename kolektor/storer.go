// Copyright (c) 2022, Geert JM Vanderkelen

package kolektor

// Storer defines methods which must be implemented by data
// stores types.
type Storer interface {
	Name() string
	GetObject(obj Modeler, field string, value any) error
	StoreObject(obj Modeler) (*Meta, error)
	RemoveCollection(model Modeler) error
	InitCollection(model Modeler) error
}
