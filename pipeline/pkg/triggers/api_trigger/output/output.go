package output

import (
	"github.com/unchain/pipeline/pkg/domain"
	"net/http"
)

func NewOutput(req *http.Request) (domain.Output, error) {
	header := req.Header
	contentType := header.Get("Content-Type")

	body, err := newBody(req, contentType)
	if err != nil {
		return nil, err
	}

	return domain.Output{
		"header": header,
		"body":   body,
		"raw":    req,
	}, nil
}
