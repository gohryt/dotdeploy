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
	Action struct {
		Follow string

		Name string
		Data any

		Next Do
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

func (deploy *Deploy) Process(action *Action) Result {
	result := Result{
		Next: action.Next,
	}

	switch action.Data.(type) {
	case *Copy:
		result.Error = action.Copy(deploy.Folder)
	case *Move:
		result.Error = action.Move(deploy.Folder)
	case *Run:
		result.Error = action.Run()
	default:
		result.Error = ErrUnknowActionType
	}

	return result
}

func (action *Action) Copy(path string) error {
	copy := action.Data.(*Copy)

	source, err := os.Open(copy.From)
	if err != nil {
		return err
	}
	defer source.Close()

	if copy.To == "" {
		copy.To = filepath.Join(path, source.Name())
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

func (action *Action) Move(path string) error {
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
		move.To = filepath.Join(path, source.Name())
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
