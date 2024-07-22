/*
Copyright Hyperledger and its contributors.

SPDX-License-Identifier: Apache-2.0
*/

package vault

import (
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
)

// VaultUserStore is a vault implementation of UserStore
type VaultUserStore struct {
	client *vault.Client
}

// NewVaultUserStore creates a new VaultUserStore instance
func NewVaultUserStore(address, token string) (*VaultUserStore, error) {
	vaultConfig := &vault.Config{
		Address: address,
	}

	client, err := vault.NewClient(vaultConfig)

	if err != nil {
		return nil, errors.Wrapf(err, "could not initialize Vault BCCSP for address: %s", vaultConfig.Address)
	}

	client.SetToken(token)

	return &VaultUserStore{client: client}, nil
}

// Store stores a user into store
func (s *VaultUserStore) Store(user *msp.UserData) error {
	_, err := s.client.Logical().Write(
		"fabric/kv/users/"+strings.ToLower(user.ID)+"@"+strings.ToLower(user.MSPID),

		map[string]interface{}{
			"value": string(user.EnrollmentCertificate),
		})

	if err != nil {
		return errors.Wrap(err, "failed to store user data into vault")
	}

	return nil
}

// Load loads a user from store
func (s *VaultUserStore) Load(id msp.IdentityIdentifier) (*msp.UserData, error) {
	secret, err := s.client.Logical().Read("fabric/kv/users/" + strings.ToLower(id.ID) + "@" + strings.ToLower(id.MSPID))

	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, msp.ErrUserNotFound
	}

	value, ok := secret.Data["value"]

	if !ok {
		return nil, msp.ErrUserNotFound
	}

	certString, ok := value.(string)

	if !ok {
		return nil, msp.ErrUserNotFound
	}

	certBytes := []byte(certString)

	if err != nil {
		return nil, errors.Wrap(err, "failed to hex decode cert bytes")
	}

	userData := msp.UserData{
		ID:    id.ID,
		MSPID: id.MSPID,
		EnrollmentCertificate: certBytes,
	}

	return &userData, nil
}
