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

func (s *Service) CreateActionVersion(params *ares.CreateActionVersionRequest) (*dto.GetComponentVersionResponse, *apperr.Error) {
	var org *orm.Organization
	var action *orm.Action
	var actionVersion *orm.ActionVersion

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		if params.Public != nil && !s.service.enforcer.IsSuperAdmin(params.Principal.ID) {
			return ares.ErrForbiddenToSetPublic(errors.New(""))
		}

		org, action, actionVersion, appErr = s.createActionVersionTx(ctx, tx, params)
		if appErr != nil {
			return appErr
		}

		err = actionVersion.L.LoadCreatedBy(ctx, tx, true, actionVersion, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), params.Name)
		}

		err = actionVersion.L.LoadUpdatedBy(ctx, tx, true, actionVersion, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), params.Name)
		}

		file, err := TarGz(actionVersion.FileName, params.ActionFile)
		if err != nil {
			return apperr.New().Wrap(err)
		}

		err = s.service.store.PutObject(actionVersion.FileID, file)
		if err != nil {
			return apperr.New().Wrap(err)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return GetActionVersionDTO(org, action, actionVersion)
}

func GetActionVersionIOSchema(actionVersion *orm.ActionVersion) (input []string, output []string, appErr *apperr.Error) {
	if actionVersion.InputSchema.Valid {
		err := json.Unmarshal(actionVersion.InputSchema.JSON, &input)
		if err != nil {
			return nil, nil, ares.ErrParseInputSchema(errors.Wrap(err, ""))
		}
	}

	if actionVersion.OutputSchema.Valid {
		err := json.Unmarshal(actionVersion.OutputSchema.JSON, &output)
		if err != nil {
			return nil, nil, ares.ErrParseOutputSchema(errors.Wrap(err, ""))
		}
	}

	return input, output, nil
}

func (s *Service) createActionVersionTx(ctx context.Context, tx *sql.Tx, params *ares.CreateActionVersionRequest) (*orm.Organization, *orm.Action, *orm.ActionVersion, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, params.OrgName)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	action, appErr := upsertActionTx(ctx, tx, params.Name, params.Principal.ID, org, params.Description)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	if err := must(canEditAction(org, action)); err != nil {
		return nil, nil, nil, ares.ErrForbiddenEdit(err)
	}

	fileName := ActionVersionFileName(params.Name, params.Version, org.Name)
	fileID := ActionVersionFileID(params.Name, params.Version, org.Name, fileName)

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
			action.Public = true
		}
	}

	actionVersion := &orm.ActionVersion{
		CreatedByID:   params.Principal.ID,
		UpdatedByID:   params.Principal.ID,
		ActionID:      action.ID,
		Version:       params.Version,
		Public:        publicVal,
		ExampleConfig: params.ExampleConfig,
		Description:   params.Description,
		Readme:        params.Readme,
		FileName:      fileName,
		FileID:        fileID,
		InputSchema:   null.JSONFrom(input),
		OutputSchema:  null.JSONFrom(output),
	}

	err = actionVersion.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return nil, nil, nil, ares.ErrDuplicateComponentVersion(err, ares.ComponentTypeAction, params.Name, params.Version)
		default:
			return nil, nil, nil, err
		}
	}

	_, err = action.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, nil, nil, ares.ParsePQErr(err)
	}

	return org, action, actionVersion, nil
}

func upsertActionTx(ctx context.Context, tx *sql.Tx, name string, userID int64, org *orm.Organization, description string) (*orm.Action, *apperr.Error) {
	var action *orm.Action
	var appErr *apperr.Error

	exists, err := orm.Actions(orm.ActionWhere.Name.EQ(name)).Exists(ctx, tx)
	if err != nil {
		return nil, ares.ParsePQErr(err)
	}

	if exists {
		action, appErr = xorm.GetActionTx(ctx, tx, org.Name, name)
		if appErr != nil {
			return nil, appErr
		}

		action.Description = description

		_, err = action.Update(ctx, tx, boil.Infer())
		if err != nil {
			return nil, ares.ParsePQErr(err)
		}

		return action, nil
	}

	action = &orm.Action{
		DeveloperID: org.ID,
		CreatedByID: userID,
		UpdatedByID: userID,
		Name:        name,
		DisplayName: name,
		Description: description,
		Public:      false,
	}

	appErr = xorm.CreateActionTx(ctx, tx, action)
	if appErr != nil {
		return nil, appErr
	}

	return action, nil
}
