package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/unchainio/pkg/xtls"
)

type Config struct {
	Host    string
	Version string
	TLS     *xtls.KeyPairConfig
	Auth    *types.AuthConfig
}
