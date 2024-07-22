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

func (s *Service) GetAllBases(orgName string, available *bool) ([]*dto.GetComponentResponse, *apperr.Error) {
	var bases []*orm.Base
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err *apperr.Error

		org, err = xorm.GetOrganizationTx(ctx, tx, orgName)
		if err != nil {
			return err
		}

		bases, err = xorm.GetAllBasesTx(ctx, tx, org.Name,
			qm.Load(orm.BaseRels.BaseVersions),
			qm.Load(orm.BaseRels.CreatedBy),
			qm.Load(orm.BaseRels.UpdatedBy),
			orm.BaseWhere.Public.EQ(true), qm.Or2(orm.BaseWhere.DeveloperID.EQ(org.ID)),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	var res []*dto.GetComponentResponse

	for _, a := range bases {
		// If the available filter is enabled, only include results that match it
		if available != nil {
			if canUseBase(org, a) != *available {
				continue
			}
		}

		base, err := GetBaseDTO(org, a)
		if err != nil {
			return nil, err
		}

		res = append(res, base)
	}

	return res, nil
}
