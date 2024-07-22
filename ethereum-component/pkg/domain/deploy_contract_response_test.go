package domain_test

import (
	"strings"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/ethereum/go-ethereum/common"
)

func (s *TestSuite) TestDeployContractResponse_AddContract() {
	hexAddress1 := "0x5d7de81c4a2009e60d4fea85fe88fdc659ae5cad"
	hexAddress2 := "0xc67206229c1f20ce80f3e13fbc9199e329dea1e1"

	address1 := common.HexToAddress(hexAddress1)
	address2 := common.HexToAddress(strings.ToUpper(hexAddress2))

	d := make(domain.DeployContractResponse)
	d.AddContract(&address1, nil, "")
	d.AddContract(&address2, nil, "")

	s.Suite.Require().Equal(2, len(d))

	_, ok := d[hexAddress1]

	s.Suite.Require().True(ok)

	_, ok = d[hexAddress2]

	s.Suite.Require().True(ok) // this verifies conversion to lower case
}
