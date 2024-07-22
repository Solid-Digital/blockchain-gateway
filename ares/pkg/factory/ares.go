package factory

import (
	"sync"

	"bitbucket.org/unchain/ares/gen/wire"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/unchainio/pkg/xconfig"
)

var aresInstance *ares.Server
var cleanupInstance func()
var once sync.Once

func (f *Factory) AresFactory() (*ares.Server, func()) {
	once.Do(func() {
		aresInstance, cleanupInstance = f.aresFactory()
	})

	return aresInstance, cleanupInstance
}

func (f *Factory) aresFactory() (*ares.Server, func()) {
	// use test config, this will allow access to all services
	var err error
	cfg := new(wire.Config)

	err = xconfig.Load(
		cfg,
		xconfig.FromPaths("../../config/test/config.toml"),
	)

	f.suite.Require().NoError(err)

	// create ares server, this will initialize all services
	var ares *ares.Server
	var cleanup func()

	// FIXME(e-nikolov) Disable the kubernetes tests in the pipelines.
	//   For the love of all that is holy, please move away from the bitbucket pipelines.
	if testhelper.InBitBucket() {
		ares, cleanup, err = wire.AresForBitBucket(f.suite.T(), f.Metadata(), cfg, f.Logger())
	} else {
		ares, cleanup, err = wire.AresForTests(f.suite.T(), f.Metadata(), cfg, f.Logger())
	}

	f.suite.Require().NoError(err)

	return ares, cleanup
}
