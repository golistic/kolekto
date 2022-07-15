// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package kolekto

import (
	"github.com/golistic/kolekto/internal/ytest"
	"github.com/golistic/kolekto/kolektor"
)

func init() {
	prepareStore[kolektor.PgSQL] = ytest.PreparePostgreSQL
}
