package api_trigger_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func (s *TestSuite) TestTrigger_NextMessage() {
	t := s.test_helpers.InitializedTrigger(s.helper.BytesFromFile("./testdata/config/config.toml"))

	cases := map[string]struct {
		Trigger      *api_trigger.Trigger
		Request      *http.Request
		ExpectedBody interface{}
		Success      bool
	}{
		"send without content type": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("some test"), map[string]string{"Authorization": "Bearer foo"}),
			"some test",
			true,
		},
		"content type application/json": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte(`{"foo": "bar"}`), map[string]string{"Authorization": "Bearer foo", "Content-Type": "application/json"}),
			map[string]interface{}{"foo": "bar"},
			true,
		},
		"content type application/xml": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte(`{"foo": "bar"}`), map[string]string{"Authorization": "Bearer foo", "Content-Type": "application/xml"}),
			`{"foo": "bar"}`,
			true,
		},
		"invalid json": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("invalid"), map[string]string{"Authorization": "Bearer foo", "Content-Type": "application/json"}),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			s.test_helpers.DelayedHttpRequest(tc.Request)
			tag, response, err := s.test_helpers.TriggerResponse(tc.Trigger, 3)

			if tc.Success {
				require.NotEmpty(t, tag)
				require.NotNil(t, response["header"])
				require.Equal(t, tc.ExpectedBody, response["body"])
				require.NoError(t, err)
			} else {
				require.NotEmpty(t, tag)
				require.Nil(t, response)
				require.Error(t, err)
			}
		})
	}

	t.Close()
}

func (s *TestSuite) TestTrigger_UploadFile() {
	fileDir, _ := os.Getwd()
	fileName := "./testdata/resources/example.txt"
	filePath := path.Join(fileDir, fileName)

	file, _ := os.Open(filePath)
	defer file.Close()

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	_, err := io.Copy(part, file)
	require.NoError(s.T(), err)
	err = writer.Close()
	require.NoError(s.T(), err)

	req := s.test_helpers.HttpPostRequestWithHeaders(buff.Bytes(), map[string]string{
		"Content-Type": writer.FormDataContentType(),
	})

	t := s.test_helpers.InitializedTrigger(s.helper.BytesFromFile("./testdata/config/config.toml"))

	s.test_helpers.DelayedHttpRequest(req)

	tag, response, err := s.test_helpers.TriggerResponse(t, 3)

	require.NotEmpty(s.T(), tag)
	require.NotNil(s.T(), response["header"])
	body, ok := response["body"].([]byte)
	require.True(s.T(), ok)
	require.Equal(s.T(), "This is a test file with small content", string(body))
}

func (s *TestSuite) TestTrigger_CORS() {
	tr := s.test_helpers.InitializedTrigger(s.helper.BytesFromFile("./testdata/config/cors_config.toml"))

	cases := map[string]struct {
		Origin  string
		Success bool
	}{
		"success": {
			"http://example.com",
			true,
		},
		"CORS wrong origin failure": {
			"http://non-listed-origin.com",
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			req := s.test_helpers.HttpPostRequestWithHeaders([]byte(`CORS request`), map[string]string{
				"Content-Type": "text/plain",
			})

			req.Header.Set("Origin", tc.Origin)
			s.test_helpers.DelayedHttpRequest(req)

			tag, response, err := s.test_helpers.TriggerResponse(tr, 3)

			fmt.Println(response)
			if tc.Success {
				require.NoError(t, err)
			} else {
				require.NotEmpty(t, tag)
				require.Error(t, err)
			}

		})
	}

	tr.Close()
}

func (s *TestSuite) TestTrigger_TriggerAPIKey() {
	t := s.test_helpers.InitializedTrigger(s.helper.BytesFromFile("./testdata/config/config_apikey.toml"))

	cases := map[string]struct {
		Trigger      *api_trigger.Trigger
		Request      *http.Request
		ExpectedBody interface{}
		Success      bool
	}{
		"valid api key": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("some test"), map[string]string{"Authorization": "Bearer foo"}),
			"some test",
			true,
		},
		"second api key": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("some test"), map[string]string{"Authorization": "Bearer bar"}),
			"some test",
			true,
		},
		"invalid api key": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("some test"), map[string]string{"Authorization": "Bearer invalid"}),
			nil,
			false,
		},
		"no api key": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("some test"), map[string]string{"Authorization": ""}),
			nil,
			false,
		},
		"basic auth": {
			t,
			s.test_helpers.HttpPostRequestWithBasicAuth([]byte("some test"), "tim", "secret"),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			s.test_helpers.DelayedHttpRequest(tc.Request)
			tag, response, err := s.test_helpers.TriggerResponse(tc.Trigger, 3)

			if tc.Success {
				require.NotEmpty(t, tag)
				require.NotNil(t, response["header"])
				require.Equal(t, tc.ExpectedBody, response["body"])
				require.NoError(t, err)
			} else {
				require.NotEmpty(t, tag)
				require.Nil(t, response)
				require.Error(t, err)
			}
		})
	}

	t.Close()
}

func (s *TestSuite) TestTrigger_TriggerBasicAuth() {
	t := s.test_helpers.InitializedTrigger(s.helper.BytesFromFile("./testdata/config/config_basicauth.toml"))

	cases := map[string]struct {
		Trigger      *api_trigger.Trigger
		Request      *http.Request
		ExpectedBody interface{}
		Success      bool
	}{
		"valid basic auth": {
			t,
			s.test_helpers.HttpPostRequestWithBasicAuth([]byte("some test"), "jim", "secret"),
			"some test",
			true,
		},
		"second basic auth": {
			t,
			s.test_helpers.HttpPostRequestWithBasicAuth([]byte("some test"), "jake", "super-secret"),
			"some test",
			true,
		},
		"invalid basicauth": {
			t,
			s.test_helpers.HttpPostRequestWithBasicAuth([]byte("some test"), "john", "not-so-secret"),
			nil,
			false,
		},
		"invalid auth": {
			t,
			s.test_helpers.HttpPostRequestWithHeaders([]byte("some test"), map[string]string{"Authorization": "foo"}),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			s.test_helpers.DelayedHttpRequest(tc.Request)
			tag, response, err := s.test_helpers.TriggerResponse(tc.Trigger, 3)

			if tc.Success {
				require.NotEmpty(t, tag)
				require.NotNil(t, response["header"])
				require.Equal(t, tc.ExpectedBody, response["body"])
				require.NoError(t, err)
			} else {
				require.NotEmpty(t, tag)
				require.Nil(t, response)
				require.Error(t, err)
			}
		})
	}

	t.Close()
}

func (s *TestSuite) TestTrigger_TriggerInvalidAPIKeyConfig() {
	_, err := s.test_helpers.InitializedTriggerWithError(s.helper.BytesFromFile("./testdata/config/config_invalid_apikey.toml"))

	s.Require().Error(err)
}

func (s *TestSuite) TestTrigger_TriggerInvalidBasicAuthConfig() {
	_, err := s.test_helpers.InitializedTriggerWithError(s.helper.BytesFromFile("./testdata/config/config_invalid_basicauth.toml"))

	s.Require().Error(err)
}
