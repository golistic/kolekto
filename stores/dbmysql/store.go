// Copyright (c) 2022, Geert JM Vanderkelen

package dbmysql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-sql-driver/mysql"
	"github.com/golistic/kolekto/kolektor"
	"github.com/golistic/kolekto/stores"
	"github.com/golistic/xstrings"
)

// Store defines the MySQL backed data store.
type Store struct {
	pool *sql.DB
}

var _ kolektor.Storer = &Store{}

func init() {
	stores.Register(kolektor.MySQL, New)
}

// New instantiates a PostgreSQL backed data store.
func New(dsn string) (kolektor.Storer, error) {
	var err error

	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed checking store connection (%w)", err)
	}
	if config.Params == nil {
		config.Params = map[string]string{}
	}
	config.Params["parseTime"] = "true"

	s := &Store{}
	s.pool, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed checking store connection (%w)", err)
	}

	if err := s.pool.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed checking store connection (%w)", err)
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("failed checking store connection (%w)", err)
	}

	return s, nil
}

// Connection returns a connection to the store. The caller is responsible
// for type asserting the result to the appropriated type for this store,
// namely *sql.Conn.
func (s *Store) Connection(ctx context.Context) (any, error) {
	return s.connection(ctx)
}

func (s *Store) connection(ctx context.Context) (*sql.Conn, error) {
	conn, err := s.pool.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting store collection (%w)", err)
	}
	return conn, nil
}

// mustSQLConn is mainly for testing.
// Panics on errors.
func (s *Store) mustSQLConn() *sql.Conn {
	conn, err := s.connection(context.Background())
	if err != nil {
		panic(err)
	}
	return conn
}

// Name returns the name of the data store.
func (s *Store) Name() string {
	return "MySQL"
}

// GetObject retrieves a stored object and stores it in obj.
func (s *Store) GetObject(obj kolektor.Modeler, fieldMap kolektor.FieldMap) error {
	if len(fieldMap) == 0 {
		return fmt.Errorf("need at least one field to filter on")
	}

	var ands []string
	var values []any
	for name, value := range fieldMap {
		if !(strings.HasPrefix(name, "(") || xstrings.Search(stores.ReservedFields, name) != -1) {
			name = fmt.Sprintf("data->>'$.%s'", name)
		}
		ands = append(ands, name+" = ?")
		values = append(values, value)
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		mysqlMergeDataMeta, obj.CollectionName(), strings.Join(ands, " AND "))

	var data []byte
	if err := s.pool.QueryRowContext(context.Background(), q, values...).Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			return stores.ErrNoObject{Name: obj.CollectionName()}
		}
		return fmt.Errorf("failed getting object (%w)", err)
	}

	if err := json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("failed getting object (%w)", err)
	}

	return nil
}

// StoreObject stores obj into the collection of the object's model.
func (s *Store) StoreObject(obj kolektor.Modeler) (*kolektor.Meta, error) {
	objID := obj.GetID()
	objUID := obj.GetUID()
	obj.SetMeta(nil) // we do not save Meta in the JSON document

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed storing object (%w)", err)
	}

	var res sql.Result
	if objID == 0 {
		q := fmt.Sprintf("INSERT INTO %s (data, uid) VALUES (?, ?)", obj.CollectionName())
		var err error
		res, err = s.pool.ExecContext(context.Background(), q, data, objUID)
		if err != nil {
			return nil, fmt.Errorf("failed storing object (%w)", err)
		}
		objID, err = res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed storing object (%w)", err)
		}
	} else {
		q := fmt.Sprintf("UPDATE %s SET data = ?, uid = ? WHERE id = ?",
			obj.CollectionName())
		var err error
		res, err = s.pool.ExecContext(context.Background(), q, data, objUID, objID)
		if err != nil {
			return nil, fmt.Errorf("failed storing object (%w)", err)
		}
	}

	// second round-trip to fetch meta
	meta := &kolektor.Meta{}
	q := "SELECT " + dmlReturningMeta + " FROM " + obj.CollectionName() + " WHERE id = ?"
	row := s.pool.QueryRowContext(context.Background(), q, objID)
	if err := row.Scan(&meta.ID, &meta.UID, &meta.Created, &meta.Updated); err != nil {
		return nil, fmt.Errorf("failed storing object (%w)", err)
	}

	return meta, nil
}

func (s *Store) init() error {
	conn, err := s.pool.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("init MySQL store failed (%w)", err)
	}
	defer func() { _ = conn.Close() }()

	for name, r := range mysqlRoutines {
		v, err := mysqlRoutineVersion(conn, name)
		switch {
		case err == sql.ErrNoRows || v < r.version:
			if _, err := conn.ExecContext(context.Background(), "DROP FUNCTION IF EXISTS "+name); err != nil {
				return fmt.Errorf("init MySQL store failed (%w)", err)
			}

			tmpl, err := template.New("sql").Parse(r.ddl)
			if err != nil {
				return fmt.Errorf("init MySQL store failed (%w)", err)
			}

			var ddl bytes.Buffer
			if err := tmpl.Execute(&ddl, struct {
				Name    string
				Version int
			}{
				Name:    name,
				Version: r.version,
			}); err != nil {
				return fmt.Errorf("init MySQL store failed (%w)", err)
			}

			if _, err := conn.ExecContext(context.Background(), ddl.String()); err != nil {
				return fmt.Errorf("init MySQL store failed (%w)", err)
			}
		case err != nil:
			return fmt.Errorf("init MySQL store failed (%w)", err)
		}
	}

	return nil
}

// InitCollection initializes the model's collection.
func (s *Store) InitCollection(model kolektor.Modeler) error {
	conn, err := s.connection(context.Background())
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	tableName := model.CollectionName()

	// default for uid is set using trigger
	ddl := ddlTable(tableName)

	// CREATE TABLE
	if _, err := conn.ExecContext(context.Background(), ddl); err != nil {
		return fmt.Errorf("failed initializing collection (%w)", err)
	}

	// CREATE TRIGGERs
	tr := fmt.Sprintf(`CREATE TRIGGER IF NOT EXISTS tr_%s_updated
BEFORE INSERT ON %s FOR EACH ROW SET new.uid = IF(new.uid='', default_uid(), new.uid)`,
		tableName, tableName)
	if _, err := conn.ExecContext(context.Background(), tr); err != nil {
		return fmt.Errorf("failed initializing collection (%w)", err)
	}

	tr = fmt.Sprintf(`CREATE TRIGGER IF NOT EXISTS tr_%s_updated
BEFORE UPDATE ON %s FOR EACH ROW SET new.uid = IF(new.uid='', default_uid(), new.uid)`,
		tableName, tableName)
	if _, err := conn.ExecContext(context.Background(), tr); err != nil {
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

	if _, err := s.pool.ExecContext(context.Background(), ddl); err != nil {
		return fmt.Errorf("failed removing collection (%w)", err)
	}

	return nil
}
