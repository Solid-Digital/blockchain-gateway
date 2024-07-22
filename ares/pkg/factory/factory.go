package factory

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null"
)

type Factory struct {
	ares  *ares.Server
	suite *suite.Suite
}

func NewFactory(s *suite.Suite) *Factory {
	return &Factory{suite: s}
}

func (f *Factory) SetAres(ares *ares.Server) {
	f.ares = ares
}

func (f *Factory) IOSchemaJSON() null.JSON {
	return null.JSONFrom([]byte("[]"))
}
