package mail

import (
	"fmt"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/matcornic/hermes"
)

func getReviewInviteMailTemplate(registrationURL, orgName string) (string, error) {
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
			Intros: []string{
				"Welcome to The Blockchain Gateway!",
				fmt.Sprintf("You have been invited to join organization %s.", orgName),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to create an account and gain access:",
					Button: hermes.Button{
						Color:     "#f36d21",
						Text:      "Accept invitation",
						TextColor: "#FFFFFF",
						Link:      registrationURL,
					},
				},
			},
			Outros: []string{
				fmt.Sprintf("For support please go to %s or reach out to us via %s or via e-mail at %s.",
					supportHost, slackHost, supportMail),
			},
		},
	}

	return h.GenerateHTML(email)
}
