package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBPipelineExists(pipelineID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Pipelines(orm.PipelineWhere.ID.EQ(pipelineID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetPipeline(pipelineID int64) *orm.Pipeline {
	if !h.DBPipelineExists(pipelineID) {
		return nil
	}

	var pipelineFromDB *orm.Pipeline
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		pipelineFromDB, err = orm.Pipelines(orm.PipelineWhere.ID.EQ(pipelineID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return pipelineFromDB
}

func (h *Helper) DBPipelineByNameExists(name string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Pipelines(orm.PipelineWhere.Name.EQ(name)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetPipelineByName(name string) *orm.Pipeline {
	if !h.DBPipelineByNameExists(name) {
		return nil
	}

	var pipelineFromDB *orm.Pipeline
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		pipelineFromDB, err = orm.Pipelines(orm.PipelineWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return pipelineFromDB
}
