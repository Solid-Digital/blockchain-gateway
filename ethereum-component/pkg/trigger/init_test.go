package trigger_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/trigger"

	"github.com/stretchr/testify/require"
	"github.com/unchainio/interfaces/adapter"
)

func (s *TestSuite) TestTrigger_Init() {
	cases := map[string]struct {
		Stub    adapter.Stub
		Config  []byte
		Success bool
	}{
		"valid config": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/valid.toml"),
			true,
		},
		"host not configured": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/no_host.toml"),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			trigger := new(trigger.Trigger)
			err := trigger.Init(tc.Stub, tc.Config)

			if tc.Success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
