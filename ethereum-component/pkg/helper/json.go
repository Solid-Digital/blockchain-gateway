package helper

import "encoding/json"

func (h *Helper) ToJSON(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)

	h.suite.Require().NoError(err)

	return bytes
}
