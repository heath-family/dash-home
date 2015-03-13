package web

import (
	"html/template"
	"log"
	"net/http"
	"os"

	httplog "github.com/99designs/goodies/http/log"
	"github.com/99designs/goodies/http/panichandler"
	"github.com/heath-family/dash-home/bom"
	"github.com/heath-family/dash-home/debugging"
)

type IndexContent struct {
	Weather  []bom.Forecast
	Calendar []interface{}
}

func Serve(input chan IndexContent) http.Handler {
	ic := <-input
	go func() {
		// Read updated content until the input channel closes, then quit
		for ic = range input {
		}
		os.Exit(0)
	}()
	return panichandler.Decorate(nil, log.New(os.Stderr, "Panic: ", log.LstdFlags), httplog.CommonLogHandler(
		nil,
		"",
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			content := ic // Copy it so it can't be modified while executing the template
			if r.URL.Path == "/" {
				t, err := template.ParseFiles("templates/index.html")
				if err != nil {
					debugging.RegisterError(err)
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				} else {
					t.Execute(rw, content)
				}
			} else {
				http.ServeFile(rw, r, "templates"+r.URL.Path)
			}
		}),
	))
}
