/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"crypto/sha256"
	"crypto/sha512"

	vault "github.com/hashicorp/vault/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric/bccsp"

	"github.com/pkg/errors"
	"github.com/unchainio/fabric-sdk-ext/vault/internal"
	"golang.org/x/crypto/sha3"
)

// Constants describing encryption algorithms
const (
	ECDSAP256 = "ecdsa-p256"
	RSA2048   = "rsa-2048"
	RSA4096   = "rsa-4096"
)

var logger = logging.NewLogger("fabsdk/core")

// CryptoSuite is a vault implementation of the core.CryptoSuite interface
type CryptoSuite struct {
	hashers map[string]Hasher
	client  *vault.Client
}

// options configure a new CryptoSuite. options are set by the OptionFunc values passed to NewCryptoSuite.
type options struct {
	client            *vault.Client
	hashers           map[string]Hasher
	address           string
	token             string
	securityAlgorithm string
	securityLevel     int
}

// OptionFunc configures how the CryptoSuite is set up.
type OptionFunc func(*options) error

// NewCryptoSuite constructs a new CryptoSuite, configured via provided OptionFuncs
func NewCryptoSuite(optFuncs ...OptionFunc) (*CryptoSuite, error) {
	var err error
	opts := &options{}

	for _, optFunc := range optFuncs {
		err = optFunc(opts)

		if err != nil {
			return nil, err
		}
	}

	if opts.client == nil {
		opts.client, err = getVaultClient(opts.address, opts.token)

		if err != nil {
			return nil, err
		}
	}

	hashers := getHashers(opts.securityAlgorithm, opts.securityLevel)

	for key, hasher := range opts.hashers {
		hashers[key] = hasher
	}

	logger.Debug("Initialized the vault CryptoSuite")

	return &CryptoSuite{
		client:  opts.client,
		hashers: hashers,
	}, nil
}

func getHashers(algorithm string, level int) map[string]Hasher {
	defaultHasher := parseHasher(algorithm, level)

	// Set the hashers
	hashers := make(map[string]Hasher)

	if defaultHasher != nil {
		hashers[bccsp.SHA] = defaultHasher
	}

	hashers[bccsp.SHA256] = &internal.Hasher{HashFunc: sha256.New}
	hashers[bccsp.SHA384] = &internal.Hasher{HashFunc: sha512.New384}
	hashers[bccsp.SHA3_256] = &internal.Hasher{HashFunc: sha3.New256}
	hashers[bccsp.SHA3_384] = &internal.Hasher{HashFunc: sha3.New384}

	return hashers
}

func parseHasher(algorithm string, level int) *internal.Hasher {
	switch {
	case algorithm == "SHA2" && level == 256:
		return &internal.Hasher{HashFunc: sha256.New}
	case algorithm == "SHA2" && level == 384:
		return &internal.Hasher{HashFunc: sha512.New384}
	case algorithm == "SHA3" && level == 256:
		return &internal.Hasher{HashFunc: sha3.New256}
	case algorithm == "SHA3" && level == 384:
		return &internal.Hasher{HashFunc: sha3.New384}

	default:
		return nil
	}
}

// WithClient allows to set the vault client of the CryptoSuite
func WithClient(client *vault.Client) OptionFunc {
	return func(o *options) error {
		o.client = client
		return nil
	}
}

// WithHashers allows to provide additional hashers to the CryptoSuite
func WithHashers(hashers map[string]Hasher) OptionFunc {
	return func(o *options) error {
		o.hashers = hashers
		return nil
	}
}

// WithAddress allows to specify the address of the vault server to be used by the CryptoSuite
func WithAddress(address string) OptionFunc {
	return func(o *options) error {
		o.address = address
		return nil
	}
}

// WithToken allows to specify the auth token of the vault server to be used by the CryptoSuite
func WithToken(token string) OptionFunc {
	return func(o *options) error {
		o.token = token
		return nil
	}
}

// WithSecurityAlgorithm allows to specify the security algorithm to be used by the CryptoSuite
func WithSecurityAlgorithm(securityAlgorithm string) OptionFunc {
	return func(o *options) error {
		o.securityAlgorithm = securityAlgorithm
		return nil
	}
}

// WithSecurityLevel allows to specify the security level to be used by the CryptoSuite
func WithSecurityLevel(securityLevel int) OptionFunc {
	return func(o *options) error {
		o.securityLevel = securityLevel
		return nil
	}
}

// FromConfig uses a core.CryptoSuiteConfig to configure the vault client of the CryptoSuite
func FromConfig(config core.CryptoSuiteConfig) OptionFunc {
	return func(o *options) error {
		//if config.SecurityProviderAddress() != "" {
		//	o.address = config.SecurityProviderAddress()
		//}
		//
		//if config.SecurityProviderToken() != "" {
		//	o.token = config.SecurityProviderToken()
		//}

		if config.SecurityAlgorithm() != "" {
			o.securityAlgorithm = config.SecurityAlgorithm()
		}

		if config.SecurityLevel() != 0 {
			o.securityLevel = config.SecurityLevel()
		}

		return nil
	}
}

func getVaultClient(address, token string) (*vault.Client, error) {
	vaultConfig := &vault.Config{
		Address: address,
	}

	client, err := vault.NewClient(vaultConfig)

	if err != nil {
		return nil, errors.Wrapf(err, "could not initialize Vault BCCSP for address: %s", vaultConfig.Address)
	}

	client.SetToken(token)

	return client, nil
}
