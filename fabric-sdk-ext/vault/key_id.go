/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"encoding/hex"

	"github.com/unchainio/fabric-sdk-ext/vault/internal"
)

func (csp *CryptoSuite) storeKeyID(ski []byte, keyID string) error {
	_, err := csp.client.Logical().Write(
		"fabric/kv/ski/"+hex.EncodeToString(ski),

		map[string]interface{}{
			"value": keyID,
		},
	)

	return err
}

func (csp *CryptoSuite) loadKeyID(ski []byte) (string, error) {
	secret, err := csp.client.Logical().Read("fabric/kv/ski/" + hex.EncodeToString(ski))

	if err != nil {
		return "", err
	}

	return internal.NewSecretWrapper(secret).ParseValue(), nil
}
