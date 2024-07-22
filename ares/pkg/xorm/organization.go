package xorm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

func GetOrganizationTx(ctx context.Context, tx *sql.Tx, name string, mods ...qm.QueryMod) (*orm.Organization, *apperr.Error) {
	org, err := orm.Organizations(
		append(mods, orm.OrganizationWhere.Name.EQ(name))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case errors.Is(err, apperr.NotFound):
			return nil, ares.ErrOrgNotFound(err, name)
		default:
			return nil, err
		}
	}

	return org, nil
}

func ExistsOrganizationTx(ctx context.Context, tx *sql.Tx, orgName string, mods ...qm.QueryMod) (bool, *apperr.Error) {
	exists, err := orm.Organizations(append(
		mods,
		orm.OrganizationWhere.Name.EQ(orgName))...,
	).Exists(ctx, tx)

	if err != nil {
		return false, ares.ParsePQErr(err)
	}

	return exists, nil
}

func CreateOrganizationTx(ctx context.Context, tx *sql.Tx, orgName string, orgDisplayName string, creator *orm.User) (*orm.Organization, *apperr.Error) {
	org := &orm.Organization{
		DisplayName: orgDisplayName,
		Name:        orgName,
		CreatedByID: creator.ID,
		UpdatedByID: creator.ID,
	}

	err := org.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case errors.Is(err, apperr.Conflict):
			return nil, ares.ErrDuplicateOrg(err, org.Name)
		default:
			return nil, err
		}
	}

	appErr := AddUserOrganizationTx(ctx, tx, creator, org)
	if appErr != nil {
		return nil, appErr
	}

	err = creator.SetDefaultOrganization(ctx, tx, false, org)
	if err != nil {
		return nil, ares.ParsePQErr(err)
	}

	// Add default environments to the newly created organization
	envs, err := orm.DefaultEnvironments().All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case errors.Is(err, apperr.NotFound):
			// Is this actually possible?
			return nil, ares.ErrDefaultEnvsNotFound(err)
		default:
			return nil, err
		}
	}

	for _, env := range envs {
		err = org.AddEnvironments(ctx, tx, true, &orm.Environment{
			Index:       env.Index,
			Name:        env.Name,
			MaxReplicas: env.MaxReplicas,
			CreatedByID: creator.ID,
			UpdatedByID: creator.ID,
		})

		if err != nil {
			return nil, ares.ParsePQErr(err)
		}
	}

	return org, nil
}
