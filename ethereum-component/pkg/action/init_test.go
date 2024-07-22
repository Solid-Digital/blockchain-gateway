package action_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/action"

	"github.com/stretchr/testify/require"
	"github.com/unchainio/interfaces/adapter"
)

func (s *TestSuite) TestAction_Init() {
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
		"host is not specified": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/no_host.toml"),
			false,
		},
		"with redis": {
			s.logger,
			s.helper.BytesFromFile("./testdata/config/with_redis.toml"),
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			action := new(action.Action)
			err := action.Init(tc.Stub, tc.Config)

			if tc.Success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
