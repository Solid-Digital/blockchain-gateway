module github.com/unchain/tbg-nodes-auth

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/davecgh/go-spew v1.1.1
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/go-chi/chi v4.0.4+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/google/wire v0.4.0
	github.com/lib/pq v1.3.0
	github.com/spf13/cast v1.3.1 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/coreos/etcd v3.3.15+incompatible // indirect
	github.com/unchainio/interfaces v0.2.1
	github.com/unchainio/pkg v0.22.1
	github.com/volatiletech/inflect v0.0.0-20170731032912-e7201282ae8d // indirect
	github.com/volatiletech/null v8.0.0+incompatible
	github.com/volatiletech/sqlboiler v3.6.1+incompatible
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/spf13/viper v1.2.2 => github.com/unchainio/viper v1.2.2-0.20190712174521-9bf201c29832

replace github.com/BurntSushi/toml v0.3.1 => github.com/unchain/toml v0.4.0
