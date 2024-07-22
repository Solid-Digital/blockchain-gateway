package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/unchainio/pkg/iferr"
	"github.com/unchainio/pkg/xlogger"

	_ "github.com/lib/pq"

	"github.com/unchain/tbg-nodes-auth/gen/wire"
	"github.com/unchainio/pkg/xconfig"

	"github.com/unchain/tbg-nodes-auth/pkg/nodeauth"
)

var version string
var branch string
var builder string
var buildDate string

func main() {
	meta := &nodeauth.Metadata{
		Name:      "tbg-nodes-auth",
		Version:   version,
		Branch:    branch,
		Builder:   builder,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
	}

	var err error

	cfg := new(wire.Config)
	info := new(xconfig.Info)

	errs := xconfig.Load(
		cfg,
		xconfig.FromPathFlag("cfg", "config/dev/config.toml"),
		xconfig.FromEnv(),
		xconfig.GetInfo(info),
	)

	fmt.Printf("%s\n", meta)

	log, err := xlogger.New(cfg.Logger)
	iferr.Exit(err)

	iferrlog, err := xlogger.New(&xlogger.Config{
		Level:       cfg.Logger.Level,
		Format:      cfg.Logger.Format,
		CallerDepth: 4,
	})
	iferr.Exit(err)

	iferr.Default, err = iferr.New(iferr.WithLogger(iferrlog))
	iferr.Exit(err)

	log.Printf("Attempted to load configs from %+v", info.Paths)
	iferr.Warn(errs)

	nodeAuth, cleanup, err := wire.NodeAuth(meta, cfg, log)
	iferr.Exit(err)
	defer cleanup()

	fmt.Printf("Running tbg node auth server\n")
	log.Fatal(http.ListenAndServe(":8080", nodeAuth.Handler()))
}
