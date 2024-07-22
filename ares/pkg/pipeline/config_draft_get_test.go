package pipeline_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetDraftConfiguratio() {
	err := s.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		v, err := orm.ActionVersions(
			qm.Load(qm.Rels(orm.ActionVersionRels.Action, orm.ActionRels.ActionVersions)),
			orm.ActionVersionWhere.ID.EQ(149),
		).All(ctx, tx)
		if err != nil {
			return err
		}
		for _, v := range v {
			for _, av := range v.R.Action.R.ActionVersions {
				fmt.Println(av.Version)
			}
			fmt.Println("----------------------")
		}

		return nil
	})

	s.Suite.Require().NoError(err)

}
func (s *TestSuite) TestService_GetDraftConfiguration() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1, draftConfig1 := s.factory.PipelineWithDraft(true, org1, user1)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2, draftConfig2 := s.factory.PipelineWithDraft(true, org2, user2)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3, draftConfig3 := s.factory.PipelineWithDraft(true, org3, user3)
	baseDraftConfig3 := s.factory.BaseDraftConfiguration(false, true, draftConfig3, org3)

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4, draftConfig4 := s.factory.PipelineWithDraft(true, org4, user4)
	triggerDraftConfig4 := s.factory.TriggerDraftConfiguration(false, true, draftConfig4, org4)

	org5, user5 := s.factory.OrganizationAndUser(true)
	pipeline5, draftConfig5 := s.factory.PipelineWithDraft(true, org5, user5)
	cfgs5 := s.factory.ManyActionDraftConfigurations(draftConfig5, org5, 10)

	cases := map[string]struct {
		Organization       *orm.Organization
		Pipeline           *orm.Pipeline
		DraftConfig        *orm.DraftConfiguration
		BaseDraftConfig    *orm.BaseDraftConfiguration
		TriggerDraftConfig *orm.TriggerDraftConfiguration
		ActionDraftConfigs []*orm.ActionDraftConfiguration
		Success            bool
	}{
		"get draft config": {
			org1,
			pipeline1,
			draftConfig1,
			nil,
			nil,
			nil,
			true,
		},
		"pipeline not from organization": {
			s.factory.Organization(true),
			pipeline2,
			draftConfig2,
			nil,
			nil,
			nil,
			false,
		},
		"with base": {
			org3,
			pipeline3,
			draftConfig3,
			baseDraftConfig3,
			nil,
			nil,
			true,
		},
		"with trigger": {
			org4,
			pipeline4,
			draftConfig4,
			nil,
			triggerDraftConfig4,
			nil,
			true,
		},
		"with actions": {
			org5,
			pipeline5,
			draftConfig5,
			nil,
			nil,
			cfgs5,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetDraftConfiguration(tc.Organization.Name, tc.Pipeline.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.DraftConfig.ID, *response.ID)

				if tc.BaseDraftConfig != nil {
					require.Equal(t, tc.BaseDraftConfig.Config, *response.Base.InitConfiguration)
				} else {
					require.Nil(t, response.Base)
				}

				if tc.TriggerDraftConfig != nil {
					require.Equal(t, tc.TriggerDraftConfig.Config, *response.Trigger.InitConfiguration)
					require.True(t, tc.TriggerDraftConfig.MessageConfig.Valid)

					require.Equal(t,
						s.helper.UnmarshalTOML(*response.Trigger.MessageConfiguration),
						s.helper.UnmarshalJSON(tc.TriggerDraftConfig.MessageConfig.JSON),
					)
				} else {
					require.Nil(t, response.Trigger)
				}
				if tc.ActionDraftConfigs != nil {
					require.Equal(t, len(tc.ActionDraftConfigs), len(response.Actions))

					for i, cfg := range tc.ActionDraftConfigs {
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
