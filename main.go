package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

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
		From string `json:"from"`
		To   string `json:"to"`
	}

	Run struct {
		Path string `json:"path"`

		Environment  []string `json:"Environment"`
		ArgumentList []string `json:"ArgumentList"`
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

	group := new(sync.WaitGroup)

	for i := range deploy.ActionList {
		action := deploy.ActionList[i]

		if action.Parallel == false {
			group.Wait()

			err = deploy.Do(action)
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		group.Add(1)

		go func() {
			err := deploy.Do(action)
			if err != nil {
				log.Fatal(err)
			}

			group.Add(-1)
		}()
	}

	group.Wait()
}

func (deploy *Deploy) Do(action Action) error {
	data := action.Data

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
		Env:  append(os.Environ(), run.Environment...),
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
	source, err := os.Open(copy.From)
	if err != nil {
		return err
	}

	if copy.To == "" {
		copy.To = filepath.Join(deploy.Folder, source.Name())
	}

	target, err := os.Create(copy.To)
	if err != nil {
		return err
	}

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}
