package test_helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (t *TestHelpers) HttpPostRequest(content []byte) *http.Request {
	method := http.MethodPost
	url := fmt.Sprintf("http://localhost:%s", TestPort)
	body := bytes.NewReader(content)
	request, err := http.NewRequest(method, url, body)
	t.suite.Require().NoError(err)

	return request
}

func (t *TestHelpers) HttpPostRequestWithHeaders(content []byte, headers map[string]string) *http.Request {
	request := t.HttpPostRequest(content)

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return request
}

func (t *TestHelpers) HttpPostRequestWithBasicAuth(content []byte, username, password string) *http.Request {
	request := t.HttpPostRequest(content)

	request.SetBasicAuth(username, password)

	return request

}

func (t *TestHelpers) HttpRequest(request *http.Request) *http.Response {
	client := &http.Client{}
	response, err := client.Do(request)

	t.suite.Require().NoError(err)

	return response
}

func (t *TestHelpers) DelayedHttpRequest(request *http.Request) chan *http.Response {
	responseCh := make(chan *http.Response)

	go func() {
		time.Sleep(time.Second)
		responseCh <- t.HttpRequest(request)
	}()

	return responseCh
}

func (t *TestHelpers) ReadBody(response *http.Response) string {
	b, err := ioutil.ReadAll(response.Body)

	t.suite.Require().NoError(err)

	return string(bytes.TrimSpace(b))
}
