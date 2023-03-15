package dotdeploy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type (
	Do []*Action

	Action struct {
		Follow string

		Name string
		Data any

		Next  Do
		Base  string
		Error error
	}

	Copy struct {
		From string `validate:"required"`
		To   string
	}

	Move struct {
		From string `validate:"required"`
		To   string
	}

	Run struct {
		Path    string `validate:"required"`
		Timeout int

		Environment []string
		Query       []string
	}

	Result struct {
		Next  Do
		Error error
	}
)

var (
	ErrUnknowActionType = errors.New("unknown action type")
)

func Process(action *Action) *Action {
	switch action.Data.(type) {
	case *Copy:
		action.Error = action.Copy()
	case *Move:
		action.Error = action.Move()
	case *Run:
		action.Error = action.Run()
	default:
		action.Error = ErrUnknowActionType
	}

	return action
}

func (do Do) Find(name string) (action *Action, ok bool) {
	for i, value := range do {
		if value.Name == name {
			return do[i], true
		}
	}

	return
}

func (action *Action) Copy() error {
	copy := action.Data.(*Copy)

	source, err := os.Open(copy.From)
	if err != nil {
		return err
	}
	defer source.Close()

	if copy.To == "" {
		copy.To = filepath.Join(action.Base, source.Name())
	}

	folder := strings.LastIndex(copy.To, "/")

	if folder > 0 {
		err = os.MkdirAll(copy.To[:folder], os.ModePerm)
		if err != nil {
			return err
		}
	}

	target, err := os.Create(copy.To)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func (action *Action) Move() error {
	move := action.Data.(*Move)

	source, err := os.Open(move.From)
	if err != nil {
		return err
	}
	defer source.Close()

	folder := strings.LastIndex(move.To, "/")

	if folder > 0 {
		err = os.MkdirAll(move.To[:folder], os.ModePerm)
		if err != nil {
			return err
		}
	}

	if move.To == "" {
		move.To = filepath.Join(action.Base, source.Name())
	}

	return os.Rename(move.From, move.To)
}

func (action *Action) Run() error {
	run := action.Data.(*Run)

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
