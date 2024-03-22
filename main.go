package main

import (
	"log"
	"os"
	"ps2manager/manager"
	"ps2manager/tui"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("games path not set")
	}
	if err := manager.InitManager(os.Args[1]); err != nil {
		log.Fatalf("failed to init game manager: %v\n", err)
	}
	ui := tui.Init()
	defer ui.Stop()
	if err := ui.Run(); err != nil {
		panic(err)
	}
}
