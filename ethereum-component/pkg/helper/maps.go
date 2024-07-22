package helper

import (
	"github.com/mitchellh/mapstructure"
	"github.com/unchainio/pkg/errors"
)

func (h *Helper) KeyFromMap(i interface{}) string {
	m := map[string]interface{}{}
	err := mapstructure.Decode(i, &m)

	h.suite.Require().NoError(err)

	for key, _ := range m {
		return key
	}

	h.suite.Require().NoError(errors.New("empty map"))

	return ""
}
