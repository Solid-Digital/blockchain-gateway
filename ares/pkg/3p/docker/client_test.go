package docker_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/docker/docker/api/types"

	"github.com/unchainio/pkg/xlogger"

	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/xconfig"

	"bitbucket.org/unchain/ares/pkg/3p/docker"
)

var defaultConfigPath = "../../../config/test/config.toml"

func TestImagePush(t *testing.T) {
	testhelper.SkipInBitbucket(t)

	cfg := getConfig(t)
	d, cleanup := getClient(t, cfg)
	defer cleanup()

	getImages(t, d)

	cases := map[string]struct {
		ref     string
		success bool
	}{
		"success": {
			"localhost:5000/janus-v2:latest",
			true,
		},
		"fail with wrong image name": {
			"localhost:5000/janus-v2z:latest",
			false,
		},
		"fail with unavailable registry": {
			"localhostz:5000/janus-v2:latest",
			false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := d.ImagePush(context.Background(), tc.ref, types.ImagePushOptions{
				RegistryAuth: "asdf",
			})

			fmt.Printf("err: %+v\n", err)
			if tc.success {
				xrequire.NoError(t, err)

				b, err := ioutil.ReadAll(out)
				xrequire.NoError(t, err)
				fmt.Printf("b: %s\n", string(b))

			} else {
				xrequire.Error(t, err)
			}
		})
	}
}

func TestImagePull(t *testing.T) {
	testhelper.SkipInBitbucket(t)

	cfg := getConfig(t)
	d, cleanup := getClient(t, cfg)
	defer cleanup()

	getImages(t, d)

	cases := map[string]struct {
		ref     string
		success bool
	}{
		"success": {
			"localhost:5000/janus-v2:latest",
			true,
		},
		"fail with wrong image name": {
			"localhost:5000/janus-v2z:latest",
			false,
		},
		"fail with unavailable registry": {
			"localhostz:5000/janus-v2:latest",
			false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := d.ImagePull(context.Background(), tc.ref, types.ImagePullOptions{
				RegistryAuth: "asdf",
			})

			fmt.Printf("err: %+v\n", err)
			if tc.success {
				xrequire.NoError(t, err)

				b, err := ioutil.ReadAll(out)
				xrequire.NoError(t, err)
				fmt.Printf("b: %s\n", string(b))

			} else {
				xrequire.Error(t, err)
			}
		})
	}
}

func getClient(t *testing.T, cfg *docker.Config) (*docker.Client, func()) {
	c, cleanup, err := docker.NewClient(xlogger.NewSimpleLogger(), cfg)
	xrequire.NoError(t, err)
	require.NotNil(t, c)

	return c, cleanup
}

func getConfig(t *testing.T) *docker.Config {
	cfg := &struct {
		Docker *docker.Config
	}{}

	err := xconfig.Load(cfg, xconfig.FromPaths(defaultConfigPath))
	xrequire.NoError(t, err)
	require.NotNil(t, cfg.Docker)

	return cfg.Docker
}

func getImages(t *testing.T, d *docker.Client) {
	_, err := d.ImagePull(context.Background(), "registry.unchain.io/unchainio/janus-v2", types.ImagePullOptions{})
	xrequire.NoError(t, err)
	err = d.ImageTag(context.Background(), "registry.unchain.io/unchainio/janus-v2", "localhost:5000/janus-v2:latest")
	xrequire.NoError(t, err)
	err = d.ImageTag(context.Background(), "registry.unchain.io/unchainio/janus-v2", "localhostz:5000/janus-v2:latest")
	xrequire.NoError(t, err)
}
