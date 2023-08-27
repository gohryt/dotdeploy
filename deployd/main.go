package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/bytedance/sonic"

	"github.com/gohryt/dotdeploy/contract"
	"github.com/gohryt/dotdeploy/deployd/command"
)

func serve(group *sync.WaitGroup, descriptor *net.UnixConn) {
	_eturn, err := (*contract.Return)(nil), error(nil)

	for {
		command := new(command.Command)

		err = sonic.ConfigStd.NewDecoder(descriptor).Decode(command)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}

			break
		}

		_eturn, err = command.Do()
		if err != nil {
			log.Println(err)
			break
		}

		err = sonic.ConfigFastest.NewEncoder(descriptor).Encode(_eturn)
		if err != nil {
			log.Println(err)
		}
	}

	err = descriptor.Close()
	if err != nil {
		log.Println(err)
	}

	group.Done()
}

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("this program should be started as systemd daemon")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	group := new(sync.WaitGroup)

	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: "/tmp/deployd"})
	if err != nil {
		log.Fatal(err)
	}

	listener.SetUnlinkOnClose(true)

	go func() {
		<-shutdown
		listener.Close()
	}()

	descriptor := (*net.UnixConn)(nil)

	for {
		descriptor, err = listener.AcceptUnix()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println(err)
			}

			break
		}

		group.Add(1)
		go serve(group, descriptor)
	}

	group.Wait()
}
