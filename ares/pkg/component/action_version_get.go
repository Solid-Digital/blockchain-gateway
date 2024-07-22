package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetActionVersion(orgName string, name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var action *orm.Action
	var actionVersion *orm.ActionVersion
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		action, appErr = xorm.GetActionTx(ctx, tx, orgName, name)
		if appErr != nil {
			return appErr
		}

		actionVersion, appErr = xorm.GetActionVersionTx(ctx, tx, orgName, action, version,
			qm.Load(orm.ActionVersionRels.CreatedBy),
			qm.Load(orm.ActionVersionRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	if err := must(canViewActionVersion(org, action, actionVersion)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	return GetActionVersionDTO(org, action, actionVersion)
}

func (s *Service) GetPublicActionVersion(name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var action *orm.Action
	var actionVersion *orm.ActionVersion
	var emptyOrg *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		action, appErr = xorm.GetActionTx(ctx, tx, "public", name)
		if appErr != nil {
			return appErr
		}

		actionVersion, appErr = xorm.GetActionVersionTx(ctx, tx, "public", action, version,
			qm.Load(orm.ActionVersionRels.CreatedBy),
			qm.Load(orm.ActionVersionRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	if err := must(canViewActionVersion(emptyOrg, action, actionVersion)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	return GetActionVersionDTO(emptyOrg, action, actionVersion)
}

func GetActionVersionDTO(org *orm.Organization, component *orm.Action, version *orm.ActionVersion) (*dto.GetComponentVersionResponse, *apperr.Error) {
	input, output, err := GetActionVersionIOSchema(version)
	if err != nil {
		return nil, err
	}

	return &dto.GetComponentVersionResponse{
		ID:             &version.ID,
		Version:        &version.Version,
		ReleaseMessage: &version.Description,
		Readme:         &version.Readme,
		ExampleConfig:  &version.ExampleConfig,
		InputSchema:    input,
		OutputSchema:   output,
		Public:         &version.Public,
		Available:      testhelper.BoolPtr(canUseActionVersion(org, component, version)),
		CreatedAt:      (*strfmt.DateTime)(&version.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       version.R.CreatedBy.ID,
			FullName: version.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&version.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       version.R.UpdatedBy.ID,
			FullName: version.R.UpdatedBy.FullName,
		},
	}, nil
}
