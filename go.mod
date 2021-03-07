module github.com/jeffrom/job-manager

go 1.13

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-playground/form v3.1.4+incompatible
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/go-cmp v0.5.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/imdario/mergo v0.3.11
	github.com/jackc/pgx/v4 v4.10.1
	github.com/jeffrom/job-manager/mjob v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.9.0 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/qri-io/jsonschema v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.1.3
	google.golang.org/protobuf v1.25.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/jeffrom/job-manager/mjob => ./mjob
