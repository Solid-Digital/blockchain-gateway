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

// KeyImport imports a key from its raw representation using opts.
// The opts argument should be appropriate for the primitive used.
func (csp *CryptoSuite) KeyImport(raw interface{}, opts core.KeyImportOpts) (k core.Key, err error) {
	if opts.Ephemeral() {
		ski, err := internal.GetSKI(raw)

		if err != nil {
			return nil, err
		}

		return internal.NewPrivateKey(ski, raw), nil
	}

	return nil, errors.New("only ephemeral imports allowed")
}
