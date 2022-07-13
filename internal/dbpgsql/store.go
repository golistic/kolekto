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
	stores.Register(stores.PgSQL, New)
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
		return nil, fmt.Errorf("failed checking store connection (%s)", err)
	}

	if _, err := s.pool.Exec(context.Background(), PostgreSQLFunctions); err != nil {
		return nil, err
	}

	return s, nil
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
		return err
	}
	if err := pgxscan.Get(context.Background(), conn, &obj, q, value); err != nil {
		return fmt.Errorf("failed getting object (%s)", err)
	}

	return err
}

// StoreObject stores obj into the collection of the object's model.
func (s *Store) StoreObject(obj kolektor.Modeler) (*kolektor.Meta, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	objID := obj.GetID()
	objUID := obj.GetUID()
	obj.SetMeta(nil) // we do not save Meta in the JSON document

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var row pgx.Row

	var action string
	if objID == 0 {
		q := fmt.Sprintf("INSERT INTO %s (data, uid) VALUES ($1, NULLIF($2, '')) "+
			"RETURNING "+dmlReturningMeta,
			obj.CollectionName())
		row = conn.QueryRow(context.Background(), q, data, objUID)
		action = "inserting"
	} else {
		q := fmt.Sprintf("UPDATE %s SET data = $1, uid = NULLIF($2, '') "+
			"WHERE id = $3 RETURNING "+dmlReturningMeta,
			obj.CollectionName())
		row = conn.QueryRow(context.Background(), q, data, objUID, objID)
		action = "updating"
	}

	meta := &kolektor.Meta{}
	if err := row.Scan(&meta.ID, &meta.UID, &meta.Created, &meta.Updated); err != nil {
		return nil, fmt.Errorf("failed %s object (%s)", action, err)
	}

	return meta, err
}

// InitCollection initializes the model's collection.
func (s *Store) InitCollection(model kolektor.Modeler) error {
	tableName := model.CollectionName()

	ddl := ddlTable(tableName)

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed acquiring data store connection (%s)", err)
	}

	// CREATE TABLE
	if _, err := conn.Exec(context.Background(), ddl); err != nil {
		return err
	}

	// CREATE TRIGGERs
	tr := fmt.Sprintf(`CREATE OR REPLACE TRIGGER tr_%s_updated
BEFORE UPDATE ON %s FOR EACH ROW EXECUTE PROCEDURE updated_now()`,
		tableName, tableName)
	if _, err := conn.Exec(context.Background(), tr); err != nil {
		return err
	}

	tr = fmt.Sprintf(`CREATE OR REPLACE TRIGGER tr_%s_uid
BEFORE INSERT OR UPDATE ON %s FOR EACH ROW EXECUTE PROCEDURE default_uid()`,
		tableName, tableName)
	if _, err := conn.Exec(context.Background(), tr); err != nil {
		return err
	}

	return nil
}

// RemoveCollection removes the model's collection.
func (s *Store) RemoveCollection(model kolektor.Modeler) error {
	ddl := fmt.Sprintf("DROP TABLE IF EXISTS %s", model.CollectionName())

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed acquiring data store connection (%s)", err)
	}

	if _, err := conn.Exec(context.Background(), ddl); err != nil {
		return err
	}

	return nil
}
