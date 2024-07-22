package auth

import (
	"github.com/unchainio/pkg/xtls"
)

type Config struct {
	TLS             *xtls.KeyPairConfig
	ExpirationDelta int
	Issuer          string
	ConnectURL      string
}
