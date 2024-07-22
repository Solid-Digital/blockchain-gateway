package fileparser_action_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/fileparser_action"
	"github.com/unchain/pipeline/pkg/domain"
	"testing"
)

func (s *TestSuite) TestFileParser_ParseJson() {
	cases := map[string]struct{
		Stub          domain.Stub
		RawFile       []byte
		Success       bool
		ExpectedValue map[string]interface{}
	}{
		"converts json successfully": {
			s.logger,
			s.helper.BytesFromFile("./testdata/json/test.json"),
			true,
			map[string]interface{}{
				"key1": "value1",
				"key2": float64(2),
				"key3": map[string]interface{}{
					"key3_key1": "value3",
				},
			},
		},
		"fails when rawfile is not valid json": {
			s.logger,
			s.helper.BytesFromFile("./testdata/json/invalid.json"),
			false,
			nil,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := fileparser_action.ParseJson(tc.Stub, tc.RawFile)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, output)
			} else {
				require.Error(t, err)
			}
		})
	}
}