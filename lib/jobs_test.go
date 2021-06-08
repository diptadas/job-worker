package lib

import (
	"strings"
	"testing"
	"time"
)

func TestCreateJobError(t *testing.T) {
	request := CreateJobRequest{
		Command: "unknown",
	}

	expectedError := "executable file not found"

	if _, err := CreateJob(request); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error: %v, found: %v", expectedError, err)
	}
}

func TestJobNotFoundError(t *testing.T) {
	expectedError := "not found"

	if _, err := GetJobStatus("0"); err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error: %v, found: %v", expectedError, err)
	}
}

func TestJobOutput(t *testing.T) {
	request := CreateJobRequest{
		Command: "echo",
		Args:    []string{"hello world"},
	}

	job, err := CreateJob(request)
	if err != nil {
		t.Errorf("found error: %v", err)
	}

	time.Sleep(time.Second)

	expectedStatus := JobExited
	expectedOutput := "hello world\n"

	job, err = GetJobStatus(job.ID)
	if err != nil {
		t.Errorf("found error: %v", err)
	} else if job.Status != expectedStatus {
		t.Errorf("expected status: %v, found: %v", expectedStatus, job.Status)
	} else if job.Output != expectedOutput {
		t.Errorf("expected output: %v, found: %v", expectedOutput, job.Output)
	}
}

func TestForceStop(t *testing.T) {
	request := CreateJobRequest{
		Command: "sleep",
		Args:    []string{"300"},
	}

	job, err := CreateJob(request)
	if err != nil {
		t.Errorf("found error: %v", err)
	}

	time.Sleep(time.Second)

	err = StopJob(job.ID)
	if err != nil { // failed to terminate job
		expectedStatus := JobRunning

		job, err = GetJobStatus(job.ID)
		if err != nil {
			t.Errorf("found error: %v", err)
		} else if job.Status != expectedStatus {
			t.Errorf("expected status: %v, found: %v", expectedStatus, job.Status)
		}
	} else { // job terminated
		expectedStatus := JobExited
		expectedError := "signal: killed"

		job, err = GetJobStatus(job.ID)
		if err != nil {
			t.Errorf("found error: %v", err)
		} else if job.Status != expectedStatus {
			t.Errorf("expected status: %v, found: %v", expectedStatus, job.Status)
		} else if job.Error == nil || !strings.Contains(job.Error.Error(), expectedError) {
			t.Errorf("expected error")
		}
	}
}
