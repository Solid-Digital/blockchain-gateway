package mail

import (
	"strconv"

	"github.com/google/wire"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/unchainio/pkg/errors"
	"gopkg.in/gomail.v2"
)

var MailerSet = wire.NewSet(NewMailer, wire.Bind(new(ares.Mailer), new(Mailer)))

// todo create interfaces
type Mailer struct {
	Dialer *gomail.Dialer
	cfg    *Config
}

func NewMailer(cfg *Config) (mail *Mailer, err error) {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse port nr for mail server")
	}
	d := gomail.NewPlainDialer(cfg.Host, port, cfg.Username, cfg.Password)

	return &Mailer{
		Dialer: d,
		cfg:    cfg,
	}, nil

}

func (m *Mailer) SendMessage(to string, subject string, body string) (err error) {
	mes := gomail.NewMessage()
	mes.SetHeader("From", m.cfg.From)
	mes.SetHeader("To", to)
	mes.SetHeader("Subject", subject)
	mes.SetBody("text/html", body)

	if err := m.Dialer.DialAndSend(mes); err != nil {
		return errors.Wrap(err, "Could not send mail")
	}

	return nil
}

func (m *Mailer) SendRecoveryMessage(to, requestID, connectURL, recoveryCode, username string) (err error) {
	mes := gomail.NewMessage()
	mes.SetHeader("From", m.cfg.From)
	mes.SetHeader("To", to)
	mes.SetHeader("Subject", "Blockchain Gateway Password Recovery")
	mes.SetHeader("X-RecoveryRequestID", requestID)
	mes.SetHeader("X-Code", recoveryCode)

	emailBody, err := getRecoveryMailTemplate(connectURL, recoveryCode, username)
	if err != nil {
		return errors.Wrap(err, "could not generate email body")
	}

	mes.SetBody("text/html", emailBody)

	if err := m.Dialer.DialAndSend(mes); err != nil {
		return errors.Wrap(err, "could not send mail")
	}

	return nil
}

func (m *Mailer) SendSignUpMessage(to, requestID, connectURL string) (err error) {
	mes := gomail.NewMessage()
	mes.SetHeader("From", m.cfg.From)
	mes.SetHeader("To", to)
	mes.SetHeader("Subject", "Blockchain Gateway Account Confirmation")
	mes.SetHeader("X-SignUpRequestID", requestID)

	emailBody, err := getSignUpMailTemplate(connectURL)
	if err != nil {
		return errors.Wrap(err, "could not generate email body")
	}

	mes.SetBody("text/html", emailBody)

	if err := m.Dialer.DialAndSend(mes); err != nil {
		return errors.Wrap(err, "could not send mail")
	}

	return nil
}

func (m *Mailer) SendInviteMessage(to, requestID, signupUrl, orgName string) (err error) {
	mes := gomail.NewMessage()
	mes.SetHeader("From", m.cfg.From)
	mes.SetHeader("To", to)
	mes.SetHeader("Subject", "Blockchain Gateway Invitation")
	mes.SetHeader("X-SignUpRequestID", requestID)

	emailBody, err := getReviewInviteMailTemplate(signupUrl, orgName)
	if err != nil {
		return errors.Wrap(err, "could not generate email body")
	}

	mes.SetBody("text/html", emailBody)

	if err := m.Dialer.DialAndSend(mes); err != nil {
		return errors.Wrap(err, "could not send mail")
	}

	return nil
}
