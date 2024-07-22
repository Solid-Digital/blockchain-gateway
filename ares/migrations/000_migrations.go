package migration

import (
	"fmt"
)

func init() {
	fmt.Println("This is ares-migrate")
	fmt.Println(`Create new migrations with: "ares-migrate create add_some_column sql"`)
	fmt.Println(`See: https://github.com/pressly/goose`)
	fmt.Println()
}
