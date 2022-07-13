// Copyright (c) 2022, Geert JM Vanderkelen

package stores

type StoreKind int

// Supported data stores.
const (
	MySQL StoreKind = iota + 1
	PgSQL
)
