package ares

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

//OrganizationService
type OrganizationService interface {
	InviteMember(params *dto.InviteMemberRequest, orgName string) (*dto.InviteMemberResponse, *apperr.Error)
	InviteMemberTx(ctx context.Context, tx *sql.Tx, user *orm.User, orgName string, roles map[string]bool) (*dto.GetMemberResponse, *apperr.Error)
	RemoveMember(email string, orgName string) *apperr.Error
	SetMemberRoles(params *dto.SetMemberRolesRequest, email string, orgName string, u *dto.User) *apperr.Error
	GetMember(email string, orgName string, u *dto.User) (*dto.GetMemberResponse, *apperr.Error)
	GetAllMembers(orgName string, principal *dto.User) ([]*dto.GetMemberResponse, *apperr.Error)
	CreateOrganization(params *dto.CreateOrganizationRequest, principal *dto.User) (*dto.GetOrganizationResponse, *apperr.Error)
	GetAllOrganizations(principal *dto.User) ([]*dto.GetOrganizationResponse, *apperr.Error)
	GetOrganization(orgName string, principal *dto.User) (*dto.GetOrganizationResponse, *apperr.Error)
	UpdateOrganization(params *dto.UpdateOrganizationRequest, orgName string) (*dto.GetOrganizationResponse, *apperr.Error)
}
