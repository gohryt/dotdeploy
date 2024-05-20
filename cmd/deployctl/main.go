package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bytedance/sonic"

	"github.com/gohryt/dotdeploy/internal/deployctl"
	"github.com/gohryt/dotdeploy/internal/models"
)

func main() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, os.Interrupt)

	ctx, err := deployctl.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	file, err := os.Open(".deploy")
	if err != nil {
		log.Fatal(err)
	}

	deploy := new(models.Deploy)

	err = sonic.ConfigDefault.NewDecoder(file).Decode(deploy)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(deploy)

	select {
	case <-signalC:
	case <-ctx.Done():
	}
}
