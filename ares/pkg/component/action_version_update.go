package component

import (
	"context"
	"database/sql"
	"encoding/json"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/unchainio/pkg/errors"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) UpdateActionVersion(params *dto.UpdateComponentVersionRequest, orgName string, name string, version string, principal *dto.User) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var org *orm.Organization
	var action *orm.Action
	var actionVersion *orm.ActionVersion

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, action, actionVersion, appErr = s.updateActionVersionTx(ctx, tx, params, orgName, name, version, principal)
		if appErr != nil {
			return appErr
		}

		err = actionVersion.L.LoadCreatedBy(ctx, tx, true, actionVersion, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), name)
		}

		err = actionVersion.L.LoadUpdatedBy(ctx, tx, true, actionVersion, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetActionVersionDTO(org, action, actionVersion)
}

func (s *Service) updateActionVersionTx(ctx context.Context, tx *sql.Tx, params *dto.UpdateComponentVersionRequest, orgName string, name string, version string, principal *dto.User) (*orm.Organization, *orm.Action, *orm.ActionVersion, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	action, appErr := xorm.GetActionTx(ctx, tx, orgName, name)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	if err := must(canEditAction(org, action)); err != nil {
		return nil, nil, nil, ares.ErrForbiddenEdit(err)
	}

	actionVersion, appErr := xorm.GetActionVersionTx(ctx, tx, orgName, action, version,
		qm.Load(orm.ActionVersionRels.CreatedBy),
		qm.Load(orm.ActionVersionRels.UpdatedBy),
	)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	actionVersion.UpdatedByID = principal.ID
	action.UpdatedByID = principal.ID

	if params.Public != nil {
		if !s.service.enforcer.IsSuperAdmin(principal.ID) {
			return nil, nil, nil, ares.ErrForbiddenToSetPublic(errors.New(""))
		}

		actionVersion.Public = *params.Public

		// Once a component version is public, the component is also public
		if actionVersion.Public {
			action.Public = true
		}
	}

	if params.ExampleConfig != nil {
		actionVersion.ExampleConfig = *params.ExampleConfig
	}

	if params.Readme != nil {
		actionVersion.Readme = *params.Readme
	}

	if params.Description != nil {
		actionVersion.Description = *params.Description
		action.Description = *params.Description
	}

	if len(params.InputSchema) != 0 {
		input, err := json.Marshal(params.InputSchema)
		if err != nil {
			return nil, nil, nil, ares.ErrParseInputSchema(errors.Wrap(err, ""))
		}

		actionVersion.InputSchema = null.JSONFrom(input)
	}

	if len(params.OutputSchema) != 0 {
		output, err := json.Marshal(params.OutputSchema)
		if err != nil {
			return nil, nil, nil, ares.ErrParseOutputSchema(errors.Wrap(err, ""))
		}

		actionVersion.OutputSchema = null.JSONFrom(output)
	}

	_, err := actionVersion.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	_, err = action.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	return org, action, actionVersion, nil
}
