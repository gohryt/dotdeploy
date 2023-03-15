package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/gohryt/dotdeploy"
	"gopkg.in/yaml.v3"
)

type (
	Deploy struct {
		Folder string `yaml:"folder"`
		Keep   bool   `yaml:"keep"`

		Remote []Connection `yaml:"Remote"`
		Do     []Action     `yaml:"Do"`
	}
)

var (
	flagVersion = flag.Bool("version", false, "print version")
	flagFile    = flag.String("file", ".deploy", "set .deploy filepath")

	version = "unknown"
)

func init() {
	flag.Parse()

	info, ok := debug.ReadBuildInfo()

	if ok == false {
		return
	}

	for i := range info.Settings {
		if info.Settings[i].Key == "vcs.revision" {
			version = info.Settings[i].Value
		}
	}
}

func main() {
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if *flagVersion {
		log.Println(version)
		return
	}

	deploy := new(Deploy)

	file, err := os.Open(*flagFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(deploy)
	if err != nil {
		log.Fatal(err)
	}

	err = dotdeploy.Work(shutdown, deploy.Deploy())
	if err != nil {
		log.Fatal(err)
	}
}

func (deploy *Deploy) Deploy() *dotdeploy.Deploy {
	new := &dotdeploy.Deploy{
		Folder: deploy.Folder,
		Keep:   deploy.Keep,

		Remote: make(dotdeploy.Remote, len(deploy.Remote)),
		Do:     make(dotdeploy.Do, len(deploy.Do)),
	}

	for i, connection := range deploy.Remote {
		new.Remote[i] = connection.Connection()
	}

	for i, action := range deploy.Do {
		new.Do[i] = action.Action()
	}

	return new
}
