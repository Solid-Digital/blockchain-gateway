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

func (s *Service) UpdateTrigger(params *dto.UpdateComponentRequest, orgName string, name string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error) {
	var trigger *orm.Trigger
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}
		trigger, appErr = xorm.GetTriggerTx(ctx, tx, orgName, name,
			qm.Load(orm.TriggerRels.TriggerVersions),
		)
		if appErr != nil {
			return appErr
		}

		if err = must(canEditTrigger(org, trigger)); err != nil {
			return ares.ErrForbiddenEdit(err)
		}

		trigger.DisplayName = params.DisplayName
		trigger.Description = params.Description
		trigger.UpdatedByID = principal.ID

		_, err = trigger.Update(ctx, tx, boil.Whitelist(orm.TriggerColumns.DisplayName, orm.TriggerColumns.Description, orm.TriggerColumns.UpdatedByID))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = trigger.L.LoadCreatedBy(ctx, tx, true, trigger, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), name)
		}

		err = trigger.L.LoadUpdatedBy(ctx, tx, true, trigger, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetTriggerDTO(org, trigger)
}
