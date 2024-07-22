/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault_test

import (
	"crypto/sha256"
	"testing"

	"reflect"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/stretchr/testify/assert"
	"github.com/unchainio/fabric-sdk-ext/cryptosuite"
	"github.com/unchainio/fabric-sdk-ext/vault"
)

func TestVaultCryptoSuiteImplementsInterface(t *testing.T) {
	cspi := reflect.TypeOf((*core.CryptoSuite)(nil)).Elem()

	ok := reflect.PtrTo(reflect.TypeOf(vault.CryptoSuite{})).Implements(cspi)

	if !ok {
		t.Fatalf("vault.CryptoSuite does not implement core.CryptoSuite")
	}
}

var digest = []byte("blabla")

func TestKeyGenECDSAP256(t *testing.T) {
	csp, closer := testVaultCryptoSuite(t)
	defer closer()

	key, err := csp.KeyGen(cryptosuite.GetECDSAP256KeyGenOpts(false))
	assert.NoError(t, err)

	testVerificationFlow(t, csp, key.SKI(), hash(digest))
}

func TestKeyGenRSA2048(t *testing.T) {
	csp, closer := testVaultCryptoSuite(t)
	defer closer()

	publicKeyBytes, err := csp.KeyGen(cryptosuite.GetRSA2048KeyGenOpts(false))
	assert.NoError(t, err)

	testVerificationFlow(t, csp, publicKeyBytes.SKI(), hash(digest))
}

func TestKeyGenRSA4096(t *testing.T) {
	csp, closer := testVaultCryptoSuite(t)
	defer closer()

	publicKeyBytes, err := csp.KeyGen(cryptosuite.GetRSA4096KeyGenOpts(false))
	assert.NoError(t, err)

	testVerificationFlow(t, csp, publicKeyBytes.SKI(), hash(digest))
}

func hash(digest []byte) []byte {
	hf := sha256.New()
	hf.Write(digest)

	return hf.Sum(nil)
}
