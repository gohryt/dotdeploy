package dotdeploy

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/gohryt/dotdeploy/multithread"
)

type (
	Validable interface {
		Validate() error
	}

	Executor interface {
		Start() error
		Close() error
		Wait() error
	}
)

type (
	CopyMeta struct {
		Path string
		From *Connection
		To   *Connection
	}

	MoveMeta struct {
		Path string
		From *Connection
	}

	ExecuteMeta struct {
		Path *Connection
	}
)

type (
	Command struct {
		*exec.Cmd
	}
)

func (Command) Close() error { return nil }

func Work(shutdown context.Context, deploy *Deploy) error {
	join := []error(nil)

	err := os.MkdirAll(deploy.Folder, os.ModePerm)
	if err != nil {
		return err
	}

	if deploy.Keep == false {
		defer os.RemoveAll(deploy.Folder)
	}

	connectionReceiver := make(chan *Connection)

	multithread.Go(connectionReceiver, deploy.Remote, Connect)

	for i := 0; i < len(deploy.Remote); i += 1 {
		select {
		case <-shutdown.Done():
			return nil
		case connection := <-connectionReceiver:
			join = append(join, connection.Error)
		}
	}

	err = errors.Join(join...)
	if err != nil {
		return err
	}

	actionReceiver := make(chan *Action)

	multithread.Go(actionReceiver, deploy.Do, Process)

	do := len(deploy.Do)

	for i := 0; i < do; i += 1 {
		action := <-actionReceiver

		if action.Error != nil {
			join = append(join, action.Error)
		} else {
			do += len(action.Next)

			multithread.Go(actionReceiver, action.Next, Process)
		}

	}

	return errors.Join(join...)
}
