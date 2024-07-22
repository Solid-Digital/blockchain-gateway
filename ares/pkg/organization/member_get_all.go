package organization

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetAllMembers(orgName string, principal *dto.User) ([]*dto.GetMemberResponse, *apperr.Error) {
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName,
			qm.Load(orm.OrganizationRels.Users),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	ret := make([]*dto.GetMemberResponse, len(org.R.Users))
	for i, user := range org.R.Users {
		ret[i] = getMember(user, s.enforcer.GetRolesForUserInOrganization(user.ID, orgName))
	}

	return ret, nil
}
