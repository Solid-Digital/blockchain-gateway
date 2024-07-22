package http_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/api/operations/hello"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/http"
	"github.com/stretchr/testify/suite"
)

type HelloHandlerTestSuite struct {
	suite.Suite
	factory *factory.Factory
	handler *http.HelloHandler
}

// This runs just once before all tests
func (s *HelloHandlerTestSuite) SetupSuite() {
	s.T().Skip()
	s.factory = factory.NewFactory(&s.Suite)
	s.handler = s.factory.HelloHandler()
}

func (s *HelloHandlerTestSuite) TestHelloHandler_Hello() {
	response := s.handler.Hello(hello.NewHelloParams())
	helloOK, ok := response.(*hello.HelloOK)

	s.Require().True(ok)
	s.Require().Contains(helloOK.ContentType, "text/plain")
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestHelloHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HelloHandlerTestSuite))
}
