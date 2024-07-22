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

func (s *Service) GetBase(orgName string, name string) (*dto.GetComponentResponse, *apperr.Error) {
	var base *orm.Base
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		base, appErr = xorm.GetBaseTx(ctx, tx, orgName, name,
			qm.Load(orm.BaseRels.BaseVersions),
			qm.Load(orm.BaseRels.CreatedBy),
			qm.Load(orm.BaseRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetBaseDTO(org, base)
}

func GetBaseDTO(org *orm.Organization, base *orm.Base) (*dto.GetComponentResponse, *apperr.Error) {
	if err := must(canViewBase(org, base)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	var versions []string

	if base.R != nil && base.R.BaseVersions != nil {
		versions = make([]string, 0)

		for _, v := range base.R.BaseVersions {
			if canViewBaseVersion(org, base, v) {
				versions = append(versions, v.Version)
			}
		}

		sortVersions(versions)
	}

	return &dto.GetComponentResponse{
		ID:          &base.ID,
		DeveloperID: &base.DeveloperID,
		Name:        &base.Name,
		DisplayName: &base.DisplayName,
		Description: &base.Description,
		Public:      &base.Public,
		Available:   testhelper.BoolPtr(canUseBase(org, base)),
		Versions:    versions,
		CreatedAt:   (*strfmt.DateTime)(&base.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       base.R.CreatedBy.ID,
			FullName: base.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&base.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       base.R.UpdatedBy.ID,
			FullName: base.R.UpdatedBy.FullName,
		},
	}, nil
}
