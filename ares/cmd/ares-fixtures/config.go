package main

import "bitbucket.org/unchain/ares/pkg/3p/sql"

// Config contains the configuration needed for an ares server
type Config struct {
	SQL *sql.Config
}
