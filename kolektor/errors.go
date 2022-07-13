// Copyright (c) 2022, Geert JM Vanderkelen

package kolektor

import "reflect"

// InvalidObjectError describes an invalid argument passed to functions
// that require a kolektor.Modeler that must be a non-nil pointer.
type InvalidObjectError struct {
	Type reflect.Type
}

func (e *InvalidObjectError) Error() string {
	if e.Type == nil {
		return "kolekto: object must not be nil"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "kolekto: object is non-pointer " + e.Type.String()
	}
	return "kolekto: object must not be nil " + e.Type.String()
}
