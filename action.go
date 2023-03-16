package dotdeploy

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"

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

	err = os.MkdirAll(filepath.Dir(copy.To), os.ModePerm)
	if err != nil {
		return err
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

	err = os.MkdirAll(filepath.Dir(move.To), os.ModePerm)
	if err != nil {
		return err
	}

	if move.To == "" {
		move.To = filepath.Join(action.Base, source.Name())
	}

	return os.Rename(move.From, move.To)
}

func (action *Action) Upload() error {
	upload := action.Data.(*Upload)

	if upload.To == "" {
		upload.To = filepath.Base(upload.From)
	}

	return action.Connection.Client.Upload(upload.From, upload.To)
}

func (action *Action) Download() error {
	download := action.Data.(*Download)

	if download.To == "" {
		download.To = filepath.Base(download.From)
	}

	return action.Connection.Client.Download(download.From, download.To)
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
