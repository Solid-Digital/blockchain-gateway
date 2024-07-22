package imap_action_test

import (
	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/require"
	"github.com/unchain/pipeline/pkg/actions/imap_action"
	"github.com/unchain/pipeline/pkg/domain"
	"os"
	"testing"
)

func (s *TestSuite) TestClient_MarkMessageAsRead() {
	imapDomain := os.Getenv("MAIL_DOMAIN")
	imapPort := os.Getenv("IMAP_PORT")
	imapUsername := os.Getenv("MAIL_USERNAME")
	imapPassword := os.Getenv("MAIL_PASSWORD")

	config := &imap_action.Config{
		Username: imapUsername,
		Password: imapPassword,
		Port:     imapPort,
		Domain:   imapDomain,
	}
	client, err := imap_action.NewClient(s.logger, config)
	s.NoError(err, "could not start client")

	cases := map[string]struct {
		stub        domain.Stub
		client      *imap_action.Client
		success     bool
		seqNum      uint32
		removeFlags bool
	}{
		"successfully update status to read": {
			s.logger,
			client,
			true,
			1,
			true,
		},
		"fails with unknown seq num": {
			s.logger,
			client,
			false,
			999,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := tc.client.MarkMessageAsRead(tc.seqNum)

			if tc.success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			// 	remove seen flag
			if tc.removeFlags {
				seqSet := new(imap.SeqSet)
				seqSet.AddNum(tc.seqNum)

				item := imap.FormatFlagsOp(imap.RemoveFlags, true)
				flags := []interface{}{imap.SeenFlag}

				err := tc.client.Client.Store(seqSet, item, flags, nil)
				require.NoError(t, err, "could not remove seen flag")
			}
		})
	}
}

func (s *TestSuite) TestClient_MoveFailedMessage() {
	imapDomain := os.Getenv("IMAP_DOMAIN")
	imapPort := os.Getenv("IMAP_PORT")
	imapUsername := os.Getenv("IMAP_USERNAME")
	imapPassword := os.Getenv("IMAP_PASSWORD")

	config := &imap_action.Config{
		Username: imapUsername,
		Password: imapPassword,
		Port:     imapPort,
		Domain:   imapDomain,
	}
	client, err := imap_action.NewClient(s.logger, config)
	s.NoError(err, "could not start client")

	cases := map[string]struct {
		stub    domain.Stub
		client  *imap_action.Client
		success bool
		seqNum  uint32
		reset   bool
	}{
		"successfully moved message": {
			s.logger,
			client,
			true,
			1,
			true,
		},
		"fails with unknown seq num": {
			s.logger,
			client,
			false,
			999,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := tc.client.MoveFailedMessage(tc.seqNum)

			if tc.success {
				_, err := tc.client.Client.Select("INBOX", false)
				status := tc.client.Client.Mailbox()
				require.Equal(t, uint32(0), status.Messages)
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			// Move message back to inbox
			if tc.reset {
				_, err := tc.client.Client.Select("Failed", false)
				require.NoError(t, err)

				// get message
				seqSet := new(imap.SeqSet)
				seqSet.AddNum(tc.seqNum)

				// move message to INBOX
				err = tc.client.Client.MoveWithFallback(seqSet, "INBOX")
				require.NoError(t, err)

				// Check message moved back successfully
				_, err = tc.client.Client.Select("INBOX", false)
				status := tc.client.Client.Mailbox()
				require.Equal(t, uint32(1), status.Messages)

			}
		})
	}
}
