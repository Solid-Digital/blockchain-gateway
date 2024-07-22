package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"github.com/go-openapi/strfmt"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetAction(orgName string, name string) (*dto.GetComponentResponse, *apperr.Error) {
	var action *orm.Action
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		action, appErr = xorm.GetActionTx(ctx, tx, orgName, name,
			qm.Load(orm.ActionRels.ActionVersions),
			qm.Load(orm.ActionRels.CreatedBy),
			qm.Load(orm.ActionRels.UpdatedBy),
		)
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

func (s *Service) GetPublicAction(name string) (*dto.GetComponentResponse, *apperr.Error) {
	var action *orm.Action

	err := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		action, appErr = xorm.GetActionTx(ctx, tx, "public", name,
			qm.Load(orm.ActionRels.ActionVersions),
			qm.Load(orm.ActionRels.CreatedBy),
			qm.Load(orm.ActionRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return GetPublicActionDTO(action)
}

func GetActionDTO(org *orm.Organization, action *orm.Action) (*dto.GetComponentResponse, *apperr.Error) {
	if err := must(canViewAction(org, action)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	var versions []string

	if action.R != nil && action.R.ActionVersions != nil {
		versions = make([]string, 0)

		for _, v := range action.R.ActionVersions {
			if canViewActionVersion(org, action, v) {
				versions = append(versions, v.Version)
			}
		}

		sortVersions(versions)

		versions = deduplicateVersions(versions)
	}

	return &dto.GetComponentResponse{
		ID:          &action.ID,
		DeveloperID: &action.DeveloperID,
		Name:        &action.Name,
		DisplayName: &action.DisplayName,
		Description: &action.Description,
		Public:      &action.Public,
		Available:   testhelper.BoolPtr(canUseAction(org, action)),
		Versions:    versions,
		CreatedAt:   (*strfmt.DateTime)(&action.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       action.R.CreatedBy.ID,
			FullName: action.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&action.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       action.R.UpdatedBy.ID,
			FullName: action.R.UpdatedBy.FullName,
		},
	}, nil
}

func GetPublicActionDTO(action *orm.Action) (*dto.GetComponentResponse, *apperr.Error) {
	if err := must(action.Public); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	var versions []string

	if action.R != nil && action.R.ActionVersions != nil {
		versions = make([]string, 0)

		for _, v := range action.R.ActionVersions {
			if v.Public {
				versions = append(versions, v.Version)
			}
		}

		sortVersions(versions)

		versions = deduplicateVersions(versions)
	}

	return &dto.GetComponentResponse{
		ID:          &action.ID,
		DeveloperID: &action.DeveloperID,
		Name:        &action.Name,
		DisplayName: &action.DisplayName,
		Description: &action.Description,
		Public:      &action.Public,
		Available:   testhelper.BoolPtr(action.Public),
		Versions:    versions,
		CreatedAt:   (*strfmt.DateTime)(&action.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       action.R.CreatedBy.ID,
			FullName: action.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&action.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       action.R.UpdatedBy.ID,
			FullName: action.R.UpdatedBy.FullName,
		},
	}, nil
}
