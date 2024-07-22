package sql

import (
	"context"
	"database/sql"

	"github.com/unchainio/pkg/iferr"

	"github.com/unchainio/interfaces/logger"

	"github.com/unchainio/pkg/errors"
)

type Driver string

const DriverPostgres Driver = "postgres"
const DriverMySQL Driver = "mysql"

type DB struct {
	Log logger.Logger
	*sql.DB
}

func NewDB(log logger.Logger, cfg *Config) (*DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.ConnectionString)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &DB{
		DB:  db,
		Log: log,
	}, nil
}

func (r *DB) BeginForRole(role string) (context.Context, *sql.Tx, error) {
	tx, err := r.DB.Begin()

	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	_, err = tx.Exec("SET LOCAL SESSION AUTHORIZATION $1", role)

	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	return context.Background(), tx, nil
}

func (r *DB) Begin() (context.Context, *sql.Tx, error) {
	tx, err := r.DB.Begin()

	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	return context.Background(), tx, nil
}

type TxFn func(ctx context.Context, tx *sql.Tx) error

func (r *DB) WrapTx(fn TxFn) (err error) {
	ctx, tx, err := r.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil { // if fn(ctx, tx) panicked, rollback and re-panic
			err2 := tx.Rollback()
			iferr.Warn(errors.Wrap(err2, ""))
			panic(p)
		} else if err != nil { // the original error is unchanged and returned
			err2 := tx.Rollback()
			iferr.Warn(errors.Wrap(err2, ""))
		} else {
			err = tx.Commit() // if tx.Commit() fails, set err to its error so that WrapTx returns it
			err = errors.Wrap(err, "")
		}
	}()
	// We can't just `return fn(ctx, tx)` because we need to check the value of err in the defer
	err = fn(ctx, tx)

	return err
}

func (r *DB) BeginTxForRole(ctx context.Context, opts *sql.TxOptions, role string) (*sql.Tx, error) {
	tx, err := r.DB.BeginTx(ctx, opts)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	_, err = tx.Exec("SET LOCAL SESSION AUTHORIZATION $1", role)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return tx, nil
}

func (r *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	tx, err := r.DB.BeginTx(ctx, opts)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return tx, nil
}
