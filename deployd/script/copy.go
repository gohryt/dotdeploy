package script

import (
	"io"
	"log"

	"github.com/gohryt/dotdeploy/contract"
)

func Copy(_opy *contract.Copy, writer io.Writer) error {
	log.Println("cp")
	writer.Write([]byte("cp\n"))
	return nil
}
