package main

import (
	"log"
	"os"
	"ps2manager/manager"
	"ps2manager/tui"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("config file not set!")
	}
	if err := manager.InitManager(os.Args[1]); err != nil {
		log.Fatalf("failed to init game manager: %v\n", err)
	}
	if err := tui.TUI().Run(); err != nil {
		panic(err)
	}
}
