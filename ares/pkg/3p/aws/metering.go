package aws

import (
	"github.com/aws/aws-sdk-go/service/marketplacemetering"
)

func (c *Client) ResolveCustomer(token string) (customerId string, productCode string, err error) {
	input := &marketplacemetering.ResolveCustomerInput{
		RegistrationToken: &token,
	}
	output, err := c.mpm.ResolveCustomer(input)
	if err != nil {
		return "", "", err
	}

	return *output.CustomerIdentifier, *output.ProductCode, nil
}
