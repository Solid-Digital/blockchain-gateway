package mail

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"fmt"
	"github.com/matcornic/hermes"
)

func getRecoveryMailTemplate(recoveryURL, recoveryCode, username string) (string, error) {
	supportHost := fmt.Sprintf("%s",
		"https://"+ares.SupportURL)

	slackHost := fmt.Sprintf("%s",
		"https://"+ares.SlackURL)

	supportMail := fmt.Sprintf("%s",
		ares.SupportMail)

	h := hermes.Hermes{
		Theme: new(hermes.Default),
		Product: hermes.Product{
			Name:      "unchain.io",
			Link:      "https://unchain.io",
			Copyright: "Copyright © 2019 unchain.io. All rights reserved.",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Title: "Hi",
			Name:  username,
			Intros: []string{
				"You have received this email because a password reset request for your account was received. ",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the link below to reset your password.",
					Button: hermes.Button{
						Color:     "#f36d21",
						TextColor: "#ffffff",
						Text:      "Reset password",
						Link:      fmt.Sprintf("%s/reset-password/%s", recoveryURL, recoveryCode),
					},
				},
			},
			Outros: []string{
				"Please do not share this e-mail with anybody.",
				fmt.Sprintf("For support please go to %s or reach out to us via %s or via e-mail at %s.",
					supportHost, slackHost, supportMail),
			},
		},
	}

	return h.GenerateHTML(email)
}
