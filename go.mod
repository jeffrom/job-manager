module github.com/jeffrom/job-manager

go 1.17

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-playground/form v3.1.4+incompatible
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/imdario/mergo v0.3.12
	github.com/jackc/pgx/v4 v4.15.0
	github.com/jeffrom/job-manager/mjob v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.3.4
	github.com/jteeuwen/go-bindata v3.0.7+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.14.0
	github.com/qri-io/jsonschema v0.2.1
	github.com/rs/zerolog v1.26.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.4.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tdewolff/minify/v2 v2.10.0 // indirect
	github.com/tdewolff/parse/v2 v2.5.27 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/jeffrom/job-manager/mjob => ./mjob
