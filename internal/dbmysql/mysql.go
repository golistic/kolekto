// Copyright (c) 2022, Geert JM Vanderkelen

package dbmysql

import (
	"context"
	"database/sql"
	"fmt"
	"kolekto/internal/stores"
)

const dmlReturningMeta = "id, uid, created, updated"

const mysqlMetaAsJson = "JSON_OBJECT('Meta', JSON_OBJECT(" +
	"'id', id, " +
	"'uid', uid, " +
	"'created', DATE_FORMAT('%Y-%m-%dT%H:%m:%i.%fZ', created), " +
	"'updated', DATE_FORMAT('%Y-%m-%dT%H:%m:%i.%fZ', updated)))"

const mysqlMergeDataMeta = "JSON_MERGE(data, " + mysqlMetaAsJson + ")"

func ddlTable(name string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
uid VARCHAR(%d) NOT NULL,
created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
updated TIMESTAMP(6) NULL ON UPDATE CURRENT_TIMESTAMP(6),
data JSON
)`, name, stores.SizeUID)
}

func mysqlRoutineVersion(db *sql.Conn, routine string) (int, error) {
	q := "SELECT JSON_EXTRACT(ROUTINE_COMMENT, '$.version') " +
		"FROM information_schema.ROUTINES " +
		"WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_NAME = ?"

	var version int

	if err := db.QueryRowContext(context.Background(), q, routine).Scan(&version); err != nil {
		return 0, err
	}

	return version, nil
}

type sqlRoutine struct {
	version int
	ddl     string
}

var mysqlRoutines = map[string]sqlRoutine{
	"uuid_generate_v4": {
		version: 1,
		// based on https://stackoverflow.com/questions/32965743/how-to-generate-a-uuidv4-in-mysql
		ddl: `
CREATE FUNCTION IF NOT EXISTS {{.Name}}() RETURNS CHAR(36) NO SQL
	COMMENT '{"version": {{.Version}}}'
BEGIN
    SET @p3 = CONCAT('4', SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3));
    SET @p4 = CONCAT(HEX(FLOOR(ASCII(RANDOM_BYTES(1)) / 64) + 8), SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3));
    RETURN LOWER(CONCAT(
            HEX(RANDOM_BYTES(4)), '-', HEX(RANDOM_BYTES(2)), '-', @p3, '-', @p4, '-', HEX(RANDOM_BYTES(6))
        ));
END`,
	},
	"default_uid": {
		version: 1,
		ddl: `
CREATE FUNCTION IF NOT EXISTS {{.Name}}() RETURNS CHAR(40) NO SQL
	COMMENT '{"version": {{.Version}}}'
BEGIN
	RETURN uuid_generate_v4();
END`,
	},
}
