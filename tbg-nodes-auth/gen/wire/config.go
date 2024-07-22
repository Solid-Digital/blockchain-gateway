package wire

import (
	"github.com/google/wire"
	"github.com/unchain/tbg-nodes-auth/pkg/3p/sql"
	"github.com/unchainio/pkg/xlogger"
)

var ConfigSet = wire.FieldsOf(new(*Config), "SQL")

// Config contains the configuration needed for an ares server
type Config struct {
	URL    string
	SQL    *sql.Config
	Logger *xlogger.Config
}
