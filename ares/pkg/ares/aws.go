package ares

import "time"

type AWSClient interface {
	ResolveCustomer(token string) (customerId string, productCode string, err error)
	DeleteSQSMessage(Handle *string) error
	ReceiveMarketplaceNotification() AWSMarketplaceNotificationMessage
	GetEntitlements(customerID string, productCode string) ([]*AWSEntitlement, error)
}

type AWSMarketplaceNotificationMessage struct {
	Body   AWSMarketplaceNotificationMessageBody
	Handle *string
	Error  error
}

type AWSMarketplaceNotificationMessageBody struct {
	Action             string
	CustomerIdentifier string
	ProductCode        string
}

type AWSEntitlement struct {
	CustomerIdentifier *string
	Dimension          *string
	ExpirationDate     *time.Time
	ProductCode        *string
	Value              *string
}

type AWSBillingInfo struct {
	CustomerIdentifier string `json:"awsCustomerId"`
	ProductCode        string
}
