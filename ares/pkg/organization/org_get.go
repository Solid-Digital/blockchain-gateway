package organization

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetOrganization(orgName string, u *dto.User) (*dto.GetOrganizationResponse, *apperr.Error) {
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return getOrg(org), nil
}

func getOrg(org *orm.Organization) *dto.GetOrganizationResponse {
	return &dto.GetOrganizationResponse{
		DisplayName: org.DisplayName,
		ID:          org.ID,
		Name:        org.Name,
	}
}
