package hello

import (
	"fmt"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.HelloService), new(Service)))

type Service struct {
	meta *ares.Metadata
}

func (s *Service) Hello() string {
	return fmt.Sprintf("%s\n%s\n", ares.Logo, s.meta)
}

func NewService(meta *ares.Metadata) *Service {
	return &Service{meta: meta}
}
