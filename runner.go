package chatterbox

import (
	"bytes"

	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/schnapper79/chatterbox/types"
)

type Runner struct {
	cmd       *exec.Cmd
	Cancel    context.CancelFunc
	ErrorChan chan error
	LogChan   chan string
	Config    *types.Model_Request
}

func NewRunner(ctx context.Context, Cancel context.CancelFunc, llamaPath, ModelPath string, config *types.Model_Request) *Runner {
	// Convert args map to string slice
	args := config.ToMap()
	args["--model"] = ModelPath + "/" + config.Model
	var argSlice []string
	for k, v := range args {
		if v == "" {
			argSlice = append(argSlice, k)
		} else {
			argSlice = append(argSlice, k, v)
		}
	}

	cmd := exec.CommandContext(ctx, fmt.Sprintf("%s/server", llamaPath), argSlice...)
	logger.Info("Starting server with args: ", strings.Join(cmd.Args, " "))
	return &Runner{
		cmd:       cmd,
		Cancel:    Cancel,
		ErrorChan: make(chan error, 1),
		LogChan:   make(chan string, 100), // Buffer of 100, adjust as needed
		Config:    config,
	}
}

func (r *Runner) Run() error {
	var stdout, stderr bytes.Buffer
	r.cmd.Stdout = &stdout
	r.cmd.Stderr = &stderr

	if err := r.cmd.Start(); err != nil {
		return err
	}

	go func() {
		for {
			line, err := stdout.ReadString('\n')
			if err != nil {
				break
			}
			r.LogChan <- strings.TrimSpace(line)
		}
	}()

	go func() {
		for {
			line, err := stderr.ReadString('\n')
			if err != nil {
				break
			}
			r.LogChan <- strings.TrimSpace(line)
		}
	}()

	go func() {
		if err := r.cmd.Wait(); err != nil {
			r.ErrorChan <- err
			close(r.LogChan)
			close(r.ErrorChan)
			return
		}
	}()
	return nil
}
