package fileparser_action_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/fileparser_action"
	"github.com/unchain/pipeline/pkg/domain"
	"testing"
)

func (s *TestSuite) TestFileParser_ParseCsv() {
	cases := map[string]struct {
		Stub          domain.Stub
		RawFile       []byte
		Header        bool
		Delimiter     rune
		Success       bool
		ExpectedValue map[string]interface{}
	}{
		"converts csv with header successfully": {
			s.logger,
			s.helper.BytesFromFile("./testdata/csv/test_headers.csv"),
			true,
			',',
			true,
			map[string]interface{}{
				"messages": []map[string]interface{}{
					map[string]interface{}{
						"header_value1": "row1_value1",
						"header_value2": "row1_value2",
						"header_value3": "row1_value3",
					},
					map[string]interface{}{
						"header_value1": "row2_value1",
						"header_value2": "row2_value2",
						"header_value3": "row2_value3",
					},
				},
			},
		},
		"converts csv without header successfully": {
			s.logger,
			s.helper.BytesFromFile("./testdata/csv/test_no_headers.csv"),
			false,
			',',
			true,
			map[string]interface{}{
				"messages": []map[string]interface{}{
					map[string]interface{}{
						"col-0": "row1_value1",
						"col-1": "row1_value2",
						"col-2": "row1_value3",
					},
					map[string]interface{}{
						"col-0": "row2_value1",
						"col-1": "row2_value2",
						"col-2": "row2_value3",
					},
				},
			},
		},
		"fails when headers don't match values": {
			s.logger,
			s.helper.BytesFromFile("./testdata/csv/test_bad_headers.csv"),
			true,
			',',
			false,
			nil,
		},
		"fails when rawfile is not a valid csv": {
			s.logger,
			s.helper.BytesFromFile("./testdata/csv/invalid.csv"),
			false,
			',',
			false,
			nil,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := fileparser_action.ParseCsv(tc.Stub, tc.RawFile, tc.Header, tc.Delimiter)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, output)
			} else {
				require.Error(t, err)
			}
		})
	}
}
