package listener

import (
	"github.com/unchain/pipeline/pkg/domain"
)

func (s *Server) responseListener(req *domain.Request) *domain.Response {
	responseChannel := make(chan *domain.Response)
	s.ResponseChannelMap.Store(req.Tag, responseChannel)
	defer close(responseChannel)

	s.RequestChannel <- req

	return <-responseChannel
}
