package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

func launch() {
	go startBackend()
	fmt.Printf("Running at http://127.0.0.1:%s\n", port)
	exec.Command("xdg-open", fmt.Sprintf("http://127.0.0.1:%s", port)).Run()
	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal, 5)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
	fmt.Println("Exiting")

}
