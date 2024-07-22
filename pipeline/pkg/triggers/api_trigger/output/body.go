package output

import (
	"github.com/unchainio/pkg/errors"
	"io/ioutil"
	"net/http"
)

func newBody(req *http.Request, contentType string) (body interface{}, err error) {
	outputType := outputType(contentType)

	switch outputType {
	case JSON:
		bodyBytes, err := getBody(req)
		if err != nil {
			return nil, err
		}

		body, err = jsonBody(bodyBytes)
		if err != nil {
			return nil, err
		}
	case FORMDATA:
		body, err = formdataBody(req)
		if err != nil {
			return nil, err
		}

	default:
		bodyBytes, err := getBody(req)
		if err != nil {
			return nil, err
		}

		body, err = textBody(bodyBytes)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

func getBody(req *http.Request) ([]byte, error) {
	content := req.Body
	defer content.Close()

	bodyBytes, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	return bodyBytes, nil
}