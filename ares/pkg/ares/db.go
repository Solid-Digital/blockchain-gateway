package ares

import (
	"context"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	stdsql "database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"github.com/lib/pq"
)

func wrapErr(err error) *apperr.Error {
	if err == nil {
		return nil
	}

	// If it's already an app error, return it
	appErr := &apperr.Error{}
	if stderr.As(err, &appErr) {
		return appErr
	}

	return apperr.Internal.Wrap(err)
}

func WrapTx(db *sql.DB, fn func(ctx context.Context, tx *stdsql.Tx) *apperr.Error) *apperr.Error {
	return wrapErr(db.WrapTx(WrapTxFn(fn)))
}

func WrapTxFn(fn func(ctx context.Context, tx *stdsql.Tx) *apperr.Error) func(ctx context.Context, tx *stdsql.Tx) error {
	return func(ctx context.Context, tx *stdsql.Tx) error {
		err := fn(ctx, tx)

		if err != nil {
			return err
		}

		return nil
	}
}

func ParsePQErr(err error) *apperr.Error {
	if err == nil {
		return nil
	}

	pqErr := pqErr(err)

	switch {
	case stderr.Is(err, stdsql.ErrNoRows):
		return apperr.NotFound.Wrap(err)
	case pqErr.Constraint != "":
		return apperr.Conflict.Wrap(err)
	default:
		return apperr.Internal.Wrap(err)
	}
}

func pqErr(err error) *pq.Error {
	pqErr := &pq.Error{}

	if stderr.As(err, &pqErr) {
		return pqErr
	}

	return nil
}
