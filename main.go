package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goccy/go-json"
)

type (
	Deploy struct {
		Folder string

		ActionList []Action `json:"Do"`
	}

	Action struct {
		Data     any
		Parallel bool
	}

	Type struct {
		Type     string `json:"type"`
		Parallel bool   `json:"parallel"`
	}

	Copy struct {
		File string `json:"file"`
	}

	Run struct {
		Path string `json:"path"`

		EnvironmentList []string `json:"Environments"`
		ArgumentList    []string `json:"Arguments"`
	}
)

func (action *Action) UnmarshalJSON(source []byte) error {
	t := new(Type)

	err := json.Unmarshal(source, t)
	if err != nil {
		return err
	}

	action.Parallel = t.Parallel

	switch t.Type {
	case "run":
		action.Data = new(Run)
	case "copy":
		action.Data = new(Copy)
	default:
		action.Data = "empty"
		return nil
	}

	return json.Unmarshal(source, action.Data)
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
	case *Run:
		return deploy.Run(data.(*Run))
	case *Copy:
		return deploy.Copy(data.(*Copy))
	default:
		log.Println("undefiden action:", data)
	}

	return nil
}

func (deploy *Deploy) Run(run *Run) error {
	if filepath.Base(run.Path) == run.Path {
		path, err := exec.LookPath(run.Path)
		if err != nil {
			return err
		}

		run.Path = path
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	command := &exec.Cmd{
		Path: run.Path,
		Env:  append(os.Environ(), run.EnvironmentList...),
		Args: append([]string{run.Path}, run.ArgumentList...),

		Stdout: stdout,
		Stderr: stderr,
	}

	err := command.Run()
	if err != nil {
		return err
	}

	_, err = os.Stdout.ReadFrom(stdout)
	if err != nil {
		return err
	}

	_, err = os.Stderr.ReadFrom(stderr)
	return err
}

func (deploy *Deploy) Copy(copy *Copy) error {
	source, err := os.Open(copy.File)
	if err != nil {
		return err
	}

	target, err := os.Create(filepath.Join(deploy.Folder, source.Name()))
	if err != nil {
		return err
	}

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}
