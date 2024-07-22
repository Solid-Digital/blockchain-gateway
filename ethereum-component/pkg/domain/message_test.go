package domain_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
)

func (s *TestSuite) TestNewMessage() {
	input1 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"to":       NonExistingContractAddress,
		"function": "set",
		"params":   map[string]interface{}{"_age": 12, "_name": "john"},
	}
	expected1 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       NonExistingContractAddress,
		Function: "set",
		Params:   map[string]interface{}{"_age": 12, "_name": "john"},
	}

	input2 := map[string]interface{}{
		"type":              domain.DeployContract,
		"from":              DefaultAccount,
		"solidity":          "foo",
		"constructorParams": map[string]map[string]interface{}{"myContract": {"_age": 12, "_name": "john"}},
	}
	expected2 := &domain.MessageDeployContract{
		From:              DefaultAccount,
		Solidity:          "foo",
		ConstructorParams: map[string]map[string]interface{}{"myContract": {"_age": 12, "_name": "john"}},
	}

	input3 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"to":       NonExistingContractAddress,
		"function": "set",
		"params":   map[string]interface{}{"_age": 12, "_name": "john"},
	}

	input4 := map[string]interface{}{
		"type":              domain.DeployContract,
		"solidity":          "foo",
		"constructorParams": map[string]map[string]interface{}{"myContract": {"_age": 12, "_name": "john"}},
	}

	input5 := map[string]interface{}{
		"type": "invalid",
	}

	input6 := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     DefaultAccount,
		"to":       NonExistingContractAddress,
		"function": "set",
		"params":   map[string]interface{}{"_age": 12, "_name": "john"},
		"nonce":    123456,
	}

	expected6 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       NonExistingContractAddress,
		Function: "set",
		Params:   map[string]interface{}{"_age": 12, "_name": "john"},
		Nonce:    uint64(123456),
	}

	cases := map[string]struct {
		Input           map[string]interface{}
		ExpectedMessage interface{}
		Success         bool
	}{
		"valid callContractFunction": {
			input1,
			expected1,
			true,
		},
		"valid deployContract": {
			input2,
			expected2,
			true,
		},
		"invalid callContractFunction - from field missing": {
			input3,
			nil,
			false,
		},
		"invalid deployContract - from field missing": {
			input4,
			nil,
			false,
		},
		"invalid type": {
			input5,
			nil,
			false,
		},
		"with nonce": {
			input6,
			expected6,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			msg, err := domain.NewMessage(tc.Input)

			if tc.Success {
				s.Require().NoError(err)
				s.Require().Equal(tc.ExpectedMessage, msg)
			} else {
				s.Require().Error(err)
				s.Require().Equal(tc.ExpectedMessage, msg)
			}
		})
	}
}
