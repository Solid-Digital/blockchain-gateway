package api_trigger_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger"
	"testing"
)

func (s *TestSuite) TestTrigger_Init() {
	cases := map[string]struct {
		Stub     domain.Stub
		Config   []byte
		Success  bool
		Shutdown bool
	}{
		"init trigger with empty config": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/config.toml"),
			true,
			true,
		},
		"init trigger with conflicting config": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/config_conflicting_auth.toml"),
			false,
			false, // since the config is invalid the server is not started so nothing to shutdown
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			trigger := new(api_trigger.Trigger)
			err := trigger.Init(tc.Stub, tc.Config)

			if tc.Success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			if tc.Shutdown {
				trigger.Close()
			}
		})
	}
}
