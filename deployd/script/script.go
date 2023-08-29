package script

import (
	"io"

	"github.com/gohryt/dotdeploy/contract"
)

func Script(item *contract.ItemScript, writer io.Writer) error {
	switch item.EnumScript.Type {
	case contract.ScriptCopy:
		return Copy((*contract.Copy)(item.Data), writer)
	case contract.ScriptMove:
		return Move((*contract.Move)(item.Data), writer)
	}
	return nil
}
