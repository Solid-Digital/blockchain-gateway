/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"encoding/base64"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/pkg/errors"
	"github.com/unchainio/fabric-sdk-ext/vault/internal"
)

// Sign signs digest using key k.
// The opts argument should be appropriate for the algorithm used.
//
// Note that when a signature of a hash of a larger message is needed,
// the caller is responsible for hashing the larger message and passing
// the hash (as digest).
func (csp *CryptoSuite) Sign(k core.Key, digest []byte, opts core.SignerOpts) (signature []byte, err error) {
	keyID, err := csp.loadKeyID(k.SKI())

	if err != nil {
		return nil, err
	}

	secret, err := csp.client.Logical().Write(
		"fabric/transit/sign/"+keyID,

		map[string]interface{}{
			"input":     base64.StdEncoding.EncodeToString(digest),
			"prehashed": true,
		},
	)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to sign the digest")
	}

	return internal.NewSecretWrapper(secret).ParseSignature(k)
}
