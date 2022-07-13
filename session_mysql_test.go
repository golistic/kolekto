// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"testing"

	"github.com/golistic/kolekto/stores"
	"github.com/golistic/kolekto/stores/dbmysql"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNew(t *testing.T) {
	t.Run("test MySQL", func(t *testing.T) {
		session, err := NewSession(stores.MySQL, testAllDSN[stores.MySQL])
		xt.OK(t, err)
		_, ok := session.store.(*dbmysql.Store)
		xt.Assert(t, ok, "expected *dbmysql.Store")
	})
}

func TestCollection_Store_mysql(t *testing.T) {
	session, err := NewSession(stores.MySQL, testAllDSN[stores.MySQL])
	xt.OK(t, err)

	testCollection_Store(t, session)
}
