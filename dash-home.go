package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/heath-family/dash-home/bom"
	"github.com/heath-family/dash-home/debugging"
	"github.com/heath-family/dash-home/web"
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

func getIndexUpdates() web.IndexContent {
	forecasts, err := bom.Latest()
	debugging.RegisterError(err)

	return web.IndexContent{
		Weather: forecasts,
	}
}

func main() {
	c := make(chan web.IndexContent)

	go func() {
		for {
			c <- getIndexUpdates()
			time.Sleep(time.Minute)
		}
	}()

	http.Handle("/", web.Serve(c))
	http.Handle("/debug", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(debugging.Errors())
	}))

	fmt.Printf("Server is listening on %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
