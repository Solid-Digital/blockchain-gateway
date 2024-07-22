package sql_test

import (
	"context"
	stdsql "database/sql"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/iferr"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"github.com/unchainio/pkg/xlogger"
)

func TestDB_WrapTx(t *testing.T) {
	cases := map[string]struct {
		WrapTxFn         func(db *sql.DB) func(fn sql.TxFn) error
		Fn               sql.TxFn
		ExpectedRollback bool
		ExpectedCommit   bool
		ExpectedPanic    bool
		ExpectedError    bool
	}{
		"Commit after no errors and no panics in the wrapped function": {
			WrapTxFn,
			SuccessFn,
			false,
			true,
			false,
			false,
		},
		"Rollback after an error in the wrapped function": {
			WrapTxFn,
			ErrorFn,
			true,
			false,
			false,
			true,
		},
		"Rollback after a panic in the wrapped function": {
			WrapTxFn,
			PanicFn,
			true,
			false,
			true,
			false,
		},
		"Example of how a bad implementation of WrapTx would not rollback after a panic in the wrapped function": {
			BadWrapTxFn,
			PanicFn,
			false,
			false,
			true,
			false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectBegin()

			if tc.ExpectedRollback {
				mock.ExpectRollback()
			}

			if tc.ExpectedCommit {
				mock.ExpectCommit()
			}

			defer func() {
				p := recover()
				require.Equal(t, tc.ExpectedPanic, p != nil, "a panic should be observed only if tc.ExpectedPanic is true")

				if tc.ExpectedError {
					require.Error(t, err, "an error should be observed if tc.ExpectedError is true")
				} else {
					require.NoError(t, err, "an error should not be observed if tc.ExpectedError is false")
				}
				require.NoError(t, mock.ExpectationsWereMet())
			}()

			err = tc.WrapTxFn(&sql.DB{
				Log: xlogger.NewSimpleLogger(),
				DB:  db,
			})(tc.Fn)
		})
	}
}

func SuccessFn(ctx context.Context, tx *stdsql.Tx) error {
	return nil
}

func ErrorFn(ctx context.Context, tx *stdsql.Tx) error {
	return errors.New("error")
}

func PanicFn(ctx context.Context, tx *stdsql.Tx) error {
	panic("panic")
}

func WrapTxFn(db *sql.DB) func(fn sql.TxFn) error {
	return func(fn sql.TxFn) error {
		return db.WrapTx(fn)
	}
}

// This implementation does not rollback on panics
func BadWrapTxFn(db *sql.DB) func(fn sql.TxFn) error {
	return func(fn sql.TxFn) error {
		ctx, tx, err := db.Begin()
		if err != nil {
			return err
		}

		err = fn(ctx, tx)

		if err != nil { // the original error is unchanged and returned
			err2 := tx.Rollback()
			iferr.Warn(errors.Wrap(err2, ""))
			return err
		}

		err = tx.Commit() // if tx.Commit() fails, set err to its error so that WrapTx returns it
		err = errors.Wrap(err, "")

		return err
	}
}

func TestDB_Begin(t *testing.T) {
	cases := map[string]struct {
	}{
		"default": {},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_ = tc
			mockdb, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectBegin()

			db := &sql.DB{
				Log: xlogger.NewSimpleLogger(),
				DB:  mockdb,
			}

			ctx, tx, err := db.Begin()
			require.NoError(t, err)
			require.NotNil(t, ctx)
			require.NotNil(t, tx)
		})
	}
}

func TestDB_BeginTx(t *testing.T) {
	cases := map[string]struct {
		CTX  context.Context
		Opts *stdsql.TxOptions
	}{
		"default": {
			context.Background(),
			nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockdb, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectBegin()

			db := &sql.DB{
				Log: xlogger.NewSimpleLogger(),
				DB:  mockdb,
			}

			tx, err := db.BeginTx(tc.CTX, tc.Opts)
			require.NoError(t, err)
			require.NotNil(t, tx)
		})
	}
}

func TestDB_BeginForRole(t *testing.T) {
	cases := map[string]struct {
		CTX  context.Context
		Opts *stdsql.TxOptions
		Role string
	}{
		"default": {
			context.Background(),
			nil,
			"role",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockdb, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectBegin()
			mock.ExpectExec("SET LOCAL SESSION AUTHORIZATION \\$1").WithArgs(tc.Role).WillReturnResult(sqlmock.NewResult(0, 0))

			db := &sql.DB{
				Log: xlogger.NewSimpleLogger(),
				DB:  mockdb,
			}

			ctx, tx, err := db.BeginForRole(tc.Role)
			require.NoError(t, err)
			require.NotNil(t, tx)
			require.NotNil(t, ctx)
		})
	}
}

func TestDB_BeginTxForRole(t *testing.T) {
	cases := map[string]struct {
		CTX  context.Context
		Opts *stdsql.TxOptions
		Role string
	}{
		"default": {
			context.Background(),
			nil,
			"role",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockdb, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectBegin()
			mock.ExpectExec("SET LOCAL SESSION AUTHORIZATION \\$1").WithArgs(tc.Role).WillReturnResult(sqlmock.NewResult(0, 0))

			db := &sql.DB{
				Log: xlogger.NewSimpleLogger(),
				DB:  mockdb,
			}

			tx, err := db.BeginTxForRole(tc.CTX, tc.Opts, tc.Role)
			require.NoError(t, err)
			require.NotNil(t, tx)
		})
	}
}
