package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jchv/go-webview2"
)

func launch() {
	debug := false
	if os.Getenv("SAENUMA_DEVELOPER") == "true" {
		debug = true
	}

	go startBackend()

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     debug,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title: "Bregana: A 3d reference image tool",
		},
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()
	w.SetSize(1400, 700, webview2.HintNone)
	w.Navigate(fmt.Sprintf("http://127.0.0.1:%s", port))
	w.Run()
}
