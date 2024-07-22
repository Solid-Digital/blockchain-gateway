package imap_action_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/imap_action"
	"github.com/unchain/pipeline/pkg/domain"
	"os"
	"testing"
)

func (s *TestSuite) TestTrigger_Invoke() {
	imapDomain := os.Getenv("IMAP_DOMAIN")
	imapPort := os.Getenv("IMAP_PORT")
	imapUsername := os.Getenv("IMAP_USERNAME")
	imapPassword := os.Getenv("IMAP_PASSWORD")

	cases := map[string]struct {
		stub    domain.Stub
		input   map[string]interface{}
		success bool
		expectedOutput map[string]interface{}
	}{
		"invoke action can load valid config successfully": {
			s.logger,
			map[string]interface{}{
				"config": fmt.Sprintf(`
					port = %s
					domain = %s
					username = %s
					password = %s
					fromFilter = ""
					subjectFilter = ""
				`, imapPort, imapDomain, imapUsername, imapPassword),
				"function": "NoOp",
				"params": map[string]interface{}{},
			},
			true,
			nil,
		},
		"invoke action can fetch email attachment successfully": {
			s.logger,
			map[string]interface{}{
				"config": fmt.Sprintf(`
					port = %s
					domain = %s
					username = %s
					password = %s
					fromFilter = ""
					subjectFilter = ""
				`, imapPort, imapDomain, imapUsername, imapPassword),
				"function": "GetNewMessageAttachments",
				"params": map[string]interface{}{},
			},

			true,
			map[string]interface{}{
				"messages": map[uint32]interface{}{
					0x1:[]uint8{0xef, 0xbb, 0xbf, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x31, 0x2c, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x32, 0xd, 0xa, 0x72, 0x6f, 0x77, 0x5f, 0x31, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x31, 0x2c, 0x72, 0x6f, 0x77, 0x5f, 0x31, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x32, 0xd, 0xa, 0x72, 0x6f, 0x77, 0x5f, 0x32, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x31, 0x2c, 0x72, 0x6f, 0x77, 0x5f, 0x32, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x32},
				},
			},
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			output, err := imap_action.Invoke(tc.stub, tc.input)

			if tc.success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			if tc.expectedOutput != nil {
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}
