// Copyright (c) 2022, Geert JM Vanderkelen

package dbmysql

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/xstrings"
)

func md5sum[T string | []byte](value T) string {
	sum := md5.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

func addIndexes(conn *sql.Conn, idxer kolektor.Indexer, tableName string) error {
	var alters []string

	haveIndexes, err := getIndexes(conn, tableName)
	if err != nil {
		return err
	}

	var wantIndexes []string
	for _, idx := range idxer.Indexes(kolektor.MySQL) {
		wantIndexes = append(wantIndexes, idx.Name)
		unique := ""
		if idx.Unique {
			unique = "UNIQUE"
		}

		exprSum := md5sum(idx.Expression)
		if haveHash, have := haveIndexes[idx.Name]; have {
			if exprSum == haveHash {
				// index did not change; skip
				continue
			} else {
				// index changed; recreate it by dropping it first
				alters = append(alters, fmt.Sprintf("DROP INDEX %s", idx.Name))
			}
		}

		alters = append(alters, fmt.Sprintf("ADD %s INDEX %s %s COMMENT 'kolekto#%s'",
			unique, idx.Name, idx.Expression, exprSum))
	}

	for name := range haveIndexes {
		if xstrings.Search(wantIndexes, name) == -1 {
			alters = append(alters, fmt.Sprintf("DROP INDEX %s", name))
		}
	}

	if len(alters) > 0 {
		dml := "ALTER TABLE " + tableName + " " + strings.Join(alters, ", ")
		if _, err := conn.ExecContext(context.Background(), dml); err != nil {
			return fmt.Errorf("failed creating indexes for %s (%w)", tableName, err)
		}
	}
	return nil
}

func getIndexes(conn *sql.Conn, tableName string) (map[string]string, error) {
	q := "SELECT INDEX_NAME, INDEX_COMMENT FROM INFORMATION_SCHEMA.STATISTICS" +
		" WHERE TABLE_SCHEMA = DATABASE() AND" +
		" TABLE_NAME = ? AND INDEX_NAME <> 'PRIMARY' AND INDEX_COMMENT LIKE 'kolekto#%'"

	rows, err := conn.QueryContext(context.Background(), q, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed getting indexes (%w)", err)
	}

	indexes := map[string]string{}

	for rows.Next() {
		var name string
		var comment string
		if err := rows.Scan(&name, &comment); err != nil {
			return nil, fmt.Errorf("failed getting indexes (%w)", err)
		}
		parts := strings.Split(comment, "#")
		if len(parts) != 2 {
			return nil, fmt.Errorf("failed getting indexes (bad index comment)")
		}
		indexes[name] = parts[1]
	}

	return indexes, nil
}
