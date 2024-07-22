package docker_test

import (
	"fmt"
	"io"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/unchainio/pkg/xlogger"

	"bitbucket.org/unchain/ares/pkg/3p/docker"
)

func TestPrepareDockerImage(t *testing.T) {
	testhelper.SkipInBitbucket(t)

	cfg := getConfig(t)
	d, cleanup := getClient(t, cfg)
	defer cleanup()
	s := docker.NewContainerService(d, xlogger.NewSimpleLogger(), cfg)
	getImages(t, d)

	cases := map[string]struct {
		ref       string
		baseRef   string
		artifacts []io.Reader
		success   bool
	}{
		"success": {
			"localhost:5000/janus-v2:prepared",
			"registry.unchain.io/unchainio/janus-v2",
			nil,
			true,
		},
		"fail with unavailable registry for the resulting image": {
			"localhostz:5000/janus-v2:prepared",
			"registry.unchain.io/unchainio/janus-v2",
			nil,
			false,
		},
		"fail with non-existent image for the base image": {
			"localhost:5000/janus-v2:prepared",
			"localhost:5000/janus-v2:non-existent",
			nil,
			false,
		},
		"fail with unavailable registry for the base image": {
			"localhost:5000/janus-v2:prepared",
			"localhostz:5000/janus-v2",
			nil,
			false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := s.PrepareImage(tc.ref, tc.baseRef, "", tc.artifacts...)

			fmt.Printf("err: %+v\n", err)
			if tc.success {
				xrequire.NoError(t, err)

			} else {
				xrequire.Error(t, err)
			}
		})
	}
}
