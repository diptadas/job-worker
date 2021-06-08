package lib

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

// allJobs stores the details of all jobs in memory.
// It maps Job object using the job ID for faster lookup.
var allJobs = make(map[string]*Job)

// CreateJob takes a CreateJobRequest, starts the job, and returns the Job object.
// It returns error if job fails to start.
// On success, it creates a new Job object, assign a unique Job.ID, and sets Job.Status to JobRunning.
// It stores the job in memory against the Job.ID.
// It creates a new go routine that waits for job to exit.
func CreateJob(request CreateJobRequest) (Job, error) {
	// combines stderr and stdout of the job process
	var outputBuffer bytes.Buffer

	cmd := exec.Command(request.Command, request.Args...)
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	// no need to save the job if it fails to start
	if err := cmd.Start(); err != nil {
		err = fmt.Errorf("job failed to start, reason: %v", err)
		log.Errorf(err.Error())
		return Job{}, err
	}

	// job started, create job object, assign ID and save the job
	job := Job{
		ID:           getNextJobID(),
		Request:      request,
		Status:       JobRunning,
		cmd:          cmd,
		outputBuffer: &outputBuffer,
	}

	log.Infof("job %v: started", job.ID)

	// save the job in memory for future reference
	allJobs[job.ID] = &job

	// wait for job to finish
	go handleFinish(&job)

	return job, nil
}

// StopJob takes a job ID, fetch the job from memory, and kills the job process.
// It returns error if job is not found or if job can not be terminated.
// On success, it waits for cmd.Wait to finish which sets the Job.Status to JobExited.
func StopJob(id string) error {
	job, ok := allJobs[id]
	if !ok {
		return fmt.Errorf("job %v: not found", id)
	}
	if job.Status == JobExited {
		log.Infof("job %v: already exited", job.ID)
		return nil
	}

	if err := job.cmd.Process.Kill(); err != nil {
		err = fmt.Errorf("job %v: failed to terminate, reason: %v", job.ID, err)
		log.Errorf(err.Error())
		return err
	}

	// wait for cmd.Wait to finish
	job.waitForExit.Wait()
	log.Infof("job %v: terminated", job.ID)
	return nil
}

// GetJobStatus takes a job ID and returns the details of the job.
// It converts the output buffer to string that indicates current output of the job process.
// It returns error if job is not found.
func GetJobStatus(id string) (Job, error) {
	if job, ok := allJobs[id]; ok {
		job.Output = string(job.outputBuffer.Bytes())
		return *job, nil
	} else {
		return Job{}, fmt.Errorf("job %v: not found", id)
	}
}

// getNextJobID generates a new globally unique identifier.
func getNextJobID() string {
	return xid.New().String()
}

// handleFinish is waits for a job process to finish.
// It is blocking, so it should be called in a separate go routine.
// It sets the Job.Status to JobExited and sets Job.Error if any.
// In case of exec.ExitError, exit code is wrapped in the error.
func handleFinish(job *Job) {
	// initialize the WaitGroup and release on exit
	job.waitForExit = &sync.WaitGroup{}
	job.waitForExit.Add(1)
	defer job.waitForExit.Done()

	if err := job.cmd.Wait(); err != nil {
		job.Error = err
		log.Infof("job %v: finished with error, reason: %v", job.ID, err)
	} else {
		log.Infof("job %v: finished successfully", job.ID)
	}

	job.Status = JobExited
}
