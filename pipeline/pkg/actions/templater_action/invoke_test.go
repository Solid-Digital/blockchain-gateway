package templater_action_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/templater_action"
	"github.com/unchain/pipeline/pkg/domain"
	"testing"
)

func (s *TestSuite) TestAction_Invoke() {
	templ1 := "This is a template with name {{.TEMPLATE_NAME}} and creation date {{.TEMPLATE_DATE}}"
	templ2 := "{\"key\":{{.VALUE}}}"
	templ3 := "I am invalid {{.TEMPLATE_NAME}} and creation date {{.TEMPLATE_DATE}"
	templ4 := "This is a template with {{.ROOT.LEAF}}"

	vars1 := map[string]interface{}{
		"TEMPLATE_NAME": "TEST TEMPLATE",
		"TEMPLATE_DATE": "20th of July, 2019",
	}
	output1 := `This is a template with name TEST TEMPLATE and creation date 20th of July, 2019`

	vars2 := map[string]interface{}{
		"TEMPLATE_NAME": "TEST TEMPLATE",
		"TEMPLATE_DATE": 2019,
	}
	output2 := `This is a template with name TEST TEMPLATE and creation date 2019`

	vars3 := map[string]interface{}{
		"TEMPLATE_NAME": true,
		"TEMPLATE_DATE": 2019,
	}
	output3 := `This is a template with name true and creation date 2019`

	vars4 := map[string]interface{}{
		"ROOT": map[string]interface{}{
			"LEAF": "nested input",
		},
	}
	output4 := `This is a template with nested input`

	vars5 := map[string]interface{}{
		"VALUE": "500",
	}
	output5 := `{"key":500}`

	cases := map[string]struct {
		Stub           domain.Stub
		Input          map[string]interface{}
		Success        bool
		ExpectedOutput map[string]interface{}
	}{
		"string input": {
			s.logger,
			map[string]interface{}{
				templater_action.InputTemplate:  templ1,
				templater_action.InputVariables: vars1,
			},
			true,
			map[string]interface{}{
				templater_action.TemplateResult: output1,
			},
		},
		"integer input": {
			s.logger,
			map[string]interface{}{
				templater_action.InputTemplate:  templ1,
				templater_action.InputVariables: vars2,
			},
			true,
			map[string]interface{}{
				templater_action.TemplateResult: output2,
			},
		},
		"boolean input": {
			s.logger,
			map[string]interface{}{
				templater_action.InputTemplate:  templ1,
				templater_action.InputVariables: vars3,
			},
			true,
			map[string]interface{}{
				templater_action.TemplateResult: output3,
			},
		},
		"nested input": {
			s.logger,
			map[string]interface{}{
				templater_action.InputTemplate:  templ4,
				templater_action.InputVariables: vars4,
			},
			true,
			map[string]interface{}{
				templater_action.TemplateResult: output4,
			},
		},
		"json input": {
			s.logger,
			map[string]interface{}{
				templater_action.InputTemplate:  templ2,
				templater_action.InputVariables: vars5,
			},
			true,
			map[string]interface{}{
				templater_action.TemplateResult: output5,
			},
		},
		"invalid template": {
			s.logger,
			map[string]interface{}{
				templater_action.InputTemplate:  templ3,
				templater_action.InputVariables: vars5,
			},
			false,
			map[string]interface{}{
				templater_action.TemplateResult: output5,
			},
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			res, err := templater_action.Invoke(tc.Stub, tc.Input)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tc.ExpectedOutput, res)

			} else {
				require.Error(t, err)
			}
		})
	}
}

func (s *TestSuite) TestAction_SprigFunctionsAvailable() {
	res, err := templater_action.Invoke(s.logger, map[string]interface{}{
		"template": `{{now | date "2006-01-02T15:04:05.000"}}`,
	})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	random, ok := res[templater_action.TemplateResult].(string)
	require.True(s.T(), ok, "random should be string")
	require.Equal(s.T(), 23, len(random), "random should be of length 23 ")
}

func (s *TestSuite) TestAction_SprigDateFunctionParsesSuccessful() {
	res, err := templater_action.Invoke(s.logger, map[string]interface{}{
		"template": `{{ toDate "20060102" .date | date "2006-01-02"}}`,
		"variables": map[string]interface{}{
			"date": "20200106",
		},
	})

	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	random, ok := res[templater_action.TemplateResult].(string)
	require.True(s.T(), ok, "random should be string")
	require.Equal(s.T(), "2020-01-06", random)
}
