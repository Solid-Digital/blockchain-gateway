package factory

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) PipelineWithNameAndDraft(name string, create bool, org *orm.Organization, user *orm.User) (*orm.Pipeline, *orm.DraftConfiguration) {
	draftConfiguration := f.DraftConfiguration(create, org, user)

	pipeline := &orm.Pipeline{
		OrganizationID: org.ID,
		CreatedByID:    user.ID,
		UpdatedByID:    user.ID,
		DisplayName:    testhelper.Randumb(fmt.Sprintf("%s Pipeline", org.DisplayName)),
		Name:           name,
		DraftConfigurationID: null.Int64{
			Int64: draftConfiguration.ID,
			Valid: true,
		},
	}

	if !create {
		return pipeline, draftConfiguration
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return pipeline.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return pipeline, draftConfiguration
}
func (f *Factory) PipelineWithDraft(create bool, org *orm.Organization, user *orm.User) (*orm.Pipeline, *orm.DraftConfiguration) {
	return f.PipelineWithNameAndDraft(testhelper.Randumb(fmt.Sprintf("%s-pipeline", org.Name)), create, org, user)
}

func (f *Factory) Pipeline(create bool, org *orm.Organization, user *orm.User) *orm.Pipeline {
	pipeline, _ := f.PipelineWithDraft(create, org, user)

	return pipeline
}

func (f *Factory) PipelineWithName(name string, create bool, org *orm.Organization, user *orm.User) *orm.Pipeline {
	pipeline, _ := f.PipelineWithNameAndDraft(name, create, org, user)

	return pipeline
}
