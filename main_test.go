// Copyright (c) 2022, Geert JM Vanderkelen

package kolekto

import (
	"fmt"
	"os"
	"testing"

	"github.com/golistic/kolekto/stores"
)

// following defaults assume Docker containers are started using the
// Docker Compose configuration found in _support/docker-compose
const defaultPgSQLDSN = "postgres://postgres:postgres@localhost:5438/kolekto_test"
const defaultMySQLDSN = "root:mysql@tcp(localhost:3360)/kolekto_test"

var (
	testExitCode int
	testErr      error
	testAllDSN   = map[stores.StoreKind]string{}
	prepareStore = map[stores.StoreKind]func(dsn string) error{}
)

func testTearDown() {
	if testErr != nil {
		fmt.Println(testErr)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	defer func() { os.Exit(testExitCode) }()
	defer testTearDown()

	if v, have := os.LookupEnv("TEST_PGSQL_DSN"); have {
		testAllDSN[stores.PgSQL] = v
	} else {
		testAllDSN[stores.PgSQL] = defaultPgSQLDSN
	}

	if v, have := os.LookupEnv("TEST_MYSQL_DSN"); have {
		testAllDSN[stores.MySQL] = v
	} else {
		testAllDSN[stores.MySQL] = defaultMySQLDSN
	}

	for kind, prep := range prepareStore {
		if testErr = prep(testAllDSN[kind]); testErr != nil {
			return
		}
	}

	testExitCode = m.Run()
}

func testDBPrefix(ses *Session) string {
	return ses.store.Name()
}
