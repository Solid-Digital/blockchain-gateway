package helper

func (h *Helper) StringPtr(s string) *string {
	return &s
}

func (h *Helper) Uint8Ptr(i uint8) *uint8 {
	return &i
}
