package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/component"
	"github.com/stretchr/testify/suite"
)

type HelperTestSuite struct {
	suite.Suite
}

func (s *HelperTestSuite) TestActionVersionFileName() {
	fileName := component.ActionVersionFileName("myAction", "0.1", "myOrg")

	s.Require().Equal("myAction.action.0.1.myOrg.so", fileName)
}

func (s *HelperTestSuite) TestActionVersionFileID() {
	fileID := component.ActionVersionFileID("myAction", "0.1", "myOrg", "myFile")

	s.Require().Equal("myOrg/action/myAction/0.1/myFile.tar.gz", fileID)
}

func (s *HelperTestSuite) TestTriggerVersionFileName() {
	fileName := component.TriggerVersionFileName("myTrigger", "0.1", "myOrg")

	s.Require().Equal("myTrigger.trigger.0.1.myOrg.so", fileName)
}

func (s *HelperTestSuite) TestTriggerVersionFileID() {
	fileID := component.TriggerVersionFileID("myTrigger", "0.1", "myOrg", "myFile")

	s.Require().Equal("myOrg/trigger/myTrigger/0.1/myFile.tar.gz", fileID)
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestHelperTestSuite(t *testing.T) {
	suite.Run(t, new(HelperTestSuite))
}
