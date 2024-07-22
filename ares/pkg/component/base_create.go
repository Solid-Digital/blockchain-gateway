package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) CreateBase(params *dto.CreateComponentRequest, orgName string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error) {
	var base *orm.Base
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		base = &orm.Base{
			DeveloperID: org.ID,
			CreatedByID: principal.ID,
			UpdatedByID: principal.ID,
			Name:        params.Name,
			DisplayName: params.DisplayName,
			Description: params.Description,
			Public:      false,
		}

		appErr = xorm.CreateBaseTx(ctx, tx, base)
		if appErr != nil {
			return appErr
		}

		err = base.L.LoadCreatedBy(ctx, tx, true, base, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), params.Name)
		}

		err = base.L.LoadUpdatedBy(ctx, tx, true, base, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), params.Name)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetBaseDTO(org, base)
}
