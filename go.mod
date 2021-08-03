module github.com/jeffrom/job-manager

go 1.16

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-playground/form v3.1.4+incompatible
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/imdario/mergo v0.3.12
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jeffrom/job-manager/mjob v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.3.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.11.0
	github.com/qri-io/jsonschema v0.2.1
	github.com/rs/zerolog v1.23.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.2.1
	google.golang.org/protobuf v1.27.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/jeffrom/job-manager/mjob => ./mjob
