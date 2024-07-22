package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) UpdateBase(params *dto.UpdateComponentRequest, orgName string, name string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error) {
	var base *orm.Base
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		base, appErr = xorm.GetBaseTx(ctx, tx, orgName, name,
			qm.Load(orm.BaseRels.BaseVersions),
		)
		if appErr != nil {
			return appErr
		}

		if err = must(canEditBase(org, base)); err != nil {
			return ares.ErrForbiddenEdit(err)
		}

		base.DisplayName = params.DisplayName
		base.Description = params.Description
		base.UpdatedByID = principal.ID

		_, err = base.Update(ctx, tx, boil.Whitelist(orm.BaseColumns.DisplayName, orm.BaseColumns.Description, orm.BaseColumns.UpdatedByID))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = base.L.LoadCreatedBy(ctx, tx, true, base, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), name)
		}

		err = base.L.LoadUpdatedBy(ctx, tx, true, base, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetBaseDTO(org, base)
}
