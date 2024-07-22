package ares

import (
	"io"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/dto"
)

type ComponentType string

const (
	ComponentTypeBase    ComponentType = "base"
	ComponentTypeTrigger ComponentType = "trigger"
	ComponentTypeAction  ComponentType = "action"
)

type ComponentService interface {
	CreateAction(params *dto.CreateComponentRequest, orgName string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error)
	CreateActionVersion(params *CreateActionVersionRequest) (*dto.GetComponentVersionResponse, *apperr.Error)
	GetPublicAction(name string) (*dto.GetComponentResponse, *apperr.Error)
	GetAction(orgName string, name string) (*dto.GetComponentResponse, *apperr.Error)
	GetActionVersion(orgName string, name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error)
	GetPublicActionVersion(name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error)
	CreateBase(params *dto.CreateComponentRequest, orgName string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error)
	CreateBaseVersion(params *dto.CreateBaseVersionRequest, orgName string, name string, principal *dto.User) (*dto.GetBaseVersionResponse, *apperr.Error)
	CreateTrigger(params *dto.CreateComponentRequest, orgName string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error)
	CreateTriggerVersion(params *CreateTriggerVersionRequest) (*dto.GetComponentVersionResponse, *apperr.Error)
	GetBase(orgName string, name string) (*dto.GetComponentResponse, *apperr.Error)
	GetBaseVersion(orgName string, name string, version string) (*dto.GetBaseVersionResponse, *apperr.Error)
	GetPublicTrigger(name string) (*dto.GetComponentResponse, *apperr.Error)
	GetTrigger(orgName string, name string) (*dto.GetComponentResponse, *apperr.Error)
	GetTriggerVersion(orgName string, name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error)
	GetPublicTriggerVersion(name string, version string) (*dto.GetComponentVersionResponse, *apperr.Error)
	UpdateAction(params *dto.UpdateComponentRequest, orgName string, name string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error)
	UpdateBase(params *dto.UpdateComponentRequest, orgName string, name string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error)
	UpdateTrigger(params *dto.UpdateComponentRequest, orgName string, name string, principal *dto.User) (*dto.GetComponentResponse, *apperr.Error)
	GetAllBases(orgName string, available *bool) ([]*dto.GetComponentResponse, *apperr.Error)
	GetAllTriggers(orgName string, available *bool) ([]*dto.GetComponentResponse, *apperr.Error)
	GetAllPublicTriggers() ([]*dto.GetComponentResponse, *apperr.Error)
	GetAllActions(orgName string, available *bool) ([]*dto.GetComponentResponse, *apperr.Error)
	GetAllPublicActions() ([]*dto.GetComponentResponse, *apperr.Error)
	UpdateTriggerVersion(params *dto.UpdateComponentVersionRequest, orgName string, name string, version string, principal *dto.User) (*dto.GetComponentVersionResponse, *apperr.Error)
	UpdateActionVersion(params *dto.UpdateComponentVersionRequest, orgName string, name string, version string, principal *dto.User) (*dto.GetComponentVersionResponse, *apperr.Error)
	UpdateBaseVersion(params *dto.UpdateBaseVersionRequest, orgName string, name string, version string, principal *dto.User) (*dto.GetBaseVersionResponse, *apperr.Error)
}

type CreateTriggerVersionRequest struct {
	/*
	  Required: true
	  In: path
	*/
	OrgName string

	/*
	  Required: true
	  In: path
	*/
	Name string

	/*version string for this action version
	  In: formData
	*/
	Version string

	/*short description of this action version
	  In: formData
	*/
	Description string

	/*readme for this action version
	  In: formData
	*/
	Readme string

	/*default config for this action version
	  In: formData
	*/
	ExampleConfig string

	/*json encoded string containing the input specification of the action
	  In: formData
	*/
	InputSchema []string

	/*json encoded string containing the output specification of the action
	  In: formData
	*/
	OutputSchema []string

	/*describes whether or not this action version is public
	  In: formData
	*/
	Public *bool

	/*the trigger version file
	  In: formData
	*/
	TriggerFile io.ReadCloser

	Principal *dto.User
}

type CreateActionVersionRequest struct {
	/*
	  Required: true
	  In: path
	*/
	OrgName string

	/*
	  Required: true
	  In: path
	*/
	Name string

	/*version string for this action version
	  In: formData
	*/
	Version string

	/*short description of this action version
	  In: formData
	*/
	Description string

	/*readme for this action version
	  In: formData
	*/
	Readme string

	/*default config for this action version
	  In: formData
	*/
	ExampleConfig string

	/*json encoded string containing the input specification of the action
	  In: formData
	*/
	InputSchema []string

	/*json encoded string containing the output specification of the action
	  In: formData
	*/
	OutputSchema []string

	/*describes whether or not this action version is public
	  In: formData
	*/
	Public *bool

	/*the action version file
	  In: formData
	*/
	ActionFile io.ReadCloser

	Principal *dto.User
}
