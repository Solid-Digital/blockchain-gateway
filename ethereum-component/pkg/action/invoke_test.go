package action_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/domain"

	"bitbucket.org/unchain/ethereum2/pkg/action"

	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestAction_Invoke() {
	simpleContractSolidity := s.helper.StringFromFile("./testdata/contract/simple_contract.sol")
	withConstructorParamsSolidity := s.helper.StringFromFile("./testdata/contract/with_constructor_params.sol")
	doesNotCompileSolidity := s.helper.StringFromFile("./testdata/contract/does_not_compile.sol")
	deployContractAction := s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml"))
	actionCfgFmt := s.helper.StringFromFile("./testdata/config/action_fmt.toml")
	actionCfgSingleAccountFmt := s.helper.StringFromFile("./testdata/config/action_fmt_single_account.toml")
	actionCfgMultipleContractsFmt := s.helper.StringFromFile("./testdata/config/action_fmt_multiple_contracts.toml")

	address1, ABI1 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action1 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgFmt, address1, ABI1))
	input1 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"to":       address1,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 1},
	}

	address2, ABI2 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action2 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgFmt, address2, ABI2))
	input2 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"to":       address2,
		"function": "set",
		"params":   map[string]interface{}{"x": 5},
	}

	action3 := s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml"))
	input3 := map[string]interface{}{
		"type":              domain.DeployContract,
		"from":              DefaultAccount,
		"solidity":          withConstructorParamsSolidity,
		"constructorParams": map[string]map[string]interface{}{"SimpleStorage": {"_foo": 34}},
	}

	input4 := map[string]interface{}{
		"from":     DefaultAccount,
		"to":       address1,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	input5 := map[string]interface{}{
		"type":     "invalidType",
		"from":     DefaultAccount,
		"to":       address1,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	input6 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"to":       address1,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	input7 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"to":       address1,
		"function": "invalidFunction",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	input8 := map[string]interface{}{
		"type":     domain.DeployContract,
		"solidity": simpleContractSolidity,
	}

	input9 := map[string]interface{}{
		"type":     domain.DeployContract,
		"from":     DefaultAccount,
		"solidity": doesNotCompileSolidity,
	}

	address10, ABI10 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action10 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgSingleAccountFmt, address10, ABI10))
	input10 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"to":       address10,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	address11, ABI11 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action11 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgFmt, address11, ABI11))
	input11 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"to":       address11,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	address12, ABI12 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action12 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgSingleAccountFmt, address12, ABI12))
	input12 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	address13, ABI13 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action13 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgMultipleContractsFmt, address13, ABI13))
	input13 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	address14, ABI14 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action14 := s.factory.InitializedAction(s.helper.FillInActionConfig(actionCfgFmt, address14, ABI14))
	input14 := map[string]interface{}{
		"type":     1,
		"from":     DefaultAccount,
		"to":       address14,
		"function": "get",
		"params":   map[string]interface{}{"_myInt": 12},
	}

	cases := map[string]struct {
		Action  *action.Action
		Input   map[string]interface{}
		Success bool
	}{
		"call contract call function": {
			action1,
			input1,
			true,
		},
		"call contract transaction function": {
			action2,
			input2,
			true,
		},
		"deploy contract": {
			action3,
			input3,
			true,
		},
		"missing type": {
			s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml")),
			input4,
			false,
		},
		"invalid type": {
			s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml")),
			input5,
			false,
		},
		"from field missing for contract call": {
			s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml")),
			input6,
			false,
		},
		"invalid function": {
			s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml")),
			input7,
			false,
		},
		"from field missing for deploy contract": {
			s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml")),
			input8,
			false,
		},
		"solidity code does not compile": {
			s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/test.toml")),
			input9,
			false,
		},
		"call contract call function without from": {
			action10,
			input10,
			false,
		},
		"call contract call function without from multiple accounts": {
			action11,
			input11,
			false,
		},
		"call contract call function without to": {
			action12,
			input12,
			false,
		},
		"call contract call function without to multiple contracts": {
			action13,
			input13,
			false,
		},
		"type is not a string": {
			action14,
			input14,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := tc.Action.Invoke(tc.Input)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.NoError(t, err)

			} else {
				require.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
