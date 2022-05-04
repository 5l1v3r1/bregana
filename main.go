package main

import (
  "os/exec"
	"runtime"
	"os/signal"
  "fmt"
  "github.com/gorilla/mux"
	"net/http"
  "os"
  "html/template"
)


func main() {
  port := "15267"

  go func() {
    r := mux.NewRouter()

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "ok bregana")
    })

    r.HandleFunc("/gs/{obj}", func (w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			rawObj, err := contentStatics.ReadFile("statics/" + vars["obj"])
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Disposition", "attachment; filename=" + vars["obj"])
			contentType := http.DetectContentType(rawObj)
			w.Header().Set("Content-Type", contentType)
			w.Write(rawObj)
		})

    r.HandleFunc("/xdg/", func (w http.ResponseWriter, r *http.Request) {
      if runtime.GOOS == "windows" {
        exec.Command("cmd", "/C", "start", r.FormValue("p")).Run()
      } else if runtime.GOOS == "linux" {
        exec.Command("xdg-open", r.FormValue("p")).Run()
      }
    })


    r.HandleFunc("/make_ref", func(w http.ResponseWriter, r *http.Request) {
      tmpl := template.Must(template.ParseFS(content, "templates/make_ref.html"))
      tmpl.Execute(w, nil)
    })

    err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
  }()

	fmt.Printf("Running at http://127.0.0.1:%s\n", port)
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/C", "start", fmt.Sprintf("http://127.0.0.1:%s", port)).Output()
	} else if runtime.GOOS == "linux" {
		exec.Command("xdg-open", fmt.Sprintf("http://127.0.0.1:%s/make_ref", port) ).Run()
	}
	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	}
	fmt.Println("Exiting")
}
