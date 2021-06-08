package lib

import (
	"bytes"
	"os/exec"
	"sync"
)

const (
	JobRunning = "RUNNING" // job running
	JobExited  = "EXITED"  // job finished successfully or finished with error or force stopped
)

// Job contains the details of a job.
type Job struct {
	ID      string           `json:"id"`      // unique autogenerated identifier of the job
	Status  string           `json:"status"`  // status of the job
	Output  string           `json:"output"`  // populated from OutputBuffer when job status is queried
	Request CreateJobRequest `json:"request"` // store CreateJobRequest for future reference

	// Error is nil if job finishes successfully
	// Error wraps the exit code if job finishes with error
	// Error contains "signal: killed" message if job is force stopped
	Error error `json:"error"`

	// for internal use only
	cmd          *exec.Cmd       // save the reference to the execution process to handle force stop
	outputBuffer *bytes.Buffer   // contains stdout and stderr of the job process
	waitForExit  *sync.WaitGroup // used by StopJob to wait for cmd.Wait to finish
}

// CreateJobRequest defines a request to create new job.
type CreateJobRequest struct {
	Command string   `json:"command"` // linux command to be executed
	Args    []string `json:"args"`    // arguments to the command
}
