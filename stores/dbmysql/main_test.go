// Copyright (c) 2022, Geert JM Vanderkelen

package dbmysql

import (
	"fmt"
	"os"
	"testing"

	"github.com/golistic/kolekto/internal/ytest"
)

// following defaults assume Docker containers are started using the
// Docker Compose configuration found in _support/docker-compose
const defaultMySQLDSN = "root:mysql@tcp(localhost:3360)/kolekto_test_store"

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

	if v, have := os.LookupEnv("TEST_MYSQL_DSN"); have {
		testDSN = v
	} else {
		testDSN = defaultMySQLDSN
	}

	testDSN, testErr = ytest.PrepareMySQL(testDSN)
	if testErr != nil {
		return
	}

	testExitCode = m.Run()
}
