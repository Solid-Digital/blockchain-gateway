package amqp_trigger_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/amqp_trigger"
	"testing"
)

func (s *TestSuite) TestTrigger_Init() {
	cases := map[string]struct {
		Stub    domain.Stub
		Config  []byte
		Success bool
	}{
		"init trigger with valid config triggers as expected": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/config.toml"),
			true,
		},
		"init trigger with empty config fails": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/empty_config.toml"),
			false,
		},
		"init trigger with false cron spec fails": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/invalid_config.toml"),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			trigger := amqp_trigger.NewTrigger()
			err := trigger.Init(tc.Stub, tc.Config)

			if tc.Success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}



