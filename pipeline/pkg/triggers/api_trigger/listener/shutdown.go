package listener

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (s *Server) Shutdown() error {
	close(s.RequestChannel)

	//FIXME: For some reason the server will not shutdown gracefully, that's why the
	// context will timeout after 3 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}

	return nil
}
