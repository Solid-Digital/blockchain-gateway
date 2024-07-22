package http_action_test

import (
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/http_action"
	"github.com/unchain/pipeline/pkg/domain"
	"net/http"
	"testing"
)

func (s *TestSuite) TestHttpAction_Invoke() {
	port := ":8888"

	cases := map[string]struct {
		Stub               domain.Stub
		Input              map[string]interface{}
		Success            bool
		ExpectedBody       interface{}
		ExpectedHeaders    map[string][]string
		ExpectedStatusCode int
		ExpectedError      string
	}{
		"successfully calls api": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://localhost:8888/api/success",
				http_action.ContentType: "application/json",
				http_action.Method:      "POST",
				http_action.RequestBody: []byte("sup"),
			},
			true,
			"{\"key\":\"value\"}\n",
			nil,
			200,
			"",
		},
		"bad request returns error and 400 status code": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://localhost:8888/api/badrequest",
				http_action.ContentType: "application/json",
				http_action.Method:      "POST",
				http_action.RequestBody: []byte("sup"),
			},
			true,
			"{\"key\":\"value\"}\n",
			nil,
			400,
			"",
		},
		"failed api gives back error and 500 status code": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://localhost:8888/api/fail",
				http_action.ContentType: "application/json",
				http_action.Method:      "POST",
				http_action.RequestBody: []byte("sup"),
			},
			true,
			nil,
			nil,
			500,
			"",
		},
		"fails without url": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         nil,
				http_action.ContentType: "application/json",
				http_action.Method:      "POST",
				http_action.RequestBody: []byte("sup"),
			},
			false,
			nil,
			nil,
			0,
			"could not cast url",
		},
		"fails with unknown url": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://gibberish:8888/api/fail",
				http_action.ContentType: "application/json",
				http_action.Method:      "POST",
				http_action.RequestBody: []byte("sup"),
			},
			false,
			nil,
			nil,
			0,
			"Post \"http://gibberish:8888/api/fail\": dial tcp: lookup gibberish: no such host",
		},
		"fails without request body": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://localhost:8888/api/success",
				http_action.ContentType: "application/json",
				http_action.Method:      "POST",
				http_action.RequestBody: nil,
			},
			false,
			nil,
			nil,
			0,
			"could not cast request body to []byte",
		},
		"fails without method": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://localhost:8888/api/success",
				http_action.ContentType: "application/json",
				http_action.Method:      nil,
				http_action.RequestBody: []byte("sup"),
			},
			false,
			nil,
			nil,
			0,
			"could not cast method to string",
		},
		"fails with unknown method": {
			s.logger,
			map[string]interface{}{
				http_action.Url:         "http://localhost:8888/api/success",
				http_action.ContentType: "application/json",
				http_action.Method:      "INVALID_METHOD",
				http_action.RequestBody: []byte("sup"),
			},
			false,
			nil,
			nil,
			0,
			"no valid method",
		},
	}

	// start local api
	http.HandleFunc("/api/success", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]interface{}{
			"key": "value",
		})
	})
	http.HandleFunc("/api/fail", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		render.JSON(w, r, map[string]interface{}{})
	})
	http.HandleFunc("/api/badrequest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		render.JSON(w, r, map[string]interface{}{
			"key": "value",
		})
	})

	server := &http.Server{
		Addr: port,
	}
	go server.ListenAndServe()

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			// time.Sleep(1 * time.Second)
			res, err := http_action.Invoke(tc.Stub, tc.Input)

			if tc.ExpectedHeaders != nil {
				headers := res[http_action.ResponseHeaders].(map[string][]string)
				require.Equal(t, tc.ExpectedHeaders, headers)
			}
			if tc.ExpectedStatusCode != 0 {
				statusCode := res[http_action.ResponseStatusCode].(int)
				require.Equal(t, tc.ExpectedStatusCode, statusCode)
			}
			if tc.ExpectedBody != nil {
				responseBody := res[http_action.ResponseBody].([]byte)
				require.Equal(t, tc.ExpectedBody, string(responseBody))
			}
			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, res)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedError, err.Error())
			}
		})
	}

}
