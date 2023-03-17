package dotdeploy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gohryt/dotdeploy/unsafe"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
)

type (
	Do []*Action

	Action struct {
		Follow string

		Name string
		Data any

		Meta any
		Next Do

		Error error
	}

	Copy struct {
		From Path
		To   Path
	}

	Move struct {
		From Path
		To   string
	}

	Execute struct {
		Timeout int

		Path        Path
		Environment []string
		Query       []string
	}
)

var (
	ErrUnknowActionType = errors.New("unknown action type")

	ErrCopyFromPathEmpty    = errors.New(`copy.From.Path == ""`)
	ErrMoveFromPathEmpty    = errors.New(`move.From.Path == ""`)
	ErrExecutePathPathEmpty = errors.New(`execute.Path.Path == ""`)

	ErrActionNotFount = errors.New("action not found")
)

var (
	ActionTypeCopy    = unsafe.Type(new(Copy))
	ActionTypeMove    = unsafe.Type(new(Move))
	ActionTypeExecute = unsafe.Type(new(Execute))
)

func Process(action *Action) *Action {
	switch action.Data.(type) {
	case *Copy:
		action.Error = action.Copy()
	case *Move:
		action.Error = action.Move()
	case *Execute:
		action.Error = action.Execute()
	default:
		action.Error = ErrUnknowActionType
	}

	return action
}

func (do Do) Find(name string) (action *Action, err error) {
	for i, value := range do {
		if value.Name == name {
			action = do[i]
			return
		}
	}

	err = ErrActionNotFount
	return
}

func (action *Action) Copy() error {
	copy := action.Data.(*Copy)
	copyMeta := action.Meta.(*CopyMeta)

	if copy.To.Connection != "" {
		if copy.To.Path == "" {
			copy.To.Path = filepath.Base(copy.From.Path)
		}

		if copy.From.Connection != "" {
			return copyConnectionToConnection(copy, copyMeta)
		} else {
			return copyLocalToConnection(copy, copyMeta)
		}
	}

	if copy.To.Path == "" {
		copy.To.Path = filepath.Join(copyMeta.Path, filepath.Base(copy.From.Path))
	}

	if copy.From.Connection != "" {
		return copyConnectionToLocal(copy, copyMeta)
	} else {
		return copyLocalToLocal(copy, copyMeta)
	}
}

func copyConnectionToConnection(copy *Copy, copyMeta *CopyMeta) error {
	from, err := copyMeta.From.Client.NewSftp()
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := copyMeta.To.Client.NewSftp()
	if err != nil {
		return err
	}
	defer to.Close()

	path := filepath.Dir(copy.To.Path)

	err = to.Remove(path)
	if err != nil {
		return err
	}

	err = to.MkdirAll(path)
	if err != nil {
		return err
	}

	source, err := from.Stat(copy.From.Path)
	if err != nil {
		return err
	}

	if source.IsDir() {
		return copyDirectoryConnectionToConnection(from, to, copy.From.Path, copy.To.Path)
	}

	return copyFileConnectionToConnection(from, to, copy.From.Path, copy.To.Path)
}

