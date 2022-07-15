// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package kolekto

import (
	"context"
	"strings"
	"testing"

	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/kolekto/stores/dbpgsql"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNew_pgsql(t *testing.T) {
	t.Run("test PostgreSQL using pgxpool", func(t *testing.T) {
		session, err := NewSession(kolektor.PgSQL, testAllDSN[kolektor.PgSQL])
		xt.OK(t, err)
		_, ok := session.store.(*dbpgsql.Store)
		xt.Assert(t, ok, "expected *dbpgsql.Store")
	})
}

func TestCollection_Store_pgsql(t *testing.T) {
	session, err := NewSession(kolektor.PgSQL, testAllDSN[kolektor.PgSQL])
	xt.OK(t, err)

	testCollection_Store(t, session)
}

func TestSession_Connection_pgsql(t *testing.T) {
	session, err := NewSession(kolektor.PgSQL, testAllDSN[kolektor.PgSQL])
	xt.OK(t, err)
	c, err := session.Connection(context.Background())
	xt.OK(t, err)
	conn, ok := c.(*pgxpool.Conn)
	xt.Assert(t, ok, "expected *pgxpool.Conn")

	var version string
	xt.OK(t, conn.QueryRow(context.Background(), "SELECT version()").Scan(&version))
	xt.Assert(t, strings.Contains(version, "PostgreSQL"))
}
