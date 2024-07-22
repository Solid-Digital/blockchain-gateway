package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdateDraftConfiguration() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1, draftConfig1 := s.factory.PipelineWithDraft(true, org1, user1)
	baseVersion1 := s.factory.BaseVersionForOrg(org1, false, true)
	params1 := &dto.UpdateDraftConfigurationRequest{
		Base: &dto.UpdateComponentConfiguration{
			ID:                s.helper.Int64Ptr(baseVersion1.ID),
			InitConfiguration: testhelper.StringPtr(testhelper.Randumb("my new base config")),
		},
	}

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2, draftConfig2 := s.factory.PipelineWithDraft(true, org2, user2)
	triggerVersion2 := s.factory.TriggerVersionForOrg(org2, false, true)
	params2 := &dto.UpdateDraftConfigurationRequest{
		Trigger: &dto.UpdateComponentConfiguration{
			ID:                   s.helper.Int64Ptr(triggerVersion2.ID),
			InitConfiguration:    testhelper.StringPtr(testhelper.Randumb("my new trigger config")),
			Name:                 "trigger",
			MessageConfiguration: s.helper.RandomTOML(),
		},
	}

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3, draftConfig3 := s.factory.PipelineWithDraft(true, org3, user3)
	actionVersion3 := s.factory.ActionVersionForOrg(org3, false, true)
	params3 := &dto.UpdateDraftConfigurationRequest{
		Actions: []*dto.UpdateComponentConfiguration{
			{
				ID:                   s.helper.Int64Ptr(actionVersion3.ID),
				InitConfiguration:    testhelper.StringPtr(testhelper.Randumb("my new action config")),
				Name:                 testhelper.Randumb(randomdata.SillyName()),
				MessageConfiguration: s.helper.RandomTOML(),
			},
			{
				ID:                   s.helper.Int64Ptr(actionVersion3.ID),
				InitConfiguration:    testhelper.StringPtr(testhelper.Randumb("my other action config")),
				Name:                 testhelper.Randumb(randomdata.SillyName()),
				MessageConfiguration: s.helper.RandomTOML(),
			},
		},
	}

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4, draftConfig4 := s.factory.PipelineWithDraft(true, org4, user4)
	baseVersion4 := s.factory.BaseVersionForOrg(org4, false, true)
	triggerVersion4 := s.factory.TriggerVersionForOrg(org4, false, true)
	actionVersion4 := s.factory.ActionVersionForOrg(org4, false, true)
	params4 := &dto.UpdateDraftConfigurationRequest{
		Base: &dto.UpdateComponentConfiguration{
			ID:                s.helper.Int64Ptr(baseVersion4.ID),
			InitConfiguration: testhelper.StringPtr(testhelper.Randumb("should not update base config")),
		},
		Trigger: &dto.UpdateComponentConfiguration{
			ID:                s.helper.Int64Ptr(triggerVersion4.ID),
			InitConfiguration: testhelper.StringPtr("should not update trigger config"),
		},
		Actions: []*dto.UpdateComponentConfiguration{
			{
				ID:                s.helper.Int64Ptr(actionVersion4.ID),
				InitConfiguration: testhelper.StringPtr(testhelper.Randumb("should not update action config")),
				Name:              testhelper.Randumb(randomdata.SillyName()),
			},
			{
				ID:                s.helper.Int64Ptr(actionVersion4.ID),
				InitConfiguration: s.helper.StringPtr(testhelper.Randumb("should not update action config")),
				Name:              testhelper.Randumb(randomdata.SillyName()),
			},
		},
	}

	org5, user5 := s.factory.OrganizationAndUser(true)
	pipeline5, draftConfig5 := s.factory.PipelineWithDraft(true, org5, user5)
	baseVersion5 := s.factory.BaseVersionForOrg(org5, false, true)
	triggerVersion5 := s.factory.TriggerVersionForOrg(org5, false, true)
	actionVersion5 := s.factory.ActionVersionForOrg(org5, false, true)
	params5 := &dto.UpdateDraftConfigurationRequest{
		Base: &dto.UpdateComponentConfiguration{
			ID:                s.helper.Int64Ptr(baseVersion5.ID),
			InitConfiguration: testhelper.StringPtr(testhelper.Randumb("should not update base config")),
		},
		Trigger: &dto.UpdateComponentConfiguration{
			ID:                s.helper.Int64Ptr(triggerVersion5.ID),
			InitConfiguration: testhelper.StringPtr("should not update trigger config"),
		},
		Actions: []*dto.UpdateComponentConfiguration{
			{
				ID:                s.helper.Int64Ptr(actionVersion5.ID),
				InitConfiguration: testhelper.StringPtr(testhelper.Randumb("should not update action config")),
				Name:              testhelper.Randumb(randomdata.SillyName()),
			},
			{
				ID:                s.helper.Int64Ptr(actionVersion5.ID),
				InitConfiguration: testhelper.StringPtr(testhelper.Randumb("should not update action config")),
				Name:              testhelper.Randumb(randomdata.SillyName()),
			},
		},
	}

	cases := map[string]struct {
		Params        *dto.UpdateDraftConfigurationRequest
		Organization  *orm.Organization
		Pipeline      *orm.Pipeline
		Principal     *dto.User
		DraftConfig   *orm.DraftConfiguration
		BaseConfig    string
		TriggerConfig string
		ActionConfigs []*dto.UpdateComponentConfiguration
		Success       bool
	}{
		"update base": {
			params1,
			org1,
			pipeline1,
			s.factory.DTOUser(true),
			draftConfig1,
			*params1.Base.InitConfiguration,
			"",
			nil,
			true,
		},
		"update trigger": {
			params2,
			org2,
			pipeline2,
			s.factory.DTOUser(true),
			draftConfig2,
			"",
			*params2.Trigger.InitConfiguration,
			nil,
			true,
		},
		"update action": {
			params3,
			org3,
			pipeline3,
			s.factory.DTOUser(true),
			draftConfig3,
			"",
			"",
			params3.Actions,
			true,
		},
		"pipeline not from organization": {
			params4,
			s.factory.Organization(true),
			pipeline4,
			s.factory.DTOUser(true),
			draftConfig4,
			"",
			"",
			nil,
			false,
		},
		"principal does not exist": {
			params5,
			org5,
			pipeline5,
			s.factory.DTOUser(false),
			draftConfig5,
			"",
			"",
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateDraftConfiguration(tc.Params, tc.Organization.Name, tc.Pipeline.Name, tc.Principal)
			draftConfigFromDB := s.helper.DBGetDraftConfiguration(tc.DraftConfig.ID)
			require.NotNil(t, draftConfigFromDB)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Principal.ID, response.UpdatedBy.ID)

				if tc.Params.Trigger != nil && tc.Params.Trigger.MessageConfiguration != "" {
					require.True(t, draftConfigFromDB.R.TriggerDraftConfiguration.MessageConfig.Valid)

					require.Equal(t,
						s.helper.UnmarshalJSON(draftConfigFromDB.R.TriggerDraftConfiguration.MessageConfig.JSON),
						s.helper.UnmarshalTOML(tc.Params.Trigger.MessageConfiguration),
					)
				}

				if tc.BaseConfig != "" {
					require.Equal(t, tc.BaseConfig, *response.Base.InitConfiguration)
					require.Equal(t, tc.BaseConfig, draftConfigFromDB.R.BaseDraftConfiguration.Config)
				} else {
					require.Nil(t, response.Base)
				}

				if tc.TriggerConfig != "" {
					require.Equal(t, tc.TriggerConfig, *response.Trigger.InitConfiguration)
					require.Equal(t, tc.TriggerConfig, draftConfigFromDB.R.TriggerDraftConfiguration.Config)
				} else {
					require.Nil(t, response.Trigger)
				}

				require.Equal(t, len(tc.Params.Actions), len(response.Actions))

				if tc.ActionConfigs != nil {
					require.Equal(t, len(tc.ActionConfigs), len(response.Actions))
					require.Equal(t, len(tc.ActionConfigs), len(draftConfigFromDB.R.ActionDraftConfigurations))

					for i, a := range tc.ActionConfigs {
						require.Equal(t, *a.InitConfiguration, *response.Actions[i].InitConfiguration)
						require.Equal(t, *a.InitConfiguration, draftConfigFromDB.R.ActionDraftConfigurations[i].Config)

						if a.MessageConfiguration != "" {
							require.True(t, draftConfigFromDB.R.ActionDraftConfigurations[i].MessageConfig.Valid)

							require.Equal(t,
								s.helper.UnmarshalJSON(draftConfigFromDB.R.ActionDraftConfigurations[i].MessageConfig.JSON),
								s.helper.UnmarshalTOML(a.MessageConfiguration),
							)

							require.Equal(t,
								s.helper.UnmarshalTOML(*response.Actions[i].MessageConfiguration),
								s.helper.UnmarshalTOML(a.MessageConfiguration),
							)
						}
					}
				} else {
					require.Equal(t, []*dto.ComponentConfiguration{}, response.Actions)
				}

			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
				require.Nil(t, draftConfigFromDB.R.BaseDraftConfiguration)
				require.Nil(t, draftConfigFromDB.R.TriggerDraftConfiguration)
				require.Nil(t, draftConfigFromDB.R.ActionDraftConfigurations)
			}
		})
	}
}
