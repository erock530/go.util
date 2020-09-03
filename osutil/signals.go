package osutil

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/erock530/go.logging"
)

// List of functions to call before exiting
var (
	exitFunctions    = make([]func(), 0)
	exitFunctionLock = sync.RWMutex{}
	exitFunctionOnce = sync.Once{}
)

func init() {
	RegisterForSignal()
}

//RegisterForSignal Register for the common process killing signals
//Exported only so util package init can call it
func RegisterForSignal() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		_ = <-sig
		exitFunctionLock.RLock()
		if len(exitFunctions) > 0 {
			logging.Warn("Catching exit signal...")
			RunExitSignals()
		}
		exitFunctionLock.RUnlock()
		os.Exit(1)
	}()
}

// ExitOnOrphan exits the program when it becomes orphaned.
// When a process is orphaned, its parent becomes init, which is always PID 1.
func ExitOnOrphan() {
	for range time.Tick(time.Minute) {
		if os.Getppid() == 1 {
			RunExitSignals()
			logging.Fatalf("Exiting because process was orphaned.")
		}
	}
}

// RunExitSignals fires each registered exit function when an exit signal is
// received.
func RunExitSignals() {
	exitFunctionOnce.Do(runExitSignals)
}

func runExitSignals() {
	exitFunctionLock.RLock()
	for _, f := range exitFunctions {
		f()
	}
	exitFunctionLock.RUnlock()
}

// CatchExitSignal is an external interface to register your exit function.
func CatchExitSignal(f func()) {
	exitFunctionLock.Lock()
	exitFunctions = append(exitFunctions, f)
	exitFunctionLock.Unlock()
}
