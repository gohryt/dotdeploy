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

	Run struct {
		Path    string
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

	ErrCopyFromEmpty           = errors.New(`copy.From == ""`)
	ErrMoveFromEmpty           = errors.New(`move.From == ""`)
	ErrUploadFromEmpty         = errors.New(`upload.From == ""`)
	ErrUploadConnectionEmpty   = errors.New(`upload.Connection == ""`)
	ErrDownloadConnectionEmpty = errors.New(`download.Connection == ""`)
	ErrDownloadFromEmpty       = errors.New(`download.From == ""`)
	ErrRunPathEmpty            = errors.New(`run.Path == ""`)
)

var (
	ActionTypeCopy     = unsafe.Type(new(Copy))
	ActionTypeMove     = unsafe.Type(new(Move))
	ActionTypeUpload   = unsafe.Type(new(Upload))
	ActionTypeDownload = unsafe.Type(new(Download))
	ActionTypeRun      = unsafe.Type(new(Run))
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

func (run *Run) Validate() error {
	if run.Path == "" {
		return ErrRunPathEmpty
	}

	return nil

}
