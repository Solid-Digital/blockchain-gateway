package helper

import (
	"bitbucket.org/unchain/ethereum2/pkg/action"
	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/unchainio/pkg/errors"
)

func (h *Helper) CallContractFunction(a *action.Action, from, to, functionName string, params map[string]interface{}) {
	input := map[string]interface{}{
		"type":     domain.CallContractFunction,
		"from":     from,
		"to":       to,
		"function": functionName,
		"params":   params,
	}

	_, err := a.Invoke(input)

	h.suite.Require().NoError(err)
}

func (h *Helper) DeploySingleContractAction(a *action.Action, from, solidity string) (contractAddress, ABI string) {
	input := map[string]interface{}{
		"type":     domain.DeployContract,
		"from":     from,
		"solidity": solidity,
	}

	response, err := a.Invoke(input)

	h.suite.Require().NoError(err)

	output := response[action.DeployContractOutput]

	for address, deployedContract := range output.(map[string]interface{}) {
		return address, deployedContract.(map[string]interface{})["ABI"].(string)
	}

	h.suite.Require().NoError(errors.New("no contract deployed"))

	return "", ""
}
