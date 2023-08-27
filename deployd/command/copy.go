package command

import "github.com/gohryt/dotdeploy/contract"

func Copy(_opy *contract.Copy) (_eturn *contract.Return, err error) {
	return &contract.Return{Type: "copy"}, nil
}
