// Copyright (c) 2022, Geert JM Vanderkelen

package dbpgsql

import (
	"fmt"
	"kolekto/internal/stores"
)

const dmlReturningMeta = "id, uid, created, updated"

const pgsqlMetaAsJson = "jsonb_build_object('Meta', jsonb_build_object(" +
	"'id', id::numeric, " +
	"'uid', uid, " +
	"'created', created, " +
	"'updated', updated))"

const pgsqlMergeDataMeta = "data || " + pgsqlMetaAsJson

const PostgreSQLFunctions = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION updated_now() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = NOW();
    RETURN NEW; 
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION default_uid() RETURNS TRIGGER AS $$
BEGIN
	NEW.uid := COALESCE(NEW.uid, uuid_generate_v4()::text);
	RETURN NEW;
END;
$$ language 'plpgsql';
`

func ddlTable(name string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
id serial NOT NULL PRIMARY KEY,
uid VARCHAR(%d) NOT NULL,
created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated TIMESTAMPTZ DEFAULT NULL,
data JSONB
)`, name, stores.SizeUID)
}
