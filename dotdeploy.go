package dotdeploy

import (
	"context"
	"errors"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/gohryt/dotdeploy/multithread"
)

type (
	Deploy struct {
		Folder string `validate:"required"`
		Keep   bool

		Remote Remote
		Do     Do
	}
)

var (
	ErrUnknownFollowReference = errors.New("unknown follow reference")
)

func Work(shutdown context.Context, deploy *Deploy) error {
	err := validator.New().Struct(deploy)
	if err != nil {
		return err
	}

	join := []error(nil)

	err = os.MkdirAll(deploy.Folder, os.ModePerm)
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
	base := Do(nil)

	for i := range deploy.Do {
		action := deploy.Do[i]

		action.Base = deploy.Folder

		if action.Follow == "" {
			base = append(base, action)
		} else {
			follow, ok := deploy.Do.Find(action.Follow)

			if ok == false {
				return ErrUnknownFollowReference
			}

			follow.Next = append(follow.Next, action)
		}
	}

	deploy.Do = base
	return nil
}
