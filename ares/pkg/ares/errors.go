package ares

import (
	"fmt"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
)

func ErrDuplicatePipeline(err *apperr.Error, orgName string, name string) *apperr.Error {
	return err.AddNamedErrors("name", fmt.Sprintf("pipeline %q already exists in organization %q", name, orgName))
}

func ErrPipelineNotFound(err *apperr.Error, orgName string, pipelineName string) *apperr.Error {
	return err.WithMessagef("pipeline %q not found in organization %q", pipelineName, orgName)
}

func ErrDeploymentNotFound(err *apperr.Error, orgName string, pipelineName string, envName string) *apperr.Error {
	return err.WithMessagef("deployment not found in environment %q of pipeline %q in organization %q", envName, pipelineName, orgName)
}

func ErrEnvVarsNotFound(err *apperr.Error, orgName string, pipelineName string, envName string) *apperr.Error {
	return err.WithMessagef("no environment variables were found for environment %q of pipeline %q in organization %q", envName, pipelineName, orgName)
}

func ErrEnvVarNotFound(err *apperr.Error, orgName string, pipelineName string, envName string, varID int64) *apperr.Error {
	return err.WithMessagef("environment variable with id %q not found for environment %q of pipeline %q in organization %q", varID, envName, pipelineName, orgName)
}

func ErrDuplicateEnvVar(err *apperr.Error, orgName, pipelineName, envName, key string) *apperr.Error {
	return err.AddNamedErrors("key", fmt.Sprintf("environment variable with key %q already exists in environment %q of pipeline %q in organization %q", key, envName, pipelineName, orgName))
}

func ErrOrgNotFound(err *apperr.Error, orgName string) *apperr.Error {
	return err.WithMessagef("organization %q not found", orgName)
}

func ErrOrgsNotFound(err *apperr.Error) *apperr.Error {
	return err.WithMessage("no organizations found")
}

func ErrDuplicateOrg(err *apperr.Error, orgName string) *apperr.Error {
	return err.WithMessagef("organization %q already exists", orgName)
}

func ErrDuplicateMembership(err *apperr.Error, userEmail string, orgName string) *apperr.Error {
	return err.WithMessagef("user %q is already a member of organization %q", userEmail, orgName)
}

func ErrDuplicateUser(err *apperr.Error, userEmail string) *apperr.Error {
	return err.WithMessagef("user with email %q already exists", userEmail)
}

func ErrDuplicateComponent(err *apperr.Error, componentType ComponentType, name string) *apperr.Error {
	return err.AddNamedErrors("name", fmt.Sprintf("component %q of type %q already exists", name, componentType))
}

func ErrDuplicateComponentVersion(err *apperr.Error, componentType ComponentType, name string, version string) *apperr.Error {
	return err.AddNamedErrors("version", fmt.Sprintf("version \"%s@%s\" of component type %q already exists", name, version, componentType))
}

func ErrForbiddenToSetPublic(err error) *apperr.Error {
	return apperr.Forbidden.Wrap(err).WithMessage("only admins can set the Public field")
}

func ErrForbiddenView(err error) *apperr.Error {
	return apperr.Forbidden.Wrap(err).WithMessage("viewing this resource is not allowed")
}

func ErrForbiddenEdit(err error) *apperr.Error {
	return apperr.Forbidden.Wrap(err).WithMessage("editing this resource is not allowed")
}

func ErrLoadingCreatedBy(err *apperr.Error, name string) *apperr.Error {
	return err.WithMessagef("failed to load the user who created %q", name)
}

func ErrLoadingUpdatedBy(err *apperr.Error, name string) *apperr.Error {
	return err.WithMessagef("failed to load the user who last updated %q", name)
}

func ErrComponentVersionNotFound(err *apperr.Error, componentType ComponentType, orgName, name string, version string) *apperr.Error {
	return err.WithMessagef("version \"%s@%s\" of component type %q not found in organization %q", name, version, componentType, orgName)
}

func ErrComponentNotFound(err *apperr.Error, componentType ComponentType, orgName, name string) *apperr.Error {
	return err.WithMessagef("component %q of type %q not found in organization %q", name, componentType, orgName)
}

func ErrPipelineConfigNotFound(err *apperr.Error, orgName, pipelineName string, revision int64) *apperr.Error {
	return err.WithMessagef("configuration revision %d not found for pipeline %q in organization %q", revision, pipelineName, orgName)
}

func ErrEnvironmentNotFound(err *apperr.Error, orgName, envName string) *apperr.Error {
	return err.WithMessagef("environment %q not found in organization %q", envName, orgName)
}

func ErrEnvironmentsNotFound(err *apperr.Error, orgName string) *apperr.Error {
	return err.WithMessagef("no environments were found in organization %q", orgName)
}

func ErrTooManyReplicas(err *apperr.Error, requested int64, maximum int64) *apperr.Error {
	return err.AddNamedErrors("replicas", fmt.Sprintf("you requested %d replicas - maximum is %d", requested, maximum))
}

func ErrComponentsNotFound(err *apperr.Error, componentType ComponentType, orgName string) *apperr.Error {
	return err.WithMessagef("no components of type %q were found in organization %q", componentType, orgName)
}

func ErrParseInputSchema(err error) *apperr.Error {
	return apperr.Internal.Wrap(err).AddNamedErrors("inputSchema", "failed to parse the input schema")
}

func ErrParseOutputSchema(err error) *apperr.Error {
	return apperr.Internal.Wrap(err).AddNamedErrors("outputSchema", "failed to parse the output schema")
}

func ErrDefaultEnvsNotFound(err *apperr.Error) *apperr.Error {
	return err.WithMessage("no default environments found")
}

func ErrUserIDNotFound(err *apperr.Error, userID int64) *apperr.Error {
	return err.WithMessagef("user with ID %d not found", userID)
}

func ErrUserEmailNotFound(err *apperr.Error, userEmail string) *apperr.Error {
	return err.WithMessagef("user with email %q not found", userEmail)
}

func ErrNotMember(err *apperr.Error, userEmail string, orgName string) *apperr.Error {
	return err.WithMessagef("user with email %q is not a member of organization %q", userEmail, orgName)
}
