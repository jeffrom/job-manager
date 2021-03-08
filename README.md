# job-manager

job-manager is a server that coordinates the execution of distributed jobs over HTTP.

The gist is that consumers poll the server for jobs, and send job status back upon completion. Checkins can be used to annotate the job with data pre-completion. A command-line tool, jobctl, is provided for queue administration. Labels, JSON Schema validation, and a gitops-style queue administration are supported to better support complex workflows across multiple development teams.

I made this initially to experiment with some higher-level features in a job queue, but it should handle a decent amount of scale by now. In most cases, the job server itself should rarely be a bottleneck, but multiple replicas can be run for high availability. Typically, the backend (currently postgresql) is going to be the bottleneck in most workloads.

*WARNING* it's still buggy, and I think I'd want to rework the data model before taking this to 1.0.

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

Migrations are implemented using [golang-migrate](https://github.com/golang-migrate/migrate).

## features

* cli controller
* job server and consumers can be scaled horizontally
* claim windows: only dequeue to consumers with matching claims for a configurable duration
* check ins
* store result data
* queue labels
* versioned queue configuration
* easily update queue configurations via cicd with `jobctl apply`
* json schema validation for job arguments, data, results
* exponential backoff
* graceful shutdown
* not very fast, and probably will never be as fast as the average job system

## clients

For now, just go:

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
