//+build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/unchain/tbg-nodes-auth/pkg/3p/sql"
	"github.com/unchain/tbg-nodes-auth/pkg/nodeauth"
	"github.com/unchainio/interfaces/logger"
)

// Wire initializer for NodeAuth. This is the spec that `wire ./pkg/ares` uses to generate the dependency injection code
// nolint[golint]
func NodeAuth(
	meta *nodeauth.Metadata,
	cfg *Config,
	log logger.Logger,
) (*nodeauth.Server, func(), error) {
	panic(wire.Build(
		ConfigSet,
		sql.NewDB,
		nodeauth.NewServer,
	))
}
