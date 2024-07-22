//+build wireinject

package wire

import (
	"bitbucket.org/unchain/ares/pkg/3p/casbin"
	"bitbucket.org/unchain/ares/pkg/3p/harbor"
	"bitbucket.org/unchain/ares/pkg/3p/mail"
	"bitbucket.org/unchain/ares/pkg/3p/redis"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/auth"
	"bitbucket.org/unchain/ares/pkg/bootstrap"
	"bitbucket.org/unchain/ares/pkg/organization"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
)

// Wire initializer for Ares. This is the spec that `wire ./pkg/ares` uses to generate the dependency injection code
// nolint[golint]
func Bootstrap(
	meta *ares.Metadata,
	cfg *Config,
	log logger.Logger,
) (*bootstrap.Service, func(), error) {
	panic(wire.Build(
		BootstrapConfigSet,
		redis.ClientSet,
		casbin.Set,
		sql.NewDB,
		mail.MailerSet,
		harbor.NewClient,
		organization.ServiceSet,
		NilAWSProvider,
		auth.ServiceSet,
		bootstrap.New,
	))
}
