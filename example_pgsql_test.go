// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package kolekto_test

import (
	"github.com/golistic/kolekto/kolektor"
	_ "github.com/golistic/kolekto/stores/dbpgsql" // register store
)

func Example_pgsql() {
	// note: the connection string reflect the PostgreSQL instance running for
	// testing the Kolekto package. To start: change into _support/docker-compose and
	// do `docker compose up -d`.
	// We assume the 'music' database exists.
	dsn := "user=postgres password=postgres host=localhost dbname=music port=5438"

	// actual example is same for each store; please check example_common_test.go
	exampleStoreRetrieveBand(kolektor.MySQL, dsn)

	// Output:
	// UID    : f5dea144-caac-4735-a521-34a82b12f20b
	// Band   : A Tribe Called Quest
	// Members:
	//  - Q-Tip
	//  - Phife Dwag
	//  - Ali Shaheed Muhammad
	//  - Jarobi White
}
