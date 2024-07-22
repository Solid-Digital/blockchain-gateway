package component

import (
	"context"
	"database/sql"
	"encoding/json"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"github.com/go-openapi/strfmt"
	"github.com/unchainio/pkg/errors"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetTriggerVersion(orgName string, name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var trigger *orm.Trigger
	var triggerVersion *orm.TriggerVersion
	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		trigger, appErr = xorm.GetTriggerTx(ctx, tx, orgName, name)
		if appErr != nil {
			return appErr
		}

		triggerVersion, appErr = xorm.GetTriggerVersionTx(ctx, tx, orgName, trigger, version,
			qm.Load(orm.TriggerVersionRels.CreatedBy),
			qm.Load(orm.TriggerVersionRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	if err := must(canViewTriggerVersion(org, trigger, triggerVersion)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	return GetTriggerVersionDTO(org, trigger, triggerVersion)
}

func (s *Service) GetPublicTriggerVersion(name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var trigger *orm.Trigger
	var triggerVersion *orm.TriggerVersion
	var emptyOrg *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		trigger, appErr = xorm.GetTriggerTx(ctx, tx, "public", name)
		if appErr != nil {
			return appErr
		}

		triggerVersion, appErr = xorm.GetTriggerVersionTx(ctx, tx, "public", trigger, version,
			qm.Load(orm.TriggerVersionRels.CreatedBy),
			qm.Load(orm.TriggerVersionRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	if err := must(canViewTriggerVersion(emptyOrg, trigger, triggerVersion)); err != nil {
		return nil, ares.ErrForbiddenView(err)
	}

	return GetTriggerVersionDTO(emptyOrg, trigger, triggerVersion)
}

func GetTriggerVersionDTO(org *orm.Organization, component *orm.Trigger, version *orm.TriggerVersion) (*dto.GetComponentVersionResponse, *apperr.Error) {
	input, output, err := GetTriggerVersionIOSchema(version)
	if err != nil {
		return nil, err
	}

	return &dto.GetComponentVersionResponse{
		ID:             &version.ID,
		Version:        &version.Version,
		ReleaseMessage: &version.Description,
		Readme:         &version.Readme,
		ExampleConfig:  &version.ExampleConfig,
		InputSchema:    input,
		OutputSchema:   output,
		Public:         &version.Public,
		Available:      testhelper.BoolPtr(canUseTriggerVersion(org, component, version)),
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

func GetTriggerVersionIOSchema(triggerVersion *orm.TriggerVersion) (input []string, output []string, appErr *apperr.Error) {
	if triggerVersion.InputSchema.Valid {
		err := json.Unmarshal(triggerVersion.InputSchema.JSON, &input)
		if err != nil {
			return nil, nil, ares.ErrParseInputSchema(errors.Wrap(err, ""))
		}
	}

	if triggerVersion.OutputSchema.Valid {
		err := json.Unmarshal(triggerVersion.OutputSchema.JSON, &output)
		if err != nil {
			return nil, nil, ares.ErrParseOutputSchema(errors.Wrap(err, ""))
		}
	}

	return input, output, nil
}
