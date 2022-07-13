// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package kolekto

import (
	"context"
	"fmt"
	"strings"

	"github.com/golistic/kolekto/stores"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func init() {
	prepareStore[stores.PgSQL] = preparePostgreSQL
}

func preparePostgreSQL(dsn string) error {
	pgConfig, err := pgconn.ParseConfig(dsn)
	if testErr != nil {
		return err
	}
	if pgConfig.Database == "" {
		return fmt.Errorf("database in DSN must not be empty")
	}

	baseDSN := strings.Replace(dsn, "/"+pgConfig.Database, "/postgres", -1)

	conn, err := pgx.Connect(context.Background(), baseDSN)
	if err != nil {
		return err
	}

	if _, err = conn.Exec(context.Background(), "DROP DATABASE IF EXISTS "+pgConfig.Database); err != nil {
		return fmt.Errorf("pgsql: %s", err)
	}

	if _, err = conn.Exec(context.Background(), "CREATE DATABASE "+pgConfig.Database); err != nil {
		return err
	}

	testAllDSN[stores.PgSQL] = dsn
	return nil
}
