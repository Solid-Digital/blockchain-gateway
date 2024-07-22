package testhelper

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/stretchr/testify/suite"
)

type Helper struct {
	suite *suite.Suite
	ares  *ares.Server
}

func NewHelper(s *suite.Suite, ares *ares.Server) *Helper {
	return &Helper{suite: s, ares: ares}
}
