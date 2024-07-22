/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/pkg/errors"
	"github.com/unchainio/fabric-sdk-ext/vault/internal"
)

// KeyGen generates a key using opts.
func (csp *CryptoSuite) KeyGen(opts core.KeyGenOpts) (core.Key, error) {
	var err error

	// Validate arguments
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}

	if opts.Ephemeral() {
		return nil, errors.New("vault does not support ephemeral keys")
	}

	keyType, err := parseKeyType(opts.Algorithm())

	if err != nil {
		return nil, err
	}

	keyID, err := csp.keyGen(keyType)

	if err != nil {
		return nil, err
	}

	sw, err := csp.getKey(keyID)

	if err != nil {
		return nil, err
	}

	key, err := sw.ParseKey()

	if err != nil {
		return nil, err
	}

	err = csp.storeKeyID(key.SKI(), sw.KeyID())

	if err != nil {
		return nil, err
	}

	return key, nil
}

func parseKeyType(algorithm string) (string, error) {
	switch algorithm {
	case bccsp.ECDSAP256:
		return ECDSAP256, nil
	case bccsp.RSA2048:
		return RSA2048, nil
	case bccsp.RSA4096:
		return RSA4096, nil

	default:
		return "", errors.Errorf("the algorithm %s is not supported.", algorithm)
	}
}

func (csp *CryptoSuite) keyGen(keyType string) (string, error) {
	var err error

	keyID := internal.RandomString(24)

	_, err = csp.client.Logical().Write(
		"fabric/transit/keys/"+keyID,

		map[string]interface{}{
			"type": keyType,
		},
	)

	if err != nil {
		return "", errors.Wrapf(err, "failed to generate a key of type %s", keyType)
	}

	return keyID, nil
}
