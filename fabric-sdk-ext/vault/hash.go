/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"hash"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/pkg/errors"
)

// Hasher is a BCCSP-like interface that provides HashFunc algorithms
type Hasher interface {

	// Hash hashes messages msg using options opts.
	// If opts is nil, the default HashFunc function will be used.
	Hash(msg []byte, opts core.HashOpts) (hash []byte, err error)

	// GetHash returns and instance of HashFunc.Hash using options opts.
	// If opts is nil, the default HashFunc function will be returned.
	GetHash(opts core.HashOpts) (h hash.Hash, err error)
}

// Hash hashes messages msg using options opts.
func (csp *CryptoSuite) Hash(msg []byte, opts core.HashOpts) (digest []byte, err error) {
	// Validate arguments
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}

	hasher, found := csp.hashers[opts.Algorithm()]
	if !found {
		return nil, errors.Errorf("unsupported 'HashOpt' provided [%v]", opts)
	}

	digest, err = hasher.Hash(msg, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed hashing with opts [%v]", opts)
	}

	return
}

// GetHash returns and instance of hash.Hash using options opts.
// If opts is nil then the default hash function is returned.
func (csp *CryptoSuite) GetHash(opts core.HashOpts) (h hash.Hash, err error) {
	// Validate arguments
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}

	hasher, found := csp.hashers[opts.Algorithm()]
	if !found {
		return nil, errors.Errorf("unsupported 'HashOpt' provided [%v]", opts)
	}

	h, err = hasher.GetHash(opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting hash function with opts [%v]", opts)
	}

	return
}
