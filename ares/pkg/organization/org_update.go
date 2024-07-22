package organization

import (
	"context"
	"database/sql"
	"errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
)

func (s *Service) UpdateOrganization(params *dto.UpdateOrganizationRequest, orgName string) (*dto.GetOrganizationResponse, *apperr.Error) {
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		org.DisplayName = params.DisplayName

		_, err = org.Update(ctx, tx, boil.Infer())
		if err != nil {
			err := ares.ParsePQErr(err)
			switch {
			case errors.Is(err, apperr.Conflict):
				return ares.ErrDuplicateOrg(err, org.Name)
			default:
				return err
			}
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return getOrg(org), nil
}
