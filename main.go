package main

import (
	"bufio"
	"log"
	"os"
	"path"

	"github.com/goccy/go-json"
)

type (
	Deploy struct {
		Folder string

		ActionList []Action `json:"Do"`
	}

	Action struct {
		Data any
	}

	Type struct {
		Type string `json:"type"`
	}

	Copy struct {
		File string `json:"file"`
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
		err := os.RemoveAll(deploy.Folder)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for i := range deploy.ActionList {
		err = deploy.Do(i)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (deploy *Deploy) Do(action int) error {
	data := deploy.ActionList[action].Data

	switch data.(type) {
	case *Copy:
		return deploy.Copy(data.(*Copy))
	default:
		log.Println("undefiden action:", data)
	}

	return nil
}

func (deploy *Deploy) Copy(copy *Copy) error {
	source, err := os.Open(copy.File)
	if err != nil {
		return err
	}

	target, err := os.Create(path.Join(deploy.Folder, source.Name()))
	if err != nil {
		return err
	}

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}
