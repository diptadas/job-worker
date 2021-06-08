# Job Worker

A prototype job worker service that provides an API to run arbitrary Linux processes.

## Library usage example

- Create a job

```go
request := CreateJobRequest{
    Command: "echo",
    Args:    []string{"hello world"},
}

job, err := CreateJob(request)
```

- Stop a job

```go
job, err := StopJob("c2vvjbucie6ju0ou7j5g")
```

- Get job status

```go
job, err := GetJobStatus("c2vvjbucie6ju0ou7j5g")
```

## Test the library

```shell
$ go test -v ./...
```
