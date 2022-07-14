// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package ytest

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func PreparePostgreSQL(dsn string) (string, error) {
	pgConfig, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return "", err
	}
	if pgConfig.Database == "" {
		return "", fmt.Errorf("database in DSN must not be empty")
	}

	baseDSN := strings.Replace(dsn, "/"+pgConfig.Database, "/postgres", -1)

	conn, err := pgx.Connect(context.Background(), baseDSN)
	if err != nil {
		return "", err
	}

	if _, err = conn.Exec(context.Background(), "DROP DATABASE IF EXISTS "+pgConfig.Database); err != nil {
		return "", fmt.Errorf("pgsql: %s", err)
	}

	if _, err = conn.Exec(context.Background(), "CREATE DATABASE "+pgConfig.Database); err != nil {
		return "", err
	}

	return dsn, nil
}
