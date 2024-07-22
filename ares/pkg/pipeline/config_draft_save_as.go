package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (s *Service) SaveDraftConfigurationAsNew(params *dto.SaveDraftConfigurationAsNewRequest, orgName string, pipelineName string, principal *dto.User) (*dto.GetConfigurationResponse, *apperr.Error) {
	var ret *dto.GetConfigurationResponse
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error

		org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		_, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName,
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.BaseDraftConfiguration)),
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.TriggerDraftConfiguration)),
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.ActionDraftConfigurations)),
		)
		if appErr != nil {
			return appErr
		}

		draft := pipeline.R.DraftConfiguration

		config := &orm.Configuration{
			CreatedByID:    principal.ID,
			UpdatedByID:    principal.ID,
			PipelineID:     pipeline.ID,
			OrganizationID: org.ID,
			Revision:       draft.Revision,
			CommitMessage:  params.CommitMessage,
		}

		err = config.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}

		draft.Revision += 1

		_, err = draft.Update(ctx, tx, boil.Whitelist(orm.ConfigurationColumns.Revision))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = config.SetBaseConfiguration(ctx, tx, true, &orm.BaseConfiguration{
			ConfigurationID: config.ID,
			VersionID:       draft.R.BaseDraftConfiguration.VersionID,
			Config:          draft.R.BaseDraftConfiguration.Config,
		})
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = config.SetTriggerConfiguration(ctx, tx, true, &orm.TriggerConfiguration{
			ConfigurationID: config.ID,
			VersionID:       draft.R.TriggerDraftConfiguration.VersionID,
			Config:          draft.R.TriggerDraftConfiguration.Config,
			MessageConfig:   draft.R.TriggerDraftConfiguration.MessageConfig, // TODO(e-nikolov) test me
			Name:            draft.R.TriggerDraftConfiguration.Name,
		})
		if err != nil {
			return ares.ParsePQErr(err)
		}

		for _, a := range draft.R.ActionDraftConfigurations {
			err = config.AddActionConfigurations(ctx, tx, true, &orm.ActionConfiguration{
				ConfigurationID: config.ID,
				VersionID:       a.VersionID,
				Index:           a.Index,
				Config:          a.Config,
				MessageConfig:   a.MessageConfig,
				Name:            a.Name,
			})
			if err != nil {
				return ares.ParsePQErr(err)
			}
		}

		ret, appErr = getConfigurationTx(ctx, tx, orgName, pipelineName, config.Revision)
		if appErr != nil {
			return appErr
		}
		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return ret, nil
}
