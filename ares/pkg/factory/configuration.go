package factory

import (
	"context"
	"database/sql"
	"math/rand"
	"sort"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"github.com/volatiletech/null"

	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) Configuration(create bool, org *orm.Organization, user *orm.User, pipeline *orm.Pipeline) *orm.Configuration {
	configuration := &orm.Configuration{
		CreatedByID:    user.ID,
		UpdatedByID:    user.ID,
		PipelineID:     pipeline.ID,
		OrganizationID: org.ID,
		Revision:       int64(randomdata.Number(1000)),
	}

	if !create {
		return configuration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return configuration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return configuration
}

func (f *Factory) DraftConfiguration(create bool, org *orm.Organization, user *orm.User) *orm.DraftConfiguration {
	draftConfiguration := &orm.DraftConfiguration{
		CreatedByID:    user.ID,
		UpdatedByID:    user.ID,
		OrganizationID: org.ID,
		Revision:       int64(randomdata.Number(1000)),
	}

	if !create {
		return draftConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return draftConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return draftConfiguration
}

func (f *Factory) BaseConfiguration(public, create bool, configuration *orm.Configuration, org *orm.Organization) *orm.BaseConfiguration {
	baseVersion := f.BaseVersionForOrg(org, public, create)

	baseConfiguration := &orm.BaseConfiguration{
		ConfigurationID: configuration.ID,
		VersionID: null.Int64{
			Int64: baseVersion.ID,
			Valid: true,
		},
		Config: testhelper.RandomTOML(),
	}

	if !create {
		return baseConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return baseConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return baseConfiguration
}

func (f *Factory) TriggerConfiguration(public, create bool, configuration *orm.Configuration, org *orm.Organization) *orm.TriggerConfiguration {
	triggerVersion := f.TriggerVersionForOrg(org, public, create)

	triggerConfiguration := &orm.TriggerConfiguration{
		ConfigurationID: configuration.ID,
		VersionID: null.Int64{
			Int64: triggerVersion.ID,
			Valid: true,
		},
		Config:        testhelper.RandomTOML(),
		MessageConfig: null.JSONFrom([]byte(testhelper.RandomJSON())),
	}

	if !create {
		return triggerConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return triggerConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return triggerConfiguration
}

func (f *Factory) TriggerConfigurationWithFile(public, create bool, configuration *orm.Configuration, org *orm.Organization, fileName string) *orm.TriggerConfiguration {
	triggerVersion := f.TriggerVersionForOrgWithFile(org, public, create, fileName)

	triggerConfiguration := &orm.TriggerConfiguration{
		ConfigurationID: configuration.ID,
		VersionID: null.Int64{
			Int64: triggerVersion.ID,
			Valid: true,
		},
		Config:        testhelper.RandomTOML(),
		MessageConfig: null.JSONFrom([]byte(testhelper.RandomJSON())),
	}

	if !create {
		return triggerConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return triggerConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return triggerConfiguration
}

func (f *Factory) ActionConfiguration(public, create bool, configuration *orm.Configuration, org *orm.Organization) *orm.ActionConfiguration {
	actionVersion := f.ActionVersionForOrg(org, public, create)

	actionConfiguration := &orm.ActionConfiguration{
		ConfigurationID: configuration.ID,
		VersionID:       actionVersion.ID,
		Index:           int64(randomdata.Number(10)),
		Config:          testhelper.RandomTOML(), // this should be randumbized
		MessageConfig:   null.JSONFrom([]byte(testhelper.RandomJSON())),
		Name:            testhelper.Randumb(randomdata.SillyName()),
	}

	if !create {
		return actionConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return actionConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return actionConfiguration
}

func (f *Factory) ManyActionCongigurations(configuration *orm.Configuration, org *orm.Organization, amount int) []*orm.ActionConfiguration {
	f.suite.Require().NotNil(f.ares)

	cfgs := make([]*orm.ActionConfiguration, amount)

	indices := rand.Perm(amount)
	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {

		for i := 0; i < amount; i++ {
			actionConfiguration := f.ActionConfiguration(false, true, configuration, org)
			actionConfiguration.Index = int64(indices[i])

			_, err := actionConfiguration.Update(ctx, tx, boil.Infer())
			if err != nil {
				return err
			}

			cfgs[i] = actionConfiguration
		}

		return nil

	})

	f.suite.Require().NoError(err)

	sort.Slice(cfgs, func(i, j int) bool {
		return cfgs[i].Index < cfgs[j].Index
	})

	return cfgs
}

func (f *Factory) BaseDraftConfiguration(public, create bool, draftConfiguration *orm.DraftConfiguration, org *orm.Organization) *orm.BaseDraftConfiguration {
	baseVersion := f.BaseVersionForOrg(org, public, create)

	baseDraftConfiguration := &orm.BaseDraftConfiguration{
		DraftConfigurationID: draftConfiguration.ID,
		VersionID: null.Int64{
			Int64: baseVersion.ID,
			Valid: true,
		},
		Config: testhelper.RandomTOML(),
	}

	if !create {
		return baseDraftConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return baseDraftConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return baseDraftConfiguration
}

func (f *Factory) TriggerDraftConfiguration(public, create bool, draftConfiguration *orm.DraftConfiguration, org *orm.Organization) *orm.TriggerDraftConfiguration {
	triggerVersion := f.TriggerVersionForOrg(org, public, create)

	triggerDraftConfiguration := &orm.TriggerDraftConfiguration{
		DraftConfigurationID: draftConfiguration.ID,
		VersionID: null.Int64{
			Int64: triggerVersion.ID,
			Valid: true,
		},
		Config:        testhelper.RandomTOML(),
		MessageConfig: null.JSONFrom([]byte(testhelper.RandomJSON())),
	}

	if !create {
		return triggerDraftConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return triggerDraftConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return triggerDraftConfiguration
}

func (f *Factory) ActionDraftConfiguration(public, create bool, draftConfiguration *orm.DraftConfiguration, org *orm.Organization) *orm.ActionDraftConfiguration {
	actionVersion := f.ActionVersionForOrg(org, public, create)

	actionDraftConfiguration := &orm.ActionDraftConfiguration{
		DraftConfigurationID: draftConfiguration.ID,
		VersionID:            actionVersion.ID,
		Name:                 testhelper.Randumb(randomdata.SillyName()),
		Config:               testhelper.RandomTOML(),
		MessageConfig:        null.JSONFrom([]byte(testhelper.RandomJSON())),
		Index:                int64(randomdata.Number(10)),
	}

	if !create {
		return actionDraftConfiguration
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return actionDraftConfiguration.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return actionDraftConfiguration
}

func (f *Factory) ManyActionDraftConfigurations(draftConfiguration *orm.DraftConfiguration, org *orm.Organization, amount int) []*orm.ActionDraftConfiguration {
	f.suite.Require().NotNil(f.ares)

	cfgs := make([]*orm.ActionDraftConfiguration, amount)

	indices := rand.Perm(amount)
	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {

		for i := 0; i < amount; i++ {
			actionDraftConfiguration := f.ActionDraftConfiguration(false, true, draftConfiguration, org)
			actionDraftConfiguration.Index = int64(indices[i])

			_, err := actionDraftConfiguration.Update(ctx, tx, boil.Infer())
			if err != nil {
				return err
			}

			cfgs[i] = actionDraftConfiguration
		}

		return nil

	})

	f.suite.Require().NoError(err)

	sort.Slice(cfgs, func(i, j int) bool {
		return cfgs[i].Index < cfgs[j].Index

	})

	return cfgs
}
