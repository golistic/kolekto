// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package dbpgsql

import (
	"fmt"
	"os"
	"testing"

	"github.com/golistic/kolekto/internal/ytest"
)

// following defaults assume Docker containers are started using the
// Docker Compose configuration found in _support/docker-compose
const defaultPgSQLDSN = "postgres://postgres:postgres@localhost:5438/kolekto_test_store"

var (
	testExitCode int
	testErr      error
	testDSN      string
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
		testDSN = v
	} else {
		testDSN = defaultPgSQLDSN
	}

	testDSN, testErr = ytest.PreparePostgreSQL(testDSN)
	if testErr != nil {
		return
	}

	testExitCode = m.Run()
}
