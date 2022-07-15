// Copyright (c) 2022, Geert JM Vanderkelen

//go:build !nopgsql

package dbpgsql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/kolekto/stores"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store defines the PostgreSQL backed data store.
type Store struct {
	pool *pgxpool.Pool
}

var _ kolektor.Storer = &Store{}

func init() {
	stores.Register(kolektor.PgSQL, New)
}

// New instantiates a PostgreSQL backed data store.
func New(dsn string) (kolektor.Storer, error) {
	var err error
	s := &Store{}

	s.pool, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if _, err := s.pool.Acquire(context.Background()); err != nil {
		return nil, fmt.Errorf("failed checking store connection (%w)", err)
	}

	if _, err := s.pool.Exec(context.Background(), PostgreSQLFunctions); err != nil {
		return nil, fmt.Errorf("failed checking store connection (%w)", err)
	}

	return s, nil
}

// mustSQLConn is mainly for testing.
// Panics on errors.
func (s *Store) mustConn() *pgxpool.Conn {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	return conn
}

// Name returns the name of the data store.
func (s *Store) Name() string {
	return "PostgreSQL"
}

// GetObject retrieves a stored object and stores it in obj.
func (s *Store) GetObject(obj kolektor.Modeler, field string, value any) error {
	collName := obj.CollectionName()

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", pgsqlMergeDataMeta, collName, field)

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed getting object (%w)", err)
	}
	if err := pgxscan.Get(context.Background(), conn, &obj, q, value); err != nil {
		if err == pgx.ErrNoRows {
			return stores.ErrNoObject{Name: obj.CollectionName()}
		}
		return fmt.Errorf("failed getting object (%w)", err)
	}

	return err
}

// StoreObject stores obj into the collection of the object's model.
func (s *Store) StoreObject(obj kolektor.Modeler) (*kolektor.Meta, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed storing object (%w)", err)
	}

	objID := obj.GetID()
	objUID := obj.GetUID()
	obj.SetMeta(nil) // we do not save Meta in the JSON document

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed storing object (%w)", err)
	}

	var row pgx.Row

	if objID == 0 {
		q := fmt.Sprintf("INSERT INTO %s (data, uid) VALUES ($1, NULLIF($2, '')) "+
			"RETURNING "+dmlReturningMeta,
			obj.CollectionName())
		row = conn.QueryRow(context.Background(), q, data, objUID)
	} else {
		q := fmt.Sprintf("UPDATE %s SET data = $1, uid = NULLIF($2, '') "+
			"WHERE id = $3 RETURNING "+dmlReturningMeta,
			obj.CollectionName())
		row = conn.QueryRow(context.Background(), q, data, objUID, objID)
	}

	meta := &kolektor.Meta{}
	if err := row.Scan(&meta.ID, &meta.UID, &meta.Created, &meta.Updated); err != nil {
		return nil, fmt.Errorf("failed storing object (%w)", err)
	}

	return meta, err
}

// InitCollection initializes the model's collection.
func (s *Store) InitCollection(model kolektor.Modeler) error {
	tableName := model.CollectionName()

	ddl := ddlTable(tableName)

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed initializing collection (%w)", err)
	}

	// CREATE TABLE
	if _, err := conn.Exec(context.Background(), ddl); err != nil {
		return fmt.Errorf("failed initializing collection (%w)", err)
	}

	// CREATE TRIGGERs
	tr := fmt.Sprintf(`CREATE OR REPLACE TRIGGER tr_%s_updated
BEFORE UPDATE ON %s FOR EACH ROW EXECUTE PROCEDURE updated_now()`,
		tableName, tableName)
	if _, err := conn.Exec(context.Background(), tr); err != nil {
		return fmt.Errorf("failed initializing collection (%w)", err)
	}

	tr = fmt.Sprintf(`CREATE OR REPLACE TRIGGER tr_%s_uid
BEFORE INSERT OR UPDATE ON %s FOR EACH ROW EXECUTE PROCEDURE default_uid()`,
		tableName, tableName)
	if _, err := conn.Exec(context.Background(), tr); err != nil {
		return fmt.Errorf("failed initializing collection (%w)", err)
	}

	// INDEXING
	if idxer, ok := model.(kolektor.Indexer); ok {
		if err := addIndexes(conn, idxer, tableName); err != nil {
			return err
		}
	}

	return nil
}

// RemoveCollection removes the model's collection.
func (s *Store) RemoveCollection(model kolektor.Modeler) error {
	ddl := fmt.Sprintf("DROP TABLE IF EXISTS %s", model.CollectionName())

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed removing collection (%w)", err)
	}

	if _, err := conn.Exec(context.Background(), ddl); err != nil {
		return fmt.Errorf("failed removing collection (%w)", err)
	}

	return nil
}
