package mail

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"fmt"
	"github.com/matcornic/hermes"
)

func getSignUpMailTemplate(signupUrl string) (string, error) {
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
			Title: "Hi!",
			Name:  "",
			Intros: []string{
				"Welcome to The Blockchain Gateway!",
				"We have prepared an account and an organization for you. You are only one step away " +
					"from using the Blockchain Gateway.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Finalize the registration by clicking the link below",
					Button: hermes.Button{
						Color:     "#f36d21",
						TextColor: "#ffffff",
						Text:      "Confirm Your Account",
						Link:      fmt.Sprintf("%s", signupUrl),
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
