# Job Worker Design Doc

Design doc for a prototype job worker service that provides an API to run arbitrary Linux processes.

## Level 2 requirements

- Library: Worker library with methods to start/stop/query status and get the output of a job.
- API:
    - HTTPS API to start/stop/get status of a running process.
    - Use mTLS authentication and verify client certificate. Set up strong set of cipher suites for TLS and good crypto
      setup for certificates. Do not use any other authentication protocols on top of mTLS.
    - Use a simple authorization scheme.
- CLI:
    - CLI should be able to connect to worker service and schedule several jobs.
    - CLI should be able to query result of the job execution and fetch its output.

## Job object

- `ID` unique identifier of a job
- `Command` linux command to be executed
- `Status` current status of the job
- `Output` output of the job process
- `Error` exit error or, termination error
- `Timeout` timeout in seconds to stop the job

## Job status

- `Created` Received a job creation request (short-lived since we immediately start the job instead of queueing)
- `Running` Job process started
- `Running_Timeout` Failed to kill process after timeout
- `Running_Force_Stop` Failed to kill process after force stop requested
- `Exited` Job finished successfully
- `Exited_Error` Job finished with error
- `Exited_Timeout` Job exited after timeout
- `Exited_Force_Stop` Job exited after force stop requested

## Job handling

- Wrap the linux command in `bash`
- Start the process using `exec.Cmd` in separate go routine
- Utilize go channels to communicate with the go routine
- The go routine should exit when job finishes or timeout reached or force stop requested

## Library

- `CreateJob(command, timeout)`
    - Generate a unique ID
    - Initialize a job object using the ID
    - Save the reference to the job into memory for future access
    - Start the job in a separate go routine
    - Return the job object
- `StopJob(jobID)`
    - Fetch the job object from memory by job ID
    - Check the current status of the job
    - Stop the job process if it is running
    - Send confirmation on success or, error on termination failure
- `PurgeJob(jobID)`
    - Fetch the job object from memory by job ID
    - Check the current status of the job
    - Remove the job from memory if it is stopped
    - Return error if the job is not found or job is running
- `GetJobStatus(jobID)`
    - fetch the job object from memory by job ID
    - return the details of a job along with output and error

## API

- `POST /create {"command": "pwd", "timeout": "1"}`
- `GET /status/{id}`
- `PUT /stop/{id}`
- `PUT /purge/{id}`

## CLI

- `job-worker [server/client]` combined CLI for running server and client
- `job-worker server [--port] [--cacert] [--cert] [--key]` run the API server
- `job-worker client [create/stop/purge/status] [--address] [--cacert] [--cert] [--key]` run client to
  communicate with the API server
- `job-worker client create [--command] [--timeout]` request a new job
- `job-worker client stop [--id]` stop a job
- `job-worker client purge [--id]` remove a job from memory
- `job-worker client status [--id]` fetch status of a job

## mTLS

- Use `openssl` to generate self-signed certificates
- Use `RSA 2048` key and `sha256` digest
- Shell script to generate certificates
- Steps:
    - Generate a root CA certificate (`ca.crt`, `ca.key`)
    - Generate server certificate signed by the CA (`server.crt`, `server.key`)
    - Generate client certificate signed by the CA (`client.crt`, `client.key`)

## Authorization

- Simple RBAC style authorization
- Permissions:
    - Read and write: Required for creating, stopping, and purging a job
    - Read only: Required for querying status of a job
- Two hard coded users to mock database
    - Username `alice`: has read and write permissions
    - Username `bob`: has read only permission
- Pass username through client certificate using `Common Name` i.e. each user will have separate client certificate

## Tests

- Test the library
- Test RBAC authorization
- Test TLS connection

## Third party libraries

- [spf13/cobra](https://github.com/spf13/cobra) for CLI
- [gorilla/mux](https://github.com/gorilla/mux) for routing
- [mikespook/gorbac](github.com/mikespook/gorbac) for simple RBAC implementation
- [rs/xid](github.com/rs/xid) for generating globally unique job ID
- [sirupsen/logrus](github.com/sirupsen/logrus) for logging

## TODOs

- Terminate all jobs before shutting down the server
- Support graceful termination of jobs along with force stop
- Retry job termination if failed
- Use DB to store users
- Save job details in DB instead of memory
- Use job queue instead of directly starting the job then schedule job based on resource availability
- Job chaining (workflow) using dependency graph
- Check for vulnerabilities before running a job
