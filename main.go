package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type (
	Deploy struct {
		Folder string `yaml:"folder" validate:"required"`
		Keep   bool   `yaml:"keep"`

		Remote Remote `yaml:"Remote"`
		Do     Do     `yaml:"Do"`
	}
)

func main() {
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := new(Deploy).Work(shutdown)
	if err != nil {
		log.Fatal(err)
	}
}

func (deploy *Deploy) Work(shutdown context.Context) error {
	file, err := os.Open(".deploy")
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(deploy)
	if err != nil {
		return err
	}

	err = validator.New().Struct(deploy)
	if err != nil {
		return err
	}

	errorReceiver := make(chan error)
	errorSlice := []error(nil)

	remote := len(deploy.Remote)

	for i := range deploy.Remote {
		go func(connection *Connection) {
			errorReceiver <- deploy.Connect(connection)
		}(deploy.Remote[i])
	}

	for i := 0; i < remote; i += 1 {
		errorSlice = append(errorSlice, <-errorReceiver)
	}

	err = errors.Join(errorSlice...)
	if err != nil {
		return err
	}

	err = os.MkdirAll(deploy.Folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if deploy.Keep == false {
		defer func() {
			err := os.RemoveAll(deploy.Folder)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	base := Do(nil)

	for i := range deploy.Do {
		action := deploy.Do[i]

		if action.Follow == "" {
			base = append(base, action)
		} else {
			follow, ok := deploy.Do.Find(action.Follow)

			if ok == false {
				return errors.New("action has unknown follow key")
			}

			follow.Next = append(follow.Next, action)
		}
	}

	resultReceiver := make(chan Result)

	deploy.Cycle(resultReceiver, base)

	do := len(deploy.Do)

	for i := 0; i < do; i += 1 {
		result := <-resultReceiver

		if result.Error != nil {
			do -= Count(result.Next)

			errorSlice = append(errorSlice, result.Error)
		} else {
			deploy.Cycle(resultReceiver, result.Next)
		}

	}

	return errors.Join(errorSlice...)
}

func Count(do Do) int {
	c := len(do)

	for i := range do {
		c += Count(do[i].Next)
	}

	return c
}

func (deploy *Deploy) Cycle(resultReceiver chan Result, do Do) {
	for i := range do {
		go func(action *Action) {
			resultReceiver <- deploy.Process(action)
		}(do[i])
	}
}
