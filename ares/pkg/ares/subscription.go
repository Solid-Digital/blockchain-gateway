package ares

type SubscriptionService interface {
	ConsumeMarketplaceNotificationMessage() error
}
