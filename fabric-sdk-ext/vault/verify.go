/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/pkg/errors"
)

// Verify verifies signature against key k and digest
// The opts argument should be appropriate for the algorithm used.
func (csp *CryptoSuite) Verify(key core.Key, signature, digest []byte, opts core.SignerOpts) (valid bool, err error) {
	// Verify signature without vault
	pubKey, err := parsePublicKey(key)

	if err != nil {
		return false, err
	}

	switch pubKey := pubKey.(type) {
	case *ecdsa.PublicKey:
		valid, err = verifyECDSA(pubKey, signature, digest)

		if err != nil {
			return false, err
		}

	case *rsa.PublicKey:
		valid, err = verifyRSA(pubKey, signature, digest)

		if err != nil {
			return false, err
		}

	}

	if !valid {
		return false, errors.New("Signature verification failed.")
	}

	return true, nil
}

func parsePublicKey(key core.Key) (interface{}, error) {
	pub, err := key.PublicKey()

	if err != nil {
		return nil, err
	}

	pubBytes, err := pub.Bytes()

	if err != nil {
		return nil, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubBytes)

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return pubKey, nil
}

func verifyRSA(k *rsa.PublicKey, signature, digest []byte) (valid bool, err error) {
	err = rsa.VerifyPSS(k, crypto.SHA256, digest, signature, nil)

	if err != nil {
		return false, err
	}

	return true, nil
}

func verifyECDSA(k *ecdsa.PublicKey, signature, digest []byte) (valid bool, err error) {
	r, s, err := utils.UnmarshalECDSASignature(signature)
	if err != nil {
		return false, fmt.Errorf("Failed unmashalling signature [%s]", err)
	}

	lowS, err := utils.IsLowS(k, s)
	if err != nil {
		return false, err
	}

	if !lowS {
		return false, fmt.Errorf("Invalid S. Must be smaller than half the order [%s][%s].", s, utils.GetCurveHalfOrdersAt(k.Curve))
	}

	return ecdsa.Verify(k, digest, r, s), nil
}

//// Verify verifies signature against key k and digest
//// The opts argument should be appropriate for the algorithm used.
//func (csp *CryptoSuite) Verify(k core.Key, signature, digest []byte, opts core.SignerOpts) (valid bool, err error) {
//	spew.Dump(k)
//	keyID, err := csp.loadKeyID(k.SKI())
//
//	if err != nil {
//		return false, err
//	}
//
//	secret, err := csp.client.Logical().Write(
//		"fabric/transit/verify/"+keyID,
//
//		map[string]interface{}{
//			"input":     base64.StdEncoding.EncodeToString(digest),
//			"signature": "vault:v1:" + base64.StdEncoding.EncodeToString(signature),
//			"prehashed": true,
//		},
//	)
//
//	if err != nil {
//		return false, errors.Wrapf(err, "failed to verify the signature")
//	}
//
//	return internal.NewSecretWrapper(secret).ParseVerification(), nil
//}
