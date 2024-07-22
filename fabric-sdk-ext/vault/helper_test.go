/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault_test

import (
	"encoding/base64"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	vaultlib "github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"

	"crypto/ecdsa"
	"crypto/x509"

	"fmt"

	"crypto/rsa"

	"crypto"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/pkg/errors"
	"github.com/unchainio/fabric-sdk-ext/vault"
)

func testVerificationFlow(tb testing.TB, csp *vault.CryptoSuite, ski []byte, digest []byte) {
	key, err := csp.GetKey(ski)
	assert.NoError(tb, err)

	// Sign via vault
	signature, err := csp.Sign(key, digest, nil)
	assert.NoError(tb, err)

	// Verify signature via vault
	valid, err := csp.Verify(key, signature, digest, nil)
	assert.NoError(tb, err)

	if !valid {
		tb.Fatalf("Signature verification failed.")
	}

	// Verify signature without vault
	pubKey := parsePublicKey(tb, key)

	valid, err = verify(tb, pubKey, signature, digest)

	assert.NoError(tb, err)

	if !valid {
		tb.Fatalf("Signature verification failed.")
	}
}

func verify(tb testing.TB, pubKey interface{}, signature []byte, digest []byte) (valid bool, err error) {
	switch pubKey := pubKey.(type) {
	case *ecdsa.PublicKey:
		return verifyECDSA(pubKey, signature, digest)

	case *rsa.PublicKey:
		return verifyRSA(pubKey, signature, digest)
	}

	return false, errors.New("unsupported key type")
}

func parsePublicKey(tb testing.TB, key core.Key) interface{} {
	pub, err := key.PublicKey()
	assert.NoError(tb, err)

	pubBytes, err := pub.Bytes()
	assert.NoError(tb, err)

	pubKey, err := x509.ParsePKIXPublicKey(pubBytes)
	assert.NoError(tb, err)

	return pubKey
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

func testVaultCryptoSuite(tb testing.TB) (*vault.CryptoSuite, func()) {
	client, closer := testVaultServer(tb)

	csp, err := vault.NewCryptoSuite(vault.WithClient(client))

	if err != nil {
		tb.Fatalf("%+v", err)
	}

	return csp, func() {
		defer closer()
	}
}

// testVaultServer creates a test vault cluster and returns a configured API
// client and closer function.
func testVaultServer(t testing.TB) (*api.Client, func()) {
	t.Helper()

	client, _, closer := testVaultServerUnseal(t)

	err := client.Sys().Mount("fabric/transit", &api.MountInput{
		Type: "transit",
	})

	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = client.Sys().Mount("fabric/kv", &api.MountInput{
		Type: "kv",
	})

	if err != nil {
		t.Fatalf("%+v", err)
	}

	return client, closer
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(t testing.TB) (*api.Client, []string, func()) {
	t.Helper()

	return testVaultServerCoreConfig(t, &vaultlib.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"database":       database.Factory,
			"generic-leased": vaultlib.LeasedPassthroughBackendFactory,
			"pki":            pki.Factory,
			"transit":        transit.Factory,
		},
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(t testing.TB, coreConfig *vaultlib.CoreConfig) (*api.Client, []string, func()) {
	t.Helper()

	cluster := vaultlib.NewTestCluster(t, coreConfig, &vaultlib.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vaultlib.TestWaitActive(t, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Convert the unseal keys to base64 encoded, since these are how the user
	// will get them.
	unsealKeys := make([]string, len(cluster.BarrierKeys))
	for i := range unsealKeys {
		unsealKeys[i] = base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])
	}

	return client, unsealKeys, func() { defer cluster.Cleanup() }
}
