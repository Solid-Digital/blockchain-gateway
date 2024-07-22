/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package internal

import (
	"crypto/x509"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/pkg/errors"
)

type privateKey struct {
	ski []byte
	pub publicKey
}

// NewPrivateKey returns a private key implementation of the core.Key interface
func NewPrivateKey(ski []byte, pub interface{}) core.Key {
	return &privateKey{
		ski: ski,
		pub: publicKey{
			ski: ski,
			pub: pub,
		},
	}
}

// Bytes converts this key to its byte representation,
// if this operation is allowed.
func (k *privateKey) Bytes() (raw []byte, err error) {
	return nil, errors.New("not supported")
}

// SKI returns the subject key identifier of this key.
func (k *privateKey) SKI() (ski []byte) {
	return k.ski
}

// Symmetric returns true if this key is a symmetric key,
// false if this key is asymmetric
func (k *privateKey) Symmetric() bool {
	return false
}

// Private returns true if this key is a private key,
// false otherwise.
func (k *privateKey) Private() bool {
	return true
}

// PublicKey returns the corresponding public key part of an asymmetric public/private key pair.
// This method returns an error in symmetric key schemes.
func (k *privateKey) PublicKey() (core.Key, error) {
	return &k.pub, nil
}

type publicKey struct {
	ski []byte
	pub interface{}
}

// Bytes converts this key to its byte representation,
// if this operation is allowed.
func (k *publicKey) Bytes() (raw []byte, err error) {
	var kpub interface{}

	switch k.pub.(type) {
	case *x509.Certificate:
		kpub = k.pub.(*x509.Certificate).PublicKey
	default:
		kpub = k.pub
	}

	raw, err = x509.MarshalPKIXPublicKey(kpub)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal key")
	}
	return raw, nil
}

// SKI returns the subject key identifier of this key.
func (k *publicKey) SKI() (ski []byte) {
	return k.ski
}

// Symmetric returns true if this key is a symmetric key,
// false if this key is asymmetric
func (k *publicKey) Symmetric() bool {
	return false
}

// Private returns true if this key is a private key,
// false otherwise.
func (k *publicKey) Private() bool {
	return false
}

// PublicKey returns the corresponding public key part of an asymmetric public/private key pair.
// This method returns an error in symmetric key schemes.
func (k *publicKey) PublicKey() (core.Key, error) {
	return k, nil
}
