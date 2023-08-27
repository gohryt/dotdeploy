package command

import (
	"github.com/bytedance/sonic"

	"github.com/gohryt/dotdeploy/contract"
)

type (
	Command struct {
		Type string
		Data any
	}
)

func (command *Command) UnmarshalJSON(buffer []byte) error {
	_ype, err := sonic.Get(buffer, "type")
	if err != nil {
		return err
	}

	command.Type, err = _ype.String()
	if err != nil {
		return err
	}

	data, err := sonic.Get(buffer, "Data")
	if err != nil {
		return err
	}

	raw, err := data.Raw()
	if err != nil {
		return err
	}

	switch command.Type {
	case "copy":
		_opy := new(contract.Copy)

		err = sonic.UnmarshalString(raw, _opy)
		if err != nil {
			return err
		}

		command.Data = _opy
	}

	return nil
}

func (command *Command) Do() (_eturn *contract.Return, err error) {
	switch typed := command.Data.(type) {
	case *contract.Copy:
		return Copy(typed)
	default:
		return &contract.Return{Type: "unknown"}, nil
	}
}
