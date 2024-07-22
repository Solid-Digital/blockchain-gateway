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

func (s *Service) UpdateAction(params *dto.UpdateComponentRequest, orgName string, name string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error) {
	var action *orm.Action
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		action, appErr = xorm.GetActionTx(ctx, tx, orgName, name,
			qm.Load(orm.ActionRels.ActionVersions),
		)
		if appErr != nil {
			return appErr
		}

		if err = must(canEditAction(org, action)); err != nil {
			return ares.ErrForbiddenEdit(err)
		}

		action.DisplayName = params.DisplayName
		action.Description = params.Description
		action.UpdatedByID = principal.ID

		_, err = action.Update(ctx, tx, boil.Whitelist(orm.ActionColumns.DisplayName, orm.ActionColumns.Description, orm.ActionColumns.UpdatedByID))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = action.L.LoadCreatedBy(ctx, tx, true, action, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), name)
		}

		err = action.L.LoadUpdatedBy(ctx, tx, true, action, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetActionDTO(org, action)
}
