package web

import (
	"html/template"
	"net/http"
	"os"

	"github.com/99designs/goodies/http/log"
	"github.com/heath-family/dash-home/bom"
)

var port = ":9000"

func init() {
	p := os.Getenv("PORT")
	if p != "" {
		port = p
	}
	if port[0] != ':' {
		port = ":" + port
	}
}

type IndexContent struct {
	Weather  []bom.Forecast
	Calendar []interface{}
}

func Serve(ic *IndexContent) {
	http.Handle("/", log.CommonLogHandler(
		nil,
		"",
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			context := *ic // Copy it so it can't be modified while executing the template
			if r.URL.Path == "/" {
				t, err := template.ParseFiles("templates/index.html")
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				} else {
					t.Execute(rw, context)
				}
			} else {
				http.ServeFile(rw, r, "templates"+r.URL.Path)
			}
		}),
	))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
