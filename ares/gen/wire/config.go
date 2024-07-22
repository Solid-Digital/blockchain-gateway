package wire

import (
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
	"bitbucket.org/unchain/ares/pkg/auth"
	"bitbucket.org/unchain/ares/pkg/bootstrap"
	"bitbucket.org/unchain/ares/pkg/pipeline"
	"bitbucket.org/unchain/ares/pkg/subscription"
	"github.com/google/wire"
	"github.com/unchainio/pkg/xlogger"
)

var ConfigSet = wire.FieldsOf(new(*Config), "Logger", "Redis", "SQL", "Elastic", "S3", "JWT", "Casbin", "Mail", "Pipelines", "Docker", "Mario", "Harbor", "Kubernetes", "AWS", "Subscription")
var ConfigForTestsSet = wire.FieldsOf(new(*Config), "Logger", "Redis", "SQL", "Elastic", "S3", "JWT", "Casbin", "Mail", "Pipelines", "Docker", "Mario", "Harbor", "Kubernetes", "Subscription")
var ConfigForBitBucketSet = wire.FieldsOf(new(*Config), "Logger", "Redis", "SQL", "Elastic", "S3", "JWT", "Casbin", "Mail", "Pipelines", "Docker", "Mario", "Harbor", "Subscription")
var BootstrapConfigSet = wire.FieldsOf(new(*Config), "Bootstrap", "Redis", "SQL", "JWT", "Casbin", "Mail", "Harbor")

// Config contains the configuration needed for an ares server
type Config struct {
	URL          string
	Bootstrap    *bootstrap.Config
	Redis        *redis.Config
	SQL          *sql.Config
	Elastic      *elastic.Config
	S3           *s3.Config
	Logger       *xlogger.Config
	JWT          *auth.Config
	Casbin       *casbin.Config
	Mail         *mail.Config
	Pipelines    *pipeline.Config
	Docker       *docker.Config
	Mario        *mario.Config
	Harbor       *harbor.Config
	Kubernetes   *kubernetes.Config
	AWS          *aws.Config
	Subscription *subscription.Config
}
