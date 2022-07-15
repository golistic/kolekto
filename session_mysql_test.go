// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/kolekto/stores/dbmysql"
)

func TestNew(t *testing.T) {
	t.Run("test MySQL", func(t *testing.T) {
		session, err := NewSession(kolektor.MySQL, testAllDSN[kolektor.MySQL])
		xt.OK(t, err)
		_, ok := session.store.(*dbmysql.Store)
		xt.Assert(t, ok, "expected *dbmysql.Store")
	})
}

func TestCollection_Store_mysql(t *testing.T) {
	session, err := NewSession(kolektor.MySQL, testAllDSN[kolektor.MySQL])
	xt.OK(t, err)

	testCollection_Store(t, session)
}

func TestSession_Connection_mysql(t *testing.T) {
	session, err := NewSession(kolektor.MySQL, testAllDSN[kolektor.MySQL])
	xt.OK(t, err)
	c, err := session.Connection(context.Background())
	xt.OK(t, err)
	conn, ok := c.(*sql.Conn)
	xt.Assert(t, ok, "expected *sql.Conn")

	var version string
	q := "SELECT VARIABLE_VALUE FROM performance_schema.global_variables WHERE VARIABLE_NAME = 'version_comment'"
	xt.OK(t, conn.QueryRowContext(context.Background(), q).Scan(&version))
	xt.Assert(t, strings.Contains(version, "MySQL"), version)
}
