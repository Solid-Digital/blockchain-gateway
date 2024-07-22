package pipeline

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/BurntSushi/toml"
	"github.com/unchainio/pkg/errors"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"

	"github.com/volatiletech/sqlboiler/boil"
)

func (s *Service) UpdateDraftConfiguration(params *dto.UpdateDraftConfigurationRequest, orgName string, pipelineName string, principal *dto.User) (*dto.GetConfigurationResponse, *apperr.Error) {
	var ret *dto.GetConfigurationResponse
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		_, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName,
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.BaseDraftConfiguration)),
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.TriggerDraftConfiguration)),
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.ActionDraftConfigurations)),
			qm.Load(qm.Rels(orm.PipelineRels.DraftConfiguration, orm.DraftConfigurationRels.CreatedBy)),
		)
		if appErr != nil {
			return appErr
		}

		draft := pipeline.R.DraftConfiguration
		draft.UpdatedByID = principal.ID

		_, err := draft.Update(ctx, tx, boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = draft.L.LoadUpdatedBy(ctx, tx, true, draft, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), "draft")
		}

		appErr = setBaseDraftConfiguration(ctx, tx, draft, params.Base)
		if appErr != nil {
			return appErr
		}

		appErr = setTriggerDraftConfiguration(ctx, tx, draft, params.Trigger)
		if appErr != nil {
			return appErr
		}

		appErr = setActionDraftConfigurations(ctx, tx, draft, params.Actions)
		if appErr != nil {
			return appErr
		}

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

func setActionDraftConfigurations(ctx context.Context, tx *sql.Tx, draft *orm.DraftConfiguration, actionParams []*dto.UpdateComponentConfiguration) *apperr.Error {
	var err error
	if draft.R.ActionDraftConfigurations != nil {
		_, err = draft.R.ActionDraftConfigurations.DeleteAll(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	for i, a := range actionParams {
		actionMessageConfig, err := tomlToJSON([]byte(a.MessageConfiguration))
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		err = draft.AddActionDraftConfigurations(ctx, tx, true, &orm.ActionDraftConfiguration{
			DraftConfigurationID: draft.ID,
			VersionID:            *a.ID,
			Index:                int64(i),
			Name:                 a.Name,
			Config:               *a.InitConfiguration,
			MessageConfig:        null.JSONFrom(actionMessageConfig),
		})
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	return nil
}

func setBaseDraftConfiguration(ctx context.Context, tx *sql.Tx, draft *orm.DraftConfiguration, baseParam *dto.UpdateComponentConfiguration) *apperr.Error {
	if draft.R.BaseDraftConfiguration != nil {
		_, err := draft.R.BaseDraftConfiguration.Delete(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	if baseParam != nil {
		err := draft.SetBaseDraftConfiguration(ctx, tx, true, &orm.BaseDraftConfiguration{
			DraftConfigurationID: draft.ID,
			VersionID:            null.Int64From(*baseParam.ID),
			Config:               *baseParam.InitConfiguration,
		})
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	return nil
}

func setTriggerDraftConfiguration(ctx context.Context, tx *sql.Tx, draft *orm.DraftConfiguration, triggerParam *dto.UpdateComponentConfiguration) *apperr.Error {
	if draft.R.TriggerDraftConfiguration != nil {
		_, err := draft.R.TriggerDraftConfiguration.Delete(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	if triggerParam != nil {
		triggerMessageConfig, err := tomlToJSON([]byte(triggerParam.MessageConfiguration))
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		err = draft.SetTriggerDraftConfiguration(ctx, tx, true, &orm.TriggerDraftConfiguration{
			DraftConfigurationID: draft.ID,
			VersionID:            null.Int64From(*triggerParam.ID),
			Config:               *triggerParam.InitConfiguration,
			Name:                 triggerParam.Name,
			MessageConfig:        null.JSONFrom(triggerMessageConfig), // TODO(e-nikolov) test me
		})
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	return nil
}

func tomlToJSON(tomlBytes []byte) (jsonBytes []byte, err error) {
	m := make(map[string]interface{})

	err = toml.Unmarshal(tomlBytes, &m)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	jsonBytes, err = json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return jsonBytes, nil
}

func jsonToTOML(jsonBytes []byte) (tomlBytes []byte, err error) {
	m := make(map[string]interface{})

	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	buf := bytes.NewBuffer(nil)
	err = toml.NewEncoder(buf).Encode(&m)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	tomlBytes = buf.Bytes()

	return tomlBytes, nil
}
