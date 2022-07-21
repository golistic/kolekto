// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"reflect"

	"github.com/golistic/kolekto/kolektor"
)

// Collection manages a JSON collection.
type Collection struct {
	ses *Session
}

func newCollection(kol *Session) (*Collection, error) {
	if kol == nil {
		panic("ses must not be nil")
	}

	coll := &Collection{
		ses: kol,
	}

	return coll, nil
}

// Get retrieves an object from the collection.
// The uid can be either an integer (int64, int) or a string. The former will
// use the model ID field, the latter the UID field.
func (coll *Collection) Get(obj kolektor.Modeler, uid any) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &kolektor.InvalidObjectError{Type: reflect.TypeOf(obj)}
	}

	var field string
	switch uid.(type) {
	case int64, int:
		field = "id"
	case string:
		field = "uid"
	}

	return coll.GetByFields(obj, map[string]any{field: uid})
}

func (coll *Collection) GetByFields(obj kolektor.Modeler, fields map[string]any) error {
	return coll.ses.store.GetObject(obj, fields)
}

// Store stores an object into the collection.
func (coll *Collection) Store(obj kolektor.Modeler) error {
	var meta *kolektor.Meta

	var err error
	meta, err = coll.ses.store.StoreObject(obj)
	if err != nil {
		return err
	}

	obj.SetMeta(meta)

	return nil
}
