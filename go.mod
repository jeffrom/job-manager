module github.com/jeffrom/job-manager

go 1.13

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-playground/form v3.1.4+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-multierror v1.1.0
	github.com/imdario/mergo v0.3.11
	github.com/jeffrom/job-manager/jobclient v0.0.0-00010101000000-000000000000
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/qri-io/jsonschema v0.2.0
	github.com/rs/zerolog v1.19.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/tdewolff/minify/v2 v2.9.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/jeffrom/job-manager/jobclient => ./jobclient
