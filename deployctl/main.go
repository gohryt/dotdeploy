package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"unsafe"

	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"

	"github.com/gohryt/dotdeploy/contract"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("command required")
	}

	err := error(nil)

	if os.Args[1] == "deployd" {
		if len(os.Args) < 3 {
			log.Fatal("no data was set")
		}

		err = deployd(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		path, file, root := ".deploy", (*os.File)(nil), new(contract.Root)

		if len(os.Args) > 2 {
			path = os.Args[1]
		}

		file, err = os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.NewDecoder(file).Decode(root)
		if err != nil {
			log.Fatal(err)
		}

		str, err := sonic.ConfigStd.MarshalToString(root)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(str)
	}
}

func deployd(json string) error {
	slice := unsafe.Slice(unsafe.StringData(json), len(json))

	if !sonic.ConfigStd.Valid(slice) {
		return errors.New("request is not valid: " + json)
	}

	descriptor, err := net.DialUnix("unix", nil, &net.UnixAddr{
		Name: "/tmp/deployd",
	})
	if err != nil {
		return err
	}

	_, err = descriptor.Write(slice)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(descriptor)

	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
	}

	return nil
}
