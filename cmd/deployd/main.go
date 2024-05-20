package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gohryt/dotdeploy/internal/deployd"
)

func main() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, os.Interrupt)

	ctx, err := deployd.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	select {
	case <-signalC:
	case <-ctx.Done():
	}
}
