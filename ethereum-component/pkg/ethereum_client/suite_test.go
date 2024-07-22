package ethereum_client_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/helper"

	"bitbucket.org/unchain/ethereum2/pkg/ethereum_client"
	"bitbucket.org/unchain/ethereum2/pkg/factory"
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
)

// When Ganache-cli starts we make sure that the DefaultAccount has a balance. Also, we control
// the private key of this account so we can use it.
const (
	DefaultAccount       = "0xfdfa8d41f986c80904bf4825402e788f3121e7af"
	NonRegisteredAccount = "0xfac399e49f5b6867af186390270af252e683b154"

	NonExistingContractAddress = "0xd0a6e6c54dbc68db5db3a091b171a77407ff7ccf"
)

type TestSuite struct {
	suite.Suite
	logger                   logger.Logger
	helper                   *helper.Helper
	factory                  *factory.Factory
	client                   *ethereum_client.Client
	blockMiningNetworkClient *ethereum_client.Client
	syncModeClient           *ethereum_client.Client
}

// This will setup an integration test suite!
func (s *TestSuite) SetupSuite() {
	s.logger = factory.DefaultLogger(&s.Suite)
	s.helper = helper.NewHelper(&s.Suite, s.logger)
	s.factory = factory.NewFactory(&s.Suite, s.logger, s.helper)

	s.client = s.factory.EthereumClient(s.helper.BytesFromFile("./testdata/config/test.toml"))
	s.blockMiningNetworkClient = s.factory.EthereumClient(s.helper.BytesFromFile("./testdata/config/block_mining_network.toml"))
	s.syncModeClient = s.factory.EthereumClient(s.helper.BytesFromFile("./testdata/config/sync_mode_client.toml"))
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
