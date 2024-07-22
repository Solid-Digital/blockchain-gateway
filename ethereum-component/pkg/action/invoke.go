package action

import (
	"encoding/json"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const (
	ContractCallOutput        = "contractCallOutput"
	ContractTransactionOutput = "contractTransactionOutput"
	DeployContractOutput      = "deployContractOutput"
)

func (a *Action) Invoke(input map[string]interface{}) (map[string]interface{}, error) {
	msg, err := domain.NewMessage(input)
	if err != nil {
		return nil, err
	}

	var response interface{}
	switch msg.(type) {
	case *domain.MessageCallContractFunction:
		response, err = a.client.CallContractFunction(msg.(*domain.MessageCallContractFunction))
	case *domain.MessageDeployContract:
		response, err = a.client.DeployContracts(msg.(*domain.MessageDeployContract))
	}

	if err != nil {
		return nil, err
	}

	output := map[string]interface{}{}
	switch response.(type) {
	case *types.Transaction:
		output[ContractTransactionOutput], err = cleanOutput(response)
	case domain.DeployContractResponse:
		output[DeployContractOutput], err = cleanOutput(response)
	default:
		//Must be a contract call response, which is interface{} or []interface{}
		output[ContractCallOutput] = response
	}

	if err != nil {
		return nil, err
	}
	return output, nil
}

func cleanOutput(response interface{}) (map[string]interface{}, error) {
	var ret map[string]interface{}

	bytes, err := json.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed converting response to bytes")
	}

	err = json.Unmarshal(bytes, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "failed converting response to map")
	}

	return ret, nil
}
