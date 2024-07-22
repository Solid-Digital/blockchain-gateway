package aws

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/marketplaceentitlementservice"
	"github.com/unchainio/pkg/errors"
)

func (c *Client) GetEntitlements(customerID, productCode string) ([]*ares.AWSEntitlement, error) {
	entitlementsInput := &marketplaceentitlementservice.GetEntitlementsInput{
		ProductCode: aws.String(productCode),
		Filter: map[string][]*string{
			"CUSTOMER_IDENTIFIER": {aws.String(customerID)},
		},
	}
	entitlements, err := c.mes.GetEntitlements(entitlementsInput)
	if err != nil {
		return nil, errors.Wrap(err, "could not get entitlements")
	}

	var ret []*ares.AWSEntitlement
	for _, e := range entitlements.Entitlements {
		ret = append(ret, &ares.AWSEntitlement{
			ProductCode:        e.ProductCode,
			CustomerIdentifier: e.CustomerIdentifier,
			Dimension:          e.Dimension,
			ExpirationDate:     e.ExpirationDate,
			Value:              e.Value.StringValue,
		})
	}

	return ret, nil
}
