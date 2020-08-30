module github.com/jeffrom/job-manager

go 1.13

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/jeffrom/job-manager/jobclient v0.0.0-00010101000000-000000000000
	github.com/kr/pretty v0.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/jeffrom/job-manager/jobclient => ./jobclient
