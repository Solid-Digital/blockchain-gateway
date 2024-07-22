package fileparser_action_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/fileparser_action"
	"github.com/unchain/pipeline/pkg/domain"
	"testing"
)

func (s *TestSuite) TestFileParserAction_Invoke() {
	cases := map[string]struct {
		Stub          domain.Stub
		Input         map[string]interface{}
		Success       bool
		ExpectedValue map[string]interface{}
	}{
		"invoke file parser for csv": {
			s.logger,
			map[string]interface{}{
				"filetype": "csv",
				"file":     s.helper.BytesFromFile("./testdata/csv/test_headers.csv"),
				"header":   true,
			},
			true,
			map[string]interface{}{"messages": []map[string]interface{}{map[string]interface{}{"header_value1": "row1_value1", "header_value2": "row1_value2", "header_value3": "row1_value3"}, map[string]interface{}{"header_value1": "row2_value1", "header_value2": "row2_value2", "header_value3": "row2_value3"}}},
		},
		"invoke file parser for csv without header": {
			s.logger,
			map[string]interface{}{
				"filetype": "csv",
				"file":     s.helper.BytesFromFile("./testdata/csv/test_no_headers.csv"),
				"header":   false,
			},
			true,
			map[string]interface{}{"messages": []map[string]interface{}{map[string]interface{}{"col-0": "row1_value1", "col-1": "row1_value2", "col-2": "row1_value3"}, map[string]interface{}{"col-0": "row2_value1", "col-1": "row2_value2", "col-2": "row2_value3"}}},
		},
		"invoke file parser for csv falls back on no header": {
			s.logger,
			map[string]interface{}{
				"filetype": "csv",
				"file":     s.helper.BytesFromFile("./testdata/csv/test_no_headers.csv"),
			},
			true,
			map[string]interface{}{"messages": []map[string]interface{}{map[string]interface{}{"col-0": "row1_value1", "col-1": "row1_value2", "col-2": "row1_value3"}, map[string]interface{}{"col-0": "row2_value1", "col-1": "row2_value2", "col-2": "row2_value3"}}},
		},
		"invoke file parser for csv with ; delimiter": {
			s.logger,
			map[string]interface{}{
				"filetype":  "csv",
				"file":      s.helper.BytesFromFile("./testdata/csv/test_delimiter.csv"),
				"delimiter": ';',
			},
			true,
			map[string]interface{}{"messages": []map[string]interface{}{map[string]interface{}{"col-0": "row1_value1", "col-1": "row1_value2", "col-2": "row1_value3"}, map[string]interface{}{"col-0": "row2_value1", "col-1": "row2_value2", "col-2": "row2_value3"}}},
		},
		"invoke file parser for unknown filetype returns error": {
			s.logger,
			map[string]interface{}{
				"filetype": "unknown",
				"file":     s.helper.BytesFromFile("./testdata/csv/test_no_headers.csv"),
			},
			false,
			nil,
		},
		"invoke file parser for non string filetype returns error": {
			s.logger,
			map[string]interface{}{
				"filetype": 1,
				"file":     s.helper.BytesFromFile("./testdata/csv/test_no_headers.csv"),
			},
			false,
			nil,
		},
		"invoke file parser for non []byte file returns error": {
			s.logger,
			map[string]interface{}{
				"filetype": "csv",
				"file":     "string instead of []byte",
			},
			false,
			nil,
		},
		"invoke file parser for json": {
			s.logger,
			map[string]interface{}{
				"filetype": "json",
				"file":     s.helper.BytesFromFile("./testdata/json/example.json"),
			},
			true,
			map[string]interface{}{
				"key": "value",
			},
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := fileparser_action.Invoke(tc.Stub, tc.Input)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, output)
			} else {
				require.Error(t, err)
			}
		})
	}
}
