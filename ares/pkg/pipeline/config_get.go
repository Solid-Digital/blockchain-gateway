package pipeline

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetConfiguration(orgName string, pipelineName string, revision int64) (*dto.GetConfigurationResponse, *apperr.Error) {
	var ret *dto.GetConfigurationResponse

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		ret, appErr = getConfigurationTx(ctx, tx, orgName, pipelineName, revision)
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

func getConfigurationTx(ctx context.Context, tx *sql.Tx, orgName string, pipelineName string, revision int64) (*dto.GetConfigurationResponse, *apperr.Error) {
	var appErr *apperr.Error
	org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
	if appErr != nil {
		return nil, appErr
	}

	config, err := pipeline.Configurations(
		qm.Load(qm.Rels(
			orm.ConfigurationRels.BaseConfiguration,
			orm.BaseConfigurationRels.Version,
			orm.BaseVersionRels.Base,
			orm.BaseRels.BaseVersions,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.BaseConfiguration,
			orm.BaseConfigurationRels.Version,
			orm.BaseVersionRels.Base,
			orm.BaseRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.BaseConfiguration,
			orm.BaseConfigurationRels.Version,
			orm.BaseVersionRels.Base,
			orm.BaseRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.BaseConfiguration,
			orm.BaseConfigurationRels.Version,
			orm.BaseVersionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.BaseConfiguration,
			orm.BaseConfigurationRels.Version,
			orm.BaseVersionRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.TriggerConfiguration,
			orm.TriggerConfigurationRels.Version,
			orm.TriggerVersionRels.Trigger,
			orm.TriggerRels.TriggerVersions,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.TriggerConfiguration,
			orm.TriggerConfigurationRels.Version,
			orm.TriggerVersionRels.Trigger,
			orm.TriggerRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.TriggerConfiguration,
			orm.TriggerConfigurationRels.Version,
			orm.TriggerVersionRels.Trigger,
			orm.TriggerRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.TriggerConfiguration,
			orm.TriggerConfigurationRels.Version,
			orm.TriggerVersionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.TriggerConfiguration,
			orm.TriggerConfigurationRels.Version,
			orm.TriggerVersionRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.ActionConfigurations,
		),
			qm.OrderBy(orm.ActionConfigurationColumns.Index),
		),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.ActionConfigurations,
			orm.ActionConfigurationRels.Version,
			orm.ActionVersionRels.Action,
			orm.ActionRels.ActionVersions,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.ActionConfigurations,
			orm.ActionConfigurationRels.Version,
			orm.ActionVersionRels.Action,
			orm.ActionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.ActionConfigurations,
			orm.ActionConfigurationRels.Version,
			orm.ActionVersionRels.Action,
			orm.ActionRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.ActionConfigurations,
			orm.ActionConfigurationRels.Version,
			orm.ActionVersionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.ConfigurationRels.ActionConfigurations,
			orm.ActionConfigurationRels.Version,
			orm.ActionVersionRels.UpdatedBy,
		)),
		qm.Load(orm.ConfigurationRels.CreatedBy),
		qm.Load(orm.ConfigurationRels.UpdatedBy),
		orm.ConfigurationWhere.Revision.EQ(revision),
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrPipelineConfigNotFound(err, orgName, pipelineName, revision)
		default:
			return nil, err
		}
	}

	ret, appErr := getConfigResponse(org, config)
	if appErr != nil {
		return nil, appErr
	}

	return ret, nil
}

func getConfigResponse(org *orm.Organization, config *orm.Configuration) (*dto.GetConfigurationResponse, *apperr.Error) {
	base, appErr := getConfigResponseBase(org, config.R.BaseConfiguration)
	if appErr != nil {
		return nil, appErr
	}

	trigger, appErr := getConfigResponseTrigger(org, config.R.TriggerConfiguration)
	if appErr != nil {
		return nil, appErr
	}

	actions, appErr := getConfigurationActions(org, config.R.ActionConfigurations)
	if appErr != nil {
		return nil, appErr
	}

	return &dto.GetConfigurationResponse{
		ID:            &config.ID,
		Revision:      &config.Revision,
		CommitMessage: config.CommitMessage,
		Base:          base,
		Trigger:       trigger,
		Actions:       actions,
		CreatedAt:     (*strfmt.DateTime)(&config.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       config.R.CreatedBy.ID,
			FullName: config.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&config.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       config.R.UpdatedBy.ID,
			FullName: config.R.UpdatedBy.FullName,
		},
	}, nil
}

func getConfigResponseTrigger(org *orm.Organization, config *orm.TriggerConfiguration) (*dto.ComponentConfiguration, *apperr.Error) {
	if config == nil || config.R == nil || config.R.Version == nil {
		return nil, nil
	}

	return getConfigResponseTriggerInternal(
		org,
		config.Name,
		config.MessageConfig,
		config.R.Version.R.Trigger,
		config.R.Version,
		config.Config,
	)
}

func getConfigResponseBase(org *orm.Organization, cfg *orm.BaseConfiguration) (*dto.ComponentConfiguration, *apperr.Error) {
	if cfg == nil || cfg.R == nil || cfg.R.Version == nil {
		return nil, nil
	}

	return getConfigResponseBaseInternal(
		org,
		cfg.R.Version.R.Base,
		cfg.R.Version,
		cfg.Config,
	)
}

func getConfigurationActions(org *orm.Organization, actions orm.ActionConfigurationSlice) ([]*dto.ComponentConfiguration, *apperr.Error) {
	res := make([]*dto.ComponentConfiguration, len(actions))

	for i, a := range actions {
		action, appErr := getConfigurationActionInternal(org, a.Name, a.MessageConfig, a.R.Version.R.Action, a.R.Version, a.Config)
		if appErr != nil {
			return nil, appErr
		}

		res[i] = action
	}

	return res, nil
}
