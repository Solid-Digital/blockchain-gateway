package helper

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func (h *Helper) AddressToString(address *common.Address) string {
	return strings.ToLower(address.String())
}
