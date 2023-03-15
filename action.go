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

	"github.com/gohryt/dotdeploy/unsafe"
	"github.com/melbahja/goph"
)

type (
	Do []*Action

	Action struct {
		Follow string

		Name string
		Data any

		Connection *Connection

		Next  Do
		Base  string
		Error error
	}

	Copy struct {
		From string
		To   string
	}

	Move struct {
		From string
		To   string
	}

	Upload struct {
		From       string
		Connection string
		To         string
	}

	Download struct {
		Connection string
		From       string
		To         string
	}

	Execute struct {
		Connection string
		Path       string
		Timeout    int

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

	ErrCopyFromEmpty           = errors.New(`copy.From == ""`)
	ErrMoveFromEmpty           = errors.New(`move.From == ""`)
	ErrUploadFromEmpty         = errors.New(`upload.From == ""`)
	ErrUploadConnectionEmpty   = errors.New(`upload.Connection == ""`)
	ErrDownloadConnectionEmpty = errors.New(`download.Connection == ""`)
	ErrDownloadFromEmpty       = errors.New(`download.From == ""`)
	ErrExecutePathEmpty        = errors.New(`execute.Path == ""`)
)

var (
	ActionTypeCopy     = unsafe.Type(new(Copy))
	ActionTypeMove     = unsafe.Type(new(Move))
	ActionTypeUpload   = unsafe.Type(new(Upload))
	ActionTypeDownload = unsafe.Type(new(Download))
	ActionTypeExecute  = unsafe.Type(new(Execute))
)

func Process(action *Action) *Action {
	switch action.Data.(type) {
	case *Copy:
		action.Error = action.Copy()
	case *Move:
		action.Error = action.Move()
	case *Upload:
		action.Error = action.Upload()
	case *Download:
		action.Error = action.Download()
	case *Execute:
		action.Error = action.Execute()
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

	if folder > -1 {
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

func (action *Action) Upload() error {
	upload := action.Data.(*Upload)

	if upload.To == "" {
		source, err := os.Open(upload.From)
		if err != nil {
			return err
		}

		upload.To = "/" + source.Name()
		source.Close()
	}

	return action.Connection.Client.Upload(upload.From, upload.To)
}

func (action *Action) Download() error {
	download := action.Data.(*Download)

	if download.To == "" {
		folder := strings.LastIndex(download.From, "/") + 1

		if folder > 0 && folder < len(download.From) {
			download.To = download.From[folder:]
		}
	}

	return action.Connection.Client.Download(download.From, download.To)
}

func (action *Action) Execute() error {
	execute := action.Data.(*Execute)

	if action.Connection == nil {
		return Local(execute)
	}

	client := action.Connection.Client

	command := (*goph.Cmd)(nil)
	err := error(nil)

	if execute.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(execute.Timeout) * time.Second))
		defer cancel()

		command, err = client.CommandContext(ctx, execute.Path, execute.Query...)
	} else {
		command, err = client.Command(execute.Path, execute.Query...)
	}
	if err != nil {
		return err
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	command.Env = execute.Environment

	command.Stdout = stdout
	command.Stderr = stderr

	err = command.Start()
	if err != nil {
		return err
	}

	err = command.Wait()

	_, one := os.Stdout.ReadFrom(stdout)
	_, two := os.Stderr.ReadFrom(stderr)

	return errors.Join(err, one, two)
}

func Local(execute *Execute) error {
	if filepath.Base(execute.Path) == execute.Path {
		path, err := exec.LookPath(execute.Path)
		if err != nil {
			return err
		}

		execute.Path = path
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		execute.Path = filepath.Join(wd, execute.Path)
	}

	command := (*exec.Cmd)(nil)

	if execute.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(execute.Timeout) * time.Second))
		defer cancel()

		command = exec.CommandContext(ctx, execute.Path, execute.Query...)
	} else {
		command = exec.Command(execute.Path, execute.Query...)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	command.Env = append(os.Environ(), execute.Environment...)

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

func (copy *Copy) Validate() error {
	if copy.From == "" {
		return ErrCopyFromEmpty
	}

	return nil
}

func (move *Move) Validate() error {
	if move.From == "" {
		return ErrMoveFromEmpty
	}

	return nil
}

func (upload *Upload) Validate() error {
	join := []error(nil)

	if upload.From == "" {
		join = append(join, ErrUploadFromEmpty)
	}

	if upload.Connection == "" {
		join = append(join, ErrUploadConnectionEmpty)
	}

	return errors.Join(join...)
}

func (download *Download) Validate() error {
	join := []error(nil)

	if download.Connection == "" {
		join = append(join, ErrDownloadConnectionEmpty)
	}

	if download.From == "" {
		join = append(join, ErrDownloadFromEmpty)
	}

	return errors.Join(join...)
}

func (execute *Execute) Validate() error {
	if execute.Path == "" {
		return ErrExecutePathEmpty
	}

	return nil

}
