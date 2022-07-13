// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nomysql

package kolekto

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"kolekto/internal/stores"
)

func init() {
	prepareStore[stores.MySQL] = prepareMySQL
}

func prepareMySQL(dsn string) error {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}
	dbName := config.DBName
	config.DBName = ""
	if config.Params == nil {
		config.Params = map[string]string{}
	}
	config.Params["parseTime"] = "true"

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	if _, err = db.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+dbName); err != nil {
		return fmt.Errorf("mysql: %s", err)
	}

	if _, err = db.ExecContext(context.Background(), "CREATE DATABASE "+dbName); err != nil {
		return err
	}

	config.DBName = dbName
	testAllDSN[stores.MySQL] = config.FormatDSN()
	return nil
}
