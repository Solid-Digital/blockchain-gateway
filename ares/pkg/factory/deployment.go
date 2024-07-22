package factory

import (
	"context"
	"database/sql"
	"fmt"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/unchainio/pkg/errors"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/Pallinder/go-randomdata"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) Deployment(create bool, user *orm.User, pipeline *orm.Pipeline, configuration *orm.Configuration, environment *orm.Environment) *orm.Deployment {
	deployment := &orm.Deployment{
		PipelineID:      pipeline.ID,
		ConfigurationID: configuration.ID,
		EnvironmentID:   environment.ID,
		Replicas:        int64(randomdata.Number(5)),
		URL:             fmt.Sprintf("http://%s", randomdata.IpV4Address()),
		FullName:        testhelper.Randumb("deployment"),
		CreatedByID:     user.ID,
		UpdatedByID:     user.ID,
	}

	if !create {
		return deployment
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return deployment.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return deployment
}

func (f *Factory) DeploymentFromService() (*orm.Organization, *orm.Pipeline, *orm.Environment, *orm.Deployment) {
	fileName := testhelper.Randumb("http.endpoint.so")

	// only works when called from the right folder
	f.File("../../test/fixtures/binary/http.endpoint.so", fileName)

	org, user := f.OrganizationAndUser(true)
	pipeline := f.Pipeline(true, org, user)
	config := f.Configuration(true, org, user, pipeline)
	env := f.Environment(org, user)
	f.TriggerConfigurationWithFile(false, true, config, org, fileName)
	f.BaseConfiguration(false, true, config, org)

	params := &dto.DeployConfigurationRequest{
		ConfigurationRevision: testhelper.Int64Ptr(config.Revision),
		Replicas:              testhelper.Int64Ptr(3),
	}

	dtoUser := f.ORMToDTOUser(user)

	d, appErr := f.ares.PipelineService.DeployConfiguration(params, org.Name, pipeline.Name, env.Name, dtoUser)
	xrequire.NoError(f.suite.T(), appErr)

	var deployment *orm.Deployment
	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		var err error

		deployment, err = orm.Deployments(
			orm.DeploymentWhere.ID.EQ(*d.ID)).One(ctx, tx)

		if err != nil {
			return errors.Wrapf(err, "No deployment found for pipeline %q/%d and env %q/%d", pipeline.Name, pipeline.ID, env.Name, env.ID)
		}

		return nil
	})

	f.suite.Require().NoError(err)

	return org, pipeline, env, deployment
}

func (f *Factory) DeploymentFromServiceForOrg(org *orm.Organization) (*orm.Pipeline, *orm.Environment, *orm.Deployment) {
	fileName := testhelper.Randumb("http.endpoint.so")

	// only works when called from the right folder
	f.File("../../test/fixtures/binary/http.endpoint.so", fileName)

	user := f.UserFromOrg(org)
	pipeline := f.Pipeline(true, org, user)
	config := f.Configuration(true, org, user, pipeline)
	env := f.Environment(org, user)
	f.TriggerConfigurationWithFile(false, true, config, org, fileName)
	f.BaseConfiguration(false, true, config, org)

	params := &dto.DeployConfigurationRequest{
		ConfigurationRevision: testhelper.Int64Ptr(config.Revision),
		Replicas:              testhelper.Int64Ptr(3),
	}

	dtoUser := f.ORMToDTOUser(user)

	d, appErr := f.ares.PipelineService.DeployConfiguration(params, org.Name, pipeline.Name, env.Name, dtoUser)
	xrequire.NoError(f.suite.T(), appErr)

	var deployment *orm.Deployment
	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		var err error

		deployment, err = orm.Deployments(
			orm.DeploymentWhere.ID.EQ(*d.ID)).One(ctx, tx)

		if err != nil {
			return errors.Wrapf(err, "No deployment found for pipeline %q/%d and env %q/%d", pipeline.Name, pipeline.ID, env.Name, env.ID)
		}

		return nil
	})

	f.suite.Require().NoError(err)

	return pipeline, env, deployment
}
