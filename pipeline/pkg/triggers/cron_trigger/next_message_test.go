package cron_trigger_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/cron_trigger"
	"testing"
)

func (s *TestSuite) TestTrigger_NextMessage() {
	cases := map[string]struct {
		Stub          domain.Stub
		Config        []byte
		Success       bool
		ExpectedValue string
	}{
		"init trigger with valid config triggers as expected": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/trigger_test_config.toml"),
			true,
			"1",
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			trigger := cron_trigger.NewTrigger()
			err := trigger.Init(tc.Stub, tc.Config)

			tag, _, _ := trigger.NextMessage()
			require.Equal(t, tc.ExpectedValue, tag)

			//  Expect trigger every so many seconds
			if tc.Success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
