/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/pkg/errors"
	"github.com/unchainio/fabric-sdk-ext/vault/internal"
)

// GetKey returns the key this CSP associates to
// the Subject Key Identifier ski.
func (csp *CryptoSuite) GetKey(ski []byte) (core.Key, error) {
	keyID, err := csp.loadKeyID(ski)

	if keyID == "" {
		return nil, errors.New("empty ski")
	}

	if err != nil {
		return nil, err
	}

	sw, err := csp.getKey(keyID)

	if err != nil {
		return nil, err
	}

	return sw.ParseKey()
}

func (csp *CryptoSuite) getKey(keyID string) (*internal.SecretWrapper, error) {
	secret, err := csp.client.Logical().Read("fabric/transit/keys/" + keyID)

	if err != nil {
		return nil, errors.Errorf("failed to find key with id `%s`", keyID)
	}

	return internal.NewSecretWrapper(secret), nil
}
