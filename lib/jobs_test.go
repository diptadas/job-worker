package lib

import (
	"strings"
	"testing"
	"time"
)

func TestCreateJobError(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "unknown",
	}

	expectedError := "executable file not found"

	if _, err := jobWorker.CreateJob(request); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error: %v, found: %v", expectedError, err)
	}
}

func TestJobNotFoundError(t *testing.T) {
	jobWorker := NewJobWorker()

	expectedError := "not found"

	if _, err := jobWorker.GetJobStatus("0"); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error: %v, found: %v", expectedError, err)
	}
}

func TestJobOutput(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "echo",
		Args:    []string{"hello world"},
	}

	job, err := jobWorker.CreateJob(request)
	if err != nil {
		t.Fatalf("found error: %v", err)
	}

	// wait for job to finish and copy output
	time.Sleep(time.Second)

	expectedStatus := JobExited
	expectedOutput := "hello world\n"

	job, err = jobWorker.GetJobStatus(job.ID)
	if err != nil {
		t.Fatalf("found error: %v", err)
	} else if job.Status != expectedStatus {
		t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
	} else if job.Output != expectedOutput {
		t.Fatalf("expected output: %v, found: %v", expectedOutput, job.Output)
	}
}

func TestForceStop(t *testing.T) {
	jobWorker := NewJobWorker()

	request := CreateJobRequest{
		Command: "sleep",
		Args:    []string{"300"},
	}

	job, err := jobWorker.CreateJob(request)
	if err != nil {
		t.Fatalf("found error: %v", err)
	}

	time.Sleep(time.Second)

	err = jobWorker.StopJob(job.ID)
	if err != nil { // failed to terminate job
		expectedStatus := JobRunning

		job, err = jobWorker.GetJobStatus(job.ID)
		if err != nil {
			t.Fatalf("found error: %v", err)
		} else if job.Status != expectedStatus {
			t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
		}
	} else { // job terminated
		expectedStatus := JobExited
		expectedError := "signal: killed"

		job, err = jobWorker.GetJobStatus(job.ID)
		if err != nil {
			t.Fatalf("found error: %v", err)
		} else if job.Status != expectedStatus {
			t.Fatalf("expected status: %v, found: %v", expectedStatus, job.Status)
		} else if job.Error == nil || !strings.Contains(job.Error.Error(), expectedError) {
			t.Fatalf("expected error")
		}
	}
}
