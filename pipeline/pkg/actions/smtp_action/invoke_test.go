package smtp_action_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/smtp_action"
	"github.com/unchain/pipeline/pkg/domain"
	"testing"
)

func (s *TestSuite) TestSmtpAction_Invoke() {
	exampleMail := []byte("To: test@nada.net\r\n" +
		"Subject: Test mail!\r\n" +
		"\r\n" +
		"This is the test email body.\r\n")

	cases := map[string]struct {
		stub    domain.Stub
		input   map[string]interface{}
		success bool
	}{
		"invoke action success ": {
			s.logger,
			map[string]interface{}{
				"username": "user@localhost",
				"password": "passw@localhost",
				"hostname": "localhost",
				"port": ":1025",
				"from": "test@localhost",
				"recipients": []string{"recipient@localhost"},
				"message": []byte(exampleMail),
			},
			true,
		},
		"invoke action fails with invalid from": {
			s.logger,
			map[string]interface{}{
				"username": "user@localhost",
				"password": "passw@localhost",
				"hostname": "localhost",
				"port": ":1025",
				"from": "NO REPLY JuicyChain",
				"recipients": []string{"jelle@unchain.io"},
				"message": []byte(exampleMail),
			},
			false,
		},
		"invoke action fails with empty recipients": {
			s.logger,
			map[string]interface{}{
				"username": "user@localhost",
				"password": "password",
				"hostname": "localhost",
				"port": ":1025",
				"from": "test@localhost",
				"recipients": []string{},
				"message": []byte(exampleMail),
			},
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			_, err := smtp_action.Invoke(tc.stub, tc.input)

			if tc.success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
