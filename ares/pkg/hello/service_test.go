package hello_test

import (
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/hello"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type HelloServiceTestSuite struct {
	suite.Suite
	factory *factory.Factory
	service *hello.Service
}

// This runs just once before all tests
func (s *HelloServiceTestSuite) SetupSuite() {
	s.factory = factory.NewFactory(&s.Suite)
	s.service = s.factory.HelloService()
}

func (s *HelloServiceTestSuite) TestService_Hello() {
	result := s.service.Hello()

	s.Require().True(strings.Contains(result, "version: 0.1"))
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestHelloServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HelloServiceTestSuite))
}
