package component

import (
	"context"
	"database/sql"
	"encoding/json"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/unchainio/pkg/errors"
	"github.com/volatiletech/null"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) CreateTriggerVersion(params *ares.CreateTriggerVersionRequest) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var org *orm.Organization
	var trigger *orm.Trigger
	var triggerVersion *orm.TriggerVersion

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		if params.Public != nil && !s.service.enforcer.IsSuperAdmin(params.Principal.ID) {
			return ares.ErrForbiddenToSetPublic(stderr.New(""))
		}

		org, trigger, triggerVersion, appErr = s.createTriggerVersionTx(ctx, tx, params)
		if appErr != nil {
			return appErr
		}

		err = triggerVersion.L.LoadCreatedBy(ctx, tx, true, triggerVersion, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), params.Name)
		}

		err = triggerVersion.L.LoadUpdatedBy(ctx, tx, true, triggerVersion, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), params.Name)
		}

		file, err := TarGz(triggerVersion.FileName, params.TriggerFile)
		if err != nil {
			return apperr.New().Wrap(err)
		}

		err = s.service.store.PutObject(triggerVersion.FileID, file)
		if err != nil {
			return apperr.New().Wrap(err)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetTriggerVersionDTO(org, trigger, triggerVersion)
}

func (s *Service) createTriggerVersionTx(ctx context.Context, tx *sql.Tx, params *ares.CreateTriggerVersionRequest) (*orm.Organization, *orm.Trigger, *orm.TriggerVersion, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, params.OrgName)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	trigger, appErr := upsertTriggerTx(ctx, tx, params.Name, params.Principal.ID, org, params.Description)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	if err := must(canEditTrigger(org, trigger)); err != nil {
		return nil, nil, nil, ares.ErrForbiddenEdit(err)
	}

	fileName := TriggerVersionFileName(params.Name, params.Version, org.Name)
	fileID := TriggerVersionFileID(params.Name, params.Version, org.Name, fileName)

	input, err := json.Marshal(params.InputSchema)
	if err != nil {
		return nil, nil, nil, ares.ErrParseInputSchema(errors.Wrap(err, ""))
	}

	output, err := json.Marshal(params.OutputSchema)
	if err != nil {
		return nil, nil, nil, ares.ErrParseOutputSchema(errors.Wrap(err, ""))
	}

	publicVal := false

	if params.Public != nil {
		publicVal = *params.Public

		// Once a component version is public, the component is also public
		if publicVal {
			trigger.Public = true
		}
	}

	triggerVersion := &orm.TriggerVersion{
		TriggerID:     trigger.ID,
		Version:       params.Version,
		CreatedByID:   params.Principal.ID,
		UpdatedByID:   params.Principal.ID,
		Public:        publicVal,
		ExampleConfig: params.ExampleConfig,
		Description:   params.Description,
		Readme:        params.Readme,
		FileName:      fileName,
		FileID:        fileID,
		InputSchema:   null.JSONFrom(input),
		OutputSchema:  null.JSONFrom(output),
	}

	err = triggerVersion.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return nil, nil, nil, ares.ErrDuplicateComponentVersion(err, ares.ComponentTypeTrigger, params.Name, params.Version)
		default:
			return nil, nil, nil, err
		}
	}

	_, err = trigger.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	return org, trigger, triggerVersion, nil
}

func upsertTriggerTx(ctx context.Context, tx *sql.Tx, name string, userID int64, org *orm.Organization, description string) (*orm.Trigger, *apperr.Error) {
	var trigger *orm.Trigger

	exists, err := orm.Triggers(orm.TriggerWhere.Name.EQ(name)).Exists(ctx, tx)
	if err != nil {
		return nil, ares.ParsePQErr(err)
	}

	if exists {
		trigger, appErr := xorm.GetTriggerTx(ctx, tx, org.Name, name)
		if appErr != nil {
			return nil, appErr
		}

		trigger.Description = description

		_, err = trigger.Update(ctx, tx, boil.Infer())
		if err != nil {
			return nil, ares.ParsePQErr(err)
		}

		return trigger, nil
	}

	trigger = &orm.Trigger{
		DeveloperID: org.ID,
		CreatedByID: userID,
		UpdatedByID: userID,
		Name:        name,
		DisplayName: name,
		Description: description,
		Public:      false,
	}

	appErr := xorm.CreateTriggerTx(ctx, tx, trigger)
	if appErr != nil {
		return nil, appErr
	}

	return trigger, nil
}
