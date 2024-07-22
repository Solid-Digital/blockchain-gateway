package ares

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/dto"
)

const (
	StatusPendingConfirmation = "PendingConfirmation"
	StatusActive              = "Active"
	StatusInactive            = "Inactive"
	TierFreePlan              = "Free Plan"
	TierStarterPlan           = "Starter Plan"
	TierBusinessPlan          = "Business Plan"
	TierEnterprisePlan        = "Enterprise Plan"
	ProviderAWS               = "AWS"
	ProviderUnchain           = "UNCHAIN"
	URLFmt                    = "%s/confirm-registration?token=%s&email=%s&org=%s&invited=%t"
)

type AuthCode struct {
	RequestID string `json:"requestID"`
	Email     string `json:"emailAddress"`
	Code      string `json:"recoveryCode"`
}

type AuthService interface {
	CreateRegistration(params *dto.CreateRegistrationRequest) *apperr.Error
	ConfirmRegistration(params *dto.ConfirmRegistrationRequest) (*dto.LoginResponse, *apperr.Error)
	Login(ip string, params *dto.LoginRequest) (*dto.LoginResponse, *apperr.Error)
	Logout(token *dto.Token) *apperr.Error
	GetCurrentUser(principal *dto.User) (*dto.GetCurrentUserResponse, *apperr.Error)
	UpdateCurrentUser(params *dto.UpdateCurrentUserRequest, u *dto.User) (*dto.GetCurrentUserResponse, *apperr.Error)
	ChangeCurrentPassword(params *dto.ChangeCurrentPasswordRequest, user *dto.User) *apperr.Error
	ResetPassword(params *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, *apperr.Error)
	ConfirmResetPassword(params *dto.ConfirmResetPasswordRequest) *apperr.Error
	DeleteCurrentUser(principal *dto.User) *apperr.Error
	HashPassword(password string) (string, *apperr.Error)
	CompareHashAndPassword(hash, password string) *apperr.Error
	Authenticate(token string) (*dto.User, error) // This should return a regular error
	InviteUserTx(ctx context.Context, tx *sql.Tx, email string, orgName string) (inviteID string, user *orm.User, appErr *apperr.Error)
	SetOrganizationService(service OrganizationService)

	GetConnectURL() string
}
