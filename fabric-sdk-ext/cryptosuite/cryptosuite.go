package cryptosuite

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric/bccsp"
)

//GetECDSAP256KeyGenOpts returns options for ECDSA key generation with curve P-256.
func GetECDSAP256KeyGenOpts(ephemeral bool) core.KeyGenOpts {
	return &bccsp.ECDSAP256KeyGenOpts{Temporary: ephemeral}
}

//GetRSA2048KeyGenOpts returns options for 2048 bit RSA key generation.
func GetRSA2048KeyGenOpts(ephemeral bool) core.KeyGenOpts {
	return &bccsp.RSA2048KeyGenOpts{Temporary: ephemeral}
}

//GetRSA4096KeyGenOpts returns options for 4096 bit RSA key generation.
func GetRSA4096KeyGenOpts(ephemeral bool) core.KeyGenOpts {
	return &bccsp.RSA4096KeyGenOpts{Temporary: ephemeral}
}
