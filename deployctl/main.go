package main

import (
	"log"
	"net"
	"os"

	"github.com/bytedance/sonic"

	"github.com/gohryt/dotdeploy/contract"
)

func main() {
	descriptor, err := net.DialUnix("unix", nil, &net.UnixAddr{
		Name: "/tmp/deployd",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = descriptor.Write([]byte(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	_eturn := new(contract.Return)

	err = sonic.ConfigFastest.NewDecoder(descriptor).Decode(_eturn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("return: " + _eturn.Type)
}
