package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetConfiguration() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)
	config1 := s.factory.Configuration(true, org1, user1, pipeline1)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)
	config2 := s.factory.Configuration(true, org2, user2, pipeline2)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)
	config3 := s.factory.Configuration(true, org3, user3, pipeline3)
	baseConfig3 := s.factory.BaseConfiguration(false, true, config3, org3)

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4 := s.factory.Pipeline(true, org4, user4)
	config4 := s.factory.Configuration(true, org4, user4, pipeline4)
	triggerConfig4 := s.factory.TriggerConfiguration(false, true, config4, org4)

	org5, user5 := s.factory.OrganizationAndUser(true)
	pipeline5 := s.factory.Pipeline(true, org5, user5)
	config5 := s.factory.Configuration(true, org5, user5, pipeline5)
	cfgs5 := s.factory.ManyActionCongigurations(config5, org5, 10)

	cases := map[string]struct {
		Organization  *orm.Organization
		Pipeline      *orm.Pipeline
		Revision      int64
		BaseConfig    *orm.BaseConfiguration
		TriggerConfig *orm.TriggerConfiguration
		ActionConfigs []*orm.ActionConfiguration
		Success       bool
	}{
		"get config": {
			org1,
			pipeline1,
			config1.Revision,
			nil,
			nil,
			nil,
			true,
		},
		"pipeline not from organization": {
			s.factory.Organization(true),
			pipeline2,
			config2.Revision,
			nil,
			nil,
			nil,
			false,
		},
		"with base": {
			org3,
			pipeline3,
			config3.Revision,
			baseConfig3,
			nil,
			nil,
			true,
		},
		"with trigger": {
			org4,
			pipeline4,
			config4.Revision,
			nil,
			triggerConfig4,
			nil,
			true,
		},
		"with actions": {
			org5,
			pipeline5,
			config5.Revision,
			nil,
			nil,
			cfgs5,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetConfiguration(tc.Organization.Name, tc.Pipeline.Name, tc.Revision)

			if tc.Success {
				xrequire.NoError(t, err)
				if tc.BaseConfig != nil {
					require.Equal(t, tc.BaseConfig.Config, *response.Base.InitConfiguration)
				} else {
					require.Nil(t, response.Base)
				}

				if tc.TriggerConfig != nil {
					require.Equal(t, tc.TriggerConfig.Config, *response.Trigger.InitConfiguration)
					require.True(t, tc.TriggerConfig.MessageConfig.Valid)

					require.Equal(t,
						s.helper.UnmarshalJSON(tc.TriggerConfig.MessageConfig.JSON),
						s.helper.UnmarshalTOML(*response.Trigger.MessageConfiguration),
					)
				} else {
					require.Nil(t, response.Trigger)
				}

				if tc.ActionConfigs != nil {
					require.Equal(t, len(tc.ActionConfigs), len(response.Actions))

					for i, cfg := range tc.ActionConfigs {
						require.Equal(t, cfg.Config, *response.Actions[i].InitConfiguration)
						require.True(t, cfg.MessageConfig.Valid)

						require.Equal(t,
							s.helper.UnmarshalJSON(cfg.MessageConfig.JSON),
							s.helper.UnmarshalTOML(*response.Actions[i].MessageConfiguration),
						)
					}
				} else {
					require.Equal(t, 0, len(response.Actions))
				}
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
