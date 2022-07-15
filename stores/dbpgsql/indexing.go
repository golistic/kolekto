// Copyright (c) 2022, Geert JM Vanderkelen

package dbpgsql

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/xstrings"
	"github.com/jackc/pgx/v4/pgxpool"
)

func md5sum[T string | []byte](value T) string {
	sum := md5.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

func addIndexes(conn *pgxpool.Conn, idxer kolektor.Indexer, tableName string) error {
	haveIndexes, err := getIndexes(conn, tableName)
	if err != nil {
		return err
	}

	var wantIndexes []string
	for _, idx := range idxer.Indexes(kolektor.PgSQL) {
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
				dml := fmt.Sprintf("DROP INDEX %s", idx.Name)
				if _, err := conn.Exec(context.Background(), dml); err != nil {
					return fmt.Errorf("failed dropping index %s (%w)", idx.Name, err)
				}
			}
		}

		dml := fmt.Sprintf("CREATE %s INDEX %s ON %s %s; COMMENT ON INDEX %s IS 'kolekto#%s'",
			unique, idx.Name, tableName, idx.Expression, idx.Name, exprSum)

		if _, err := conn.Exec(context.Background(), dml); err != nil {
			return fmt.Errorf("failed creating index %s (%w)", idx.Name, err)
		}
	}

	for name := range haveIndexes {
		if xstrings.Search(wantIndexes, name) == -1 {
			dml := fmt.Sprintf("DROP INDEX %s", name)
			if _, err := conn.Exec(context.Background(), dml); err != nil {
				return fmt.Errorf("failed dropping index %s (%w)", name, err)
			}
		}
	}

	return nil
}

func getIndexes(conn *pgxpool.Conn, tableName string) (map[string]string, error) {
	q := "SELECT indexrelname, description" +
		" FROM pg_catalog.pg_stat_all_indexes as idx" +
		" LEFT JOIN pg_catalog.pg_description ON idx.indexrelid = pg_description.objoid" +
		" WHERE relname = $1 AND schemaname = \"current_schema\"() AND" +
		" description LIKE 'kolekto#%'"

	rows, err := conn.Query(context.Background(), q, tableName)
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
