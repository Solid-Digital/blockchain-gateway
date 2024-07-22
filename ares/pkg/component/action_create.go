package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) CreateAction(params *dto.CreateComponentRequest, orgName string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error) {
	var action *orm.Action
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, action, appErr = CreateActionTx(ctx, tx, params, orgName, principal)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetActionDTO(org, action)
}

func CreateActionTx(ctx context.Context, tx *sql.Tx, params *dto.CreateComponentRequest, orgName string, principal *dto.User) (*orm.Organization, *orm.Action, *apperr.Error) {
	var err error
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, nil, appErr
	}

	action := &orm.Action{
		DeveloperID: org.ID,
		CreatedByID: principal.ID,
		UpdatedByID: principal.ID,
		Name:        params.Name,
		DisplayName: params.DisplayName,
		Description: params.Description,
		Public:      false,
	}

	appErr = xorm.CreateActionTx(ctx, tx, action)
	if appErr != nil {
		return nil, nil, appErr
	}

	err = action.L.LoadCreatedBy(ctx, tx, true, action, nil)
	if err != nil {
		return nil, nil, ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), params.Name)
	}

	err = action.L.LoadUpdatedBy(ctx, tx, true, action, nil)
	if err != nil {
		return nil, nil, ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), params.Name)
	}

	return org, action, nil
}
