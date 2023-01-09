package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jchv/go-webview2"
)

func main() {
	debug := false
	if os.Getenv("SAENUMA_DEVELOPER") == "true" {
		debug = true
	}

	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}
	os.MkdirAll(rootPath, 0777)

	port := "15267"

	go func() {
		r := mux.NewRouter()

		r.HandleFunc("/gs/{obj}", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			rawObj, err := contentStatics.ReadFile("statics/" + vars["obj"])
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Disposition", "attachment; filename="+vars["obj"])
			contentType := http.DetectContentType(rawObj)
			w.Header().Set("Content-Type", contentType)
			w.Write(rawObj)
		})

		r.HandleFunc("/xdg/", func(w http.ResponseWriter, r *http.Request) {
			if runtime.GOOS == "windows" {
				exec.Command("cmd", "/C", "start", r.FormValue("p")).Run()
			} else if runtime.GOOS == "linux" {
				exec.Command("xdg-open", r.FormValue("p")).Run()
			}
		})

		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			type Context struct {
				OutPath string
			}
			tmpl := template.Must(template.ParseFS(content, "templates/make_ref.html"))
			tmpl.Execute(w, Context{rootPath})
		})

		r.HandleFunc("/save_canvas_content", func(w http.ResponseWriter, r *http.Request) {
			base64Str := r.FormValue("imgBase64")
			newBase64Str := strings.ReplaceAll(base64Str, "data:image/png;base64,", "")
			decoded, err := base64.StdEncoding.DecodeString(newBase64Str)
			if err != nil {
				fmt.Fprint(w, err.Error())
				return
			}
			reader := bytes.NewReader(decoded)
			img, err := png.Decode(reader)
			if err != nil {
				fmt.Fprintf(w, "bad png: %s", err.Error())
				return
			}

			imgPath := filepath.Join(rootPath, UntestedRandomString(5)+".png")
			f, err := os.Create(imgPath)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer f.Close()

			err = png.Encode(f, img)
			if err != nil {
				fmt.Fprint(w, err.Error())
				return
			}

			fmt.Fprintf(w, "ok")
		})

		err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

	}()

	if runtime.GOOS == "windows" {
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
	} else if runtime.GOOS == "linux" {
		fmt.Printf("Running at http://127.0.0.1:%s\n", port)
		exec.Command("xdg-open", fmt.Sprintf("http://127.0.0.1:%s", port)).Run()
		// Wait until the interrupt signal arrives or browser window is closed
		sigc := make(chan os.Signal, 5)
		signal.Notify(sigc, os.Interrupt)
		<-sigc
		fmt.Println("Exiting")

	}
}
