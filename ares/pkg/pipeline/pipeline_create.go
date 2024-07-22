package pipeline

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) CreatePipeline(params *dto.CreatePipelineRequest, orgName string, principal *dto.User) (*dto.GetPipelineResponse, *apperr.Error) {
	var ret *dto.GetPipelineResponse
	err := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		user, appErr := xorm.GetUserTxByID(ctx, tx, principal.ID)
		if appErr != nil {
			return appErr
		}

		pipeline := &orm.Pipeline{
			OrganizationID: org.ID,
			CreatedByID:    user.ID,
			UpdatedByID:    user.ID,
			Name:           params.Name,
			DisplayName:    params.DisplayName,
			Description:    params.Description,
		}

		err := pipeline.Insert(ctx, tx, boil.Infer())
		if err != nil {
			err := ares.ParsePQErr(err)
			switch {
			case stderr.Is(err, apperr.Conflict):
				return ares.ErrDuplicatePipeline(err, orgName, params.Name)
			default:
				return err
			}
		}

		err = pipeline.L.LoadCreatedBy(ctx, tx, true, pipeline, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), params.Name)
		}

		err = pipeline.L.LoadUpdatedBy(ctx, tx, true, pipeline, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), params.Name)
		}

		draft := &orm.DraftConfiguration{
			CreatedByID:    principal.ID,
			UpdatedByID:    principal.ID,
			OrganizationID: org.ID,
			Revision:       1,
		}

		err = pipeline.SetDraftConfiguration(ctx, tx, true, draft)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = draft.SetBaseDraftConfiguration(ctx, tx, true, &orm.BaseDraftConfiguration{})
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = draft.SetTriggerDraftConfiguration(ctx, tx, true, &orm.TriggerDraftConfiguration{})
		if err != nil {
			return ares.ParsePQErr(err)
		}

		envs, appErr := xorm.GetAllOrgEnvironmentsTx(ctx, tx, org)
		if appErr != nil {
			return appErr
		}

		ret, appErr = s.toDTOPipeline(orgName, pipeline, envs)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}
