package dotdeploy

import (
	"context"
	"errors"
	"os"

	"github.com/gohryt/dotdeploy/multithread"
)

type (
	Deploy struct {
		Folder string
		Keep   bool

		Remote Remote
		Do     Do
	}

	Validable interface {
		Validate() error
	}
)

var (
	ErrUnknownFollowReference     = errors.New("unknown follow reference")
	ErrUnknownConnectionReference = errors.New("unknown connection reference")

	ErrDeployFolderEmpty = errors.New(`deploy.Folder == ""`)
)

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

func (deploy *Deploy) Prepare() error {
	err := deploy.Validate()
	if err != nil {
		return err
	}

	base := Do(nil)

	for i := range deploy.Do {
		action := deploy.Do[i]

		action.Base = deploy.Folder

		remote := ""

		switch action.Data.(type) {
		case *Upload:
			remote = action.Data.(*Upload).Connection
		case *Download:
			remote = action.Data.(*Download).Connection
		case *Execute:
			remote = action.Data.(*Execute).Connection
		}

		if remote != "" {
			connection, ok := deploy.Remote.Find(remote)

			if ok == false {
				return ErrUnknownConnectionReference
			}

			action.Connection = connection
		}

		if action.Follow != "" {
			follow, ok := deploy.Do.Find(action.Follow)

			if ok == false {
				return ErrUnknownFollowReference
			}

			follow.Next = append(follow.Next, action)
		} else {
			base = append(base, action)
		}
	}

	deploy.Do = base

	return nil
}

func (deploy *Deploy) Validate() error {
	join := []error(nil)

	if deploy.Folder == "" {
		join = append(join, ErrDeployFolderEmpty)
	}

	for i := range deploy.Remote {
		err := deploy.Remote[i].Data.(Validable).Validate()
		if err != nil {
			join = append(join, err)
		}
	}

	for i := range deploy.Do {
		err := deploy.Do[i].Data.(Validable).Validate()
		if err != nil {
			join = append(join, err)
		}
	}

	return errors.Join(join...)
}
