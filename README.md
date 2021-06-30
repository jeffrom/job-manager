# job-manager

job-manager is a server that coordinates the execution of distributed jobs over HTTP.

The gist is that consumers poll the server for jobs, and send job status back upon completion. Checkins can be used to annotate the job with data pre-completion. A command-line tool, jobctl, is provided for queue administration. Labels, JSON Schema validation, and a gitops-style queue administration support complex workflows across multiple development teams.

I made this initially to experiment with some higher-level features in a job queue, but it should be capable of comparable performance to other postgres-backed job queues. In most cases, the job server itself is not a performance bottleneck compared to the backend, however multiple replicas can be run for high availability.

## install

Using go:

```bash
$ go get github.com/jeffrom/job-manager/...
```

Via docker: `docker pull jeffmartin1117/job-manager`

Or download a github release.

### migrations

To run postgresql migrations:

```sh
$ jobctl migrate
```

Migrations are implemented using [golang-migrate](https://github.com/golang-migrate/migrate). Note that this command requests a job-manager server to execute the migration using its configured credentials. This requires additional permissions in Postgresql.

## features

* cli controller
* straightforward Rest API (protobuf support)
* job server is stateless / scales horizontally
* implement your own backend, comes with postgresql
* in-memory backend for development & testing purposes (not really working right now)
* claim windows: only dequeue to consumers with matching claims for a configurable duration
* check ins
* store result data
* queue labels
* versioned queue configuration
* easily update queue configurations via cicd with `jobctl apply`
* json schema validation for job arguments, data, results
* the usual job queue features: retries, exponential backoff, durability
* graceful shutdown
* not fast, predictable 😎

## clients

There is a go client:

```bash
$ go get github.com/jeffrom/job-manager/mjob
```

## develop

To start a development server on your laptop:

```bash
$ make dev
```

Run tests and static analysis:

```bash
$ make test
$ make lint
```

Point jobctl at the local server dev proxy:

```
$ export HOST=:4000
$ jobctl stats
```
