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

func (s *Service) GetAllActions(orgName string, available *bool) ([]*dto.GetComponentResponse, *apperr.Error) {
	var actions []*orm.Action
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		actions, appErr = xorm.GetAllActionsTx(ctx, tx, org.Name,
			qm.Load(orm.ActionRels.ActionVersions),
			qm.Load(orm.ActionRels.CreatedBy),
			qm.Load(orm.ActionRels.UpdatedBy),
			orm.ActionWhere.Public.EQ(true), qm.Or2(orm.ActionWhere.DeveloperID.EQ(org.ID)),
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

	for _, a := range actions {
		// If the available filter is enabled, only include results that match it
		if available != nil {
			if canUseAction(org, a) != *available {
				continue
			}
		}

		action, err := GetActionDTO(org, a)
		if err != nil {
			return nil, err
		}

		res = append(res, action)
	}

	return res, nil
}

func (s *Service) GetAllPublicActions() ([]*dto.GetComponentResponse, *apperr.Error) {
	var actions []*orm.Action

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		actions, appErr = xorm.GetAllActionsTx(ctx, tx, "public",
			qm.Load(orm.ActionRels.ActionVersions),
			qm.Load(orm.ActionRels.CreatedBy),
			qm.Load(orm.ActionRels.UpdatedBy),
			orm.ActionWhere.Public.EQ(true),
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

	for _, a := range actions {
		action, err := GetPublicActionDTO(a)
		if err != nil {
			return nil, err
		}

		res = append(res, action)
	}

	return res, nil
}
