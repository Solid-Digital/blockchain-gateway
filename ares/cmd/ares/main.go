package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	_ "github.com/lib/pq"

	"bitbucket.org/unchain/ares/gen/wire"
	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/unchainio/pkg/iferr"
	"github.com/unchainio/pkg/xconfig"
	"github.com/unchainio/pkg/xlogger"
)

var version = ""
var branch = ""
var builder = ""
var buildDate = ""
var v *bool

func main() {
	meta := &ares.Metadata{
		Name:      "ares",
		Version:   version,
		Branch:    branch,
		Builder:   builder,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
	}
	flags()

	var err error

	cfg := new(wire.Config)
	info := new(xconfig.Info)

	errs := xconfig.Load(
		cfg,
		xconfig.FromPathFlag("cfg", "config/dev/config.toml"),
		xconfig.FromEnv(),
		xconfig.GetInfo(info),
	)

	fmt.Print(ares.Logo)
	fmt.Printf("%s\n", meta)

	HandleFlags()

	log, err := xlogger.New(cfg.Logger)
	iferr.Exit(err)

	iferrlog, err := xlogger.New(&xlogger.Config{
		Level:       cfg.Logger.Level,
		Format:      cfg.Logger.Format,
		HideFName:   cfg.Logger.HideFName,
		CallerDepth: 4,
	})
	iferr.Exit(err)

	iferr.Default, err = iferr.New(iferr.WithLogger(iferrlog))
	iferr.Exit(err)

	log.Printf("Attempted to load configs from %+v", info.Paths)
	iferr.Warn(errs)

	ares, cleanup, err := wire.Ares(meta, cfg, log)
	iferr.Exit(err)
	defer cleanup()

	h, err := ares.Handler()
	iferr.Exit(err)

	log.Printf("Start listening on port %s with log level: %s", cfg.URL, cfg.Logger.Level)
	log.Fatal(http.ListenAndServe(cfg.URL, h))
}

func flags() {
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
	v = flag.Bool("version", false, "Print the version of ares")

}

func HandleFlags() {
	flag.Parse()

	if *v || (len(os.Args) == 2 && os.Args[1] == "version") {
		os.Exit(0)
	}
}
