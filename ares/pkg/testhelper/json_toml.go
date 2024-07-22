package testhelper

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

func (h *Helper) EqualJSON(a, b interface{}) {
	eq, err := EqualJSON(a, b)

	h.suite.Require().NoError(err)
	h.suite.Require().True(eq)
}

func EqualJSON(a, b interface{}) (bool, error) {
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	bBytes, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return string(aBytes) == string(bBytes), nil
}

func (h *Helper) UnmarshalJSON(jsonString []byte) map[string]interface{} {
	res := make(map[string]interface{})
	err := json.Unmarshal(jsonString, &res)
	h.suite.Require().NoError(err)

	return res
}

func (h *Helper) UnmarshalTOML(tomlString string) map[string]interface{} {
	res := make(map[string]interface{})

	err := toml.Unmarshal([]byte(tomlString), &res)
	h.suite.Require().NoError(err)

	return res
}

func RandomTOML() string {
	return fmt.Sprintf(`foo = "%s"`, Randumb("bar"))
}

func (h *Helper) RandomTOML() string {
	return RandomTOML()
}

func RandomJSON() string {
	return fmt.Sprintf(`{"foo": "%s"}`, Randumb("bar"))
}

func (h *Helper) RandomJSON() string {
	return RandomJSON()
}
