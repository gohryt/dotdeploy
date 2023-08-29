package script

import (
	"io"
	"log"

	"github.com/gohryt/dotdeploy/contract"
)

func Move(_opy *contract.Move, writer io.Writer) error {
	log.Println("mv")
	writer.Write([]byte("mv\n"))
	return nil
}
