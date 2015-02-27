package web

import (
	"github.com/99designs/goodies/http/log"
	"github.com/danielheath/dash-home/bom"
	"html/template"
	"net/http"
	"os"
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

type IndexContext struct {
	Weather []bom.Forecast
}

func Serve() {
	http.Handle("/", log.CommonLogHandler(
		nil,
		"",
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				t, err := template.ParseFiles("templates/index.html")
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				} else {
					t.Execute(rw, IndexContext{bom.Sample})
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
