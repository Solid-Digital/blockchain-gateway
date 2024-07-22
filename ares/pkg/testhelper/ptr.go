package testhelper

func (h *Helper) StringPtr(s string) *string {
	return StringPtr(s)
}

func StringPtr(s string) *string {
	return &s
}

func (h *Helper) BoolPtr(b bool) *bool {
	return BoolPtr(b)
}

func BoolPtr(b bool) *bool {
	return &b
}

func (h *Helper) Int64Ptr(i int64) *int64 {
	return Int64Ptr(i)
}

func Int64Ptr(i int64) *int64 {
	return &i
}
