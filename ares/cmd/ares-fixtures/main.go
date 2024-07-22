package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/unchainio/pkg/iferr"

	"github.com/unchainio/pkg/xconfig"

	"github.com/pressly/goose"

	// Init DB drivers.
	_ "bitbucket.org/unchain/ares/fixtures"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/ziutek/mymysql/godrv"
)

const MigrationsPath = "fixtures"

var (
	flags = flag.NewFlagSet("ares-fixtures", flag.ExitOnError)
	cfg   = flags.String("cfg", "config/dev/config.toml", "which configuration file to use")
)

func main() {
	goose.SetTableName("goose_fixture_version")

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	config := new(Config)

	errs := xconfig.Load(
		config,
		xconfig.FromPathFlag("cfg", "config/dev/config.toml"),
		xconfig.FromEnv(),
	)
	iferr.Warn(errs)

	cfg := config.SQL

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	switch args[0] {
	case "create":
		if err := goose.Run("create", nil, MigrationsPath, args[1:]...); err != nil {
			log.Fatalf("ares-fixtures run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, MigrationsPath); err != nil {
			log.Fatalf("ares-fixtures run: %v", err)
		}
		return
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	switch cfg.Driver {
	case "postgres", "mysql", "sqlite3", "redshift":
		if err := goose.SetDialect(cfg.Driver); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("%q driver not supported\n", cfg.Driver)
	}

	switch cfg.ConnectionString {
	case "":
		log.Fatalf("Empty connection string not supported\n")
	default:
	}

	if cfg.Driver == "redshift" {
		cfg.Driver = "postgres"
	}

	db, err := sql.Open(cfg.Driver, cfg.ConnectionString)
	if err != nil {
		log.Fatalf("Connection=%q: %v\n", cfg.ConnectionString, err)
	}

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	if err := goose.Run(command, db, MigrationsPath, arguments...); err != nil {
		log.Fatalf("ares-fixtures run: %v", err)
	}
}

func usage() {
	log.Print(usagePrefix)
	flags.PrintDefaults()
	log.Print(usageCommands)
}

var (
	usagePrefix = `Usage: ares-fixtures [OPTIONS] DRIVER DBSTRING COMMAND

Drivers:
    postgres
    mysql
    sqlite3
    redshift

Examples:
    ares-fixtures status
    ares-fixtures create init sql
    ares-fixtures create add_some_column sql
    ares-fixtures create fetch_user_data go
    ares-fixtures up


Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
	fix                  Apply sequential ordering to migrations
`
)
