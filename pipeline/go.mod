module github.com/unchain/pipeline

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/coreos/etcd v3.3.15+incompatible // indirect
	github.com/emersion/go-imap v1.0.6
	github.com/emersion/go-imap-move v0.0.0-20190710073258-6e5a51a5b342
	github.com/emersion/go-message v0.13.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/go-plugin v1.4.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/common v0.6.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/unchainio/interfaces v0.2.1
	github.com/unchainio/pkg v0.22.1
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	google.golang.org/grpc v1.27.1 // indirect
)

replace github.com/spf13/viper v1.2.2 => github.com/unchainio/viper v1.2.2-0.20190712174521-9bf201c29832
