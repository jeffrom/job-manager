module github.com/jeffrom/job-manager

go 1.13

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-multierror v1.1.0
	github.com/jeffrom/job-manager/jobclient v0.0.0-00010101000000-000000000000
	github.com/kr/pretty v0.1.0 // indirect
	github.com/qri-io/jsonschema v0.2.0
	github.com/rs/zerolog v1.19.0
	github.com/satori/go.uuid v1.2.0
	github.com/tdewolff/minify v2.3.6+incompatible // indirect
	github.com/tdewolff/minify/v2 v2.9.1
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/jeffrom/job-manager/jobclient => ./jobclient
