package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/unchainio/pkg/errors"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) UpdateBaseVersion(params *dto.UpdateBaseVersionRequest, orgName string, name string, version string, principal *dto.User) (*dto.GetBaseVersionResponse, *apperr.Error) {
	var org *orm.Organization
	var base *orm.Base
	var baseVersion *orm.BaseVersion

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, base, baseVersion, appErr = s.updateBaseVersionTx(ctx, tx, params, orgName, name, version, principal)
		if appErr != nil {
			return appErr
		}

		err = baseVersion.L.LoadCreatedBy(ctx, tx, true, baseVersion, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), name)
		}

		err = baseVersion.L.LoadUpdatedBy(ctx, tx, true, baseVersion, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetBaseVersionDTO(org, base, baseVersion)
}

func (s *Service) updateBaseVersionTx(ctx context.Context, tx *sql.Tx, params *dto.UpdateBaseVersionRequest, orgName string, name string, version string, principal *dto.User) (*orm.Organization, *orm.Base, *orm.BaseVersion, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	base, appErr := xorm.GetBaseTx(ctx, tx, orgName, name)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	if err := must(canEditBase(org, base)); err != nil {
		return nil, nil, nil, ares.ErrForbiddenEdit(err)
	}

	baseVersion, appErr := xorm.GetBaseVersionTx(ctx, tx, orgName, base, version,
		qm.Load(orm.BaseVersionRels.CreatedBy),
		qm.Load(orm.BaseVersionRels.UpdatedBy),
	)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	baseVersion.UpdatedByID = principal.ID
	base.UpdatedByID = principal.ID

	if params.Public != nil {
		if !s.service.enforcer.IsSuperAdmin(principal.ID) {
			return nil, nil, nil, ares.ErrForbiddenToSetPublic(errors.New(""))
		}

		baseVersion.Public = *params.Public

		// Once a component version is public, the component is also public
		if baseVersion.Public {
			base.Public = true
		}
	}

	if params.Readme != nil {
		baseVersion.Readme = *params.Readme
	}

	if params.Description != nil {
		baseVersion.Description = *params.Description
		base.Description = *params.Description
	}

	_, err := baseVersion.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	_, err = base.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	return org, base, baseVersion, nil
}
