package dotdeploy

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/melbahja/goph"
)

type (
	Execute struct {
		Connection string
		Path       string
		Timeout    int

		Environment []string
		Query       []string
	}

	Executor interface {
		Start() error
		Wait() error
	}
)

func (action *Action) Execute() error {
	execute := action.Data.(*Execute)

	executor := Executor(nil)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	if action.Connection != nil {
		goph, err := Goph(execute, action.Connection.Client, stdout, stderr)
		if err != nil {
			return err
		}

		executor = goph
	} else {
		executor = Local(execute, stdout, stderr)
	}

	err := executor.Start()
	if err != nil {
		return err
	}

	err = executor.Wait()

	_, one := os.Stdout.ReadFrom(stdout)
	_, two := os.Stderr.ReadFrom(stderr)

	return errors.Join(err, one, two)
}

func Goph(execute *Execute, client *goph.Client, stdout, stderr io.Writer) (executor Executor, err error) {
	command := (*goph.Cmd)(nil)

	if execute.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(execute.Timeout) * time.Second))
		defer cancel()

		command, err = client.CommandContext(ctx, execute.Path, execute.Query...)
	} else {
		command, err = client.Command(execute.Path, execute.Query...)
	}
	if err != nil {
		return
	}

	command.Env = execute.Environment

	command.Stdout = stdout
	command.Stderr = stderr

	return command, nil
}

func Local(execute *Execute, stdout, stderr io.Writer) Executor {
	command := (*exec.Cmd)(nil)

	if execute.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(execute.Timeout) * time.Second))
		defer cancel()

		command = exec.CommandContext(ctx, execute.Path, execute.Query...)
	} else {
		command = exec.Command(execute.Path, execute.Query...)
	}

	command.Env = execute.Environment

	command.Stdout = stdout
	command.Stderr = stderr

	return command
}

func (execute *Execute) Validate() error {
	if execute.Path == "" {
		return ErrExecutePathEmpty
	}

	return nil

}
