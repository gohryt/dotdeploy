package main

import (
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
		From string `yaml:"from"`
		To   string `yaml:"to"`
	}

	Move struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
	}

	Run struct {
		Path    string `yaml:"path"`
		Timeout int    `yaml:"timeout"`

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
		mask.Type = unsafe.Type(new(dotdeploy.Copy))
	case "move":
		action.Data = new(Move)

		err = value.Decode(action.Data)
		mask.Type = unsafe.Type(new(dotdeploy.Move))
	case "run":
		action.Data = new(Run)

		err = value.Decode(action.Data)
		mask.Type = unsafe.Type(new(dotdeploy.Run))
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
