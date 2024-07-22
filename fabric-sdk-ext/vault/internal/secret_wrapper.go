/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package internal

import (
	"crypto/ecdsa"
	"strings"

	"encoding/base64"

	"crypto/rsa"

	vault "github.com/hashicorp/vault/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric/bccsp/utils"
	"github.com/pkg/errors"
)

// SecretWrapper is a wrapper around Vault secrets, used for parsing various items of interest
type SecretWrapper struct {
	secret *vault.Secret
}

// NewSecretWrapper constructs a new SecretWrapper out of a vault secret
func NewSecretWrapper(secret *vault.Secret) *SecretWrapper {
	return &SecretWrapper{
		secret: secret,
	}
}

// ParseKey parses a core.Key out of a vault secret's Data field
func (sw *SecretWrapper) ParseKey() (core.Key, error) {
	pubBytes := sw.PublicKey()
	return ParseKey(pubBytes)
}

// PublicKey parses a public key out of a vault secret's Data field
func (sw *SecretWrapper) PublicKey() []byte {
	keys, _ := sw.secret.Data["keys"].(map[string]interface{})
	onlyKey, _ := keys["1"].(map[string]interface{})
	pub, _ := onlyKey["public_key"].(string)

	return []byte(pub)
}

// KeyID parses a key's ID out of a vault secret's Data field
func (sw *SecretWrapper) KeyID() string {
	keyID, _ := sw.secret.Data["name"].(string)
	return keyID
}

// ParseValue parses the value out of a vault secret's Data field
func (sw *SecretWrapper) ParseValue() string {
	// TODO add errors?
	if sw.secret == nil {
		return ""
	}

	stringValue, ok := sw.secret.Data["value"]

	if !ok {
		return ""
	}

	value, _ := stringValue.(string)

	return value
}

// ParseSignature parses the signature out of a vault secret's Data field
func (sw *SecretWrapper) ParseSignature(key core.Key) ([]byte, error) {
	signature, _ := sw.secret.Data["signature"].(string)
	signature = strings.TrimPrefix(signature, "vault:v1:")

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse vault signature")
	}

	switch typedKey := key.(*privateKey).pub.pub.(type) {
	case *ecdsa.PublicKey:
		return ecdsaToLowS(typedKey, signatureBytes)
	case *rsa.PublicKey:
		return signatureBytes, nil

	default:
		return nil, errors.New("unsupported key type")
	}
}

func ecdsaToLowS(key *ecdsa.PublicKey, signatureBytes []byte) ([]byte, error) {
	r, s, err := utils.UnmarshalECDSASignature(signatureBytes)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal signature bytes")
	}

	s, _, err = utils.ToLowS(key, s)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert signature to low s")
	}

	return utils.MarshalECDSASignature(r, s)
}

// ParseVerification parses the verification out of a vault secret's Data field
func (sw *SecretWrapper) ParseVerification() bool {
	valid, _ := sw.secret.Data["valid"].(bool)
	return valid
}
