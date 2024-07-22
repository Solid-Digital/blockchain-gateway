package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/pkg/component"

	"github.com/go-openapi/strfmt"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetDraftConfiguration(orgName string, pipelineName string) (*dto.GetConfigurationResponse, *apperr.Error) {
	var ret *dto.GetConfigurationResponse
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		ret, appErr = getDraftConfigurationTx(ctx, tx, orgName, pipelineName)
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

func getDraftConfigurationTx(ctx context.Context, tx *sql.Tx, orgName, pipelineName string) (*dto.GetConfigurationResponse, *apperr.Error) {
	org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
	if appErr != nil {
		return nil, appErr
	}

	draft, err := pipeline.DraftConfiguration(
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.BaseDraftConfiguration,
			orm.BaseDraftConfigurationRels.Version,
			orm.BaseVersionRels.Base,
			orm.BaseRels.BaseVersions,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.TriggerDraftConfiguration,
			orm.TriggerDraftConfigurationRels.Version,
			orm.TriggerVersionRels.Trigger,
			orm.TriggerRels.TriggerVersions,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.ActionDraftConfigurations,
		),
			qm.OrderBy(orm.ActionDraftConfigurationColumns.Index),
		),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.ActionDraftConfigurations,
			orm.ActionDraftConfigurationRels.Version,
			orm.ActionVersionRels.Action,
			orm.ActionRels.ActionVersions,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.BaseDraftConfiguration,
			orm.BaseDraftConfigurationRels.Version,
			orm.BaseVersionRels.Base,
			orm.BaseRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.BaseDraftConfiguration,
			orm.BaseDraftConfigurationRels.Version,
			orm.BaseVersionRels.Base,
			orm.BaseRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.BaseDraftConfiguration,
			orm.BaseDraftConfigurationRels.Version,
			orm.BaseVersionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.BaseDraftConfiguration,
			orm.BaseDraftConfigurationRels.Version,
			orm.BaseVersionRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.TriggerDraftConfiguration,
			orm.TriggerDraftConfigurationRels.Version,
			orm.TriggerVersionRels.Trigger,
			orm.TriggerRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.TriggerDraftConfiguration,
			orm.TriggerDraftConfigurationRels.Version,
			orm.TriggerVersionRels.Trigger,
			orm.TriggerRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.TriggerDraftConfiguration,
			orm.TriggerDraftConfigurationRels.Version,
			orm.TriggerVersionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.TriggerDraftConfiguration,
			orm.TriggerDraftConfigurationRels.Version,
			orm.TriggerVersionRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.ActionDraftConfigurations,
			orm.ActionDraftConfigurationRels.Version,
			orm.ActionVersionRels.Action,
			orm.ActionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.ActionDraftConfigurations,
			orm.ActionDraftConfigurationRels.Version,
			orm.ActionVersionRels.Action,
			orm.ActionRels.UpdatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.ActionDraftConfigurations,
			orm.ActionDraftConfigurationRels.Version,
			orm.ActionVersionRels.CreatedBy,
		)),
		qm.Load(qm.Rels(
			orm.DraftConfigurationRels.ActionDraftConfigurations,
			orm.ActionDraftConfigurationRels.Version,
			orm.ActionVersionRels.UpdatedBy,
		)),
		qm.Load(orm.DraftConfigurationRels.CreatedBy),
		qm.Load(orm.DraftConfigurationRels.UpdatedBy),
	).One(ctx, tx)
	if err != nil {
		return nil, ares.ParsePQErr(err)
	}

	return getDraftConfigResponse(org, draft)
}

