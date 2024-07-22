package helper

import "io/ioutil"

func (h *Helper) BytesFromFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)

	h.suite.Require().NoError(err)

	return bytes
}

func (h *Helper) StringFromFile(path string) string {
	bytes := h.BytesFromFile(path)

	return string(bytes)
}
