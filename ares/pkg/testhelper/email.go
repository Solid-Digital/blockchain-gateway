package testhelper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jhillyerd/inbucket/pkg/rest/client"
)

// FIXME: get from config, but how?
const mailURL = "http://localhost:9000"

func (h *Helper) GetSignUpCodeFromEmail(requestID string, email string) string {
	return h.getCodeFromMail("X-SignUpRequestID", requestID, email)
}

func (h *Helper) GetRecoveryCodeFromEmail(requestID string, email string) string {
	return h.getCodeFromMail("X-RecoveryRequestID", requestID, email)
}

func (h *Helper) GetMailbox(email string) []*client.MessageHeader {
	mailClient := h.getMailClient()

	return h.getMailbox(mailClient, email)
}

func (h *Helper) getMailClient() *client.Client {
	mailClient, err := client.New(mailURL)

	h.suite.Require().NoError(err)

	return mailClient
}

func (h *Helper) getMailbox(mailClient *client.Client, email string) []*client.MessageHeader {
	mailbox, err := mailClient.ListMailbox(email)

	h.suite.Require().NoError(err)

	return mailbox
}

func (h *Helper) getCodeFromMail(codeType string, requestID string, email string) string {
	mailClient := h.getMailClient()
	mailbox := h.getMailbox(mailClient, email)

	var AuthCodeEmailRegexp = regexp.MustCompile("X-Code: ([A-Za-z0-9]{100,})")

	for _, mail := range mailbox {
		buf, err := mailClient.GetMessageSource(email, mail.ID)

		h.suite.Require().NoError(err)

		pattern := fmt.Sprintf("%s: %s", codeType, requestID)
		if strings.Contains(buf.String(), pattern) {
			matches := AuthCodeEmailRegexp.FindStringSubmatch(buf.String())
			if len(matches) < 2 {
				return ""
			}

			return matches[1]
		}
	}

	return ""
}