func getDraftConfigResponse(org *orm.Organization, config *orm.DraftConfiguration) (*dto.GetConfigurationResponse, *apperr.Error) {
	base, err := getDraftConfigResponseBase(org, config.R.BaseDraftConfiguration)
	if err != nil {
		return nil, err
	}

	trigger, err := getDraftConfigResponseTrigger(org, config.R.TriggerDraftConfiguration)
	if err != nil {
		return nil, err
	}

	actions, err := getDraftConfigurationActions(org, config.R.ActionDraftConfigurations)
	if err != nil {
		return nil, err
	}

	return &dto.GetConfigurationResponse{
		ID:        &config.ID,
		Revision:  &config.Revision,
		Base:      base,
		Trigger:   trigger,
		Actions:   actions,
		CreatedAt: (*strfmt.DateTime)(&config.CreatedAt),
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

func getDraftConfigResponseTrigger(org *orm.Organization, cfg *orm.TriggerDraftConfiguration) (*dto.ComponentConfiguration, *apperr.Error) {
	if cfg == nil || cfg.R == nil || cfg.R.Version == nil {
		return nil, nil
	}

	return getConfigResponseTriggerInternal(
		org,
		cfg.Name,
		cfg.MessageConfig,
		cfg.R.Version.R.Trigger,
		cfg.R.Version,
		cfg.Config,
	)
}

func getConfigResponseTriggerInternal(org *orm.Organization, name string, messageConfig null.JSON, trigger *orm.Trigger, version *orm.TriggerVersion, config string) (*dto.ComponentConfiguration, *apperr.Error) {
	messageConfigString, appErr := getAdapterMessageConfigTOML(messageConfig) // TODO(e-nikolov) test me
	if appErr != nil {
		return nil, appErr
	}

	c, appErr := component.GetTriggerDTO(org, trigger)
	if appErr != nil {
		return nil, appErr
	}

	cv, appErr := component.GetTriggerVersionDTO(org, trigger, version)
	if appErr != nil {
		return nil, appErr
	}

	return &dto.ComponentConfiguration{
		Name:                 name,
		Component:            c,
		ComponentVersion:     cv,
		InitConfiguration:    &config,
		MessageConfiguration: &messageConfigString,
	}, nil
}

func getConfigResponseBaseInternal(org *orm.Organization, base *orm.Base, version *orm.BaseVersion, config string) (*dto.ComponentConfiguration, *apperr.Error) {
	c, appErr := component.GetBaseDTO(org, base)
	if appErr != nil {
		return nil, appErr
	}

	// TODO(e-nikolov) Change the return type to dto.BaseConfiguration that contains a base version instead of component version
	//	 Because the base version has different fields than the component version
	//cv, err := component.GetBaseVersionDTO(org, cfg.R.Version.R.Base, cfg.R.Version)
	//if err != nil {
	//	return nil, err
	//}

	return &dto.ComponentConfiguration{
		Name:      *c.Name,
		Component: c,
		//ComponentVersion: cv,
		ComponentVersion: &dto.GetComponentVersionResponse{
			ID:             &version.ID,
			Version:        &version.Version,
			ReleaseMessage: &version.Description,
			Readme:         &version.Readme,
			Public:         &version.Public,
		},

		InitConfiguration: &config,
	}, nil
}

func getDraftConfigResponseBase(org *orm.Organization, cfg *orm.BaseDraftConfiguration) (*dto.ComponentConfiguration, *apperr.Error) {
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

func getDraftConfigurationActions(org *orm.Organization, actions orm.ActionDraftConfigurationSlice) ([]*dto.ComponentConfiguration, *apperr.Error) {
	res := make([]*dto.ComponentConfiguration, len(actions))

	for i, a := range actions {
		action, err := getConfigurationActionInternal(org, a.Name, a.MessageConfig, a.R.Version.R.Action, a.R.Version, a.Config)
		if err != nil {
			return nil, err
		}

		res[i] = action
	}

	return res, nil
}

func getConfigurationActionInternal(org *orm.Organization, name string, messageConfig null.JSON, action *orm.Action, version *orm.ActionVersion, config string) (*dto.ComponentConfiguration, *apperr.Error) {
	messageConfigString, appErr := getAdapterMessageConfigTOML(messageConfig)
	if appErr != nil {
		return nil, appErr
	}

	c, appErr := component.GetActionDTO(org, action)
	if appErr != nil {
		return nil, appErr
	}

	cv, appErr := component.GetActionVersionDTO(org, action, version)
	if appErr != nil {
		return nil, appErr
	}

	return &dto.ComponentConfiguration{
		Name:                 name,
		Component:            c,
		ComponentVersion:     cv,
		InitConfiguration:    &config,
		MessageConfiguration: &messageConfigString,
	}, nil
}
