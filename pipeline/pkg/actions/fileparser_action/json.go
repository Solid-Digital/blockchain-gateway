package fileparser_action

import (
	"encoding/json"
	"github.com/unchainio/interfaces/logger"
)

func ParseJson(logger logger.Logger, rawJson []byte) (map[string]interface{}, error){
	var res map[string]interface{}

	err := json.Unmarshal(rawJson, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
