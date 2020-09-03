package osutil

import (
	"testing"
)

//Test run until complete, should return success
func TestRunCmd(t *testing.T) {
	args := []string{"1"}
	cmd, _, _ := RunCmd("sleep", args, -1, false)
	if !cmd.ProcessState.Success() {
		t.Fatalf("cmd failed")
	}
}

//Test run until complete or timeout expires, should return success
func TestRunCmdTimeout(t *testing.T) {
	args := []string{"1"}
	cmd, _, _ := RunCmd("sleep", args, 10, false)
	if !cmd.ProcessState.Success() {
		t.Fatalf("cmd failed")
	}
}

//Test run until complete or timeout expires,
// timer should expire, kill the job and return not successful
func TestRunCmdFailedTimeout(t *testing.T) {
	args := []string{"10"}
	cmd, _, _ := RunCmd("sleep", args, 1, false)
	if cmd.ProcessState.Success() {
		t.Fatalf("cmd passed but should have failed")
	}
}
