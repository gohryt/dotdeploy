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
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type (
	Action struct {
		Data     Checkable
		Parallel bool `json:"parallel"`
	}

	ActionType struct {
		Type     string `json:"type"`
		Parallel bool   `json:"parallel"`
	}

	Copy struct {
		From Path `json:"From"`
		To   Path `json:"To"`
	}

	Move struct {
		From Path `json:"From"`
		To   Path `json:"To"`
	}

	Run struct {
		Path    string `json:"path"`
		Timeout int    `json:"timeout"`

		Environment []string `json:"Environment"`
		Query       []string `json:"Query"`
	}

	Empty string
)

func (action *Action) UnmarshalJSON(source []byte) error {
	t := new(ActionType)

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

func (copy *Copy) Check() error {
	if copy.From.Path == "" {
		return errors.New("'from' can't be empty")
	}

	return nil
}

func (copy *Copy) String() string {
	return fmt.Sprintf("copy from %s to %s", copy.From, copy.To)
}

func (move *Move) Check() error {
	if move.From.Path == "" {
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

func (deploy *Deploy) Process(action *Action) error {
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
	source, err := os.Open(copy.From.Path)
	if err != nil {
		return err
	}
	defer source.Close()

	if copy.To.Path == "" {
		copy.To.Path = filepath.Join(deploy.Folder, source.Name())
	}

	folder := strings.LastIndex(copy.To.Path, "/")

	if folder > 0 {
		err = os.MkdirAll(copy.To.Path[:folder], os.ModePerm)
		if err != nil {
			return err
		}
	}

	target, err := os.Create(copy.To.Path)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func (deploy *Deploy) Move(move *Move) error {
	source, err := os.Open(move.From.Path)
	if err != nil {
		return err
	}
	defer source.Close()

	folder := strings.LastIndex(move.To.Path, "/")

	if folder > 0 {
		err = os.MkdirAll(move.To.Path[:folder], os.ModePerm)
		if err != nil {
			return err
		}
	}

	if move.To.Path == "" {
		move.To.Path = filepath.Join(deploy.Folder, source.Name())
	}

	return os.Rename(move.From.Path, move.To.Path)
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
