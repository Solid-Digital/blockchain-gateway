package account

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/unchainio/interfaces/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/unchainio/pkg/errors"
)

type Account struct {
	logger     logger.Logger
	cfg        *Config
	Address    *common.Address
	privateKey *ecdsa.PrivateKey
}

func newAccount(logger logger.Logger, cfg *Config) (*Account, error) {
	address := common.HexToAddress(cfg.Address)

	// For some reason go-ethereum expects addresses to be prefixed with '0x' and
	// private key's not. Usually both of them are prefixed, so we cut the '0x' from
	// the private key to make it work in both cases.
	privateKey, err := crypto.HexToECDSA(cut0x(cfg.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}

	return &Account{logger: logger, cfg: cfg, Address: &address, privateKey: privateKey}, nil
}

// GetSignerFn returns a function that can sign a transaction with the private key of the Account
func (a Account) GetSignerFn() (bind.SignerFn, error) {
	signerFn := func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != *a.Address {
			return nil, errors.New("not authorized to sign this account")
		}
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a.privateKey)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}

	return signerFn, nil
}

// cut0x removes leading '0x' from string
func cut0x(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s[2:]
	}

	return s
}
