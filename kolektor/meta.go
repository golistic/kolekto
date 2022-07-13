// Copyright (c) 2022, Geert JM Vanderkelen

package kolektor

import "time"

// Meta stores metadata which each object within a collection
// must provide.
// Note that metadata is not stored within the actual JSON document.
type Meta struct {
	ID      int64      `json:"id"`
	UID     string     `json:"uid"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
}
