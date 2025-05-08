package db

import (
	"context"
	"database/sql"

	"healthmonitor/internal/db/migrations"

	"github.com/pressly/goose/v3"
)

type DB struct {
	sqldb *sql.DB
}

func New(url string) (*DB, error) {
	var err error
	db := &DB{}

	db.sqldb, err = sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	register := []*goose.Migration{
		goose.NewGoMigration(
			1,
			&goose.GoFunc{RunTx: migrations.Up00001},
			&goose.GoFunc{RunTx: migrations.Down00001},
		),
	}

	goose.SetLogger(LogAdapter{})
	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db.sqldb,
		nil,
		goose.WithGoMigrations(register...),
	)
	if err != nil {
		return nil, err
	}

	_, err = provider.Up(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.sqldb.Close()
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.sqldb.Exec(query, args...)
}
