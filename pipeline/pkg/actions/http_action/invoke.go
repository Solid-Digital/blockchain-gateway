package http_action

import (
	"bytes"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Url                = "url"
	Method             = "method"
	RequestBody        = "requestBody"
	ContentType        = "contentType"
	ResponseBody       = "responseBody"
	ResponseHeaders    = "responseHeaders"
	ResponseStatusCode = "responseStatusCode"
)

func Invoke(stub domain.Stub, input map[string]interface{}) (output map[string]interface{}, err error) {
	url, ok := input[Url].(string)
	if !ok {
		return nil, errors.New("could not cast url")
	}

	requestBody, ok := input[RequestBody].([]byte)
	if !ok {
		return nil, errors.New("could not cast request body to []byte")
	}

	method, ok := input[Method].(string)
	if !ok {
		return nil, errors.New("could not cast method to string")
	}

	contentType := input[ContentType].(string)
	if !ok {
		return nil, errors.New("could not cast content type to string")
	}

	var resp *http.Response
	switch strings.ToUpper(method) {
	case "POST":
		resp, err = http.Post(url, contentType, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}
		stub.Debugf("request send to url: %s\n", url)
	default:
		return nil, errors.New("no valid method")
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}
	stub.Debugf("returning http request response with code %v, body:\n %v", resp.StatusCode, string(bodyBytes))

	return map[string]interface{}{
		ResponseBody:       bodyBytes,
		ResponseHeaders:    resp.Header,
		ResponseStatusCode: resp.StatusCode,
	}, nil
}
