package testhelper

import "github.com/volatiletech/null"

func (h *Helper) ToInt64(i null.Int64) int64 {
	h.suite.Require().True(i.Valid)

	return i.Int64
}
