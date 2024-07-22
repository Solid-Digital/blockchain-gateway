package smtp_action

import (
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"net/smtp"
)

func Invoke(stub domain.Stub, input map[string]interface{}) (output map[string]interface{}, err error) {
	msg, err := NewInput(input)
	if err != nil {
		return nil, err
	}
	auth := smtp.PlainAuth("", msg.Username, msg.Password, msg.Hostname)

	if len(msg.Recipients) < 1 {
		return nil, errors.New("no recipients found for email")
	}

	err = smtp.SendMail(msg.Hostname+msg.Port, auth, msg.From, msg.Recipients, msg.Message)
	if err != nil {
		return nil, errors.Wrap(err, "could not send mail")
	}
	stub.Printf("mail send")

	return nil, nil
}
