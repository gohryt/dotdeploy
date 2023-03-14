package main

import (
	"errors"
	"log"
	"os"

	"github.com/goccy/go-json"
)

type (
	Deploy struct {
		Folder string `json:"folder"`
		Keep   bool   `json:"keep"`

		Remote []Connection `json:"Remote"`
		Do     []Action     `json:"Do"`
	}

	Checkable interface {
		Check() error
		String() string
	}

	Path struct {
		Connection string `json:"connection"`
		Path       string `json:"path"`
	}
)

func main() {
	name := ".deploy"

	if len(os.Args) == 2 {
		name = os.Args[1]
	}

	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	deploy := new(Deploy)

	err = json.NewDecoder(file).Decode(deploy)
	if err != nil {
		log.Fatal(err)
	}

	for i := range deploy.Remote {
		err = deploy.Connect(&deploy.Remote[i])
		if err != nil {
			log.Fatal(err)
		}
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

	receiver := make(chan error)

	in := int64(0)
	done := int64(0)

	for i, action := range deploy.Do {
		if action.Parallel == false {
			for ; done < in; done += 1 {
				err = errors.Join(err, <-receiver)
			}

			if err != nil {
				log.Fatal(err)
			}

			err = deploy.Process(&deploy.Do[i])
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		in += 1

		go func(action *Action) {
			receiver <- deploy.Process(action)
		}(&deploy.Do[i])
	}

	for ; done < in; done += 1 {
		err = errors.Join(err, <-receiver)
	}

	if err != nil {
		log.Fatal(err)
	}
}
