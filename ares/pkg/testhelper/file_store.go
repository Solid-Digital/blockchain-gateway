package testhelper

import (
	"io/ioutil"
)

func (h *Helper) FileExists(fileID string) bool {
	file, err := h.ares.FileStore.GetObject(fileID)

	h.suite.Require().NoError(err)

	_, err = ioutil.ReadAll(file)

	if err != nil {
		return false
	}

	return true
}
