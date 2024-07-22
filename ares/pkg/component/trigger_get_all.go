package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetAllTriggers(orgName string, available *bool) ([]*dto.GetComponentResponse, *apperr.Error) {
	var triggers []*orm.Trigger
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		triggers, appErr = xorm.GetAllTriggersTx(ctx, tx, org.Name,
			qm.Load(orm.TriggerRels.TriggerVersions),
			qm.Load(orm.TriggerRels.CreatedBy),
			qm.Load(orm.TriggerRels.UpdatedBy),
			orm.TriggerWhere.Public.EQ(true), qm.Or2(orm.TriggerWhere.DeveloperID.EQ(org.ID)),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	var res []*dto.GetComponentResponse

	for _, a := range triggers {
		// If the available filter is enabled, only include results that match it
		if available != nil {
			if canUseTrigger(org, a) != *available {
				continue
			}
		}

		trigger, err := GetTriggerDTO(org, a)
		if err != nil {
			return nil, err
		}

		res = append(res, trigger)
	}

	return res, nil
}

func (s *Service) GetAllPublicTriggers() ([]*dto.GetComponentResponse, *apperr.Error) {
	var triggers []*orm.Trigger

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		triggers, appErr = xorm.GetAllTriggersTx(ctx, tx, "public",
			qm.Load(orm.TriggerRels.TriggerVersions),
			qm.Load(orm.TriggerRels.CreatedBy),
			qm.Load(orm.TriggerRels.UpdatedBy),
			orm.TriggerWhere.Public.EQ(true),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	var res []*dto.GetComponentResponse

	for _, a := range triggers {
		trigger, err := GetPublicTriggerDTO(a)
		if err != nil {
			return nil, err
		}

		res = append(res, trigger)
	}

	return res, nil
}
