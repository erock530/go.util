package osutil

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/erock530/go.logging"
)

// GetSetEnv Gets and potentially sets environment variable to the fallback value.
// Returns environment-set value if present, fallback otherwise
func GetSetEnv(envVar, fallback string) string {
	envString := os.Getenv(envVar)
	if envString == "" {
		envString = fallback
		if err := os.Setenv(envVar, envString); err != nil {
			logging.Errorf("Unable to set env %s = %s: %v", envVar, envString, err)
		}
	}

	return envString
}

// RunCmd runs the given executable with the given arguments.
//  executable: executable name(or full path if its not in the path)
//  args: list of arguments as strings
//  timeout: if < 0 wait for the commmand to return before exiting
//           if == 0 return  immediately
//           if > 0 max number of seconds process will run before killing process
//  verbose: whether to print information about how the call goes
func RunCmd(executable string, args []string, timeout int, verbose bool) (*exec.Cmd, bytes.Buffer, bytes.Buffer) {
	cmd := exec.Command(executable, args...)

	if verbose {
		logging.Info("%+v", cmd)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	var err error
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if timeout < 0 {
		err = cmd.Run()
	} else {
		err = cmd.Start()
		//running in timeout mode
		if err == nil && timeout > 0 {
			//kill running process if timeout expires
			kill := time.AfterFunc(time.Second*time.Duration(timeout), func() { cmd.Process.Kill() })

			//wait until process returns,
			//process will return if timeout expires and kill signal is sent
			err = cmd.Wait()

			//stop time if it's still running,
			//we don't want to send a signal to a process in the future with a recovered pid
			kill.Stop()
		}
	}

	if verbose {
		if err != nil {
			logging.Errorf("Error running %s %s", executable,
				fmt.Sprint(err)+": "+stderr.String())
		} else {
			logging.Info("No Error: %v", stderr.String())
		}
		logging.Info("Result: %v", out.String())
	}

	return cmd, out, stderr
}

//RunCmdSafely runs the command and logs errors
func RunCmdSafely(cmdStr string, args []string, workdir string, env []string) error {
	return RunCmdSafelyStdOut(cmdStr, nil, args, workdir, env)
}

//RunCmdSafelyStdOut runs the command and logs errors and captures std out into buffer
func RunCmdSafelyStdOut(cmdStr string, out *bytes.Buffer, args []string, workdir string, env []string) error {
	// run the command and return err
	cmd := exec.Command(cmdStr, args...)
	if len(env) > 0 {
		cmd.Env = env
	}
	if workdir != "" {
		cmd.Dir = workdir
	}
	if out != nil {
		cmd.Stdout = out
	}

	// Run the command and pipe errors to log
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	errPipeScanner := bufio.NewScanner(stderr)
	for errPipeScanner.Scan() {
		logging.Errorf(errPipeScanner.Text())
	}
	if err := errPipeScanner.Err(); err != nil {
		return err
	}

	return cmd.Wait()
}

