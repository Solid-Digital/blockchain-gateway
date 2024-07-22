//+build wireinject

package wire

import (
	"bitbucket.org/unchain/ares/gen/api"
	"bitbucket.org/unchain/ares/pkg/3p/aws"
	"bitbucket.org/unchain/ares/pkg/3p/casbin"
	"bitbucket.org/unchain/ares/pkg/3p/docker"
	"bitbucket.org/unchain/ares/pkg/3p/elastic"
	"bitbucket.org/unchain/ares/pkg/3p/harbor"
	"bitbucket.org/unchain/ares/pkg/3p/kubernetes"
	"bitbucket.org/unchain/ares/pkg/3p/mail"
	"bitbucket.org/unchain/ares/pkg/3p/mario"
	"bitbucket.org/unchain/ares/pkg/3p/redis"
	"bitbucket.org/unchain/ares/pkg/3p/s3"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/auth"
	"bitbucket.org/unchain/ares/pkg/component"
	"bitbucket.org/unchain/ares/pkg/hello"
	"bitbucket.org/unchain/ares/pkg/http"
	"bitbucket.org/unchain/ares/pkg/organization"
	"bitbucket.org/unchain/ares/pkg/pipeline"
	"bitbucket.org/unchain/ares/pkg/subscription"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
)

// Wire initializer for Ares. This is the spec that `wire ./pkg/ares` uses to generate the dependency injection code
// nolint[golint]
func Ares(
	meta *ares.Metadata,
	cfg *Config,
	log logger.Logger,
) (*ares.Server, func(), error) {
	panic(wire.Build(
		ConfigSet,
		redis.ClientSet,
		elastic.NewClient,
		casbin.Set,
		sql.NewDB,
		s3.FileStoreSet,
		mail.MailerSet,
		docker.ContainerServiceSet,
		mario.ImageBuilderSet,
		kubernetes.ServiceSet,
		harbor.NewClient,
		aws.ClientSet,
		http.AuthenticationMiddleware,
		http.AuthorizationMiddleware,
		pipeline.ServiceSet,
		organization.ServiceSet,
		hello.ServiceSet,
		auth.ServiceSet,
		component.ServiceSet,
		subscription.ServiceSet,
		http.HelloSet,
		http.AuthSet,
		http.PipelineSet,
		http.OrganizationSet,
		http.ComponentSet,
		http.InnerMiddlewareProvider,
		http.MiddlewareProvider,
		api.Handlers{},
		ares.NewServer,
	))
}
