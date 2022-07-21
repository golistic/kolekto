// Copyright (c) 2022, Geert JM Vanderkelen

package kolektor

import "fmt"

type StoreKind int

// Supported data stores.
const (
	MySQL StoreKind = iota + 1
	PgSQL
)

func (sk StoreKind) String() string {
	switch sk {
	case MySQL:
		return "MySQL"
	case PgSQL:
		return "PostgreSQL"
	default:
		return fmt.Sprintf("StoreNameMissing{%d}", sk)
	}
}
