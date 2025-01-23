package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.Int("port", 9070, "--port 9070")
	respond := flag.String("respond", "", "--respond hello")
	flag.Parse()
	if len(*respond) == 0 {
		*respond = fmt.Sprintf("hello %d", *port)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, *respond)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
