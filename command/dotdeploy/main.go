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

	deploy := new(dotdeploy.Deploy)

	file, err := os.Open(*flagFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(deploy)
	if err != nil {
		log.Fatal(err)
	}

	err = deploy.Work(shutdown)
	if err != nil {
		log.Fatal(err)
	}
}
