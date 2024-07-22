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

func (s *Service) GetBaseVersion(orgName string, name string, version string) (*dto.GetBaseVersionResponse, *apperr.Error) {
	var base *orm.Base
	var baseVersion *orm.BaseVersion
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		base, appErr = xorm.GetBaseTx(ctx, tx, orgName, name)
		if appErr != nil {
			return appErr
		}

		baseVersion, appErr = xorm.GetBaseVersionTx(ctx, tx, orgName, base, version,
			qm.Load(orm.BaseVersionRels.CreatedBy),
			qm.Load(orm.BaseVersionRels.UpdatedBy),
			orm.BaseVersionWhere.Version.EQ(version),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	if err := must(canViewBaseVersion(org, base, baseVersion)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	return GetBaseVersionDTO(org, base, baseVersion)
}

func GetBaseVersionDTO(org *orm.Organization, base *orm.Base, version *orm.BaseVersion) (*dto.GetBaseVersionResponse, *apperr.Error) {
	return &dto.GetBaseVersionResponse{
		ID:             &version.ID,
		Version:        &version.Version,
		Description:    &version.Description,
		Readme:         &version.Readme,
		DockerImageRef: &version.DockerImageRef,
		Entrypoint:     &version.Entrypoint,
		Public:         &version.Public,
		Available:      testhelper.BoolPtr(canUseBaseVersion(org, base, version)),
		CreatedAt:      (*strfmt.DateTime)(&version.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       version.R.CreatedBy.ID,
			FullName: version.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&version.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       version.R.UpdatedBy.ID,
			FullName: version.R.UpdatedBy.FullName,
		},
	}, nil
}
