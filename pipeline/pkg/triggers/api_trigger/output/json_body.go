package output

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func jsonBody(body []byte) (output map[string]interface{}, err error) {
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return output, nil
}
