# job-manager

job-manager is a server that coordinates the execution of distributed jobs.

## install

Using go:

```bash
$ go get github.com/jeffrom/job-manager/...
```

Via docker: `docker pull jeffmartin1117/job-manager`

Or download a github release.

### migrations

```sh
$ jobctl migrate
```

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
* not very fast :)

## clients

For now, just go:

```bash
$ go get github.com/jeffrom/job-manager/mjob
```
