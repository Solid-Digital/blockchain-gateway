package testhelper

import (
	"os"
	"testing"
)

func InBitBucket() bool {
	return os.Getenv("BITBUCKET_BUILD_NUMBER") != ""
}

func (h *Helper) InBitBucket() bool {
	return InBitBucket()
}

func SkipInBitbucket(t *testing.T) {
	if InBitBucket() {
		t.Skip("bitbucket pipelines suck")
	}
}

func (h *Helper) SkipInBitbucket() {
	SkipInBitbucket(h.suite.T())
}