func copyDirectoryConnectionToConnection(clientFrom, clientTo *sftp.Client, from, to string) error {
	source, err := clientFrom.ReadDir(from)
	if err != nil {
		return err
	}

	err = clientTo.MkdirAll(to)
	if err != nil {
		return err
	}

	for i := range source {
		entry := source[i]

		name := entry.Name()

		from := filepath.Join(from, name)
		to := filepath.Join(to, name)

		if entry.IsDir() {
			err = copyDirectoryConnectionToConnection(clientFrom, clientTo, from, to)
		} else {
			err = copyFileConnectionToConnection(clientFrom, clientTo, from, to)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFileConnectionToConnection(clientFrom, clientTo *sftp.Client, from, to string) error {
	source, err := clientFrom.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := clientTo.Create(to)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func copyLocalToConnection(copy *Copy, copyMeta *CopyMeta) error {
	to, err := copyMeta.To.Client.NewSftp()
	if err != nil {
		return err
	}
	defer to.Close()

	path := filepath.Dir(copy.To.Path)

	err = to.Remove(path)
	if err != nil {
		return err
	}

	err = to.MkdirAll(path)
	if err != nil {
		return err
	}

	source, err := os.Stat(copy.From.Path)
	if err != nil {
		return err
	}

	if source.IsDir() {
		return copyDirectoryLocalToConnection(to, copy.From.Path, copy.To.Path)
	}

	return copyFileLocalToConnection(to, copy.From.Path, copy.To.Path)
}

func copyDirectoryLocalToConnection(client *sftp.Client, from, to string) error {
	source, err := os.ReadDir(from)
	if err != nil {
		return err
	}

	err = client.MkdirAll(to)
	if err != nil {
		return err
	}

	for i := range source {
		entry := source[i]

		name := entry.Name()

		from := filepath.Join(from, name)
		to := filepath.Join(to, name)

		if entry.IsDir() {
			err = copyDirectoryLocalToConnection(client, from, to)
		} else {
			err = copyFileLocalToConnection(client, from, to)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFileLocalToConnection(client *sftp.Client, from, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := client.Create(to)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func copyConnectionToLocal(copy *Copy, copyMeta *CopyMeta) error {
	from, err := copyMeta.From.Client.NewSftp()
	if err != nil {
		return err
	}
	defer from.Close()

	path := filepath.Dir(copy.To.Path)

	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	source, err := from.Stat(copy.From.Path)
	if err != nil {
		return err
	}

	if source.IsDir() {
		return copyDirectoryConnectionToLocal(from, copy.From.Path, copy.To.Path)
	}

	return copyFileConnectionToLocal(from, copy.From.Path, copy.To.Path)
}

func copyDirectoryConnectionToLocal(client *sftp.Client, from, to string) error {
	source, err := client.ReadDir(from)
	if err != nil {
		return err
	}

	err = os.MkdirAll(to, os.ModePerm)
	if err != nil {
		return err
	}

	for i := range source {
		entry := source[i]

		name := entry.Name()

		from := filepath.Join(from, name)
		to := filepath.Join(to, name)

		if entry.IsDir() {
			err = copyDirectoryConnectionToLocal(client, from, to)
		} else {
			err = copyFileConnectionToLocal(client, from, to)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFileConnectionToLocal(client *sftp.Client, from, to string) error {
	source, err := client.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(to)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func copyLocalToLocal(copy *Copy, copyMeta *CopyMeta) error {
	path := filepath.Dir(copy.To.Path)

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	source, err := os.Stat(copy.From.Path)
	if err != nil {
		return err
	}

	if source.IsDir() {
		return copyDirectoryLocalToLocal(copy.From.Path, copy.To.Path)
	}

	return copyFileLocalToLocal(copy.From.Path, copy.To.Path)
}

func copyDirectoryLocalToLocal(from, to string) error {
	source, err := os.ReadDir(from)
	if err != nil {
		return err
	}

	err = os.MkdirAll(to, os.ModePerm)
	if err != nil {
		return err
	}

	for i := range source {
		entry := source[i]

		name := entry.Name()

		from := filepath.Join(from, name)
		to := filepath.Join(to, name)

		if entry.IsDir() {
			err = copyDirectoryLocalToLocal(from, to)
		} else {
			err = copyFileLocalToLocal(from, to)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyFileLocalToLocal(from, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(to)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = bufio.NewWriter(target).ReadFrom(source)
	return err
}

func (action *Action) Move() error {
	move := action.Data.(*Move)
	moveMeta := action.Meta.(*MoveMeta)

	if move.From.Connection != "" {
		return moveConnection(move, moveMeta)
	} else {
		return moveLocal(move, moveMeta)
	}
}

func moveConnection(move *Move, moveMeta *MoveMeta) error {
	sftp, err := moveMeta.From.Client.NewSftp()
	if err != nil {
		return err
	}
	defer sftp.Close()

	source, err := sftp.Stat(move.From.Path)
	if err != nil {
		return err
	}

	if move.To == "" {
		move.To = source.Name()
	}

	err = sftp.MkdirAll(filepath.Dir(move.To))
	if err != nil {
		return err
	}

	return sftp.Rename(move.From.Path, move.To)
}

func moveLocal(move *Move, moveMeta *MoveMeta) error {
	source, err := os.Open(move.From.Path)
	if err != nil {
		return err
	}
	defer source.Close()

	if move.To == "" {
		move.To = filepath.Join(moveMeta.Path, source.Name())
	}

	err = os.MkdirAll(filepath.Dir(move.To), os.ModePerm)
	if err != nil {
		return err
	}

	return os.Rename(move.From.Path, move.To)
}

func (action *Action) Execute() error {
	execute := action.Data.(*Execute)
	executeMeta := action.Meta.(*ExecuteMeta)

	executor := Executor(nil)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	if execute.Path.Connection != "" {
		goph, err := executeConnection(execute, executeMeta.Path.Client, stdout, stderr)
		if err != nil {
			return err
		}

		executor = goph
	} else {
		executor = executeLocal(execute, stdout, stderr)
	}

	err := executor.Start()
	if err != nil {
		return err
	}

	err = executor.Wait()

	_, one := os.Stdout.ReadFrom(stdout)
	_, two := os.Stderr.ReadFrom(stderr)

	return errors.Join(err, one, two)
}

func executeConnection(execute *Execute, client *goph.Client, stdout, stderr io.Writer) (executor Executor, err error) {
	command := (*goph.Cmd)(nil)

	if execute.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(execute.Timeout) * time.Second))
		defer cancel()

		command, err = client.CommandContext(ctx, execute.Path.Path, execute.Query...)
	} else {
		command, err = client.Command(execute.Path.Path, execute.Query...)
	}
	if err != nil {
		return
	}

	command.Env = execute.Environment

	command.Stdout = stdout
	command.Stderr = stderr

	return command, nil
}

func executeLocal(execute *Execute, stdout, stderr io.Writer) Executor {
	command := Command{}

	if execute.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(execute.Timeout) * time.Second))
		defer cancel()

		command.Cmd = exec.CommandContext(ctx, execute.Path.Path, execute.Query...)
	} else {
		command.Cmd = exec.Command(execute.Path.Path, execute.Query...)
	}

	command.Env = append(os.Environ(), execute.Environment...)

	command.Stdout = stdout
	command.Stderr = stderr

	return command
}

func (copy *Copy) Validate() error {
	if copy.From.Path == "" {
		return ErrCopyFromPathEmpty
	}

	return nil
}

func (move *Move) Validate() error {
	if move.From.Path == "" {
		return ErrMoveFromPathEmpty
	}

	return nil
}

func (execute *Execute) Validate() error {
	if execute.Path.Path == "" {
		return ErrExecutePathPathEmpty
	}

	return nil

}
