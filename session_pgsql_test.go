// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package kolekto

import (
	"testing"

	"github.com/golistic/kolekto/stores"
	"github.com/golistic/kolekto/stores/dbpgsql"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNew_pgsql(t *testing.T) {
	t.Run("test PostgreSQL using pgxpool", func(t *testing.T) {
		session, err := NewSession(stores.PgSQL, testAllDSN[stores.PgSQL])
		xt.OK(t, err)
		_, ok := session.store.(*dbpgsql.Store)
		xt.Assert(t, ok, "expected *dbpgsql.Store")
	})
}

func TestCollection_Store_pgsql(t *testing.T) {
	session, err := NewSession(stores.PgSQL, testAllDSN[stores.PgSQL])
	xt.OK(t, err)

	testCollection_Store(t, session)
}
