// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package harness

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	"sync"

	"github.com/revel/cmd/model"
	"github.com/revel/cmd/utils"
	"runtime"
)

// App contains the configuration for running a Revel app.  (Not for the app itself)
// Its only purpose is constructing the command to execute.
type App struct {
	BinaryPath     string            // Path to the app executable
	Port           int               // Port to pass as a command line argument.
	cmd            AppCmd            // The last cmd returned.
	PackagePathMap map[string]string // Package to directory path map
	Paths          *model.RevelContainer
}

// NewApp returns app instance with binary path in it
func NewApp(binPath string, paths *model.RevelContainer, packagePathMap map[string]string) *App {
	return &App{BinaryPath: binPath, Paths: paths, Port: paths.HTTPPort, PackagePathMap:packagePathMap}
}

// Cmd returns a command to run the app server using the current configuration.
func (a *App) Cmd(runMode string) AppCmd {
	a.cmd = NewAppCmd(a.BinaryPath, a.Port, runMode, a.Paths)
	return a.cmd
}

// Kill the last app command returned.
func (a *App) Kill() {
	a.cmd.Kill()
}

// AppCmd manages the running of a Revel app server.
// It requires revel.Init to have been called previously.
type AppCmd struct {
	*exec.Cmd
}

// NewAppCmd returns the AppCmd with parameters initialized for running app
func NewAppCmd(binPath string, port int, runMode string, paths *model.RevelContainer) AppCmd {
	cmd := exec.Command(binPath,
		fmt.Sprintf("-port=%d", port),
		fmt.Sprintf("-importPath=%s", paths.ImportPath),
		fmt.Sprintf("-runMode=%s", runMode))
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return AppCmd{cmd}
}

// Start the app server, and wait until it is ready to serve requests.
func (cmd AppCmd) Start(c *model.CommandConfig) error {
	listeningWriter := &startupListeningWriter{os.Stdout, make(chan bool), c, &bytes.Buffer{}}
	cmd.Stdout = listeningWriter
	utils.CmdInit(cmd.Cmd, !c.Vendored, c.AppPath)
	utils.Logger.Info("Exec app:", "path", cmd.Path, "args", cmd.Args, "dir", cmd.Dir, "env", cmd.Env)
	if err := cmd.Cmd.Start(); err != nil {
		utils.Logger.Fatal("Error running:", "error", err)
	}

	select {
	case exitState := <-cmd.waitChan():
		fmt.Println("Startup failure view previous messages, \n Proxy is listening :", c.Run.Port)
		err := utils.NewError("", "Revel Run Error", "starting your application there was an exception. See terminal output, " + exitState, "")
	// TODO pretiffy command line output
	// err.MetaError = listeningWriter.getLastOutput()
		return err

	case <-time.After(60 * time.Second):
		println("Revel proxy is listening, point your browser to :", c.Run.Port)
		utils.Logger.Error("Killing revel server process did not respond after wait timeout.", "processid", cmd.Process.Pid)
		cmd.Kill()
		return errors.New("revel/harness: app timed out")

	case <-listeningWriter.notifyReady:
		println("Revel proxy is listening, point your browser to :", c.Run.Port)
		return nil
	}

}

// Run the app server inline.  Never returns.
func (cmd AppCmd) Run(c *model.CommandConfig) {
	utils.CmdInit(cmd.Cmd, !c.Vendored, c.AppPath)
	utils.Logger.Info("Exec app:", "path", cmd.Path, "args", cmd.Args)
	if err := cmd.Cmd.Run(); err != nil {
		utils.Logger.Fatal("Error running:", "error", err)
	}
}

// Kill terminates the app server if it's running.
func (cmd AppCmd) Kill() {

	if cmd.Cmd != nil && (cmd.ProcessState == nil || !cmd.ProcessState.Exited()) {
		// Windows appears to send the kill to all threads, shutting down the
		// server before this can, this check will ensure the process is still running
		if _, err := os.FindProcess(int(cmd.Process.Pid)); err != nil {
			// Server has already exited
			utils.Logger.Info("Server not running revel server pid", "pid", cmd.Process.Pid)
			return
		}

		// Wait for the shutdown channel
		waitMutex := &sync.WaitGroup{}
		waitMutex.Add(1)
		ch := make(chan bool, 1)
		go func() {
			waitMutex.Done()
			s, err := cmd.Process.Wait()
			defer func() {
				ch <- true
			}()
			if err != nil {
				utils.Logger.Info("Wait failed for process ", "error", err)
			}
			if s != nil {
				utils.Logger.Info("Revel App exited", "state", s.String())
			}
		}()
		// Wait for the channel to begin waiting
		waitMutex.Wait()

		// Send an interrupt signal to allow for a graceful shutdown
		utils.Logger.Info("Killing revel server pid", "pid", cmd.Process.Pid)
		var err error
		if runtime.GOOS != "windows" {
			// os.Interrupt is not available on windows
			err = cmd.Process.Signal(os.Interrupt)
		}

		if err != nil {
			utils.Logger.Info(
				"Revel app already exited.",
				"processid", cmd.Process.Pid, "error", err,
				"killerror", cmd.Process.Kill())
			return
		}


		// Use a timer to ensure that the process exits
		utils.Logger.Info("Waiting to exit")
		select {
		case <-ch:
			return
		case <-time.After(60 * time.Second):
		// Kill the process
			utils.Logger.Error(
				"Revel app failed to exit in 60 seconds - killing.",
				"processid", cmd.Process.Pid,
				"killerror", cmd.Process.Kill())
		}

		utils.Logger.Info("Done Waiting to exit")
	}
}

// Return a channel that is notified when Wait() returns.
func (cmd AppCmd) waitChan() <-chan string {
	ch := make(chan string, 1)
	go func() {
		_ = cmd.Wait()
		state := cmd.ProcessState
		exitStatus := " unknown "
		if state != nil {
			exitStatus = state.String()
		}

		ch <- exitStatus
	}()
	return ch
}

// A io.Writer that copies to the destination, and listens for "Revel engine is listening on.."
// in the stream.  (Which tells us when the revel server has finished starting up)
// This is super ghetto, but by far the simplest thing that should work.
type startupListeningWriter struct {
	dest        io.Writer
	notifyReady chan bool
	c           *model.CommandConfig
	buffer      *bytes.Buffer
}

// Writes to this output stream
func (w *startupListeningWriter) Write(p []byte) (int, error) {
	if w.notifyReady != nil && bytes.Contains(p, []byte("Revel engine is listening on")) {
		w.notifyReady <- true
		w.notifyReady = nil
	}
	if w.c.HistoricMode {
		if w.notifyReady != nil && bytes.Contains(p, []byte("Listening on")) {
			w.notifyReady <- true
			w.notifyReady = nil
		}
	}
	if w.notifyReady != nil {
		w.buffer.Write(p)
	}
	return w.dest.Write(p)
}

// Returns the cleaned output from the response
// TODO clean the response more
func (w *startupListeningWriter) getLastOutput() string {
	return w.buffer.String()
}

