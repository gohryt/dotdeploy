package main

import (
	"log"
	"os"

	"github.com/goccy/go-json"
)

type (
	Deploy struct {
		Folder string

		Do []Action
	}

	Action struct {
		Data any
	}

	Type struct {
		Type string
	}

	Copy struct {
		File string
	}
)

func (action *Action) UnmarshalJSON(source []byte) error {
	t := new(Type)

	err := json.Unmarshal(source, t)
	if err != nil {
		return err
	}

	switch t.Type {
	case "copy":
		action.Data = new(Copy)

		err = json.Unmarshal(source, action.Data)
	default:
		action.Data = "empty"
	}

	if err != nil {
		return err
	}

	return nil
}

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

	err = os.MkdirAll(deploy.Folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := os.Remove(deploy.Folder)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for i := range deploy.Do {
		err = deploy.do(i)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (deploy *Deploy) do(action int) error {
	data := deploy.Do[action].Data

	switch data.(type) {
	case *Copy:
		log.Println(data)
	default:
		log.Println("undefiden action:", data)
	}

	return nil
}
