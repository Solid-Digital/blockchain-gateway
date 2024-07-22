package pipeline_test

import (
	"context"
	"database/sql"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_SaveDraftConfigurationAsNew() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1, draftConfig1 := s.factory.PipelineWithDraft(true, org1, user1)
	baseDraftConfig1 := s.factory.BaseDraftConfiguration(false, true, draftConfig1, org1)
	triggerDraftConfig1 := s.factory.TriggerDraftConfiguration(false, true, draftConfig1, org1)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2, draftConfig2 := s.factory.PipelineWithDraft(true, org2, user2)
	baseDraftConfig2 := s.factory.BaseDraftConfiguration(false, true, draftConfig2, org2)
	triggerDraftConfig2 := s.factory.TriggerDraftConfiguration(false, true, draftConfig2, org2)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3, draftConfig3 := s.factory.PipelineWithDraft(true, org3, user3)
	baseDraftConfig3 := s.factory.BaseDraftConfiguration(false, true, draftConfig3, org3)
	triggerDraftConfig3 := s.factory.TriggerDraftConfiguration(false, true, draftConfig3, org3)
	for i := 0; i < 10; i++ {
		s.factory.ActionDraftConfiguration(false, true, draftConfig3, org3)
	}

	cases := map[string]struct {
		Params             *dto.SaveDraftConfigurationAsNewRequest
		Organization       *orm.Organization
		Pipeline           *orm.Pipeline
		DraftConfig        *orm.DraftConfiguration
		BaseDraftConfig    *orm.BaseDraftConfiguration
		TriggerDraftConfig *orm.TriggerDraftConfiguration
		Principal          *dto.User
		Success            bool
	}{
		"save config draft as new config": {
			&dto.SaveDraftConfigurationAsNewRequest{
				CommitMessage: testhelper.Randumb(randomdata.Paragraph()),
			},
			org1,
			pipeline1,
			draftConfig1,
			baseDraftConfig1,
			triggerDraftConfig1,
			s.factory.DTOUser(true),
			true,
		},
		"pipeline not from organization": {
			&dto.SaveDraftConfigurationAsNewRequest{
				CommitMessage: testhelper.Randumb(randomdata.Paragraph()),
			},
			s.factory.Organization(true),
			pipeline2,
			draftConfig2,
			baseDraftConfig2,
			triggerDraftConfig2,
			s.factory.DTOUser(true),
			false,
		},
		"preserve order of action components": {
			&dto.SaveDraftConfigurationAsNewRequest{
				CommitMessage: testhelper.Randumb(randomdata.Paragraph()),
			},
			org3,
			pipeline3,
			draftConfig3,
			baseDraftConfig3,
			triggerDraftConfig3,
			s.factory.DTOUser(true),
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.SaveDraftConfigurationAsNew(tc.Params, tc.Organization.Name, tc.Pipeline.Name, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Principal.ID, response.CreatedBy.ID)
				require.Equal(t, tc.Principal.ID, response.UpdatedBy.ID)
				require.Equal(t, tc.BaseDraftConfig.Config, *response.Base.InitConfiguration)

				configFromDB := s.helper.DBGetConfiguration(*response.ID)

				require.NotNil(t, configFromDB)
				require.Equal(t, tc.BaseDraftConfig.Config, configFromDB.R.BaseConfiguration.Config)
				require.Equal(t, tc.BaseDraftConfig.Config, configFromDB.R.BaseConfiguration.Config)

				if tc.TriggerDraftConfig != nil {
					require.Equal(t, tc.TriggerDraftConfig.MessageConfig, configFromDB.R.TriggerConfiguration.MessageConfig)
				}

				s.verifyActions(tc.DraftConfig)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) verifyActions(draftConfig *orm.DraftConfiguration) {
	var draftActionConfigs []*orm.ActionDraftConfiguration
	var err error

	err = s.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		draftActionConfigs, err = draftConfig.ActionDraftConfigurations().All(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	s.Require().NoError(err)

	for _, draftActionConfig := range draftActionConfigs {
		actionConfig := s.helper.GetActionConfiguration(draftActionConfig.Name)
		s.Require().NotNil(actionConfig)

		s.Require().Equal(draftActionConfig.Index, actionConfig.Index)
		s.Require().Equal(draftActionConfig.MessageConfig, actionConfig.MessageConfig)
	}
}
