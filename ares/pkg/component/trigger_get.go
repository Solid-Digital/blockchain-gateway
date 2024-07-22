package component

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/testhelper"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (s *Service) GetTrigger(orgName string, name string) (*dto.GetComponentResponse, *apperr.Error) {
	var trigger *orm.Trigger
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		trigger, appErr = xorm.GetTriggerTx(ctx, tx, orgName, name,
			qm.Load(orm.TriggerRels.TriggerVersions),
			qm.Load(orm.TriggerRels.CreatedBy),
			qm.Load(orm.TriggerRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetTriggerDTO(org, trigger)
}

func (s *Service) GetPublicTrigger(name string) (*dto.GetComponentResponse, *apperr.Error) {
	var trigger *orm.Trigger

	err := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		trigger, appErr = xorm.GetTriggerTx(ctx, tx, "public", name,
			qm.Load(orm.TriggerRels.TriggerVersions),
			qm.Load(orm.TriggerRels.CreatedBy),
			qm.Load(orm.TriggerRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return GetPublicTriggerDTO(trigger)
}

func GetTriggerDTO(org *orm.Organization, trigger *orm.Trigger) (*dto.GetComponentResponse, *apperr.Error) {
	if err := must(canViewTrigger(org, trigger)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	var versions []string

	if trigger.R != nil && trigger.R.TriggerVersions != nil {
		versions = make([]string, 0)

		for _, v := range trigger.R.TriggerVersions {
			if canViewTriggerVersion(org, trigger, v) {
				versions = append(versions, v.Version)
			}
		}

		sortVersions(versions)
	}

	return &dto.GetComponentResponse{
		ID:          &trigger.ID,
		DeveloperID: &trigger.DeveloperID,
		Name:        &trigger.Name,
		DisplayName: &trigger.DisplayName,
		Description: &trigger.Description,
		Public:      &trigger.Public,
		Available:   testhelper.BoolPtr(canUseTrigger(org, trigger)),
		Versions:    versions,
		CreatedAt:   (*strfmt.DateTime)(&trigger.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       trigger.R.CreatedBy.ID,
			FullName: trigger.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&trigger.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       trigger.R.UpdatedBy.ID,
			FullName: trigger.R.UpdatedBy.FullName,
		},
	}, nil
}

func GetPublicTriggerDTO(trigger *orm.Trigger) (*dto.GetComponentResponse, *apperr.Error) {
	if err := must(trigger.Public); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	var versions []string

	if trigger.R != nil && trigger.R.TriggerVersions != nil {
		versions = make([]string, 0)

		for _, v := range trigger.R.TriggerVersions {
			if v.Public {
				versions = append(versions, v.Version)
			}
		}

		sortVersions(versions)
	}

	return &dto.GetComponentResponse{
		ID:          &trigger.ID,
		DeveloperID: &trigger.DeveloperID,
		Name:        &trigger.Name,
		DisplayName: &trigger.DisplayName,
		Description: &trigger.Description,
		Public:      &trigger.Public,
		Available:   testhelper.BoolPtr(trigger.Public),
		Versions:    versions,
		CreatedAt:   (*strfmt.DateTime)(&trigger.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       trigger.R.CreatedBy.ID,
			FullName: trigger.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&trigger.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       trigger.R.UpdatedBy.ID,
			FullName: trigger.R.UpdatedBy.FullName,
		},
	}, nil
}
