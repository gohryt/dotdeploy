package main

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/gohryt/dotdeploy"
	"github.com/gohryt/dotdeploy/unsafe"
)

type (
	Action struct {
		Follow string
		Name   string
		Data   any
	}

	ActionType struct {
		Type   string `yaml:"type"`
		Name   string `yaml:"name"`
		Follow string `yaml:"follow"`
	}

	Copy struct {
		From Path `yaml:"From"`
		To   Path `yaml:"To"`
	}

	Move struct {
		From Path   `yaml:"From"`
		To   string `yaml:"To"`
	}

	Execute struct {
		Timeout int `yaml:"timeout"`

		Path        Path     `yaml:"Path"`
		Environment []string `yaml:"Environment"`
		Query       []string `yaml:"Query"`
	}
)

func (action *Action) UnmarshalYAML(value *yaml.Node) error {
	t := new(ActionType)

	err := value.Decode(t)
	if err != nil {
		return err
	}

	if t.Name != "" {
		action.Name = t.Name
	} else {
		action.Name = t.Type
	}

	action.Follow = t.Follow

	mask := unsafe.As[unsafe.Any](&action.Data)

	switch t.Type {
	case "copy":
		action.Data = new(Copy)

		err = value.Decode(action.Data)
		mask.Type = dotdeploy.ActionTypeCopy
	case "move":
		action.Data = new(Move)

		err = value.Decode(action.Data)
		mask.Type = dotdeploy.ActionTypeMove
	case "execute":
		action.Data = new(Execute)

		err = value.Decode(action.Data)
		mask.Type = dotdeploy.ActionTypeExecute
	case "file":
		file := new(File)

		err = value.Decode(file)
		if err != nil {
			return err
		}

		f, err := os.Open(file.Path)
		if err != nil {
			return err
		}

		return yaml.NewDecoder(f).Decode(action)
	default:
		return dotdeploy.ErrUnknowActionType
	}

	return err
}

func (action *Action) Action() *dotdeploy.Action {
	return &dotdeploy.Action{
		Follow: action.Follow,

		Name: action.Name,
		Data: action.Data,
	}
}
