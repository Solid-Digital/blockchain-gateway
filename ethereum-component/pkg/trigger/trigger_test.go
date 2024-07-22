package trigger_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/errors"

	"bitbucket.org/unchain/ethereum2/pkg/trigger"
)

func (s *TestSuite) TestTrigger_Trigger() {
	simpleContractSolidity := s.helper.StringFromFile("./testdata/contract/simple_contract.sol")
	singleValueEventSolidity := s.helper.StringFromFile("./testdata/contract/single_value_event.sol")
	singleNonIndexedValueEventSolidity := s.helper.StringFromFile("./testdata/contract/single_non_indexed_value_event.sol")
	deployContractAction := s.factory.InitializedAction(s.helper.BytesFromFile("./testdata/config/action_config.toml"))

	address1, ABI1 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, singleValueEventSolidity)
	action1 := s.factory.InitializedAction(s.helper.FillInActionConfig(s.helper.StringFromFile("./testdata/config/action_config_fmt.toml"), address1, ABI1))
	trigger1 := s.factory.InitializedTrigger(s.helper.FillInTriggerConfig(s.helper.StringFromFile("./testdata/config/set_event_fmt.toml"), address1, ABI1))
	func1 := func() {
		// this slows down tests a lot, is there a more efficient way to do this? Like with defer?
		time.Sleep(time.Second)
		s.helper.CallContractFunction(action1, DefaultAccount, address1, "set", nil)
	}
	expected1 := map[string]interface{}{
		"Event": "Set",
		"Values": []map[string]interface{}{
			{
				"sender": DefaultAccount,
			},
		},
	}

	address2, ABI2 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, singleNonIndexedValueEventSolidity)
	action2 := s.factory.InitializedAction(s.helper.FillInActionConfig(s.helper.StringFromFile("./testdata/config/action_config_fmt.toml"), address2, ABI2))
	trigger2 := s.factory.InitializedTrigger(s.helper.FillInTriggerConfig(s.helper.StringFromFile("./testdata/config/set_event_fmt.toml"), address2, ABI2))
	func2 := func() {
		// this slows down tests a lot, is there a more efficient way to do this? Like with defer?
		time.Sleep(time.Second)
		s.helper.CallContractFunction(action2, DefaultAccount, address2, "set", nil)
	}
	expected2 := map[string]interface{}{
		"Event": "Set",
		"Values": []map[string]interface{}{
			{
				"sender": DefaultAccount,
			},
		},
	}

	address3, ABI3 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action3 := s.factory.InitializedAction(s.helper.FillInActionConfig(s.helper.StringFromFile("./testdata/config/action_config_fmt.toml"), address3, ABI3))
	trigger3 := s.factory.InitializedTrigger(s.helper.FillInTriggerConfig(s.helper.StringFromFile("./testdata/config/set_event_fmt.toml"), address3, ABI3))
	func3 := func() {
		// this slows down tests a lot, is there a more efficient way to do this? Like with defer?
		time.Sleep(time.Second)
		s.helper.CallContractFunction(action3, DefaultAccount, address3, "set", nil)
	}
	expected3 := map[string]interface{}{
		"Event": "Set",
		"Values": []map[string]interface{}{
			{
				"sender": DefaultAccount,
			},
			{
				"x": 10,
			},
		},
	}

	address4, ABI4 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action4 := s.factory.InitializedAction(s.helper.FillInActionConfig(s.helper.StringFromFile("./testdata/config/action_config_fmt.toml"), address4, ABI4))
	trigger4 := s.factory.InitializedTrigger(s.helper.FillInTriggerConfig(s.helper.StringFromFile("./testdata/config/set_event_with_filters_fmt.toml"), address4, ABI4))
	func4 := func() {
		// this slows down tests a lot, is there a more efficient way to do this? Like with defer?
		time.Sleep(time.Second)
		s.helper.CallContractFunction(action4, DefaultAccount, address4, "set", nil)
	}
	expected4 := map[string]interface{}{
		"Event": "Set",
		"Values": []map[string]interface{}{
			{
				"sender": DefaultAccount,
			},
			{
				"x": 10,
			},
		},
	}

	address5, ABI5 := s.helper.DeploySingleContractAction(deployContractAction, DefaultAccount, simpleContractSolidity)
	action5 := s.factory.InitializedAction(s.helper.FillInActionConfig(s.helper.StringFromFile("./testdata/config/action_config_fmt.toml"), address5, ABI5))
	trigger5 := s.factory.InitializedTrigger(s.helper.FillInTriggerConfig(s.helper.StringFromFile("./testdata/config/set_event_with_filters_fmt.toml"), address5, ABI5))
	func5 := func() {
		// this slows down tests a lot, is there a more efficient way to do this? Like with defer?
		time.Sleep(time.Second)
		s.helper.CallContractFunction(action5, AlternativeAccount, address5, "set", nil)
	}
	expected5 := map[string]interface{}{
		"Event": "Set",
		"Values": []map[string]interface{}{
			{
				"sender": DefaultAccount,
			},
			{
				"x": 10,
			},
		},
	}

	cases := map[string]struct {
		Trigger  *trigger.Trigger
		EventFn  func()
		Expected []byte
		Success  bool
	}{
		"emit event single indexed return value": {
			trigger1,
			func1,
			s.helper.ToJSON(expected1),
			true,
		},
		"emit event single non-indexed return value": {
			trigger2,
			func2,
			s.helper.ToJSON(expected2),
			true,
		},
		"emit event multiple return values": {
			trigger3,
			func3,
			s.helper.ToJSON(expected3),
			true,
		},
		"event passes filters": {
			trigger4,
			func4,
			s.helper.ToJSON(expected4),
			true,
		},
		"event does not pass filters": {
			trigger5,
			func5,
			s.helper.ToJSON(expected5),
			false,
		},
	}

	filterMap := func(input map[string]interface{}, keys []string) map[string]interface{} {
		ret := map[string]interface{}{}
		for _, key := range keys {
			ret[key] = input[key]
		}

		return ret
	}

	listen := func(trigger *trigger.Trigger) (string, map[string]interface{}, error) {
		type triggerResponse struct {
			tag      string
			response map[string]interface{}
			err      error
		}

		responseChan := make(chan triggerResponse)

		go func() {
			tag, response, err := trigger.Trigger()
			responseChan <- triggerResponse{
				tag:      tag,
				response: response,
				err:      err,
			}
		}()

		for i := 0; i < 3; i++ {
			select {
			case response := <-responseChan:
				return response.tag, response.response, response.err
			default:
				time.Sleep(time.Second)
			}
		}

		return "", nil, errors.New("trigger did not receive anything")
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			tc.EventFn()
			tag, response, err := listen(tc.Trigger)

			if tc.Success {
				require.NoError(t, err)
				require.NotEmpty(t, tag)
				require.True(t, len(response) > 0)

				// Not a great way of testing this, and we're completely ignoring the "Log" output
				filtered := filterMap(response, []string{"Event", "Values"})

				require.Equal(t, tc.Expected, s.helper.ToJSON(filtered))
			} else {
				require.Error(t, err)
				require.Equal(t, "", tag)
				require.Equal(t, 0, len(response))
			}
		})
	}
}
