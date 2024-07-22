package organization

import (
	"context"
	"database/sql"
	"errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetAllOrganizations(principal *dto.User) ([]*dto.GetOrganizationResponse, *apperr.Error) {
	var orgs orm.OrganizationSlice

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		orgs, appErr = getAllOrganizationsTx(ctx, tx)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	ret := make([]*dto.GetOrganizationResponse, len(orgs))
	for i, org := range orgs {
		ret[i] = getOrg(org)
	}

	return ret, nil
}

func getAllOrganizationsTx(ctx context.Context, tx *sql.Tx) (orm.OrganizationSlice, *apperr.Error) {
	orgs, err := orm.Organizations().All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case errors.Is(err, apperr.NotFound):
			// Is this actually possible?
			return nil, ares.ErrOrgsNotFound(err)
		default:
			return nil, err
		}
	}

	return orgs, nil
}
