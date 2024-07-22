package component

import (
	"context"
	"database/sql"
	"encoding/json"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/unchainio/pkg/errors"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/volatiletech/null"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) UpdateTriggerVersion(params *dto.UpdateComponentVersionRequest, orgName string, name string, version string, principal *dto.User) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var org *orm.Organization
	var trigger *orm.Trigger
	var triggerVersion *orm.TriggerVersion

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, trigger, triggerVersion, appErr = s.updateTriggerVersionTx(ctx, tx, params, orgName, name, version, principal)
		if appErr != nil {
			return appErr
		}

		err = triggerVersion.L.LoadCreatedBy(ctx, tx, true, triggerVersion, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), name)
		}

		err = triggerVersion.L.LoadUpdatedBy(ctx, tx, true, triggerVersion, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetTriggerVersionDTO(org, trigger, triggerVersion)
}

func (s *Service) updateTriggerVersionTx(ctx context.Context, tx *sql.Tx, params *dto.UpdateComponentVersionRequest, orgName string, name string, version string, principal *dto.User) (*orm.Organization, *orm.Trigger, *orm.TriggerVersion, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	trigger, appErr := xorm.GetTriggerTx(ctx, tx, orgName, name)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	if err := must(canEditTrigger(org, trigger)); err != nil {
		return nil, nil, nil, ares.ErrForbiddenEdit(err)
	}

	triggerVersion, appErr := xorm.GetTriggerVersionTx(ctx, tx, orgName, trigger, version,
		qm.Load(orm.TriggerVersionRels.CreatedBy),
		qm.Load(orm.TriggerVersionRels.UpdatedBy),
	)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	triggerVersion.UpdatedByID = principal.ID
	trigger.UpdatedByID = principal.ID

	if params.Public != nil {
		if !s.service.enforcer.IsSuperAdmin(principal.ID) {
			return nil, nil, nil, ares.ErrForbiddenToSetPublic(errors.New(""))
		}

		triggerVersion.Public = *params.Public

		// Once a component version is public, the component is also public
		if triggerVersion.Public {
			trigger.Public = true
		}
	}

	if params.ExampleConfig != nil {
		triggerVersion.ExampleConfig = *params.ExampleConfig
	}

	if params.Readme != nil {
		triggerVersion.Readme = *params.Readme
	}

	if params.Description != nil {
		triggerVersion.Description = *params.Description
		trigger.Description = *params.Description
	}

	if len(params.InputSchema) != 0 {
		input, err := json.Marshal(params.InputSchema)
		if err != nil {
			return nil, nil, nil, ares.ErrParseInputSchema(errors.Wrap(err, ""))
		}

		triggerVersion.InputSchema = null.JSONFrom(input)
	}

	if len(params.OutputSchema) != 0 {
		output, err := json.Marshal(params.OutputSchema)
		if err != nil {
			return nil, nil, nil, ares.ErrParseOutputSchema(errors.Wrap(err, ""))
		}

		triggerVersion.OutputSchema = null.JSONFrom(output)
	}

	_, err := triggerVersion.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	_, err = trigger.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	return org, trigger, triggerVersion, nil
}
