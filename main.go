package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
)

type (
	Deploy struct {
		Folder string `json:"folder"`
		Keep   bool   `json:"keep"`

		Do []Action `json:"Do"`
	}

	Action struct {
		Data     Checkable
		Parallel bool `json:"parallel"`
	}

	Checkable interface {
		Check() error
		String() string
	}

	Type struct {
		Type     string `json:"type"`
		Parallel bool   `json:"parallel"`
	}

	Copy struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	Move struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	Run struct {
		Path    string `json:"path"`
		Timeout int    `json:"timeout"`

		Environment []string `json:"Environment"`
		Query       []string `json:"Query"`
	}

	Empty string
)

var (
	empty = Empty("empty")
)

func (action *Action) UnmarshalJSON(source []byte) error {
	t := new(Type)

	err := json.Unmarshal(source, t)
	if err != nil {
		return err
	}

	action.Parallel = t.Parallel

	switch t.Type {
	case "copy":
		action.Data = new(Copy)
	case "move":
		action.Data = new(Move)
	case "run":
		action.Data = new(Run)
	default:
		action.Data = &empty

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

	for i := range deploy.Do {
		action := deploy.Do[i]

		if action.Parallel == false {
			for ; done < in; done += 1 {
				err = errors.Join(err, <-receiver)
			}

			if err != nil {
				log.Fatal(err)
			}

			err = deploy.Process(action)
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		in += 1

		go func() {
			receiver <- deploy.Process(action)
		}()
	}

	for ; done < in; done += 1 {
		err = errors.Join(err, <-receiver)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func (deploy *Deploy) Process(action Action) error {
	data := action.Data

	switch data.(type) {
	case *Copy:
		return deploy.Copy(data.(*Copy))
	case *Move:
		return deploy.Move(data.(*Move))
	case *Run:
		return deploy.Run(data.(*Run))
	default:
		log.Println("undefiden action:", data.String())
	}

	return nil
}

func (deploy *Deploy) Copy(copy *Copy) error {
	source, err := os.Open(copy.From)
	if err != nil {
		return err
	}
	defer source.Close()

	if copy.To == "" {
		copy.To = filepath.Join(deploy.Folder, source.Name())
	}

	target, err := os.Create(copy.To)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func (deploy *Deploy) Move(move *Move) error {
	source, err := os.Open(move.From)
	if err != nil {
		return err
	}
	defer source.Close()

	if move.To == "" {
		move.To = filepath.Join(deploy.Folder, source.Name())
	}

	return os.Rename(move.From, move.To)
}

func (deploy *Deploy) Run(run *Run) error {
	if filepath.Base(run.Path) == run.Path {
		path, err := exec.LookPath(run.Path)
		if err != nil {
			return err
		}

		run.Path = path
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		run.Path = filepath.Join(wd, run.Path)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	command := (*exec.Cmd)(nil)

	if run.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(run.Timeout) * time.Second))
		defer cancel()

		command = exec.CommandContext(ctx, run.Path, run.Query...)
	} else {
		command = exec.Command(run.Path, run.Query...)
	}

	command.Env = append(os.Environ(), run.Environment...)

	command.Stdout = stdout
	command.Stderr = stderr

	err := command.Start()
	if err != nil {
		return err
	}

	err = command.Wait()

	_, one := os.Stdout.ReadFrom(stdout)
	_, two := os.Stderr.ReadFrom(stderr)

	return errors.Join(err, one, two)
}

func (copy *Copy) Check() error {
	if copy.From == "" {
		return errors.New("'from' can't be empty")
	}

	return nil
}

func (copy *Copy) String() string {
	return fmt.Sprintf("copy from %s to %s", copy.From, copy.To)
}

func (move *Move) Check() error {
	if move.From == "" {
		return errors.New("'from' can't be empty")
	}

	return nil
}

func (move *Move) String() string {
	return fmt.Sprintf("move from %s to %s", move.From, move.To)
}

func (run *Run) Check() error {
	if run.Path == "" {
		return errors.New("'path' can't be empty")
	}

	return nil
}

func (run *Run) String() string {
	return fmt.Sprintf("run %s with timeout %ds, environment %v and query %v", run.Path, run.Timeout, run.Environment, run.Query)
}

func (empty *Empty) Check() error {
	return nil
}

func (empty *Empty) String() string {
	return string(*empty)
}
