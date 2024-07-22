/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package internal

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"

	"encoding/pem"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/pkg/errors"
)

// ParseKey parses a core.Key out of a public key's bytes
func ParseKey(pubBytes []byte) (core.Key, error) {
	block, _ := pem.Decode(pubBytes)

	if block == nil {
		return nil, errors.Errorf("failed to decode the bytes of the vault secret")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, errors.Errorf("failed to parse the core.Key from the decoded vault secret")
	}

	ski, err := GetSKI(pub)

	if err != nil {
		return nil, err
	}

	return NewPrivateKey(ski, pub), nil
}

// GetSKI returns an SKI from a *rsa.PublicKey or *ecdsa.PublicKey
func GetSKI(pub interface{}) ([]byte, error) {
	switch tpub := pub.(type) {
	case *rsa.PublicKey:
		// Marshall the public key
		raw, err := asn1.Marshal(*tpub)

		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal key")
		}

		// Hash it
		hash := sha256.New()
		_, err = hash.Write(raw)

		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate HashFunc")
		}

		return hash.Sum(nil), nil

	case *ecdsa.PublicKey:
		// Marshall the public key
		raw := elliptic.Marshal(tpub.Curve, tpub.X, tpub.Y)

		// Hash it
		hash := sha256.New()
		_, err := hash.Write(raw)

		if err != nil {
			return nil, err
		}

		return hash.Sum(nil), nil
	case *x509.Certificate:
		return GetSKI(tpub.PublicKey)

	default:
		return nil, errors.New("unsupported key type")
	}
}
