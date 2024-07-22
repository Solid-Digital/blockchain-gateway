package api_trigger_test

import (
	"github.com/unchain/pipeline/pkg/triggers/api_trigger"
	"github.com/unchainio/pkg/errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestTrigger_Respond() {
	t := s.test_helpers.InitializedTrigger([]byte(`port="8888"`))

	cases := map[string]struct {
		Trigger        *api_trigger.Trigger
		TriggerFn      func(*api_trigger.Trigger) (string, chan *http.Response)
		RespondInput   map[string]interface{}
		RespondErr     error
		ExpectedStatus int
		ExpectedBody   string
	}{
		"send valid message": {
			t,
			s.triggerFn(),
			map[string]interface{}{"foo": "bar"},
			nil,
			http.StatusOK,
			`{"foo":"bar"}`,
		},
		"empty input": {
			t,
			s.triggerFn(),
			map[string]interface{}{},
			nil,
			http.StatusOK,
			`{}`,
		},
		"xml response": {
			t,
			s.xmlTriggerFn(),
			map[string]interface{}{"foo": "bar"},
			nil,
			http.StatusInternalServerError,
			`xml: unsupported type: map[string]interface {}`,
		},
		"with error": {
			t,
			s.triggerFn(),
			nil,
			errors.New("failed"),
			http.StatusInternalServerError,
			`error: failed`,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			tag, responseCh := tc.TriggerFn(tc.Trigger)
			err := s.test_helpers.RespondResponse(tc.Trigger, tag, tc.RespondInput, tc.RespondErr, 3)
			s.Require().NoError(err)

			response := <-responseCh

			require.Equal(t, tc.ExpectedStatus, response.StatusCode)
			require.Equal(t, tc.ExpectedBody, s.test_helpers.ReadBody(response))
		})
	}

	t.Close()
}

func (s *TestSuite) triggerFn() func(t *api_trigger.Trigger) (tag string, responseCh chan *http.Response) {
	return func(t *api_trigger.Trigger) (tag string, responseCh chan *http.Response) {
		// do we still need a delay? -> no not for this test
		responseCh = s.test_helpers.DelayedHttpRequest(s.test_helpers.HttpPostRequest([]byte(`test`)))

		tag, _, err := s.test_helpers.TriggerResponse(t, 3)
		s.Suite.Require().NoError(err)

		return tag, responseCh
	}
}

func (s *TestSuite) xmlTriggerFn() func(t *api_trigger.Trigger) (tag string, responseCh chan *http.Response) {
	return func(t *api_trigger.Trigger) (tag string, responseCh chan *http.Response) {
		// do we still need a delay? -> no not for this test
		headers := map[string]string{"Accept": "application/xml"}
		req := s.test_helpers.HttpPostRequestWithHeaders([]byte(`test`), headers)

		responseCh = s.test_helpers.DelayedHttpRequest(req)

		tag, _, err := s.test_helpers.TriggerResponse(t, 3)
		s.Suite.Require().NoError(err)

		return tag, responseCh
	}
}
