// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nomysql

package kolekto

import (
	"github.com/golistic/kolekto/internal/ytest"
	"github.com/golistic/kolekto/kolektor"
)

func init() {
	prepareStore[kolektor.MySQL] = ytest.PrepareMySQL
}
