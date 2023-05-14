package main

import (
	"fmt"
	"io"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func actionDisplayCheckPage(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, "OK")
}
