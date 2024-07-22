package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190402195034, Down20190402195034)
}

func Up20190402195034(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func Down20190402195034(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
