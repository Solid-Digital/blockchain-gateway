package component

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) CreateBaseVersion(params *dto.CreateBaseVersionRequest, orgName string, name string, principal *dto.User) (*dto.GetBaseVersionResponse, *apperr.Error) {
	var org *orm.Organization
	var base *orm.Base
	var baseVersion *orm.BaseVersion

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		if params.Public != nil && !s.service.enforcer.IsSuperAdmin(principal.ID) {
			return ares.ErrForbiddenToSetPublic(stderr.New(""))
		}

		org, base, baseVersion, appErr = s.createBaseVersionTx(ctx, tx, params, orgName, name, principal)
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

func (s *Service) createBaseVersionTx(ctx context.Context, tx *sql.Tx, params *dto.CreateBaseVersionRequest, orgName, name string, principal *dto.User) (*orm.Organization, *orm.Base, *orm.BaseVersion, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	base, appErr := upsertBaseTx(ctx, tx, name, principal.ID, org.ID, params.Description)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	if err := must(canEditBase(org, base)); err != nil {
		return nil, nil, nil, ares.ErrForbiddenEdit(err)
	}

	publicVal := false

	if params.Public != nil {
		publicVal = *params.Public

		// Once a component version is public, the component is also public
		if publicVal {
			base.Public = true
		}
	}

	baseVersion := &orm.BaseVersion{
		BaseID:         base.ID,
		CreatedByID:    principal.ID,
		UpdatedByID:    principal.ID,
		Version:        params.Version,
		Description:    params.Description,
		Readme:         params.Readme,
		DockerImageRef: params.DockerImageRef,
		Entrypoint:     params.Entrypoint,
		Public:         publicVal,
	}

	err := baseVersion.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return nil, nil, nil, ares.ErrDuplicateComponentVersion(err, ares.ComponentTypeBase, name, params.Version)
		default:
			return nil, nil, nil, err
		}
	}

	_, err = base.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	return org, base, baseVersion, nil
}

func upsertBaseTx(ctx context.Context, tx *sql.Tx, name string, userID int64, orgID int64, description string) (*orm.Base, *apperr.Error) {
	var base *orm.Base

	exists, err := orm.Bases(orm.BaseWhere.Name.EQ(name)).Exists(ctx, tx)
	if err != nil {
		return nil, ares.ParsePQErr(err)
	}

	if exists {
		base, err = orm.Bases(orm.BaseWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return nil, ares.ParsePQErr(err)
		}

		base.Description = description

		_, err = base.Update(ctx, tx, boil.Infer())
		if err != nil {
			return nil, ares.ParsePQErr(err)
		}

		return base, nil
	}

	base = &orm.Base{
		DeveloperID: orgID,
		CreatedByID: userID,
		UpdatedByID: userID,
		Name:        name,
		DisplayName: name,
		Description: description,
		Public:      false,
	}

	appErr := xorm.CreateBaseTx(ctx, tx, base)
	if appErr != nil {
		return nil, appErr
	}

	return base, nil
}
